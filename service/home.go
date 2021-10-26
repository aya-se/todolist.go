package service

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

// Home renders index.html
func Home(ctx *gin.Context) {
	t := time.Now().Local().Format("2006-01-02T15:04")
	ctx.HTML(http.StatusOK, "index.html", gin.H{"Title": "HOME", "Now": t})
}
