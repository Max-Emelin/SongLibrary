package service

import (
	"SongLibrary/pkg/model"
	"SongLibrary/pkg/repository"
)

type Song interface {
	Create(song model.Song) (int, error)
	GetLyrics(songId int, limit, offset int) ([]model.Lyrics, error)
	GetAllSongsWithFilter(filter model.SongFilter, limit, offset int) ([]model.Song, error)
	GetById(songId int) (model.Song, error)
	Delete(songId int) error
	Update(songId int, input model.UpdateSongInput) error
}

type Service struct {
	Song
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Song: NewSongService(repos.Song),
	}
}
