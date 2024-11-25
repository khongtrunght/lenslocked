package models

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Image struct {
	// TODO: Add fields to this type
	GalleryID int
	Path      string
	Filename  string
}

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type GalleryService struct {
	DB *sql.DB

	// ImagesDir is used to tell the GalleryService where to store and locate
	// images. If not set, it will default to "images" directory.
	ImagesDir string
}

func (gs *GalleryService) Create(title string, userID int) (*Gallery, error) {
	gallery := Gallery{
		Title:  title,
		UserID: userID,
	}
	err := gs.DB.QueryRow(`
    INSERT INTO galleries (title, user_id)
    VALUES ($1, $2) RETURNING id;
  `, gallery.Title, gallery.UserID).Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}
	return &gallery, nil
}

func (gs *GalleryService) ByID(id int) (*Gallery, error) {
	var gallery Gallery
	err := gs.DB.QueryRow(`
    SELECT id, title, user_id FROM galleries
    WHERE id = $1
  `, id).Scan(&gallery.ID, &gallery.Title, &gallery.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query gallery: %w", err)
	}
	return &gallery, nil
}

func (gs *GalleryService) ByUserID(userID int) ([]Gallery, error) {
	rows, err := gs.DB.Query(`
    SELECT id, title FROM galleries
    WHERE user_id = $1
    ORDER BY id
  `, userID)
	if err != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}
	defer rows.Close()
	var galleries []Gallery
	for rows.Next() {
		gallery := Gallery{UserID: userID}
		err := rows.Scan(&gallery.ID, &gallery.Title)
		if err != nil {
			return nil, fmt.Errorf("query galleries by user: %w", err)
		}
		galleries = append(galleries, gallery)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}
	return galleries, nil
}

func (gs *GalleryService) Update(gallery *Gallery) error {
	_, err := gs.DB.Exec(`
    UPDATE galleries
    SET title = $2
    WHERE id = $1
  `, gallery.ID, gallery.Title)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}
	return nil
}

func (gs *GalleryService) Delete(id int) error {
	_, err := gs.DB.Exec(`
    DELETE FROM galleries WHERE id = $1
  `, id)
	if err != nil {
		return fmt.Errorf("delete gallery: %w", err)
	}

	err = os.RemoveAll(gs.galleryDir(id))
	if err != nil {
		return fmt.Errorf("delete gallery images: %w", err)
	}
	return nil
}

func (gs *GalleryService) Images(galleryID int) ([]Image, error) {
	globPattern := filepath.Join(gs.galleryDir(galleryID), "*")

	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retrieving gallery images: %w", err)
	}

	var images []Image
	for _, file := range allFiles {
		if hasExtension(file, gs.extensions()) {
			images = append(images, Image{
				Path:      file,
				Filename:  filepath.Base(file),
				GalleryID: galleryID,
			})
		}
	}

	return images, nil
}

func (gs *GalleryService) Image(galleryID int, filename string) (Image, error) {
	imagePath := filepath.Join(gs.galleryDir(galleryID), filename)
	_, err := os.Stat(imagePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Image{}, ErrNotFound
		}
		return Image{}, fmt.Errorf("querying for image: %w", err)
	}

	return Image{
		Filename:  filename,
		GalleryID: galleryID,
		Path:      imagePath,
	}, nil
}

func (gs *GalleryService) CreateImage(galleryID int, filename string, content io.ReadSeeker) error {
	err := checkContentType(content, gs.imageContentTypes())
	if err != nil {
		return fmt.Errorf("create image %v: %w", filename, err)
	}
	err = checkExtension(filename, gs.extensions())
	if err != nil {
		return fmt.Errorf("create image %v: %w", filename, err)
	}

	galleryDir := gs.galleryDir(galleryID)
	err = os.MkdirAll(galleryDir, 0755)
	if err != nil {
		return fmt.Errorf("creating gallery-%d directory: %w", galleryID, err)
	}
	imagePath := filepath.Join(galleryDir, filename)
	file, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("creating image file: %w", err)
	}
	defer file.Close()
	_, err = file.ReadFrom(content)
	if err != nil {
		return fmt.Errorf("writing image content: %w", err)
	}

	return nil
}

func (gs *GalleryService) DeleteImage(galleryID int, filename string) error {
	image, err := gs.Image(galleryID, filename)
	if err != nil {
		return fmt.Errorf("delete image: %w", err)
	}
	err = os.Remove(image.Path)
	if err != nil {
		return fmt.Errorf("delete image: %w", err)
	}
	return nil
}

func (gs *GalleryService) extensions() []string {
	return []string{".jpg", ".png", ".jpeg", ".gif"}
}

func (gs *GalleryService) imageContentTypes() []string {
	return []string{"image/png", "image/jpeg", "image/gif"}
}

func (gs *GalleryService) galleryDir(id int) string {
	imagesDir := gs.ImagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}

	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
}

func hasExtension(file string, extensions []string) bool {
	for _, ext := range extensions {
		file = strings.ToLower(file)
		ext = strings.ToLower(ext)
		if filepath.Ext(file) == ext {
			return true
		}
	}
	return false
}
