package handlers

import (
	"encoding/json"
	"errors"
	"faceit_parser/internal/service"
	custom_errors "faceit_parser/pkg/errs"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	userService service.Service
}

func InitHandler(service service.Service) Handler {
	return Handler{
		userService: service,
	}
}

func (h Handler) GetUUID(c *gin.Context) {
	sessionID := uuid.New()

	c.SetCookie("session_id", sessionID.String(), 3600, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"session_id": sessionID})
}

func (h Handler) GetAllMatchesStats(c *gin.Context) {
	type data struct {
		Nickname string `json:"nickname"`
		Game     string `json:"game"`
	}

	body, _ := io.ReadAll(c.Request.Body)

	var values data
	if err := json.Unmarshal(body, &values); err != nil {
		// error log
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if values.Nickname == "" || values.Game == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no nickname or game info provided"})
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("user/%v", values.Nickname))

	matches, err := h.userService.GetAllMatchesStats(c.Request.Context(), values.Nickname, values.Game)
	if err != nil {
		if errors.Is(err, custom_errors.ErrWrongNickname) {
			c.JSON(http.StatusBadRequest, gin.H{"wrong nickname": values.Nickname})
			return
		}

		if errors.Is(err, custom_errors.ErrNoMatches) {
			c.JSON(http.StatusNoContent, gin.H{"game": values.Game, "message": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"requested nickname": values.Nickname, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, matches)
}
