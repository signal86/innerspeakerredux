package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"

	"github.com/signal86/innerspeakerredux/model"
)

func GetDatastore(c *gin.Context) {
	var user model.User
	// already logged in
	if session_token, err := c.Cookie("session_token");
	err == nil && session_token != "" {
		fmt.Println(session_token)
		if err := model.DB.First(&user, "session_token = ?", session_token).Error; 
		err != nil {
			c.HTML(http.StatusOK, "login.tmpl", gin.H{
				"error": "Invalid session token",
			})
		} else {
			c.HTML(http.StatusOK, "datastore.tmpl", gin.H{})
		}
	} else {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"error": "",
		})
	}
}

