package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLFiles("./views/login.html")
	r.LoadHTMLFiles("./views/oauth.html")

	// html router
	r.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", gin.H{})
	})
	r.GET("/oauth", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "oauth.html", gin.H{})
	})

	// api router
	apiGroup := r.Group("/api")
	apiGroup.POST("/oauth", func(ctx *gin.Context) {
		type oauthRequest struct {
			OauthToken string `json:"oauth_token"`
		}

		var in oauthRequest
		if err := ctx.ShouldBind(&in); err != nil {
			log.Printf("error on binding oauth-request. err=%+v", err)
		}

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"oauth_token": in.OauthToken,
		})
	})

	httpaddress := ":3000"
	log.Printf("serving %s", httpaddress)
	if err := r.Run(httpaddress); err != nil {
		log.Printf("error on serving. err=%+v", err)
	}
}
