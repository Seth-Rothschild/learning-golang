package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed index.html
var content embed.FS

type status string

const (
	Todo status = "TODO"
	Done status = "DONE"
	Wait status = "WAIT"
)

type TodoItem struct {
	Description string
	Status      status
}

func (t *TodoItem) markDone() {
	switch t.Status {
	case Todo:
		t.Status = Done
	case Done:
		t.Status = Wait
	case Wait:
		t.Status = Todo
	}
}

var todoList []TodoItem

func getValue(r *http.Request, value string) string {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	return r.FormValue(value)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFS(content, "index.html"))
	tpl.Execute(w, todoList)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	item := getValue(r, "description")
	newItem := TodoItem{item, Todo}
	todoList = append(todoList, newItem)
	http.Redirect(w, r, "/", 303)
	http.Get("/")
}

func toggleHandler(w http.ResponseWriter, r *http.Request) {
	item := getValue(r, "description")
	for i, v := range todoList {
		if v.Description == item {
			todoList[i].markDone()
		}
	}
	http.Redirect(w, r, "/", 303)
	http.Get("/")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/toggle", toggleHandler)
	fmt.Println("Serving on 8100")
	http.ListenAndServe(":8100", nil)
}
