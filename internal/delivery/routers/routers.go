package routers

import (
	"faceit_parser/internal/delivery/handlers"
	"faceit_parser/internal/repository"
	"faceit_parser/internal/service"
	"faceit_parser/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
)

func InitRouting(router *gin.Engine, logger *log.Logger) {
	cache := database.InitRedis()

	userRepo := repository.InitRepository(logger)
	userService := service.InitService(userRepo, logger, cache)
	userHandler := handlers.InitHandler(userService)

	router.Group("/")

	router.GET("", userHandler.GetUUID)
	router.POST("", userHandler.GetAllMatchesStats)
}
