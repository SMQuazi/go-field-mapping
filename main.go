package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	titles, titleErr := getTitles("spreadsheet.xlsx")

	handleError(titleErr, c)
	fmt.Print(titles)

	bindError := c.Bind(&fieldHeaders)
	handleError(bindError, c)

	settings, err := getSettings()
	handleError(err, c)
	fmt.Print(settings)
}

func main() {
	r := gin.Default()
	r.GET("/ping", handlePing)
	r.POST("/match", handleMatch)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
