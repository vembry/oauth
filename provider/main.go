package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLFiles("views/index.html")

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	httpaddress := ":3001"
	log.Printf("serving %s", httpaddress)
	if err := r.Run(httpaddress); err != nil {
		log.Printf("error on serving. err=%+v", err)
	}
}
