package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func getSettings() Settings {
	data, err := os.ReadFile("settings.json")
	handleError(err)
	var settings Settings
	json.Unmarshal(data, &settings)
	return settings
}

func getTitles(excelPath string) []string {
	f, err := excelize.OpenFile(excelPath)
	handleError(err)
	rows, err := f.GetRows("Sheet1")
	handleError(err)
	titles := rows[0]
	return titles
}

func handlePing(c *gin.Context) {
	titles := getTitles("spreadsheet.xlsx")
	fmt.Print(titles)

	c.JSON(http.StatusOK,
		gin.H{
			"message": "ok",
		})
}

func handleMatch(c *gin.Context) {
	var fieldHeaders matchApiPostData
	err := c.Bind(&fieldHeaders)
	handleError(err)
}

func main() {
	r := gin.Default()
	r.GET("/ping", handlePing)
	r.POST("/match", handleMatch)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
