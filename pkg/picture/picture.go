package picture

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Picture errors.
var (
	ErrCreateFile   = errors.New("Create file fail")
	ErrCopyFile     = errors.New("Copy file fail")
	ErrFileNotFound = errors.New("File not found")
)

// Picture is an image to save in server.
type Picture struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// GetExtension returns the file extension.
func GetExtension(fileName string) string {
	titleSlice := strings.Split(fileName, ".")
	ext := titleSlice[len(titleSlice)-1]
	return ext
}

// GetSlug returns the picture slug to save in database.
func GetSlug(clientID, fileName string) string {
	return fmt.Sprintf("%s-%s-%d", fileName, clientID, time.Now().UnixNano())
}

// Upload receives the image file and saves it on the server.
func Upload(r *http.Request, fileName, clientID string) (string, error) {
	mpf, mph, err := r.FormFile(fileName)
	if err != nil {
		return "", ErrFileNotFound
	}
	defer mpf.Close()

	ext := GetExtension(mph.Filename)
	name := GetSlug(clientID, fileName)

	f, err := os.Create(filepath.Join("public", name+"."+ext))
	if err != nil {
		return "", ErrCreateFile
	}
	defer f.Close()

	_, err = io.Copy(f, mpf)
	if err != nil {
		return "", ErrCopyFile
	}

	fileURL := fmt.Sprintf("%s/public/%s.%s", os.Getenv("XD_BASE_URL_SERVER"), name, ext)

	return fileURL, nil
}

type Pictures []Picture
