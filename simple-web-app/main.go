package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type ToDoList struct {
	// count number of todos
	ToDoCount int
	// array to hold to dos from a file
	ToDos []string
}

func errorCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Pass a string to write a message and display
func write(writer http.ResponseWriter, msg string) {
	_, err := writer.Write([]byte(msg))
	errorCheck(err)
}

// Data fetch from a file
func getStrings(fileName string) []string {
	var lines []string
	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return nil
	}
	errorCheck(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	errorCheck(scanner.Err())
	return lines
}

// Pass an english hello to a http reponse to print in a page
func englishHandler(writer http.ResponseWriter, request *http.Request) {
	write(writer, "Hello Internet")
}

// Pass a polish hello to a http reponse to print in a page
func polishHandler(writer http.ResponseWriter, request *http.Request) {
	write(writer, "Siema Internet")
}

// Handler to fetch data, map to object and display on a page
func interactHandler(writer http.ResponseWriter, request *http.Request) {

	// Data fetch [this could be switched to a db persistance]
	todoVals := getStrings("todos.txt")

	// Log retrieved data to console
	fmt.Printf("%#v\n", todoVals)

	tmpl, err := template.ParseFiles("view.html")
	errorCheck(err)

	todos := ToDoList{
		ToDoCount: len(todoVals),
		ToDos:     todoVals,
	}
	err = tmpl.Execute(writer, todos)

}

// New todo route
func newHandler(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles("new.html")
	errorCheck(err)
	err = tmpl.Execute(writer, nil)
}

// Handling POST request from a new todo view
func createHandler(writer http.ResponseWriter, request *http.Request) {

	// Value passed from the form
	todo := request.FormValue("todo")

	// Options for working with the data file
	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile("todos.txt", options, os.FileMode(0600))
	errorCheck(err)

	_, err = fmt.Fprintln(file, todo)
	errorCheck(err)

	err = file.Close()
	http.Redirect(writer, request, "/interact", http.StatusFound)
}

func main() {
	http.HandleFunc("/hello", englishHandler)
	http.HandleFunc("/siema", polishHandler)
	http.HandleFunc("/interact", interactHandler)
	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/create", createHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}
