package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"

	// "gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/signal86/innerspeakerredux/controller"
	"github.com/signal86/innerspeakerredux/model"
)

// Running the server with extra data will not start the server and will
// instead manual add a user to the db
func KeyGenerator(username string) string {
	newKey := string(rand.Text())
	model.ConnectDatabase()
	user := model.User{
		Username: os.Args[1],
		Password: model.Hash(newKey),
	}

	// key overwrite vs new user
	if err := model.DB.First(&user, "username = ?", username).Error; err == nil {
		fmt.Printf("Key already exists for %s\nOverwriting\n", username)
		model.DB.Model(&user).Where("username = ?", username).Update("password", model.Hash(newKey))
	} else {
		model.DB.Create(&user)
	}

	fmt.Printf("Key Generated: %s\n", newKey)
	return newKey
}

// Actual server
func main() {
	// Except for this conditional
	if len(os.Args) > 1 {
		fmt.Printf("Generating key for %s\n", os.Args[1])
		fmt.Printf("New key: %s\n", KeyGenerator(os.Args[1]))
		return
	}
	model.ConnectDatabase()
	router := gin.Default()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	router.LoadHTMLGlob("views/templates/*.tmpl")
	router.Static("/assets", "./views/assets")

	router.GET("/", controller.GetIndex)
	router.GET("/software", controller.GetSoftware)
	router.GET("/datastore", controller.GetDatastore)
	router.POST("/login", controller.LoginHandler)
	router.POST("/upload", controller.PostFileUpload)
	router.POST("/delete", controller.PostFileDelete)
	router.POST("/toggle-visibility", controller.PostVisibilityToggle)
	router.GET("/file/:username/:filename", controller.GetFile)
	router.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	fmt.Println("Listening on port 8080")
	router.Run(":8080")
}
