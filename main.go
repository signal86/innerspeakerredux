package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/signal86/innerspeakerredux/controller"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("views/templates/*.tmpl")
	router.Static("/assets", "./views/assets")

	router.GET("/", controller.GetIndex )
	router.GET("/software", controller.GetSoftware )
	router.GET("/datastore", controller.GetDatastore )
	router.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	fmt.Println("Listening on port 8080")
	router.Run(":8080")
}
