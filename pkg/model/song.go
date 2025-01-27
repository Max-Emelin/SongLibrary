package model

import "errors"

type Song struct {
	ID       int    `json:"id" db:"id"`
	Group    string `json:"group" binding:"required"`
	SongName string `json:"song_name" binding:"required"`
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
