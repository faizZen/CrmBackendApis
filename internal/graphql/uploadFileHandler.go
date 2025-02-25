package graphql

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	initializers "github.com/Zenithive/it-crm-backend/Initializers"
	"github.com/Zenithive/it-crm-backend/auth"
	"github.com/Zenithive/it-crm-backend/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Max file size: 10MB
const maxUploadSize = 10 << 20

// Allowed file extensions
var allowedExtensions = map[string]bool{
	".pdf":  true,
	".docx": true,
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from JWT token
	userID, err := auth.ExtractUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Println("User ID:", userID)

	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	// Parse multipart form data
	err = r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
		return
	}

	// Extract form values
	referenceID := r.FormValue("referenceID")
	referenceType := r.FormValue("referenceType")
	tags := r.MultipartForm.Value["tags"] // Get tags as an array

	fmt.Println("Reference ID:", referenceID)
	fmt.Println("Reference Type:", referenceType)
	fmt.Println("Tags:", tags)

	if referenceID == "" || referenceType == "" {
		http.Error(w, "Missing referenceID or referenceType", http.StatusBadRequest)
		return
	}

	// Get file
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if !allowedExtensions[ext] {
		http.Error(w, "Invalid file type. Only PDF and DOCX are allowed.", http.StatusBadRequest)
		return
	}

	// Check file size
	if handler.Size > maxUploadSize {
		http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
		return
	}

	// Create uploads directory if it doesn't exist
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// Generate unique file name
	fileUUID := uuid.New().String()
	filePath := filepath.Join(uploadDir, fileUUID+ext)

	// Save file to disk
	if err := saveUploadedFile(file, filePath); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Parse UUIDs
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	parsedReferenceID, err := uuid.Parse(referenceID)
	if err != nil {
		http.Error(w, "Invalid reference ID", http.StatusBadRequest)
		return
	}

	// Store document details in the database
	docDetails := models.Document{
		ID:            uuid.New(),
		Title:         handler.Filename,
		FilePath:      filePath,
		FileSize:      fmt.Sprintf("%d bytes", handler.Size),
		FileType:      ext,
		Tags:          pq.StringArray(tags),
		UserID:        parsedUserID,
		ReferenceID:   parsedReferenceID,
		ReferenceType: referenceType,
	}

	if err := initializers.DB.Create(&docDetails).Error; err != nil {
		http.Error(w, "Failed to store document in database", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "File uploaded successfully: %s", filePath)
}

// Function to save file to disk
func saveUploadedFile(file multipart.File, filePath string) error {
	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	return err
}
