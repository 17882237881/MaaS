module maas-platform

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/google/uuid v1.5.0
	go.uber.org/zap v1.26.0
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.25.5
	github.com/redis/go-redis/v9 v9.3.0
	github.com/spf13/viper v1.18.1
	github.com/segmentio/kafka-go v0.4.46
	google.golang.org/grpc v1.60.0
	google.golang.org/protobuf v1.32.0
	github.com/hibiken/asynq v0.24.1
	github.com/casbin/casbin/v2 v2.79.0
	github.com/prometheus/client_golang v1.17.0
	github.com/jaegertracing/jaeger-client-go v2.30.0+incompatible
	github.com/swaggo/swag v1.16.2
	github.com/swaggo/gin-swagger v1.6.0
)
