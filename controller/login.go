package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"

	"github.com/signal86/innerspeakerredux/model"
)

type LoginForm struct {
	Key string `form:"password"`
}

func LoginHandler(c *gin.Context) {
	var loginForm LoginForm
	if err := c.ShouldBind(&loginForm); err != nil {
		c.HTML(http.StatusUnprocessableEntity, "login.tmpl", gin.H{
			"error": "Malformed request",
		})
		return
	}
	var user model.User
	if err := model.DB.First(&user, "password = ?", model.Hash(loginForm.Key)).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
			"error": "Invalid key",
		})
		return
	}
	session_token := uuid.New().String()
	user.SessionToken = session_token
	model.DB.Save(&user)
	c.SetCookie("session_token", session_token, 3600, "/", "localhost", true, true)
	c.Redirect(http.StatusFound, "/datastore")
}
