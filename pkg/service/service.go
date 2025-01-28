package service

import (
	"SongLibrary/pkg/apiClient"
	"SongLibrary/pkg/model"
	"SongLibrary/pkg/repository"

	"github.com/sirupsen/logrus"
)

type Song interface {
	Create(song model.Song) (int, error)
	GetLyrics(songId int, limit, offset int) ([]model.Lyrics, error)
	GetAllSongsWithFilter(filter model.SongFilter) ([]model.Song, error)
	GetById(songId int) (model.Song, error)
	Delete(songId int) error
	Update(songId int, input model.UpdateSongInput) error
	FetchSongDetailsFromAPI(group, songName string) (*model.SongAPIResponse, error)
	UpdateSongWithAPIInfo(songId int, apiResponse model.SongAPIResponse) error
}

type Service struct {
	Song
}

func NewService(repos *repository.Repository, client *apiClient.Client) *Service {
	logrus.Debug("NewService - initializing service")
	service := &Service{
		Song: NewSongService(repos.Song, client),
	}
	logrus.Debug("NewService - service initialized successfully")

	return service
}
