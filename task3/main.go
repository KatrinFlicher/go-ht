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
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func parseArgs() Arguments {
	operation := flag.String("operation", "", "Operation")
	filename := flag.String("fileName", "", "Filename")
	item := flag.String("item", "", "Item of User")
	id := flag.String("id", "", "User id")
	flag.Parse()
	return Arguments{
		"operation": *operation,
		"fileName":  *filename,
		"item":      *item,
		"id":        *id,
	}
}

func check(e error) (err error) {
	if e != nil {
		err = e
	}
	return
}

func list(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}

func Perform(args Arguments, writer io.Writer) (error error) {
	operation := args["operation"]
	if operation == "" {
		return fmt.Errorf("-operation flag has to be specified")
	}
	filename := args["fileName"]
	if filename == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}
	switch operation {
	case "add":
		item := args["item"]
		if item == "" {
			return fmt.Errorf("-item flag has to be specified")
		}
		var user User
		error = json.Unmarshal([]byte(item), &user)
		var users []User
		bytes, err := list(filename)
		error = check(err)
		error = json.Unmarshal(bytes, &users)
		sameId := false
		for _, value := range users {
			if value.Id == user.Id {
				errStr := "Item with id " + user.Id + " already exists"
				writer.Write([]byte(errStr))
				//errStr :=
				//_, err := writer.Write([]byte(errStr))
				//error = check(err)
				sameId = true
			}
		}
		if !sameId {
			users = append(users, user)
			bytes, err := json.Marshal(users)
			error = check(err)
			ioutil.WriteFile(filename, bytes, 064)
		}
	case "remove":
		id := args["id"]
		if id == "" {
			return errors.New("-id flag has to be specified")
		}
		var users []User
		bytesUser, err := list(filename)
		error = check(err)
		err = json.Unmarshal(bytesUser, &users)
		var newUser []User
		for i, value := range users {
			if value.Id == id {
				newUser = append(users[:i], users[i+1:]...)
			}
		}
		if len(newUser) == 0 {
			errStr := "Item with id " + id + " not found"
			writer.Write([]byte(errStr))
		}
		bytes, err := json.Marshal(newUser)
		error = check(err)
		return ioutil.WriteFile(filename, bytes, 0)
	case "list":
		users, err2 := list(filename)
		error = check(err2)
		_, err := writer.Write(users)
		error = check(err)
	case "findById":
		id := args["id"]
		if id == "" {
			return errors.New("-id flag has to be specified")
		}
		var users []User
		bytes, err := list(filename)
		error = check(err)
		error = json.Unmarshal(bytes, &users)
		var user User
		for _, value := range users {
			if value.Id == id {
				user = value
				break
			}
		}
		userBytes := []byte("")
		if user.Id != "" {
			bytes, err := json.Marshal(user)
			error = check(err)
			userBytes = bytes
		}
		_, err = writer.Write(userBytes)
		error = check(err)
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
