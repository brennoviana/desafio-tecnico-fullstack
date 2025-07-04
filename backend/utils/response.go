package utils

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func RespondSuccess(c *gin.Context, data interface{}) {
	c.JSON(200, APIResponse{
		Status: "success",
		Data:   data,
	})
}

func RespondError(c *gin.Context, status int, errMsg string) {
	c.JSON(status, APIResponse{
		Status: "error",
		Error:  errMsg,
	})
}
