package middleware

import (
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/i2dou/sponge/pkg/gin/response"
	"github.com/i2dou/sponge/pkg/gohttp"
	"github.com/i2dou/sponge/pkg/utils"

	"github.com/gin-gonic/gin"
)

func runRateLimiterHTTPServer() string {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// e.g. (1) use default
	// r.Use(RateLimit())

	// e.g. (2) custom parameters
	r.Use(RateLimit(
		WithWindow(time.Second*10),
		WithBucket(200),
		WithCPUThreshold(500),
		WithCPUQuota(0.5),
	))

	r.GET("/hello", func(c *gin.Context) {
		if rand.Int()%2 == 0 {
			response.Output(c, http.StatusInternalServerError)
		} else {
			response.Success(c, "hello "+c.ClientIP())
		}
	})

	go func() {
		err := r.Run(serverAddr)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Millisecond * 200)
	return requestAddr
}

func TestRateLimiter(t *testing.T) {
	requestAddr := runRateLimiterHTTPServer()

	var success, failures int32
	for j := 0; j < 10; j++ {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				result := &gohttp.StdResult{}
				if err := gohttp.Get(result, requestAddr+"/hello"); err != nil {
					atomic.AddInt32(&failures, 1)
				} else {
					atomic.AddInt32(&success, 1)
				}
			}
		}()

		wg.Wait()
		t.Logf("%s   success: %d, failures: %d\n",
			time.Now().Format(time.RFC3339Nano), success, failures)
	}
}
