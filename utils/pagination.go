package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPaginationParams reads the "limit" and "skip" query parameters.
func GetPaginationParams(c *gin.Context) (limit int, skip int) {
	limitStr := c.DefaultQuery("limit", "10")
	skipStr := c.DefaultQuery("skip", "0")

	lim, err := strconv.Atoi(limitStr)
	if err != nil {
		lim = 10
	}

	off, err := strconv.Atoi(skipStr)
	if err != nil {
		off = 0
	}

	return lim, off
}