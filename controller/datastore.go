package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetDatastore(c *gin.Context) {
	c.HTML(http.StatusOK, "datastore.tmpl", gin.H{
	})
}
