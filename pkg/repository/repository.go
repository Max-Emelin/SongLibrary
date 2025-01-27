package repository

import (
	"SongLibrary/pkg/model"

	"github.com/jmoiron/sqlx"
)

type Song interface {
	Create(song model.Song) (int, error)
	GetLyrics(songId int, limit, offset int) ([]model.Lyrics, error)
	GetAllSongsWithFilter(filter model.SongFilter, limit, offset int) ([]model.Song, error)
	GetById(songId int) (model.Song, error)
	Delete(songId int) error
	Update(songId int, input model.UpdateSongInput) error
}

type Repository struct {
	Song
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Song: NewSongPostgres(db),
	}
}
