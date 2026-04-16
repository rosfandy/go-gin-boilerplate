package presenters

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Success bool    `json:"success"`
	Message *string `json:"message,omitempty"`
	Data    *any    `json:"data,omitempty"`
	Cursor  *int64  `json:"cursor,omitempty"`
	Total   *int64  `json:"total,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func SendSuccessResponse(c *gin.Context, response SuccessResponse) {
	c.JSON(http.StatusOK, response)
}

func SendErrorResponse(c *gin.Context, errCode int, err string) {
	c.JSON(errCode, ErrorResponse{
		Success: false,
		Error:   err,
	})
}
