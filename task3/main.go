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

type Arguments map[string] string

type User struct {
	id string `json:"id"`
	email string `json:"email"`
	age int `json:"age"`
}

func parseArgs() Arguments {
	operation := flag.String("operation", "", "Operation")
	filename := flag.String("filename", "", "Filename")
	item := flag.String("item", "", "Item of User")
	id := flag.String("id", "", "User id")
	flag.Parse()
	return Arguments{
		"operation": *operation,
		"filename": *filename,
		"item": *item,
		"id": *id,
	}

}

func check(e error)  {
	if e != nil {
		panic(e)
	}
}


func list(arr Arguments) (users []User) {
	filename := arr["filename"]
	if filename == "" {
		check(errors.New("-fileName flag has to be specified"))
	}
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	bytes, err := ioutil.ReadAll(f)
	check(err)
	err = json.Unmarshal(bytes, &users)
	check(err)
	return
}

func add(filename string, item string)  {
	if filename == "" {
		check(errors.New("-fileName flag has to be specified"))
	}
	if item == "" {
		check(errors.New("-item flag has to be specified"))
	}
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	byte, err := f.Write([]byte(item))
	check(err)
	fmt.Println("wrote ", byte, " bytes")
}

func remove(filename string, id string)  {
	if filename == "" {
		check(errors.New("-fileName flag has to be specified"))
	}
	if id == "" {
		check(errors.New("-id flag has to be specified"))
	}
	users := list(filename)
	for i, value := range users {
		if value.id == id {
			users = append(users[:i], users[i+1:]...)
		}
	}
	bytes, err := json.Marshal(users)
	check(err)
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	byte, err := f.Write(bytes)
	check(err)
	fmt.Println("wrote ", byte, " bytes")
}

func findById(filename string, id string) (user User) {
	if filename == "" {
		check(errors.New("-fileName flag has to be specified"))
	}
	if id == "" {
		check(errors.New("-id flag has to be specified"))
	}
	users := list(filename)
	for _, value := range users {
		if value.id == id {
			user = value
		}
	}
	return
}

func Perform(args Arguments, writer io.Writer) error {
	var mapOperation map[string] func(arguments Arguments){
		"list": list(args)

	}

}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
