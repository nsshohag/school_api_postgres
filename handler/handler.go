package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "sadat"
	DB_PASSWORD = "11235813"
	DB_NAME     = "school_db"
)

var db *sql.DB

func initDB() {

	var err error
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=localhost port=5433 sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)

	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database:", err)
	}
	fmt.Println("Successfully connected to the school_db database!")
}

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Class int    `json:"class"`
}

var wg = sync.WaitGroup{}

func Handle() {

	fmt.Println("Cool Cool")

	initDB()

	r := mux.NewRouter()
	r.HandleFunc("/students", getStudents_All).Methods("GET")
	r.HandleFunc("/students", createStudent).Methods("POST")
	r.HandleFunc("/students/{id}", updateStudent).Methods("PUT")
	r.HandleFunc("/students/{id}", deleteStudent).Methods("DELETE")

	wg.Add(1)
	go func() {
		log.Println("Server running on port 8080...")
		if err := http.ListenAndServe(":8080", r); err != nil {
			log.Fatal("Server failed to start:", err)
		}
		wg.Done()
	}()

	fmt.Println("Server Started")
	wg.Wait()
}

func getStudents_All(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query("SELECT * FROM students")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.ID, &s.Name, &s.Age, &s.Class); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		students = append(students, s)
	}

	// Log the JSON response before sending it to the client
	jsonResponse, err := json.Marshal(students)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("All students list will be sent\n", string(jsonResponse))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)

}

func createStudent(w http.ResponseWriter, r *http.Request) {
	var s Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow("INSERT INTO students (name, age, class) VALUES ($1, $2, $3) RETURNING id", s.Name, s.Age, s.Class).Scan(&s.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the JSON response before sending it to the client
	jsonResponse, err := json.Marshal(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Student created is\n", string(jsonResponse))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s) //w.Write(jsonResponse)

}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var s Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE students SET name=$1, age=$2, class=$3 WHERE id=$4", s.Name, s.Age, s.Class, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the JSON response before sending it to the client
	jsonResponse, err := json.Marshal(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(id, "id is updated with\n", string(jsonResponse))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.Exec("DELETE FROM students WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println(id, "id is deleted")
	w.WriteHeader(http.StatusNoContent)
}
