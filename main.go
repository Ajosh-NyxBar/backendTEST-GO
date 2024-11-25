package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var tasks []Task
var nextID = 1

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id:[0-9]+}", updateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id:[0-9]+}", deleteTask).Methods("DELETE")


	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	fmt.Println("Server starting on port " + port + "...")
	http.ListenAndServe(":"+port, router)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	newTask.ID = nextID
	nextID++
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var updatedTask Task
	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i] = updatedTask
			tasks[i].ID = id
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully"})
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}
