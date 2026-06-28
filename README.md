# Library Management System

A comprehensive **Library Management System** built with Go, featuring a clean **Controller-Service-DAO** architecture. This system provides complete functionality for managing library operations including book inventory, student management, and borrowing workflows.

## 🏗️ Architecture

The project follows a **3-tier architecture** pattern for better maintainability and testability:

```
┌─────────────────┐
│   Controller    │  ← HTTP Request Handling
├─────────────────┤
│    Service      │  ← Business Logic
├─────────────────┤
│      DAO        │  ← Data Access Layer
├─────────────────┤
│   GORM/MySQL    │  ← Database
└─────────────────┘
```

## ✨ Features

### 👥 User Management
- **Student Registration & Authentication**
- **Librarian Registration & Authentication**
- **Password Security** with bcrypt hashing
- **User Profile Management**

### 📚 Book Management
- **CRUD Operations** for books
- **Inventory Tracking** (quantity & availability)
- **Category-based Organization**
- **Search and Pagination**

### 📖 Borrowing System
- **Book Borrowing** with availability checks
- **Book Returns** with date tracking
- **Borrowing History** for students
- **Transaction Safety** with database locks

### 👩‍💼 Librarian Features
- **Student Management** (view, update, delete)
- **Complete Book Management**
- **Borrowing Analytics**
- **System Administration**

## 🚀 Getting Started

### Prerequisites
- **Go 1.21+**
- **MySQL** (or SQLite for testing)
- **Git**

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd library-management-system
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up database**
   - Configure your MySQL database
   - Update connection string in `utils/utils.go`

4. **Run the application**
   ```bash
   go run .
   ```

The server will start on `http://localhost:8080`

## 📖 API Documentation

### Student Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/students/register` | Register new student |
| `POST` | `/students/login` | Student login |
| `GET` | `/students/:id/history` | Get borrowing history |

### Librarian Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/librarians/register` | Register new librarian |
| `POST` | `/librarians/login` | Librarian login |
| `GET` | `/librarians/students` | Get all students |
| `GET` | `/librarians/students/:id` | Get student details |
| `PUT` | `/librarians/students/:id` | Update student info |
| `DELETE` | `/librarians/students/:id` | Delete student |

### Book Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/books` | Get all books (paginated) |
| `GET` | `/books/:id` | Get book by ID |
| `POST` | `/books` | Add new book (Librarian) |
| `PUT` | `/books/:id` | Update book (Librarian) |
| `DELETE` | `/books/:id` | Delete book (Librarian) |
| `POST` | `/books/borrow` | Borrow a book |
| `POST` | `/books/return` | Return a book |

### Example Requests

**Student Registration:**
```bash
curl -X POST http://localhost:8080/students/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com", 
    "phone": "1234567890",
    "password": "secure123"
  }'
```

**Borrow Book:**
```bash
curl -X POST http://localhost:8080/books/borrow \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 1,
    "book_id": 1
  }'
```

## 🧪 Testing

The project includes comprehensive tests covering all layers:

**Run all tests:**
```bash
go test -v ./...
```

**Run specific tests:**
```bash
go test -v -run TestDatabaseSetup
go test -v -run TestServiceLayer  
go test -v -run TestArchitectureLayers
```

**Test Coverage:**
- ✅ Database layer (DAO) tests
- ✅ Business logic (Service) tests  
- ✅ Architecture integration tests
- ✅ Error handling tests

## 📁 Project Structure

```
library-management-system/
├── controller/          # HTTP request handlers
│   ├── book_controller.go
│   ├── student_controller.go
│   └── librarian_controller.go
├── service/            # Business logic layer
│   ├── book_service.go
│   ├── student_service.go
│   └── librarian_service.go
├── dao/               # Data access objects
│   ├── book_dao.go
│   ├── student_dao.go
│   ├── librarian_dao.go
│   └── borrowed_book_dao.go
├── model/             # Data models
│   ├── book.go
│   ├── student.go
│   ├── librarian.go
│   ├── borrowed_book.go
│   ├── borrow_request.go
│   └── login.go
├── boot/              # Application bootstrap
│   └── boot.go
├── constants/         # Application constants
│   └── constants.go
├── utils/             # Utility functions
│   └── utils.go
├── main.go           # Application entry point
├── go.mod            # Go module dependencies
└── README.md         # Project documentation
```

## 🛡️ Security Features

- **Password Hashing**: Bcrypt for secure password storage
- **SQL Injection Prevention**: GORM ORM protection
- **Input Validation**: Request validation at controller level
- **Error Handling**: Secure error messages (no data leakage)

## 🔧 Technologies Used

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: MySQL (SQLite for testing)
- **Authentication**: bcrypt
- **Testing**: Testify
- **Architecture**: Controller-Service-DAO

## 📈 Performance Features

- **Database Connection Pooling**: Efficient connection management
- **Transaction Safety**: ACID compliance for critical operations
- **Pagination Support**: Efficient data retrieval
- **Row-level Locking**: Concurrent borrowing safety

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request


## 🚀 Deployment

For production deployment:

1. **Build the application**
   ```bash
   go build -o library-management-system
   ```

2. **Set environment variables**
   ```bash
   export DB_HOST=your-mysql-host
   export DB_PORT=3306
   export DB_NAME=library_db
   export DB_USER=your-username  
   export DB_PASS=your-password
   ```

3. **Run the binary**
   ```bash
   ./library-management-system
   ```

