package web

import (
	"log_analyzer/pkg/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// func ShowFilterPage(c *gin.Context) {
// 	c.HTML(http.StatusOK, "index1.html", gin.H{
// 		"Level":     []string{},
// 		"Component": []string{},
// 		"Host":      []string{},
// 		"RequestID": "",
// 		"Timestamp": "",
// 	})
// }

// func RunFilter(c *gin.Context) {

// 	//  checkboxes
// 	levels := c.PostFormArray("level")
// 	components := c.PostFormArray("component")
// 	hosts := c.PostFormArray("host")

// 	// Textboxes
// 	requestID := c.PostForm("request_id")
// 	timestamp := c.PostForm("timestamp") // eg >2025-11-17 10:00:00"

// 	entries, err := database.FilterLogsWeb(DB, levels, components, hosts, requestID, timestamp)
// 	if err != nil {
// 		c.HTML(http.StatusOK, "index.html", gin.H{
// 			"Error":     err.Error(),
// 			"Level":     levels,
// 			"Component": components,
// 			"Host":      hosts,
// 			"RequestID": requestID,
// 			"Timestamp": timestamp,
// 		})
// 		return
// 	}

// 	c.HTML(http.StatusOK, "index.html", gin.H{
// 		"Entries":   entries,
// 		"Count":     len(entries),
// 		"Level":     levels,
// 		"Component": components,
// 		"Host":      hosts,
// 		"RequestID": requestID,
// 		"Timestamp": timestamp,
// 	})
// }
// func Hello(c *gin.Context) {
// 	c.HTML(http.StatusOK, "hello.html", gin.H{
// 		"name": "fizza",
// 	})
// }

// func RunFilter1(c *gin.Context) {
// 	//  checkboxes
// 	levels := c.PostFormArray("level")
// 	components := c.PostFormArray("component")
// 	hosts := c.PostFormArray("host")

// 	// Textboxes
// 	requestID := c.PostForm("request_id")
// 	timestamp := c.PostForm("timestamp") // eg >2025-11-17 10:00:00"

// 	entries, err := database.FilterLogsWeb(DB, levels, components, hosts, requestID, timestamp)
// 	if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {

// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": err.Error(),
// 			})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{
// 			"entries": entries,
// 			"count":   len(entries),
// 		})
// 		return
// 	}

// }
// func ShowAllLogs(c *gin.Context){
// 	entries, err:=database.GetAllLogs(DB)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": err.Error(),
// 			})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{
// 			"entries": entries[0:9999],
// 		})

// }
func FilterPaginatedLogs(c *gin.Context) {
	//query params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "100"))

	offset := page * pageSize

	//json body
	var body struct {
		Levels     []string `json:"levels"`
		Components []string `json:"components"`
		Hosts      []string `json:"hosts"`
		RequestID  string   `json:"requestId"`
		TimeStamp  string   `json:"timeStamp"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filtered, err := database.FilterLogsWeb(
		DB,
		body.Levels,
		body.Components,
		body.Hosts,
		body.RequestID,
		body.TimeStamp,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	total := len(filtered)

	start := offset
	end := offset + pageSize

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	pageEntries := filtered[start:end]

	c.JSON(http.StatusOK, gin.H{
		"entries": pageEntries,
		"total":   total,
	})
}
