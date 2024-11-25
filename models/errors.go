package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

var (
	ErrEmailTaken = errors.New("models: email address is already in use")
	ErrNotFound   = errors.New("models: no resource could be found with the provided information")
)

type FileError struct {
	Issue string
}

func (e FileError) Error() string {
	return fmt.Sprintf("invalid file: %v", e.Issue)
}

func checkContentType(r io.ReadSeeker, allowedTypes []string) error {
	testBytes := make([]byte, 512)
	_, err := r.Read(testBytes)
	if err != nil {
		return fmt.Errorf("checking content type: %w", err)
	}

	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("checking content type: %w", err)
	}

	contentType := http.DetectContentType(testBytes)
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			return nil
		}
	}

	return FileError{Issue: fmt.Sprintf("invalid content type: %v", contentType)}
}

func checkExtension(filename string, allowedExtensions []string) error {
	if hasExtension(filename, allowedExtensions) {
		return nil
	}

	return FileError{Issue: fmt.Sprintf("invalid file extension: %v", filepath.Ext(filename))}
}
