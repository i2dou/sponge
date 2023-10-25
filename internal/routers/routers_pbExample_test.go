package routers

import (
	"context"
	"testing"
	"time"

	serverNameExampleV1 "github.com/i2dou/sponge/api/serverNameExample/v1"
	"github.com/i2dou/sponge/configs"
	"github.com/i2dou/sponge/internal/config"

	"github.com/i2dou/sponge/pkg/gin/middleware"
	"github.com/i2dou/sponge/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter_pbExample(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		t.Fatal(err)
	}

	config.Get().App.EnableMetrics = true
	config.Get().App.EnableTrace = true
	config.Get().App.EnableHTTPProfile = true
	config.Get().App.EnableLimit = true
	config.Get().App.EnableCircuitBreaker = true

	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		gin.SetMode(gin.ReleaseMode)
		r := NewRouter_pbExample()
		assert.NotNil(t, r)
		cancel()
	})
}

func Test_middlewareConfig(t *testing.T) {
	c := newMiddlewareConfig()

	c.setGroupPath("/api/v1", middleware.Auth())
	assert.Equal(t, 1, len(c.groupPathMiddlewares["/api/v1"]))
	c.setGroupPath("/api/v1", middleware.RateLimit(), middleware.RequestID())
	assert.Equal(t, 3, len(c.groupPathMiddlewares["/api/v1"]))

	c.setSinglePath("DELETE", "/api/v1/userExample/:id", middleware.Auth())
	assert.Equal(t, 1, len(c.singlePathMiddlewares[getSinglePathKey("DELETE", "/api/v1/userExample/:id")]))
	c.setSinglePath("POST", "/api/v1/userExample/list", middleware.RateLimit(), middleware.RequestID())
	assert.Equal(t, 2, len(c.singlePathMiddlewares[getSinglePathKey("POST", "/api/v1/userExample/list")]))
}

func Test_userExampleServiceRouter(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	c := newMiddlewareConfig()
	userExampleServiceRouter(r, c.groupPathMiddlewares, c.singlePathMiddlewares, &mockGw{})
}

type mockGw struct{}

func (m mockGw) Create(ctx context.Context, req *serverNameExampleV1.CreateUserExampleRequest) (*serverNameExampleV1.CreateUserExampleReply, error) {
	return nil, nil
}

func (m mockGw) DeleteByID(ctx context.Context, req *serverNameExampleV1.DeleteUserExampleByIDRequest) (*serverNameExampleV1.DeleteUserExampleByIDReply, error) {
	return nil, nil
}

func (m mockGw) DeleteByIDs(ctx context.Context, req *serverNameExampleV1.DeleteUserExampleByIDsRequest) (*serverNameExampleV1.DeleteUserExampleByIDsReply, error) {
	return nil, nil
}

func (m mockGw) GetByID(ctx context.Context, req *serverNameExampleV1.GetUserExampleByIDRequest) (*serverNameExampleV1.GetUserExampleByIDReply, error) {
	return nil, nil
}

func (m mockGw) GetByCondition(ctx context.Context, req *serverNameExampleV1.GetUserExampleByConditionRequest) (*serverNameExampleV1.GetUserExampleByConditionReply, error) {
	return nil, nil
}

func (m mockGw) List(ctx context.Context, req *serverNameExampleV1.ListUserExampleRequest) (*serverNameExampleV1.ListUserExampleReply, error) {
	return nil, nil
}

func (m mockGw) ListByIDs(ctx context.Context, req *serverNameExampleV1.ListUserExampleByIDsRequest) (*serverNameExampleV1.ListUserExampleByIDsReply, error) {
	return nil, nil
}

func (m mockGw) UpdateByID(ctx context.Context, req *serverNameExampleV1.UpdateUserExampleByIDRequest) (*serverNameExampleV1.UpdateUserExampleByIDReply, error) {
	return nil, nil
}
