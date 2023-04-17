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

func GetSettings() Settings {
	data, err := os.ReadFile("settings.json")
	handleError(err)
	var settings Settings
	json.Unmarshal(data, &settings)
	return settings
}

func GetTitles(excelPath string) []string {
	f, err := excelize.OpenFile(excelPath)
	handleError(err)
	rows, err := f.GetRows("Sheet1")
	handleError(err)
	titles := rows[0]
	return titles
}

func main() {
	r := gin.Default()
	titles := GetTitles("spreadsheet.xlsx")
	fmt.Print(titles)
	r.GET("/ping", func(c *gin.Context) {

		c.JSON(http.StatusOK,
			gin.H{
				"message": "pong",
			})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
