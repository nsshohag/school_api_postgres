# 📚 School Management API

![School API Logo](https://miro.medium.com/v2/resize:fit:2000/format:webp/1*OcmVkcsM5BWRHrg8GC17iw.png)

A simple RESTful API for managing students in a school database using **Go (Golang)**, **PostgreSQL**, and **Go chi**.

## 🚀 Features

- **🔄 CRUD Operations:** Create, Read, Update, Delete students.
- **📜 Pagination:** Efficiently handle large datasets.
- **📦 JSON-based API Responses:** Standardized data format for easy consumption.
- **✅ Input Validation:** Ensure data integrity before processing.
- **⚡ Bulk Insert:** Efficiently insert multiple records in one request.
- **🔒 Environment Variables:** Securely manage database connection details.
- **🗄️ PostgreSQL Database Connection:** Persistent data storage with PostgreSQL.
- **🚯 IP-Based Rate Limiting:** Prevent abuse by limiting requests per IP.
- **🛑 Graceful Shutdown:** Ensure smooth termination of the API.
- **🐳 Docker Support:** Easily deploy the application using Docker.
- **☸️ Kubernetes Deployment:** Scalable and manageable deployment in Kubernetes.
- **🌐 Docker Hub Integration:** Uploaded images are available on Docker Hub.

---

## 🏗️ Tech Stack

- **Backend:** Go (Golang), GO chi
- **Database:** PostgreSQL
- **Containerization:** Docker
- **Orchestration:** Kubernetes
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
dockerfile            # Dockerfile for building the image
docker-compose.yaml   # Docker Compose file for multi-container setup
dockerfile_web-server:1.0 
dockerfile_web-server-ubuntu:1.0
kubernetes/           # Kubernetes configuration files
  ├── postgres.yaml
  ├── postgres-config.yaml
  ├── postgres-secret.yaml
  ├── server.yaml

```

---

## 📦 Installation & Setup

### Prerequisites
- Install **Go (v1.24 or latest)**
- Install and set up **PostgreSQL**
- Install **Docker** and **Kubernetes**

### 1️⃣ Clone the Repository
```sh
git clone https://github.com/nsshohag/school_api_postgres.git
cd school_api_postgres
```

### 2️⃣ Configure Environment Variables
If you don't use docker-compose or kubernetes then, create a `.env` file in the root directory add your PostgreSQL credentials:
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

### 6️⃣ Run Docker Compose
For multi-container setup, you can use Docker Compose:
```sh
docker-compose up -d
```

---

### 7️⃣ Deploy to Kubernetes
If you have your Kubernetes cluster set up, you can deploy using:
```sh
kubectl apply -f kubernetes/postgres-config.yaml
kubectl apply -f kubernetes/postgres-secret.yaml
kubectl apply -f kubernetes/postgres.yaml
kubectl apply -f kubernetes/server.yaml
```

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

## 🔍 Additional Features

### 🗄️ PostgreSQL Database Connection
The API connects to a **PostgreSQL** database for persistent data storage. Connection parameters are managed through environment variables to ensure security and flexibility.

### 🚯 IP-Based Rate Limiting
To prevent abuse, the API enforces **IP-based rate limiting**, restricting excessive requests from the same IP within a specific timeframe.

### 🛑 Graceful Shutdown
The API supports **graceful shutdown**, ensuring proper cleanup of resources when the server is stopped, preventing issues like lingering database connections. When the server gets a shutdown request, it finishes ongoing requests for a specific time, and during that time, it does not take any new requests.

### 📜 Pagination
The API implements **pagination** for efficient handling of large datasets, allowing clients to request data in smaller chunks for better performance.

### ⚡ Bulk Insert with Batching
Bulk insert allows efficiently adding multiple records in a single request, reducing the number of database transactions and improving performance. Also did batching so that query does not exceed.

### 🌐 Docker Hub Integration
Images for the API and web server have been uploaded to Docker Hub, allowing for easy deployment and version management.

### ☸️ Kubernetes Configuration
The project includes Kubernetes configuration files for deploying the application in a cluster, ensuring scalability and manageability.

---

## 📜 License

MIT License. See `LICENSE` for more details.

---

## ⭐ Contributing

Pull requests are welcome! For major changes, please open an issue first.

---

## 🏆 Author

Developed by **Nazmus Sadat Shohag**

🔗 Connect with me on [LinkedIn](https://www.linkedin.com/in/nazmus-sadat-492bba291/)