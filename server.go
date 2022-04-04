package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed index.html
var content embed.FS

type status int

const (
	TODO status = 0
	DONE status = 1
	WAIT status = 2
)

type Todo struct {
	Description string
	Done        status
}

func (t *Todo) mark_done() {
	switch t.Done {
	case TODO:
		t.Done = DONE
	case DONE:
		t.Done = WAIT
	case WAIT:
		t.Done = TODO
	}
}

var todos []Todo

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFS(content, "index.html"))
	tpl.Execute(w, todos)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	item := r.FormValue("item")
	todo := Todo{item, TODO}
	todos = append(todos, todo)
	http.Redirect(w, r, "/", 303)
	http.Get("/")
}

func toggleHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	item := r.FormValue("todoitem")
	for i, v := range todos {
		if v.Description == item {
			todos[i].mark_done()
		}
	}
	http.Redirect(w, r, "/", 303)
	http.Get("/")
}

func main() {
	http.FS(content)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/toggle", toggleHandler)
	fmt.Println("Serving on 8100")
	http.ListenAndServe(":8100", nil)
}
