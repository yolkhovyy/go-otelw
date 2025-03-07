package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yolkhovyy/go-otelw/cmd/example/internal/domain"
)

type Controller struct {
	domain  *domain.Controller
	handler *gin.Engine
}

func New(domain *domain.Controller, middleware ...gin.HandlerFunc) *Controller {
	controller := Controller{
		domain: domain,
	}

	gin.SetMode("release")
	engine := gin.New()
	engine.RedirectTrailingSlash = false
	engine.Use(middleware...)

	engine.POST("/echo", controller.echo)

	controller.handler = engine

	return &controller
}

func (c *Controller) Handler() http.Handler {
	return c.handler
}
