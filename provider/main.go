package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

type session struct {
	Username   string    `json:"username"`
	HasConsent bool      `json:"has_consent"`
	CreatedAt  time.Time `json:"created_at"`
}

var oauthMapper = map[string]string{}
var mapper = map[string]session{}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLFiles("views/index.html")

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	apiRouter := r.Group("/api")
	{
		apiRouter.POST("/oauth/login", func(ctx *gin.Context) {
			type loginRequest struct {
				Username   string `json:"username"`
				HasConsent bool   `json:"has_consent"`
			}

			var in loginRequest
			if err := ctx.ShouldBind(&in); err != nil {
				log.Fatalf("error on parsing request. err=%+v", err)
			}

			// construct oauth session
			key := ksuid.New().String()
			mapper[key] = session{
				Username:   in.Username,
				HasConsent: in.HasConsent,
				CreatedAt:  time.Now().UTC(),
			}

			oauthToken := ksuid.New().String()
			oauthMapper[oauthToken] = key

			// return oauth key
			ctx.JSON(http.StatusOK, map[string]interface{}{
				"token": oauthToken,
			})
		})

		apiRouter.GET("/oauth/validate/:oauth_token", func(ctx *gin.Context) {
			oauthToken, ok := ctx.Params.Get("oauth_token")
			if !ok {
				ctx.JSON(http.StatusBadRequest, nil)
				return
			}
			val, ok := oauthMapper[oauthToken]
			if !ok {
				ctx.JSON(http.StatusUnauthorized, nil)
				return
			}

			// delete(oauthMapper, oauthToken)

			ctx.JSON(http.StatusOK, map[string]string{
				"oauth_session_token": val,
			})
		})

		apiRouter.GET("/oauth/session/:oauth_token", func(ctx *gin.Context) {
			oauthToken, ok := ctx.Params.Get("oauth_token")
			if !ok {
				ctx.JSON(http.StatusBadRequest, nil)
				return
			}

			session := mapper[oauthToken]

			if session.HasConsent {
				ctx.JSON(http.StatusOK, session)
				return
			}

			ctx.JSON(http.StatusOK, map[string]interface{}{})
		})

		apiRouter.GET("/oauth/active", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, oauthMapper)
		})

		apiRouter.GET("/oauth/active-session", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, mapper)
		})
	}

	httpaddress := ":3001"
	log.Printf("serving %s", httpaddress)
	if err := r.Run(httpaddress); err != nil {
		log.Printf("error on serving. err=%+v", err)
	}
}
