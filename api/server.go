package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ssksameer56/CloudIndexer/controllers"
)

func RunServer() {

	sc := controllers.SearchController{}

	router := gin.Default()
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/search", sc.Search)

	router.Run(":" + "8080")
}
