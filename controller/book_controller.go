package controller

import (
	"net/http"
	"strconv"

	"library-management-system/constants"
	"library-management-system/model"
	"library-management-system/service"

	"github.com/gin-gonic/gin"
)

// BookController handles HTTP requests for book operations
type BookController struct {
	bookService service.BookService
}

// NewBookController creates a new BookController instance
func NewBookController(bookService service.BookService) *BookController {
	return &BookController{
		bookService: bookService,
	}
}

// GetAllBooks handles GET /books - retrieves paginated list of books
func (ctrl *BookController) GetAllBooks(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 {
		limit = 20
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	books, count, err := ctrl.bookService.GetAllBooks(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Books retrieved successfully",
		"books":   books,
		"count":   count,
	})
}

// GetBookByID handles GET /books/:id - retrieves a single book by ID
func (ctrl *BookController) GetBookByID(c *gin.Context) {
	bookID := c.Param("id")
	id, err := strconv.Atoi(bookID)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	book, err := ctrl.bookService.GetBookByID(uint(id))
	if err != nil {
		if err.Error() == constants.BookNotFoundError {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book retrieved successfully",
		"book":    book,
	})
}

// AddBook handles POST /books - adds a new book (Librarian only)
func (ctrl *BookController) AddBook(c *gin.Context) {
	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.bookService.AddBook(&book); err != nil {
		if err.Error() == "quantity must be > 0" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": constants.BookAddedMessage,
		"book":    book,
	})
}

// UpdateBook handles PUT /books/:id - updates an existing book (Librarian only)
func (ctrl *BookController) UpdateBook(c *gin.Context) {
	bookID := c.Param("id")
	id, err := strconv.Atoi(bookID)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	var payload model.Book
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build update data map (only include non-empty fields)
	updateData := map[string]interface{}{}
	if payload.Title != "" {
		updateData["title"] = payload.Title
	}
	if payload.Author != "" {
		updateData["author"] = payload.Author
	}
	if payload.Category != "" {
		updateData["category"] = payload.Category
	}
	// Include numeric fields if they are provided (non-zero)
	if payload.Quantity != 0 {
		updateData["quantity"] = payload.Quantity
	}
	// Include available field (always include it as it could be zero)
	updateData["available"] = payload.Available

	book, err := ctrl.bookService.UpdateBook(uint(id), updateData)
	if err != nil {
		if err.Error() == constants.BookNotFoundError {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err.Error() == "no valid fields to update" ||
			err.Error() == "available copies cannot exceed total quantity" ||
			err.Error() == "available copies cannot be negative" ||
			err.Error() == "quantity cannot be negative" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": constants.BookUpdatedMessage,
		"book":    book,
	})
}

// DeleteBook handles DELETE /books/:id - deletes a book (Librarian only)
func (ctrl *BookController) DeleteBook(c *gin.Context) {
	bookID := c.Param("id")
	id, err := strconv.Atoi(bookID)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	err = ctrl.bookService.DeleteBook(uint(id))
	if err != nil {
		if err.Error() == constants.BookNotFoundError {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err.Error() == constants.BookCurrentlyBorrowedError {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": constants.BookDeletedMessage})
}

// BorrowBook handles POST /books/borrow - allows students to borrow books
func (ctrl *BookController) BorrowBook(c *gin.Context) {
	var req model.BorrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ctrl.bookService.BorrowBook(uint(req.GetStudentID()), uint(req.GetBookID()))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "student_id and book_id must be positive" ||
			err.Error() == constants.BookAlreadyBorrowedError ||
			err.Error() == constants.BookNotAvailableError {
			statusCode = http.StatusBadRequest
		} else if err.Error() == constants.StudentNotFoundError ||
			err.Error() == constants.BookNotFoundError {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        constants.BookBorrowedMessage,
		"borrow_details": response,
	})
}

// ReturnBook handles POST /books/return - allows students to return borrowed books
func (ctrl *BookController) ReturnBook(c *gin.Context) {
	var req model.BorrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ctrl.bookService.ReturnBook(uint(req.GetStudentID()), uint(req.GetBookID()))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "student_id and book_id must be positive" ||
			err.Error() == constants.BookNotBorrowedError {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        constants.BookReturnedMessage,
		"return_details": response,
	})
}
