package main

import (
	"fmt"
	shorten_url "go-url-shortener/api/shorten"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Shortening db works with the shorten_url package
// Need to test with a front end app
func main() {

	// Create a new instance of a router
	r := gin.New()

	fmt.Println(http.StatusOK)

	r.POST("/shorten", shorten_url.ShortenUrl)

	r.GET("/get-url/:id", shorten_url.GetOriginalUrlFromDb)

	// Start serving the application
	http.ListenAndServe(":8080", r)
}
