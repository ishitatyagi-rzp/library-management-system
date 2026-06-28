package handlers

import (
	"net/http"
	"strconv"
	"time"

	"library-management-system/constants"
	"library-management-system/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetAllBooks returns a paginated list of books.
// Query params: ?limit=20&offset=0
func GetAllBooks(c *gin.Context, db *gorm.DB) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 {
		limit = 20
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	var books []model.Book
	result := db.Order("id DESC").
		Limit(limit).
		Offset(offset).
		Find(&books)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.FetchBooksError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Books retrieved successfully",
		"books":   books,
		"count":   len(books),
	})
}

// GetBookByID retrieves a single book by id
func GetBookByID(c *gin.Context, db *gorm.DB) {
	bookID := c.Param("id")
	id, err := strconv.Atoi(bookID)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	var book model.Book
	if err := db.First(&book, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": constants.BookNotFoundError})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.FetchBookError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book retrieved successfully",
		"book":    book,
	})
}

// AddBook adds a new book (Librarian only)
func AddBook(c *gin.Context, db *gorm.DB) {
	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if book.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quantity must be > 0"})
		return
	}

	// initialize available and added_at
	book.Available = book.Quantity
	if book.AddedAt.IsZero() {
		book.AddedAt = time.Now()
	}

	if err := db.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.AddBookError})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": constants.BookAddedMessage,
		"book":    book,
	})
}

// UpdateBook updates allowed fields of an existing book (Librarian only)
func UpdateBook(c *gin.Context, db *gorm.DB) {
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

	var book model.Book
	if err := db.First(&book, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": constants.BookNotFoundError})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.UpdateBookError})
		return
	}

	// Validate business rules:
	if payload.Available > payload.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "available copies cannot exceed total quantity"})
		return
	}
	if payload.Available < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "available copies cannot be negative"})
		return
	}
	if payload.Quantity < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quantity cannot be negative"})
		return
	}

	// Build map of fields to update explicitly to avoid unintended overwrites
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
	// quantity and available are numeric: update only if provided (non-zero or explicit)
	if payload.Quantity != 0 {
		updateData["quantity"] = payload.Quantity
	}
	// if available was provided (could be zero), update it
	// to detect whether client provided it, we rely on payload.Available != book.Available OR client intent;
	// for simplicity, allow updating available when it differs
	if payload.Available != book.Available {
		updateData["available"] = payload.Available
	}

	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no valid fields to update"})
		return
	}

	if err := db.Model(&book).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.UpdateBookError})
		return
	}

	// refresh model
	if err := db.First(&book, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.UpdateBookError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": constants.BookUpdatedMessage,
		"book":    book,
	})
}

// DeleteBook deletes a book (Librarian only). Prevent deletion if currently borrowed.
func DeleteBook(c *gin.Context, db *gorm.DB) {
	bookID := c.Param("id")
	id, err := strconv.Atoi(bookID)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	var borrowedCount int64
	if err := db.Model(&model.BorrowedBook{}).
		Where("book_id = ? AND return_date IS NULL", id).
		Count(&borrowedCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.DeleteBookError})
		return
	}

	if borrowedCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.BookCurrentlyBorrowedError})
		return
	}

	result := db.Delete(&model.Book{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.DeleteBookError})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": constants.BookNotFoundError})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": constants.BookDeletedMessage})
}

// BorrowBook allows a student to borrow a book. It uses a transaction and atomic update.
func BorrowBook(c *gin.Context, db *gorm.DB) {
	var req model.BorrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// basic validation - ensure positive ids
	if req.GetStudentID() <= 0 || req.GetBookID() <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_id and book_id must be positive"})
		return
	}

	// start transaction
	tx := db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.BorrowBookError})
		return
	}
	defer tx.Rollback() // safe to call even after commit

	// check student exists
	var student model.Student
	if err := tx.First(&student, req.GetStudentID()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": constants.StudentNotFoundError})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.BorrowBookError})
		return
	}

	// lock book row (FOR UPDATE) to prevent races and read current values
	var book model.Book
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&book, req.GetBookID()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": constants.BookNotFoundError})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.BorrowBookError})
		return
	}

	// check availability
	if book.Available <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.BookNotAvailableError})
		return
	}

	// check already borrowed (not returned)
	var count int64
	if err := tx.Model(&model.BorrowedBook{}).
		Where("student_id = ? AND book_id = ? AND return_date IS NULL", req.GetStudentID(), req.GetBookID()).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.BorrowBookError})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.BookAlreadyBorrowedError})
		return
	}

	// decrement available atomically
	result := tx.Model(&book).
		Where("id = ? AND available > 0", book.ID).
		UpdateColumn("available", gorm.Expr("available - ?", 1))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.BorrowBookError})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.BookNotAvailableError})
		return
	}

	// create borrow record
	now := time.Now()
	borrowRecord := model.BorrowedBook{
		StudentID:  uint(req.GetStudentID()),
		BookID:     uint(req.GetBookID()),
		BorrowDate: now,
	}
	if err := tx.Create(&borrowRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.BorrowBookError})
		return
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.BorrowBookError})
		return
	}

	// successful response
	c.JSON(http.StatusOK, gin.H{
		"message": constants.BookBorrowedMessage,
		"borrow_details": gin.H{
			"student_id":  borrowRecord.StudentID,
			"book_id":     borrowRecord.BookID,
			"book_title":  book.Title,
			"book_author": book.Author,
			"borrow_date": borrowRecord.BorrowDate,
		},
	})
}

// ReturnBook processes returning a borrowed book
func ReturnBook(c *gin.Context, db *gorm.DB) {
	var req model.BorrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.GetStudentID() <= 0 || req.GetBookID() <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_id and book_id must be positive"})
		return
	}

	tx := db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ReturnBookError})
		return
	}
	defer tx.Rollback()

	// find active borrowed record
	var borrowRecord model.BorrowedBook
	if err := tx.Where("student_id = ? AND book_id = ? AND return_date IS NULL",
		req.GetStudentID(), req.GetBookID()).First(&borrowRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": constants.BookNotBorrowedError})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ReturnBookError})
		return
	}

	// set return date (use same timestamp for DB and response)
	now := time.Now()
	if err := tx.Model(&borrowRecord).Update("return_date", &now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ReturnBookError})
		return
	}

	// increment available but ensure it doesn't exceed quantity
	result := tx.Model(&model.Book{}).
		Where("id = ? AND available < quantity", req.GetBookID()).
		UpdateColumn("available", gorm.Expr("available + ?", 1))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ReturnBookError})
		return
	}
	// It's OK if RowsAffected == 0: that means available == quantity (shouldn't happen normally)

	// fetch book for response
	var book model.Book
	if err := tx.First(&book, req.GetBookID()).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ReturnBookError})
		return
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ReturnBookError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": constants.BookReturnedMessage,
		"return_details": gin.H{
			"student_id":  req.GetStudentID(),
			"book_id":     req.GetBookID(),
			"book_title":  book.Title,
			"book_author": book.Author,
			"return_date": now,
		},
	})
}
