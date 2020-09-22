package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	// docs is generated by Swag CLI, you have to import it.

	"github.com/gin-contrib/pprof"
	"github.com/morentharia/anothergoproxy/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Api struct {
	*gin.Engine
}

func NewApi() (*Api, error) {
	docs.SwaggerInfo.Title = "Swagger API"
	docs.SwaggerInfo.Description = ""
	docs.SwaggerInfo.Version = "1.0"
	restURL, err := url.Parse(options.RestAddr)
	if err != nil {
		return nil, err
	}
	docs.SwaggerInfo.Host = fmt.Sprintf("%s", restURL.Host)
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := &Api{gin.Default()}

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With",
		)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	pprof.Register(r.Engine, "debug/pprof")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		// The url pointing to API definition
		ginSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", options.RestAddr)),
	))
	r.GET("/config", r.configHandler)
	r.GET("/reloadPage", r.reloadPageHandler)
	r.GET("/infoPages", r.infoPagesHandler)
	r.POST("/navigatePage", r.navigatePageHandler)
	r.POST("/log", r.logHandler)

	return r, nil
}

// Config godoc
// @Accept json
// @Produce json
// @Router /config [get]
// @Success 200 {string} string "answer"
func (a Api) configHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, options)
}

// Config godoc
// @Accept json
// @Produce json
// @Router /reloadPage [get]
// @Success 200 {string} string "answer"
func (a Api) reloadPageHandler(ctx *gin.Context) {
	err := browser.ReloadPageByURL()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, struct{}{})
	}

	ctx.JSON(http.StatusOK, struct{}{})
}

// Config godoc
// @Accept json
// @Produce json
// @Router /navigatePage [post]
// @Success 200 {string} string "answer"
func (a Api) navigatePageHandler(ctx *gin.Context) {
	req := struct {
		URL      string      `json:"url" binding:"required"`
		TargetID string      `json:"targetId" type:"integer" binding:"required"`
		WaitSec  json.Number `json:"waitSec" type:"integer"`
	}{
		WaitSec: "0",
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// pp.Println(req)
	waitSec, err := req.WaitSec.Int64()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err = browser.Navigate(req.TargetID, req.URL, int(waitSec)); err != nil {
		logrus.WithError(err).Error("navigate")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, struct{}{})
}

// Config godoc
// @Accept json
// @Produce json
// @Router /infoPages [Get]
// @Success 200 {string} string "answer"
func (a Api) infoPagesHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"result": browser.PagesInfo()})
}

// Config godoc
// @Accept json
// @Produce json
// @Router /log [post]
// @Success 200 {string} string "answer"
func (a Api) logHandler(ctx *gin.Context) {
	type Request struct {
		Type   string
		Params interface{}
	}
	var req Request
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// logrus.WithField("data", req).Info(req.Type)
	// go func(req *Request) {
	rotlog.WithField("data", req).Info(req.Type)
	// }(&req)

	ctx.JSON(http.StatusOK, struct{}{})
	return
}
