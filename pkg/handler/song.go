package handler

import (
	"SongLibrary/pkg/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createSong(c *gin.Context) {
	var input model.Song
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	songId, err := h.services.Song.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	////////////////////////////////////////////

	apiResponse, err := h.services.Song.FetchSongDetailsFromAPI(input.Group, input.SongName)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to fetch song details")
		return
	}

	err = h.services.Song.UpdateSongWithAPIInfo(songId, *apiResponse)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to update song with API info")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Song added successfully",
		"song_id": songId,
		"details": apiResponse,
	})

	///////////////////////////////////////////////////////////////

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": songId,
	})
}

func (h *Handler) GetAllSongsWithFilter(c *gin.Context) {
	var filter model.SongFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid filter parameters: %v", err))
		return
	}

	if filter.Limit == 0 {
		filter.Limit = 10
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	songs, err := h.services.Song.GetAllSongsWithFilter(filter)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error fetching songs: %v", err))
		return
	}

	c.JSON(http.StatusOK, getAllSongsResponse{
		Data: songs,
	})
}

func (h *Handler) getLyrics(c *gin.Context) {
	songId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		newErrorResponse(c, http.StatusBadRequest, "invalid limit param")
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		newErrorResponse(c, http.StatusBadRequest, "invalid offset param")
		return
	}

	lyrics, err := h.services.Song.GetLyrics(songId, limit, offset)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getLyricsResponse{
		Data: lyrics,
	})
}

func (h *Handler) getById(c *gin.Context) {
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	song, err := h.services.Song.GetById(listId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, song)
}

func (h *Handler) deleteSong(c *gin.Context) {
	songId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.Song.Delete(songId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) updateSong(c *gin.Context) {
	songId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input model.UpdateSongInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Song.Update(songId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

type getAllSongsResponse struct {
	Data []model.Song `json:"data"`
}

type getLyricsResponse struct {
	Data []model.Lyrics `json:"data"`
}
