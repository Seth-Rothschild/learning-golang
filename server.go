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

func handle(action func(string)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		item := getValue(r, "description")
		action(item)
		http.Redirect(w, r, "/", 303)
		http.Get("/")
	}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFS(content, "index.html"))
	tpl.Execute(w, todoList)
}

func todoItemFromString(input string) TodoItem {
	ID := fmt.Sprint(len(todoList) + 1)
	trimmedString := strings.TrimLeft(input, " ")
	leadingSpaces := len(input) - len(trimmedString)
	newItem := TodoItem{ID, trimmedString, Todo, leadingSpaces * 40}
	return newItem
}

func add(item string) {
	todoList = append(todoList, todoItemFromString(item))
}

func toggle(item string) {
	for i, v := range todoList {
		if v.ID == item {
			todoList[i].markDone()
		}
	}
}

func indent(indentAmount int) func(string) {
	return func(item string) {
		for i, v := range todoList {
			if v.ID == item {
				todoList[i].indent(indentAmount)
			}
		}
	}
}

func order(upOrDown int) func(string) {
	return func(item string) {
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
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", handle(add))
	http.HandleFunc("/toggle", handle(toggle))
	http.HandleFunc("/indentRight", handle(indent(40)))
	http.HandleFunc("/indentLeft", handle(indent(-40)))
	http.HandleFunc("/moveUp", handle(order(-1)))
	http.HandleFunc("/moveDown", handle(order(1)))
	fmt.Println("Serving on 8100")
	http.ListenAndServe(":8100", nil)
}
