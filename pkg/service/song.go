package service

import (
	"SongLibrary/pkg/apiClient"
	"SongLibrary/pkg/model"
	"SongLibrary/pkg/repository"
	"strings"

	"github.com/sirupsen/logrus"
)

type SongService struct {
	repo   repository.Song
	client *apiClient.Client
}

func NewSongService(repo repository.Song, client *apiClient.Client) *SongService {
	logrus.Debug("Creating a new instance of SongService")
	return &SongService{
		repo:   repo,
		client: client,
	}
}

func (s *SongService) Create(song model.Song) (int, error) {
	logrus.Debugf("Creating a song: %+v", song)
	id, err := s.repo.Create(song)
	if err != nil {
		logrus.Errorf("Error creating song: %v", err)
		return 0, err
	}

	return id, nil
}

func (s *SongService) GetLyrics(songId int, limit, offset int) ([]model.Lyrics, error) {
	logrus.Debugf("Fetching lyrics for song with ID: %d, limit: %d, offset: %d", songId, limit, offset)
	lyrics, err := s.repo.GetLyrics(songId, limit, offset)
	if err != nil {
		logrus.Errorf("Error fetching lyrics for song with ID %d: %v", songId, err)
		return nil, err
	}

	return lyrics, nil
}

func (s *SongService) GetAllSongsWithFilter(filter model.SongFilter) ([]model.Song, error) {
	logrus.Debugf("Fetching songs with filter: %+v", filter)
	songs, err := s.repo.GetAllSongsWithFilter(filter)
	if err != nil {
		logrus.Errorf("Error fetching songs with filter: %v", err)
		return nil, err
	}

	return songs, nil
}

func (s *SongService) GetById(songId int) (model.Song, error) {
	logrus.Debugf("Fetching song by ID: %d", songId)
	song, err := s.repo.GetById(songId)
	if err != nil {
		logrus.Errorf("Error fetching song with ID %d: %v", songId, err)
		return model.Song{}, err
	}

	return song, nil
}

func (s *SongService) Delete(songId int) error {
	logrus.Debugf("Deleting song with ID: %d", songId)
	err := s.repo.Delete(songId)
	if err != nil {
		logrus.Errorf("Error deleting song with ID %d: %v", songId, err)
		return err
	}

	return nil
}

func (s *SongService) Update(songId int, input model.UpdateSongInput) error {
	logrus.Debugf("Updating song with ID: %d, input data: %+v", songId, input)
	if err := input.Validate(); err != nil {
		logrus.Errorf("Validation error for song update: %v", err)
		return err
	}
	err := s.repo.Update(songId, input)
	if err != nil {
		logrus.Errorf("Error updating song with ID %d: %v", songId, err)
		return err
	}

	return nil
}

func (s *SongService) FetchSongDetailsFromAPI(group, songName string) (*model.SongAPIResponse, error) {
	logrus.Debugf("Fetching song details from API for group: %s and song: %s", group, songName)
	apiResponse, err := s.client.FetchSongDetails(group, songName)
	if err != nil {
		logrus.Errorf("Error fetching details from API for group: %s and song: %s: %v", group, songName, err)
		return nil, err
	}

	return apiResponse, nil
}

func (s *SongService) UpdateSongWithAPIInfo(songId int, apiResponse model.SongAPIResponse) error {
	logrus.Debugf("Updating song with ID: %d based on API data", songId)
	err := s.repo.UpdateSongWithAPIInfo(model.UpdateSongApiData{
		SongId:      songId,
		ReleaseDate: apiResponse.ReleaseDate,
		Link:        apiResponse.Link,
		Lyrics:      splitIntoVerses(apiResponse.Text),
	})
	if err != nil {
		logrus.Errorf("Error updating song with ID: %d: %v", songId, err)
		return err
	}

	return nil
}

func splitIntoVerses(lyrics string) []string {
	return strings.Split(lyrics, "\n\n")
}
