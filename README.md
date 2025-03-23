# 📚 School Management API

![School API Logo](https://miro.medium.com/v2/resize:fit:2000/format:webp/1*OcmVkcsM5BWRHrg8GC17iw.png)

A simple RESTful API for managing students in a school database using **Go (Golang)**, **PostgreSQL**, and **Gorilla Mux**.

## 🚀 Features

- **🔄 CRUD Operations:** Create, Read, Update, Delete students.
- **📜 Pagination:** Efficiently handle large datasets.
- **📦 JSON-based API Responses:** Standardized data format for easy consumption.
- **✅ Input Validation:** Ensure data integrity before processing.
- **⚡ Bulk Insert:** Efficiently insert multiple records in one request.
- **🔒 Environment Variables:** Securely manage database connection details.
- **🗄️ PostgreSQL Database Connection:** Persistent data storage with PostgreSQL.
- **🛑 Graceful Shutdown:** Ensure smooth termination of the API.

---

## 🏗️ Tech Stack

- **Backend:** Go (Golang), Gorilla Mux
- **Database:** PostgreSQL
- **Logging:** Log Package
- **API Format:** RESTful, JSON

---

## 📂 Project Structure

```
handlers
  ├── handlers.go     # API Endpoints
models
  ├── models.go       # Student struct
validation
  ├── validation.go   # Input Validation

main.go               # Entry Point
.env                  # Environment Variables
README.md             # Project Documentation
LICENSE               # License Information
go.mod                # Module Dependencies
go.sum                # Dependency Checksums
```

---

## 📦 Installation & Setup

### Prerequisites
- Install **Go (v1.24 or latest)**
- Install and set up **PostgreSQL**

### 1️⃣ Clone the Repository
```sh
git clone https://github.com/nsshohag/school_api_postgres.git
cd school_api_postgres
```

### 2️⃣ Configure Environment Variables
Create a `.env` file in the root directory is exists then add your PostgreSQL credentials:
```env
DB_HOST=localhost
DB_PORT=5433
DB_USER=sadat
DB_PASSWORD=11235813
DB_NAME=school_db
```

### 3️⃣ Install Dependencies
```sh
go mod tidy
```

### 4️⃣ Run the Application
```sh
go run main.go
```

The server will start at `http://localhost:8080`

---

## 🛠️ Database Setup

Run the following SQL query in your PostgreSQL database:
```sql
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    age INT NOT NULL,
    class INTEGER NOT NULL CHECK (class BETWEEN 1 AND 10)
);
```

---

## 📖 API Endpoints

### Student Routes

| Method | Endpoint                       | Description                   |
|--------|--------------------------------|-------------------------------|
| GET    | `/api/v1/students`            | Get All Students              |
| POST   | `/api/v1/students`            | Create Student                |
| GET    | `/api/v1/students/{id}`       | Get Student by ID             |
| PUT    | `/api/v1/students/{id}`       | Update Student                |
| PATCH  | `/api/v1/students/{id}`       | Patch Student                 |
| DELETE | `/api/v1/students/{id}`       | Delete Student                |
| POST   | `/api/v1/students/bulk`       | Bulk Insert Students          |

### 🔍 Get All Students
```http
GET api/v1/students
```
**Response:**
```json
[
  {
    "id": 1,
    "name": "Nazmus Sadat Shohag",
    "age": 24,
    "class": 10
  },
    {
    "id": 2,
    "name": "SH Rony",
    "age": 24,
    "class": 10
  }
]
```

### 📌 Get Student by ID
```http
GET api/v1/students/{id}
```

### ➕ Create Student
```http
POST api/v1/students
```
**Request Body:**
```json
{
  "name": "Preity",
  "age": 24,
  "class": 9
}
```

### ✏️ Update Student
```http
PUT api/v1/students/{id}
```
**Request Body:**
```json
{
  "name": "Preety",
  "age": 25,
  "class": 10
}
```


### 🔄 Patch Student
```http
PATCH api/v1/students/{id}
```
**Request Body:**
```json
{
  "age": 26
}
```

### 🗑️ Delete Student
```http
DELETE api/v1/students/{id}
```

### 🔄 Bulk Insert Students
```http
POST /api/v1/students/bulk
```
**Request Body:**
```json
[
  { "name": "Abir", "age": 10, "class": 4 },
  { "name": "Dristy", "age": 9, "class": 3 }
]
```

---
<!-- 
## 📸 Screenshots

![API Example](https://via.placeholder.com/800x400?text=API+Example)

---
-->

## 📜 License

MIT License. See `LICENSE` for more details.

---

## ⭐ Contributing

Pull requests are welcome! For major changes, please open an issue first.

---

## 🏆 Author

Developed by **Nazmus Sadat Shohag**

🔗 Connect with me on [LinkedIn](https://www.linkedin.com/in/nazmus-sadat-492bba291/)