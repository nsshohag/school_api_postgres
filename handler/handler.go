package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" //
)

var db *sql.DB

func initDB() {

	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}

	// loading from .env
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT)

	var err error
	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	errPing := db.Ping()
	if errPing != nil {
		log.Fatal("Failed to ping the database:", errPing)
	}
	fmt.Println("Successfully connected to the school_db database!")
}

// Student ... that holds data and key = `json:"id"` and so on if not provided then it would be ID
type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Class int    `json:"class"`
}

var wg = sync.WaitGroup{}

// Handle ...
func Handle() {

	fmt.Println("Server initialization starting...")
	initDB()

	// routes
	r := mux.NewRouter()
	r.HandleFunc("/students", getStudentsAll).Methods("GET")
	r.HandleFunc("/students", createStudentSingle).Methods("POST")
	r.HandleFunc("/students/bulk", createStudentBulk).Methods("POST")
	r.HandleFunc("/students/{id}", updateStudent).Methods("PUT")
	r.HandleFunc("/students/{id}", deleteStudent).Methods("DELETE")
	r.HandleFunc("/students/{id}", patchStudent).Methods("PATCH")
	r.HandleFunc("/students/{id}", getStudentOne).Methods("GET")

	// starting server on a port
	wg.Add(1)
	go func() {
		log.Println("Server running on port 8080...")
		if err := http.ListenAndServe(":8080", r); err != nil {
			log.Fatal("Server failed to start:", err)
		}
		wg.Done()

	}()

	fmt.Println("This waitgroup and goroutine is used to excute code after ListenAndServe as it is blcoking and falls into infinite loop for taking request")
	wg.Wait()
}

// Fetch multiple rows --- db.Query
// Fetch a single row --- db.QueryRow
// Insert, update, delete --- db.Exec

// GET --get all students
func getStudentsAll(w http.ResponseWriter, r *http.Request) {

	// Get page and limit from query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Set default values
	page := 1
	limit := 3

	// Parse page number
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Parse limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Calculate offset
	offset := (page - 1) * limit

	rows, err := db.Query("SELECT * FROM students LIMIT $1 OFFSET $2", limit, offset)
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(students)

}

// POST --insert a student

func createStudentSingle(w http.ResponseWriter, r *http.Request) {

	var s Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch multiple rows --- db.Query
	// Fetch a single row --- db.QueryRow
	// Insert, update, delete --- db.Exec

	// kintu eitate kaj korche karon database supports RETURNING (e.g., PostgreSQL
	// but eita ar concern na karon this time I have modified the code so that multiple students can be inserted
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
	w.WriteHeader(http.StatusCreated)
	//json.NewEncoder(w).Encode(s) //w.Write(jsonResponse)

}

// func createStudent(w http.ResponseWriter, r *http.Request) {

// 	/*
// 		var students []Student
// 			// zodi multiple object pathay [[{}]
// 			err := json.NewDecoder(r.Body).Decode(&students)
// 			if err != nil {
// 				// {} - if decoding as an array fails tahole decode kortechias a single student object
// 				var singleStudent Student
// 				if err := json.NewDecoder(r.Body).Decode(&singleStudent); err != nil {
// 					http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 					return
// 				}
// 				students = append(students, singleStudent) // Convert single student to array format
// 			}
// 	*/

// 	//Read the request body once and store it because the previous won't work beacuse
// 	//When you attempt to decode the request body the first time (as an array)
// 	//it consumes the body, and the second attempt to decode it (as a single object) fails
// 	//because the body is already empty.

// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Failed to read request body", http.StatusBadRequest)
// 		return
// 	}
// 	defer r.Body.Close()

// 	var students []Student

// 	// Try to decode the stored body as an array of students
// 	err = json.Unmarshal(body, &students)
// 	if err != nil {
// 		// If decoding as an array fails, try to decode as a single student object
// 		var singleStudent Student
// 		if err := json.Unmarshal(body, &singleStudent); err != nil {
// 			http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 			return
// 		}
// 		students = append(students, singleStudent) // Convert single student to array format
// 	}

// 	if len(students) == 0 {
// 		http.Error(w, "No student data provided", http.StatusBadRequest)
// 		return
// 	}

// 	tx, err := db.Begin()
// 	if err != nil {
// 		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
// 		return
// 	}
// 	defer tx.Rollback()

// 	stmt, err := tx.Prepare("INSERT INTO students (name, age, class) VALUES ($1, $2, $3) RETURNING id")
// 	if err != nil {
// 		http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
// 		return
// 	}
// 	defer stmt.Close()

// 	var insertedStudents []Student
// 	for _, s := range students {
// 		var id int
// 		err := stmt.QueryRow(s.Name, s.Age, s.Class).Scan(&id)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Failed to insert student: %v", err), http.StatusInternalServerError)
// 			return
// 		}
// 		s.ID = id
// 		insertedStudents = append(insertedStudents, s)
// 	}

// 	if err := tx.Commit(); err != nil {
// 		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
// 		return
// 	}

// 	jsonResponse, err := json.Marshal(insertedStudents)
// 	if err != nil {
// 		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
// 		return
// 	}
// 	log.Println("Inserted students:\n", string(jsonResponse))

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	w.Write(jsonResponse)

// }

// Bulk with batch because query size matters
func createStudentBulk(w http.ResponseWriter, r *http.Request) {

	// Read the request body once and store it
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var students []Student

	err = json.Unmarshal(body, &students)
	if err != nil {
		var singleStudent Student
		if err := json.Unmarshal(body, &singleStudent); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		// Convert single student to array format
		students = append(students, singleStudent)
	}

	if len(students) == 0 {
		http.Error(w, "No student data provided", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // Ensure the transaction is rolled back if anything goes wrong

	batchSize := 3 // batch size dietchi to bulk insert
	var insertedStudents []Student

	// Insert students in batches
	for i := 0; i < len(students); i += batchSize {

		end := i + batchSize
		//
		if end > len(students) { // this is for last batch
			end = len(students)
		}

		batch := students[i:end] // Get the current batch of students

		// Build the VALUES clause dynamically for the current batch
		var valueStrings []string
		var valueArgs []interface{}

		for j, student := range batch {

			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", j*3+1, j*3+2, j*3+3)) // Create placeholders for each student example = ($1, $2, $3), ($4, $5, $6), )
			valueArgs = append(valueArgs, student.Name, student.Age, student.Class)                  // Append the student data to the valueArgs slice
		}

		// Combine the query and VALUES clause
		query := `
			INSERT INTO students (name, age, class)
			VALUES %s
			RETURNING id
		`
		finalQuery := fmt.Sprintf(query, strings.Join(valueStrings, ", "))

		// Execute the bulk insert query for the current batch
		rows, err := tx.Query(finalQuery, valueArgs...)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to execute bulk insert: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		j := 0 // batch
		// Collect the inserted IDs for the current batch
		for rows.Next() {
			var id int
			if err := rows.Scan(&id); err != nil {
				http.Error(w, fmt.Sprintf("Failed to scan inserted ID: %v", err), http.StatusInternalServerError)
				return
			}

			insertedStudents = append(insertedStudents, Student{ID: id,
				Name:  batch[j].Name,
				Age:   batch[j].Age,
				Class: batch[j].Class,
			}) // Append the inserted student to the insertedStudents slice
			j++
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// ekhane amra sudhu resposne dekhbo but pathabo na response to client
	// zodio name dekhabe na sudhu id dekhabe karon return kori nai
	jsonResponse, err := json.Marshal(insertedStudents)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	log.Println("Inserted students:\n", string(jsonResponse))

	// this not necessary in the case of insertion we dont need to send the result to the user
	// Send the response to the client
	// w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// w.Write(jsonResponse)
}

// PUT --update all the information of a student
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

	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(s)
}

// DELETE --delete a student from the database
func deleteStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := db.Exec("DELETE FROM students WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		// If no rows were affected, the ID does not exist
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	log.Println(id, "id is deleted")
	w.WriteHeader(http.StatusNoContent)
}

// PATCH --update partial information of a student
func patchStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updates []string
	var values []interface{}
	values = append(values, id)

	var s Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Build query dynamically based on provided fields
	if s.Name != "" {
		updates = append(updates, "name=$"+fmt.Sprint(len(values)+1))
		values = append(values, s.Name)
	}
	if s.Age != 0 {
		updates = append(updates, "age=$"+fmt.Sprint(len(values)+1))
		values = append(values, s.Age)
	}
	if s.Class != 0 {
		updates = append(updates, "class=$"+fmt.Sprint(len(values)+1))
		values = append(values, s.Class)
	}

	if len(updates) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf("UPDATE students SET %s WHERE id=$1",
		stringJoin(updates, ", "))

	_, err := db.Exec(query, values...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Student with id %s is updated.", id)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
}

// helper function to join strings
func stringJoin(arr []string, sep string) string {
	result := ""
	for i, s := range arr {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

// GET --get information of a single student
func getStudentOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var s Student
	err := db.QueryRow("SELECT * FROM students WHERE id=$1", id).Scan(&s.ID, &s.Name, &s.Age, &s.Class)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Log the JSON response before sending it to the client
	jsonResponse, err := json.Marshal(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Student details sent for ID", id, "\n", string(jsonResponse))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
}
