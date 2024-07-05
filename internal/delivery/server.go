package delivery

import (
	"faceit_parser/internal/delivery/routers"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func Start(logger *log.Logger) {
	server := gin.Default()

	routers.InitRouting(server, logger)

	if err := server.Run(":8080"); err != nil {
		panic(fmt.Sprintf("error while running the server: %v", err.Error()))
	}
}
