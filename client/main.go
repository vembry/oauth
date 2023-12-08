package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

type session struct {
	OauthSessionToken string
	Data              data
}

type data struct {
	Username string    `json:"username"`
	CreateAt time.Time `json:"created_at"`
}

var sessionMapper = map[string]session{}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLFiles(
		"./views/index.html",
		"./views/login.html",
		"./views/oauth.html",
	)

	// html router
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", gin.H{})
	})
	r.GET("/oauth", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "oauth.html", gin.H{})
	})

	// api router
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/oauth/:oauth_token", func(ctx *gin.Context) {
			oauthToken, ok := ctx.Params.Get("oauth_token")
			if !ok {
				ctx.JSON(http.StatusBadRequest, nil)
				return
			}

			oauthSessionToken, err := getOauthSessionToken(oauthToken)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, nil)
				return
			}

			oauthSession, err := getOauthSession(oauthSessionToken)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, nil)
				return
			}
			if oauthSession == nil {
				ctx.JSON(http.StatusUnauthorized, nil)
				return
			}

			sessionId := ksuid.New().String()

			sessionMapper[sessionId] = session{
				OauthSessionToken: oauthSessionToken,
				Data: data{
					Username: oauthSession.Username,
					CreateAt: time.Now().UTC(),
				},
			}

			ctx.JSON(http.StatusOK, map[string]interface{}{
				"session_id": sessionId,
			})
		})

		apiGroup.GET("/session/:session_id", func(ctx *gin.Context) {
			sessionId, ok := ctx.Params.Get("session_id")
			if !ok {
				ctx.JSON(http.StatusBadRequest, nil)
				return
			}

			val, ok := sessionMapper[sessionId]
			if !ok {
				ctx.JSON(http.StatusUnauthorized, nil)
				return
			}

			ctx.JSON(http.StatusOK, val.Data)
		})
	}

	httpaddress := ":3000"
	log.Printf("serving %s", httpaddress)
	if err := r.Run(httpaddress); err != nil {
		log.Printf("error on serving. err=%+v", err)
	}
}

func getOauthSessionToken(oauthToken string) (string, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://localhost:3001/api/oauth/validate/%s", oauthToken),
		nil,
	)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request return non-200 http status. status=%d", res.StatusCode)
	}

	raw, _ := io.ReadAll(res.Body)

	var out map[string]string
	json.Unmarshal(raw, &out)

	return out["oauth_session_token"], nil
}

type oauthSession struct {
	Username   string    `json:"username"`
	HasConsent bool      `json:"has_consent"`
	CreatedAt  time.Time `json:"created_at"`
}

func getOauthSession(oauthSessionToken string) (*oauthSession, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://localhost:3001/api/oauth/session/%s", oauthSessionToken),
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request return non-200 http status. status=%d", res.StatusCode)
	}

	raw, _ := io.ReadAll(res.Body)

	var out oauthSession
	json.Unmarshal(raw, &out)

	return &out, nil
}
