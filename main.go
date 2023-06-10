// Service to allow matching of titles to fields customized in a settings file
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func main() {
	r := gin.Default()
	r.GET("/ping", handlePing)
	r.POST("/match", HandleMatch)
	r.Run()
}

// Gets titles from Excel file
func GetTitles(excelPath string) ([]string, error) {
	xlData, err := excelize.OpenFile(excelPath)
	if err != nil {
		return nil, err
	}
	rows, err := xlData.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}
	titles := rows[0]
	return titles, nil
}

func HandleError(err error, c *gin.Context) {
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err})
	}
}

// Health check API
func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"message": "ok",
		})
}

// API end point that returns matches to a passed in array of titles
func HandleMatch(c *gin.Context) {
	var titles TitlesToMatch
	bindError := c.Bind(&titles)
	HandleError(bindError, c)
	suggestion := MatchFields(titles)

	c.JSON(http.StatusOK, gin.H{"data": suggestion})
}
