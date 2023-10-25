package initial

import (
	"context"
	"time"

	"github.com/i2dou/sponge/internal/config"
	"github.com/i2dou/sponge/internal/model"

	"github.com/i2dou/sponge/pkg/app"
	"github.com/i2dou/sponge/pkg/tracer"
)

// RegisterClose register for released resources
func RegisterClose(servers []app.IServer) []app.Close {
	var closes []app.Close

	// close server
	for _, s := range servers {
		closes = append(closes, s.Stop)
	}

	// close mysql
	closes = append(closes, func() error {
		return model.CloseMysql()
	})

	// close redis
	if config.Get().App.CacheType == "redis" {
		closes = append(closes, func() error {
			return model.CloseRedis()
		})
	}

	// close tracing
	if config.Get().App.EnableTrace {
		closes = append(closes, func() error {
			ctx, _ := context.WithTimeout(context.Background(), 2*time.Second) //nolint
			return tracer.Close(ctx)
		})
	}

	return closes
}
