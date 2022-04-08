package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"
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
	ID     string
	Title  string
	Status status
	Depth  int
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

func (t *TodoItem) indent(i int) {
	t.Depth += i
	if t.Depth < 0 {
		t.Depth = 0
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

func parseString(input string) TodoItem {
	ID := string(len(todoList)+1)
	trimmedString := strings.TrimLeft(input, " ")
	leadingSpaces := len(input) - len(trimmedString)
	newItem := TodoItem{ID, trimmedString, Todo, leadingSpaces * 40}
	return newItem
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFS(content, "index.html"))
	tpl.Execute(w, todoList)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	item := getValue(r, "description")
	todoList = append(todoList, parseString(item))
	http.Redirect(w, r, "/", 303)
	http.Get("/")
}

func toggleHandler(w http.ResponseWriter, r *http.Request) {
	item := getValue(r, "description")
	for i, v := range todoList {
		if v.ID == item {
			todoList[i].markDone()
		}
	}
	http.Redirect(w, r, "/", 303)
	http.Get("/")
}

func indent(indentAmount int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		item := getValue(r, "description")
		for i, v := range todoList {
			if v.ID == item {
				todoList[i].indent(indentAmount)
			}
		}
		http.Redirect(w, r, "/", 303)
		http.Get("/")
	}
}

func order(upOrDown int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		item := getValue(r, "description")
		todoListCopy := make([]TodoItem, len(todoList))
		copy(todoListCopy, todoList)
		for i, v := range todoListCopy {
			if v.ID == item {
				targetIndex := i + upOrDown

				if targetIndex > len(todoList)-1 {
					targetIndex = len(todoList) - 1
				} else if targetIndex < 0 {
					targetIndex = 0
				}
				todoList[targetIndex] = v
				todoList[i] = todoListCopy[targetIndex]

			}

		}
		http.Redirect(w, r, "/", 303)
		http.Get("/")
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/toggle", toggleHandler)
	http.HandleFunc("/indentRight", indent(40))
	http.HandleFunc("/indentLeft", indent(-40))
	http.HandleFunc("/moveUp", order(-1))
	http.HandleFunc("/moveDown", order(1))
	fmt.Println("Serving on 8100")
	http.ListenAndServe(":8100", nil)
}
