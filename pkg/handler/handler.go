package handler

import (
	"SongLibrary/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		songs := api.Group("songs")
		{
			songs.POST("/", h.createSong)
			songs.GET("/", h.GetAllSongsWithFilter)
			songs.GET("/:id", h.getById)
			songs.GET("/:id/lyrics", getLyrics(services))
			songs.DELETE("/:id", deleteSong(services))
			songs.PUT("/:id", updateSong(services))

		}
	}

	return router
}
