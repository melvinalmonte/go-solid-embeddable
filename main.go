package main

import (
	"context"
	"encoding/json"
	"errors"
	"go-solid-embed/frontend"
	"go-solid-embed/utils"
	"go.uber.org/zap"
	"strconv"

	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Todo struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Active bool   `json:"active"`
	Done   bool   `json:"done"`
}

type Todos []Todo

var todos = Todos{{1, "Learn Golang", true, false}, {2, "Learn SolidJS", true, false}}

type customFs struct {
	fs       http.FileSystem
	fallback string
}

func (c *customFs) Open(name string) (http.File, error) {
	f, err := c.fs.Open(name)

	if errors.Is(err, os.ErrNotExist) {
		return c.fs.Open(c.fallback)
	}

	return f, err
}

func WithNotFoundFallbackFileSystem(root http.FileSystem, fallback string) http.FileSystem {
	return &customFs{
		fs:       root,
		fallback: fallback,
	}
}

func FileServer(contents fs.FS) http.Handler {
	return http.FileServer(WithNotFoundFallbackFileSystem(http.FS(contents), "index.html"))
}

func handleTodos(w http.ResponseWriter, r *http.Request) {
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

func api(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode("TODO Server")
	if err != nil {
		zap.S().Errorw("failed to encode todos", "error", err)
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

func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
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

func updateTodo(w http.ResponseWriter, r *http.Request) {
	zap.S().Info("Updating todo")
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}
	zap.S().Info("Updating todo with id: ", id)
	var updatedTodo Todo
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

func main() {
	utils.InitLogger()
	zap.S().Info("Starting server...")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	contents, err := frontend.Contents()
	if err != nil {
		zap.S().Errorw("failed to find assets", "error", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", FileServer(contents))
	mux.HandleFunc("/api", api)
	mux.HandleFunc("/api/todos", handleTodos)

	server := &http.Server{
		Addr:    ":9090",
		Handler: mux,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			zap.S().Errorw("failed to shutdown server", "error", err)
		}
	}()

	zap.S().Info("Server is ready to handle requests at http://localhost:9090/")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		zap.S().Errorw("server error", "error", err)
	}
}
