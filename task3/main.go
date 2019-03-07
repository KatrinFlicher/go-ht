package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string

type User struct {
	id    string `json:"id"`
	email string `json:"email"`
	age   int    `json:"age"`
}

func parseArgs() Arguments {
	operation := flag.String("operation", "", "Operation")
	filename := flag.String("filename", "", "Filename")
	item := flag.String("item", "", "Item of User")
	id := flag.String("id", "", "User id")
	flag.Parse()
	return Arguments{
		"operation": *operation,
		"filename":  *filename,
		"item":      *item,
		"id":        *id,
	}

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func list(f *os.File) []byte {
	bytes, err := ioutil.ReadAll(f)
	check(err)
	return bytes
}

func Perform(args Arguments, writer io.Writer) (err error) {
	filename := args["filename"]
	operation := args["operation"]
	if operation == "" {
		return fmt.Errorf("-operation flag has to be specified")
	}

	if filename == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	check(err)
	defer f.Close()
	switch operation {
	case "add":
		item := args["item"]
		if item == "" {
			return fmt.Errorf("-item flag has to be specified")
		}
		var user User
		err = json.Unmarshal([]byte(item), &user)
		var users []User
		err = json.Unmarshal(list(f), &users)
		for _, value := range users {
			if value.id == user.id {
				return fmt.Errorf("Item with id %s already exists", user.id)
			}
		}
		users = append(users, user)
		bytes, err := json.Marshal(users)
		check(err)
		err = ioutil.WriteFile(filename, bytes, 0644)
		check(err)
		check(err)
	case "remove":
		id := args["id"]
		if id == "" {
			return errors.New("-id flag has to be specified")
		}
		var users []User
		err = json.Unmarshal(list(f), &users)
		for i, value := range users {
			if value.id == id {
				users = append(users[:i], users[i+1:]...)
			}
		}
		bytes, err := json.Marshal(users)
		check(err)
		err = ioutil.WriteFile(filename, bytes, 0644)
		check(err)
	case "list":
		users := list(f)
		_, err := writer.Write(users)
		check(err)
	case "findById":
		id := args["id"]
		if id == "" {
			return errors.New("-id flag has to be specified")
		}
		var users []User
		err = json.Unmarshal(list(f), &users)
		var user User
		for _, value := range users {
			if value.id == id {
				user = value
			}
		}
		bytes, err := json.Marshal(user)
		check(err)
		_, err = writer.Write(bytes)
		check(err)
	default:
		return fmt.Errorf("Operation %s not allowed!", operation)
	}
	return

}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
