package graphql

import (
	"encoding/json"
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
func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Get file ID from URL (e.g., /download?id=1234)
	fileID := r.URL.Query().Get("id")
	if fileID == "" {
		http.Error(w, "Missing file ID", http.StatusBadRequest)
		return
	}

	// Fetch file details from DB
	var document models.Document
	if err := initializers.DB.First(&document, "id = ?", fileID).Error; err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	filePath := document.FilePath
	fmt.Println("File Path:", filePath)
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found on server", http.StatusNotFound)
		return
	}

	// Set the response headers for file download
	// w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	w.Header().Set("Content-Disposition", "inline; filename="+filepath.Base(filePath))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", getFileSize(filePath)))
	fmt.Println("File Size:", getFileSize(filePath))
	// Serve the file
	http.ServeFile(w, r, filePath)
}
func getFileSize(path string) int64 {
	file, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return file.Size()
}

func listDocuments(w http.ResponseWriter, r *http.Request) {
	var documents []models.Document

	// Fetch all documents using GORM
	result := initializers.DB.Find(&documents)
	if result.Error != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	response := make([]map[string]interface{}, len(documents))
	for i, doc := range documents {
		response[i] = map[string]interface{}{
			"id":            doc.ID,
			"title":         doc.Title,
			"userId":        doc.UserID,
			"filePath":      doc.FilePath,
			"fileSize":      doc.FileSize,
			"fileType":      doc.FileType,
			"reference":     doc.ReferenceID,
			"referenceType": doc.ReferenceType,
			"tags":          doc.Tags,
			"createdAt":     doc.CreatedAt,
			"updatedAt":     doc.UpdatedAt,
			"deletedAt":     doc.DeletedAt,
		}
	}
	// Print total documents count
	fmt.Println("Total documents:", result.RowsAffected)

	// Send response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
