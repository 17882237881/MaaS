package tenant

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

type createTenantRequest struct {
	Name string `json:"name"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func RegisterRoutes(h *server.Hertz, svc *Service) {
	h.GET("/healthz", func(_ context.Context, c *app.RequestContext) {
		c.String(http.StatusOK, "ok")
	})

	h.POST("/tenants", func(_ context.Context, c *app.RequestContext) {
		var req createTenantRequest
		if err := c.Bind(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid request"})
			return
		}
		tenant, err := svc.Create(req.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, tenant)
	})

	h.GET("/tenants/:id", func(_ context.Context, c *app.RequestContext) {
		id := c.Param("id")
		tenant, ok := svc.Get(id)
		if !ok {
			c.JSON(http.StatusNotFound, errorResponse{Error: "tenant not found"})
			return
		}
		c.JSON(http.StatusOK, tenant)
	})

	h.GET("/tenants", func(_ context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, svc.List())
	})

	h.DELETE("/tenants/:id", func(_ context.Context, c *app.RequestContext) {
		id := c.Param("id")
		if ok := svc.Delete(id); !ok {
			c.JSON(http.StatusNotFound, errorResponse{Error: "tenant not found"})
			return
		}
		c.Status(http.StatusNoContent)
	})
}
