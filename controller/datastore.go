package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"

	"github.com/signal86/innerspeakerredux/model"
)

var (
	uploadDirectory = "./datastore"
)

type Entry struct {
	Name    string
	Visible bool
}

func alreadyLoggedIn(c *gin.Context, user *model.User) bool {
	if session_token, err := c.Cookie("session_token"); err == nil && session_token != "" {
		fmt.Println(session_token)
		if err := model.DB.First(user, "session_token = ?", session_token).Error; err != nil {
			return false
		} else {
			return true
		}
	}
	return false
}

func GetDatastore(c *gin.Context) {
	var user model.User
	// already logged in
	if alreadyLoggedIn(c, &user) {
		var files []model.File
		model.DB.Where("username = ?", user.Username).Find(&files)
		var userstore []Entry
		errMsg := c.Query("error")
		for _, f := range files {
			userstore = append(userstore, Entry{Name: f.Name, Visible: f.Visible})
		}
		c.HTML(http.StatusOK, "datastore.tmpl", gin.H{
			"Files":    userstore,
			"Error":    errMsg,
			"Username": user.Username,
		})
	} else {
		c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
			"error": "Invalid session token",
		})
	}
	// if session_token, err := c.Cookie("session_token"); err == nil && session_token != "" {
	// 	fmt.Println(session_token)
	// 	if err := model.DB.First(&user, "session_token = ?", session_token).Error; err != nil {
	// 		c.HTML(http.StatusOK, "login.tmpl", gin.H{
	// 			"error": "Invalid session token",
	// 		})
	// 	} else {
	// 		c.HTML(http.StatusOK, "datastore.tmpl", gin.H{})
	// 	}
	// } else {
	// 	c.HTML(http.StatusOK, "login.tmpl", gin.H{
	// 		"error": "",
	// 	})
	// }
}

func PostFileUpload(c *gin.Context) {
	var user model.User
	if !alreadyLoggedIn(c, &user) {
		c.Redirect(http.StatusFound, "/?error=session invalid")
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.Redirect(http.StatusFound, "/?error="+err.Error())
		return
	}

	files := form.File["files"]
	for _, header := range files {
		destination := filepath.Join(uploadDirectory, user.Username, header.Filename)

		var existing model.File
		if err := model.DB.Where("name = ? AND username = ?", header.Filename, user.Username).First(&existing).Error; err == nil {
			c.Redirect(http.StatusFound, "/datastore?error=File already exists: "+header.Filename)
			return
		}

		if err := c.SaveUploadedFile(header, destination); err != nil {
			c.Redirect(http.StatusFound, "/datastore?error=Internal server failure: "+header.Filename)
			return
		}

		model.DB.Create(&model.File{Name: header.Filename, Visible: false, Username: user.Username})
	}

	c.Redirect(http.StatusFound, "/datastore")
}

func PostFileDelete(c *gin.Context) {
	var user model.User
	if !alreadyLoggedIn(c, &user) {
		c.Redirect(http.StatusFound, "/?error=session invalid")
		return
	}

	filename := c.PostForm("filename")
	if filename == "" {
		c.Redirect(http.StatusFound, "/?error=no filename provided")
		return
	}

	destination := filepath.Join(uploadDirectory, user.Username, filename)

	if _, err := os.Stat(destination); os.IsNotExist(err) {
		c.Redirect(http.StatusFound, "/?error=file not found")
		return
	}

	if err := os.Remove(destination); err != nil {
		c.Redirect(http.StatusFound, "/?error="+err.Error())
		return
	}

	model.DB.Unscoped().Where("name = ? AND username = ?", filename, user.Username).Delete(&model.File{})

	c.Redirect(http.StatusFound, "/datastore")
}

func PostVisibilityToggle(c *gin.Context) {
	var user model.User
	if !alreadyLoggedIn(c, &user) {
		c.Redirect(http.StatusFound, "/?error=session invalid")
		return
	}

	filename := c.PostForm("filename")
	if filename == "" {
		c.Redirect(http.StatusFound, "/?error=no filename provided")
		return
	}

	visible := c.PostForm("visible") == "on"

	model.DB.Model(&model.File{}).Where("name = ? AND username = ?", filename, user.Username).Update("visible", visible)

	c.Redirect(http.StatusFound, "/datastore")
}

func GetFile(c *gin.Context) {
	username := c.Param("username")
	filename := c.Param("filename")

	var file model.File
	if err := model.DB.Where("name = ? AND username = ?", filename, username).First(&file).Error; err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	// if hidden, must be the logged in owner
	if !file.Visible {
		var user model.User
		if !alreadyLoggedIn(c, &user) || user.Username != username {
			c.Redirect(http.StatusFound, "/")
			return
		}
	}

	filePath := filepath.Join(uploadDirectory, username, filename)
	c.File(filePath)
}
