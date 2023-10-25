package service

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	serverNameExampleV1 "github.com/i2dou/sponge/api/serverNameExample/v1"
	"github.com/i2dou/sponge/api/types"
	"github.com/i2dou/sponge/configs"
	"github.com/i2dou/sponge/internal/config"

	"github.com/i2dou/sponge/pkg/grpc/benchmark"
)

// Test each method of userExample via the rpc client
func Test_service_userExample_methods(t *testing.T) {
	conn := getRPCClientConnForTest()
	cli := serverNameExampleV1.NewUserExampleClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	//ctx = interceptor.SetJwtTokenToCtx(ctx, "Bearer jwt-token-value")

	tests := []struct {
		name    string
		fn      func() (interface{}, error)
		wantErr bool
	}{
		// todo generate the service struct code here
		// delete the templates code start
		{
			name: "Create",
			fn: func() (interface{}, error) {
				// todo type in the parameters to test
				req := &serverNameExampleV1.CreateUserExampleRequest{
					Name:     "foo7",
					Email:    "foo7@bar.com",
					Password: "f447b20a7fcbf53a5d5be013ea0b15af",
					Phone:    "16000000000",
					Avatar:   "http://internal.com/7.jpg",
					Age:      11,
					Gender:   2,
				}
				return cli.Create(ctx, req)
			},
			wantErr: false,
		},

		{
			name: "UpdateByID",
			fn: func() (interface{}, error) {
				// todo type in the parameters to test
				req := &serverNameExampleV1.UpdateUserExampleByIDRequest{
					Id:    7,
					Phone: "16000000001",
					Age:   11,
				}
				return cli.UpdateByID(ctx, req)
			},
			wantErr: false,
		},
		// delete the templates code end
		{
			name: "DeleteByID",
			fn: func() (interface{}, error) {
				// todo type in the parameters to test
				req := &serverNameExampleV1.DeleteUserExampleByIDRequest{
					Id: 100,
				}
				return cli.DeleteByID(ctx, req)
			},
			wantErr: false,
		},

		{
			name: "DeleteByIDs",
			fn: func() (interface{}, error) {
				// todo type in the parameters to test
				req := &serverNameExampleV1.DeleteUserExampleByIDsRequest{
					Ids: []uint64{100},
				}
				return cli.DeleteByIDs(ctx, req)
			},
			wantErr: false,
		},

		{
			name: "GetByID",
			fn: func() (interface{}, error) {
				// todo type in the parameters to test
				req := &serverNameExampleV1.GetUserExampleByIDRequest{
					Id: 1,
				}
				return cli.GetByID(ctx, req)
			},
			wantErr: false,
		},

		{
			name: "GetByCondition",
			fn: func() (interface{}, error) {
				// todo type in the parameters to test
				req := &serverNameExampleV1.GetUserExampleByConditionRequest{
					Conditions: &types.Conditions{
						Columns: []*types.Column{
							{
								Name:  "email",
								Value: "foo@bar.com",
							},
						},
					},
				}
				return cli.GetByCondition(ctx, req)
			},
			wantErr: false,
		},

		{
			name: "ListByIDs",
			fn: func() (interface{}, error) {
				// todo type in the parameters to test
				req := &serverNameExampleV1.ListUserExampleByIDsRequest{
					Ids: []uint64{1, 2, 3},
				}
				return cli.ListByIDs(ctx, req)
			},
			wantErr: false,
		},

		{
			name: "List",
			fn: func() (interface{}, error) {
				// todo type in the parameters to test
				req := &serverNameExampleV1.ListUserExampleRequest{
					Params: &types.Params{
						Page:  0,
						Limit: 10,
						Sort:  "",
						Columns: []*types.Column{
							{
								Name:  "id",
								Exp:   ">=",
								Value: "1",
								Logic: "",
							},
						},
					},
				}
				return cli.List(ctx, req)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fn()
			if (err != nil) != tt.wantErr {
				t.Logf("test '%s' error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			data, _ := json.MarshalIndent(got, "", "    ")
			fmt.Println(string(data))
		})
	}
}

// Perform a stress test on {{.LowerName}}'s method and
// copy the press test report to your browser when you are finished.
func Test_service_userExample_benchmark(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		panic(err)
	}
	host := fmt.Sprintf("127.0.0.1:%d", config.Get().Grpc.Port)
	protoFile := configs.Path("../api/serverNameExample/v1/userExample.proto")
	// If third-party dependencies are missing during the press test,
	// copy them to the project's third_party directory.
	importPaths := []string{
		configs.Path("../third_party"), // third_party directory
		configs.Path(".."),             // Previous level of third_party
	}

	tests := []struct {
		name    string
		fn      func() error
		wantErr bool
	}{
		{
			name: "GetByID",
			fn: func() error {
				// todo type in the parameters to test
				message := &serverNameExampleV1.GetUserExampleByIDRequest{
					Id: 1,
				}
				var total uint = 1000 // total number of requests
				b, err := benchmark.New(host, protoFile, "GetByID", message, total, importPaths...)
				if err != nil {
					return err
				}
				return b.Run()
			},
			wantErr: false,
		},

		{
			name: "ListByIDs",
			fn: func() error {
				// todo type in the parameters to test
				message := &serverNameExampleV1.ListUserExampleByIDsRequest{
					Ids: []uint64{1, 2, 3},
				}
				var total uint = 1000 // total number of requests
				b, err := benchmark.New(host, protoFile, "ListByIDs", message, total, importPaths...)
				if err != nil {
					return err
				}
				return b.Run()
			},
			wantErr: false,
		},

		{
			name: "List",
			fn: func() error {
				// todo type in the parameters to test
				message := &serverNameExampleV1.ListUserExampleRequest{
					Params: &types.Params{
						Page:  0,
						Limit: 10,
						Sort:  "",
						Columns: []*types.Column{
							{
								Name:  "id",
								Exp:   ">=",
								Value: "1",
								Logic: "",
							},
						},
					},
				}
				var total uint = 100 // total number of requests
				b, err := benchmark.New(host, protoFile, "List", message, total, importPaths...)
				if err != nil {
					return err
				}
				return b.Run()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if (err != nil) != tt.wantErr {
				t.Errorf("test '%s' error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
		})
	}
}
