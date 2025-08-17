package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSoftware(c *gin.Context) {
	c.HTML(http.StatusOK, "software.tmpl", gin.H{
	})
}
