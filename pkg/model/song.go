package model

import (
	"errors"
	"time"
)

type Song struct {
	ID          int        `json:"id" db:"id"`
	Group       string     `json:"group" db:"group" binding:"required"`
	SongName    string     `json:"song_name" db:"song_name" binding:"required"`
	ReleaseDate *time.Time `json:"release_date" db:"release_date"`
	Link        string     `json:"link" db:"link"`
}

type UpdateSongInput struct {
	Group    *string `json:"group"`
	SongName *string `json:"song_name"`
}

func (s *UpdateSongInput) Validate() error {
	if s.Group == nil && s.SongName == nil {
		return errors.New("update structure has no values")
	}

	return nil
}

type UpdateSongApiData struct {
	SongId      int
	ReleaseDate string
	Link        string
	Lyrics      []string
}
