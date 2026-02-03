package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/17882237881/MaaS/internal/config"
	"github.com/17882237881/MaaS/internal/logging"
	"github.com/17882237881/MaaS/internal/tenant"
)

func main() {
	cfg := config.Load()
	logging.Init(cfg.LogLevel, nil)

	h := server.Default(server.WithHostPorts(cfg.HTTPAddr))
	svc := tenant.NewService(tenant.NewInMemoryStore())
	tenant.RegisterRoutes(h, svc)

	h.Spin()
}
