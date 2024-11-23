package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type GalleryService struct {
	DB *sql.DB
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
	return nil
}
