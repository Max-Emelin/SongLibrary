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

	id, err := h.services.Song.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) GetAllSongsWithFilter(c *gin.Context) {
	var filter model.SongFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid filter parameters: %v", err))
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		newErrorResponse(c, http.StatusBadRequest, "invalid limit parameter")
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		newErrorResponse(c, http.StatusBadRequest, "invalid offset parameter")
		return
	}

	songs, err := h.services.Song.GetAllSongsWithFilter(filter, limit, offset)
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
		newErrorResponse(c, http.StatusBadRequest, "invalid limit parameter")
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		newErrorResponse(c, http.StatusBadRequest, "invalid offset parameter")
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

type getAllSongsResponse struct {
	Data []model.Song `json:"data"`
}

type getLyricsResponse struct {
	Data []model.Lyrics `json:"data"`
}
