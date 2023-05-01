package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type matchApiPostData []string

func main() {
	r := gin.Default()
	r.GET("/ping", handlePing)
	r.POST("/match", handleMatch)
	r.Run()
}

func getTitles(excelPath string) ([]string, error) {
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

func handleError(err error, c *gin.Context) {
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err})
	}
}

func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"message": "ok",
		})
}

func handleMatch(c *gin.Context) {
	var fieldHeaders matchApiPostData
	bindError := c.Bind(&fieldHeaders)
	handleError(bindError, c)

	suggestions := suggestFieldsForTitles(fieldHeaders)
	c.JSON(http.StatusOK, suggestions)
}
