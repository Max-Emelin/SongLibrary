package service

import (
	"SongLibrary/pkg/model"
	"SongLibrary/pkg/repository"
)

type SongService struct {
	repo repository.Song
}

func NewSongService(repo repository.Song) *SongService {
	return &SongService{repo: repo}
}

func (s *SongService) Create(song model.Song) (int, error) {
	return s.repo.Create(song)
}

func (s *SongService) GetLyrics(songId int, limit, offset int) ([]model.Lyrics, error) {
	return s.repo.GetLyrics(songId, limit, offset)
}

func (s *SongService) GetAllSongsWithFilter(filter model.SongFilter, limit, offset int) ([]model.Song, error) {
	return s.repo.GetAllSongsWithFilter(filter, limit, offset)
}

func (s *SongService) GetById(songId int) (model.Song, error) {
	return s.repo.GetById(songId)
}

func (s *SongService) Delete(songId int) error {
	return s.repo.Delete(songId)
}

func (s *SongService) Update(songId int, input model.UpdateSongInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.Update(songId, input)
}
