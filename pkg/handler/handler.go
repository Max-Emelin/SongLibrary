package handler

import (
	"SongLibrary/pkg/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	logrus.Debug("Initializing new handler with provided service")
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	logrus.Debug("Initializing routes")

	router := gin.New()

	api := router.Group("/api")
	{
		songs := api.Group("songs")
		{
			songs.POST("/", h.createSong)
			songs.GET("/", h.GetAllSongsWithFilter)
			songs.GET("/:id", h.getById)
			songs.GET("/:id/lyrics", h.getLyrics)
			songs.DELETE("/:id", h.deleteSong)
			songs.PUT("/:id", h.updateSong)
		}
	}

	logrus.Info("Routes initialized successfully")
	return router
}
