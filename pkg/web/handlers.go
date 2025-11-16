package web

import (
	"log_analyzer/pkg/database"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ShowFilterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func RunFilter(c *gin.Context) {
	rawFilter := c.PostForm("filter")

	if strings.TrimSpace(rawFilter) == "" {
		c.HTML(http.StatusOK, "results.html", gin.H{
			"Error": "Filter cannot be empty",
		})
		return
	}

	// Split into parts by AND (space or comma)
	parts := database.SplitUserFilter(rawFilter)

	entries, err := database.QueryDB(DB, parts)
	if err != nil {
		c.HTML(http.StatusOK, "results.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "results.html", gin.H{
		"Entries": entries,
	})
}
