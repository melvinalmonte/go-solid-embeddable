package handlers

import (
	"encoding/json"
	"go-solid-embed/models"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Todos []models.Todo

var todos = Todos{{1, "Learn Golang", true, false}, {2, "Learn SolidJS", true, false}}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	todo.ID = len(todos) + 1
	todos = append(todos, todo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(todo)
	if err != nil {
		zap.S().Errorw("failed to encode todo", "error", err)
		return
	}
}

func getTodos(w http.ResponseWriter) {
	zap.S().Info("Getting list of todos")
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(todos)
	if err != nil {
		zap.S().Errorw("failed to encode todos", "error", err)
		return
	}
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	zap.S().Info("Updating todo")
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}
	zap.S().Info("Updating todo with id: ", id)
	var updatedTodo models.Todo
	err = json.NewDecoder(r.Body).Decode(&updatedTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Text = updatedTodo.Text
			todos[i].Done = updatedTodo.Done
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	http.Error(w, "Todo not found", http.StatusNotFound)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}
	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	http.Error(w, "Todo not found", http.StatusNotFound)
}
func HandleTodos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getTodos(w)
	case "POST":
		createTodo(w, r)
	case "PUT":
		updateTodo(w, r)
	case "DELETE":
		deleteTodo(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
