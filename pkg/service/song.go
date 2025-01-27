package service

import (
	"SongLibrary/pkg/apiClient"
	"SongLibrary/pkg/model"
	"SongLibrary/pkg/repository"
	"strings"
)

type SongService struct {
	repo   repository.Song
	client *apiClient.Client
}

func NewSongService(repo repository.Song, client *apiClient.Client) *SongService {
	return &SongService{
		repo:   repo,
		client: client,
	}
}

func (s *SongService) Create(song model.Song) (int, error) {
	return s.repo.Create(song)
}

func (s *SongService) GetLyrics(songId int, limit, offset int) ([]model.Lyrics, error) {
	return s.repo.GetLyrics(songId, limit, offset)
}

func (s *SongService) GetAllSongsWithFilter(filter model.SongFilter) ([]model.Song, error) {
	return s.repo.GetAllSongsWithFilter(filter)
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

func (s *SongService) FetchSongDetailsFromAPI(group, songName string) (*model.SongAPIResponse, error) {
	return s.client.FetchSongDetails(group, songName)
}

func (s *SongService) UpdateSongWithAPIInfo(songId int, apiResponse model.SongAPIResponse) error {
	return s.repo.UpdateSongWithAPIInfo(model.UpdateSongApiData{
		SongId:      songId,
		ReleaseDate: apiResponse.ReleaseDate,
		Link:        apiResponse.Link,
		Lyrics:      splitIntoVerses(apiResponse.Text),
	})
}

func splitIntoVerses(lyrics string) []string {
	return strings.Split(lyrics, "\n\n")
}
