package handler

import (
	"SongLibrary/pkg/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) createSong(c *gin.Context) {
	logrus.Debug("createSong - binding input to model")
	var input model.Song
	if err := c.BindJSON(&input); err != nil {
		logrus.Error("createSong - failed to bind input", err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logrus.Debug("createSong - calling service to create song")
	songId, err := h.services.Song.Create(input)
	if err != nil {
		logrus.Error("createSong - failed to create song", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Info("createSong - song created successfully with ID:", songId)

	////////////////////////////////////////////

	// apiResponse, err := h.services.Song.FetchSongDetailsFromAPI(input.Group, input.SongName)
	// if err != nil {
	// 	logrus.Error("createSong - failed to fetch song details", err)
	// 	newErrorResponse(c, http.StatusInternalServerError, "failed to fetch song details")
	// 	return
	// }

	// err = h.services.Song.UpdateSongWithAPIInfo(songId, *apiResponse)
	// if err != nil {
	// 	logrus.Error("createSong - failed to update song with API info", err)
	// 	newErrorResponse(c, http.StatusInternalServerError, "failed to update song with API info")
	// 	return
	// }

	///////////////////////////////////////////////////////////////

	c.JSON(http.StatusOK, map[string]interface{}{"id": songId})
}

func (h *Handler) GetAllSongsWithFilter(c *gin.Context) {
	logrus.Debug("GetAllSongsWithFilter - binding filter from query parameters")
	var filter model.SongFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		logrus.Error("GetAllSongsWithFilter - invalid filter parameters", err)
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid filter parameters: %v", err))
		return
	}

	if filter.Limit == 0 {
		filter.Limit = 10
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	logrus.Debug("GetAllSongsWithFilter - fetching songs with filter")
	songs, err := h.services.Song.GetAllSongsWithFilter(filter)
	if err != nil {
		logrus.Error("GetAllSongsWithFilter - error fetching songs", err)
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error fetching songs: %v", err))
		return
	}

	logrus.Info("GetAllSongsWithFilter - successfully fetched songs")
	c.JSON(http.StatusOK, getAllSongsResponse{Data: songs})
}

func (h *Handler) getLyrics(c *gin.Context) {
	logrus.Debug("getLyrics - extracting songId and parameters from request")
	songId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error("getLyrics - invalid songId param", err)
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		logrus.Error("getLyrics - invalid limit param", err)
		newErrorResponse(c, http.StatusBadRequest, "invalid limit param")
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		logrus.Error("getLyrics - invalid offset param", err)
		newErrorResponse(c, http.StatusBadRequest, "invalid offset param")
		return
	}

	logrus.Debug("getLyrics - fetching lyrics for songId")
	lyrics, err := h.services.Song.GetLyrics(songId, limit, offset)
	if err != nil {
		logrus.Error("getLyrics - error fetching lyrics", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Info("getLyrics - successfully fetched lyrics")
	c.JSON(http.StatusOK, getLyricsResponse{Data: lyrics})
}

func (h *Handler) getById(c *gin.Context) {
	logrus.Debug("getById - extracting songId from request")
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error("getById - invalid songId param", err)
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	logrus.Debug("getById - fetching song by ID")
	song, err := h.services.Song.GetById(listId)
	if err != nil {
		logrus.Error("getById - error fetching song", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Info("getById - successfully fetched song")
	c.JSON(http.StatusOK, song)
}

func (h *Handler) deleteSong(c *gin.Context) {
	logrus.Debug("deleteSong - extracting songId from request")
	songId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error("deleteSong - invalid songId param", err)
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	logrus.Debug("deleteSong - deleting song by ID")
	err = h.services.Song.Delete(songId)
	if err != nil {
		logrus.Error("deleteSong - error deleting song", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Info("deleteSong - successfully deleted song")
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

func (h *Handler) updateSong(c *gin.Context) {
	logrus.Debug("updateSong - extracting songId from request")
	songId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error("updateSong - invalid songId param", err)
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	logrus.Debug("updateSong - binding update input")
	var input model.UpdateSongInput
	if err := c.BindJSON(&input); err != nil {
		logrus.Error("updateSong - failed to bind update input", err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logrus.Debug("updateSong - calling service to update song")
	if err := h.services.Song.Update(songId, input); err != nil {
		logrus.Error("updateSong - failed to update song", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Info("updateSong - successfully updated song")
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

type getAllSongsResponse struct {
	Data []model.Song `json:"data"`
}

type getLyricsResponse struct {
	Data []model.Lyrics `json:"data"`
}
