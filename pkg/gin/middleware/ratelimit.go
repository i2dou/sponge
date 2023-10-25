package middleware

import (
	"net/http"
	"time"

	"github.com/i2dou/sponge/pkg/gin/response"
	rl "github.com/i2dou/sponge/pkg/shield/ratelimit"

	"github.com/gin-gonic/gin"
)

// ErrLimitExceed is returned when the rate limiter is
// triggered and the request is rejected due to limit exceeded.
var ErrLimitExceed = rl.ErrLimitExceed

// RateLimitOption set the rate limits rateLimitOptions.
type RateLimitOption func(*rateLimitOptions)

type rateLimitOptions struct {
	window       time.Duration
	bucket       int
	cpuThreshold int64
	cpuQuota     float64
}

func defaultRatelimitOptions() *rateLimitOptions {
	return &rateLimitOptions{
		window:       time.Second * 10,
		bucket:       100,
		cpuThreshold: 800,
	}
}

func (o *rateLimitOptions) apply(opts ...RateLimitOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithWindow with window size.
func WithWindow(d time.Duration) RateLimitOption {
	return func(o *rateLimitOptions) {
		o.window = d
	}
}

// WithBucket with bucket size.
func WithBucket(b int) RateLimitOption {
	return func(o *rateLimitOptions) {
		o.bucket = b
	}
}

// WithCPUThreshold with cpu threshold
func WithCPUThreshold(threshold int64) RateLimitOption {
	return func(o *rateLimitOptions) {
		o.cpuThreshold = threshold
	}
}

// WithCPUQuota with real cpu quota(if it can not collect from process correct);
func WithCPUQuota(quota float64) RateLimitOption {
	return func(o *rateLimitOptions) {
		o.cpuQuota = quota
	}
}

// RateLimit an adaptive rate limiter middleware
func RateLimit(opts ...RateLimitOption) gin.HandlerFunc {
	o := defaultRatelimitOptions()
	o.apply(opts...)
	limiter := rl.NewLimiter(
		rl.WithWindow(o.window),
		rl.WithBucket(o.bucket),
		rl.WithCPUThreshold(o.cpuThreshold),
		rl.WithCPUQuota(o.cpuQuota),
	)

	return func(c *gin.Context) {
		done, err := limiter.Allow()
		if err != nil {
			response.Output(c, http.StatusTooManyRequests, err.Error())
			c.Abort()
			return
		}

		c.Next()

		done(rl.DoneInfo{Err: c.Request.Context().Err()})
	}
}
