package service

import (
	"context"
	"errors"
	"strings"

	serverNameExampleV1 "github.com/i2dou/sponge/api/serverNameExample/v1"
	"github.com/i2dou/sponge/internal/cache"
	"github.com/i2dou/sponge/internal/dao"
	"github.com/i2dou/sponge/internal/ecode"
	"github.com/i2dou/sponge/internal/model"

	"github.com/i2dou/sponge/pkg/grpc/interceptor"
	"github.com/i2dou/sponge/pkg/logger"
	"github.com/i2dou/sponge/pkg/mysql/query"

	"github.com/jinzhu/copier"
	"google.golang.org/grpc"
)

func init() {
	registerFns = append(registerFns, func(server *grpc.Server) {
		serverNameExampleV1.RegisterUserExampleServer(server, NewUserExampleServer()) // register service to the rpc service
	})
}

var _ serverNameExampleV1.UserExampleServer = (*userExample)(nil)

type userExample struct {
	serverNameExampleV1.UnimplementedUserExampleServer

	iDao dao.UserExampleDao
}

// NewUserExampleServer create a new service
func NewUserExampleServer() serverNameExampleV1.UserExampleServer {
	return &userExample{
		iDao: dao.NewUserExampleDao(
			model.GetDB(),
			cache.NewUserExampleCache(model.GetCacheType()),
		),
	}
}

// Create a record
func (s *userExample) Create(ctx context.Context, req *serverNameExampleV1.CreateUserExampleRequest) (*serverNameExampleV1.CreateUserExampleReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	record := &model.UserExample{}
	err = copier.Copy(record, req)
	if err != nil {
		return nil, ecode.StatusCreateUserExample.Err()
	}

	ctx = context.WithValue(ctx, interceptor.ContextRequestIDKey, interceptor.ServerCtxRequestID(ctx)) //nolint
	err = s.iDao.Create(ctx, record)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("userExample", record), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	return &serverNameExampleV1.CreateUserExampleReply{Id: record.ID}, nil
}

// DeleteByID delete a record by id
func (s *userExample) DeleteByID(ctx context.Context, req *serverNameExampleV1.DeleteUserExampleByIDRequest) (*serverNameExampleV1.DeleteUserExampleByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	ctx = context.WithValue(ctx, interceptor.ContextRequestIDKey, interceptor.ServerCtxRequestID(ctx)) //nolint
	err = s.iDao.DeleteByID(ctx, req.Id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", req.Id), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	return &serverNameExampleV1.DeleteUserExampleByIDReply{}, nil
}

// DeleteByIDs delete records by batch id
func (s *userExample) DeleteByIDs(ctx context.Context, req *serverNameExampleV1.DeleteUserExampleByIDsRequest) (*serverNameExampleV1.DeleteUserExampleByIDsReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	ctx = context.WithValue(ctx, interceptor.ContextRequestIDKey, interceptor.ServerCtxRequestID(ctx)) //nolint
	err = s.iDao.DeleteByIDs(ctx, req.Ids)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("ids", req.Ids), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	return &serverNameExampleV1.DeleteUserExampleByIDsReply{}, nil
}

// UpdateByID update a record by id
func (s *userExample) UpdateByID(ctx context.Context, req *serverNameExampleV1.UpdateUserExampleByIDRequest) (*serverNameExampleV1.UpdateUserExampleByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	record := &model.UserExample{}
	err = copier.Copy(record, req)
	if err != nil {
		return nil, ecode.StatusUpdateByIDUserExample.Err()
	}
	record.ID = req.Id

	ctx = context.WithValue(ctx, interceptor.ContextRequestIDKey, interceptor.ServerCtxRequestID(ctx)) //nolint
	err = s.iDao.UpdateByID(ctx, record)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("userExample", record), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	return &serverNameExampleV1.UpdateUserExampleByIDReply{}, nil
}

// GetByID get a record by id
func (s *userExample) GetByID(ctx context.Context, req *serverNameExampleV1.GetUserExampleByIDRequest) (*serverNameExampleV1.GetUserExampleByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	ctx = context.WithValue(ctx, interceptor.ContextRequestIDKey, interceptor.ServerCtxRequestID(ctx)) //nolint
	record, err := s.iDao.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, query.ErrNotFound) {
			logger.Warn("GetByID error", logger.Err(err), logger.Any("id", req.Id), interceptor.ServerCtxRequestIDField(ctx))
			return nil, ecode.StatusNotFound.Err()
		}
		logger.Error("GetByID error", logger.Err(err), logger.Any("id", req.Id), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	data, err := convertUserExample(record)
	if err != nil {
		logger.Warn("convertUserExample error", logger.Err(err), logger.Any("userExample", record), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusGetByIDUserExample.Err()
	}

	return &serverNameExampleV1.GetUserExampleByIDReply{UserExample: data}, nil
}

// GetByCondition get a record by id
func (s *userExample) GetByCondition(ctx context.Context, req *serverNameExampleV1.GetUserExampleByConditionRequest) (*serverNameExampleV1.GetUserExampleByConditionReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	ctx = context.WithValue(ctx, interceptor.ContextRequestIDKey, interceptor.ServerCtxRequestID(ctx)) //nolint
	conditions := &query.Conditions{}
	for _, v := range req.Conditions.GetColumns() {
		column := query.Column{}
		_ = copier.Copy(&column, v)
		conditions.Columns = append(conditions.Columns, column)
	}
	err = conditions.CheckValid()
	if err != nil {
		logger.Warn("Parameters error", logger.Err(err), logger.Any("conditions", conditions), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	record, err := s.iDao.GetByCondition(ctx, conditions)
	if err != nil {
		if errors.Is(err, query.ErrNotFound) {
			logger.Warn("GetByCondition error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
			return nil, ecode.StatusNotFound.Err()
		}
		logger.Error("GetByCondition error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	data, err := convertUserExample(record)
	if err != nil {
		logger.Warn("convertUserExample error", logger.Err(err), logger.Any("userExample", record), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusGetByConditionUserExample.Err()
	}

	return &serverNameExampleV1.GetUserExampleByConditionReply{
		UserExample: data,
	}, nil
}

// ListByIDs list of records by batch id
func (s *userExample) ListByIDs(ctx context.Context, req *serverNameExampleV1.ListUserExampleByIDsRequest) (*serverNameExampleV1.ListUserExampleByIDsReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	ctx = context.WithValue(ctx, interceptor.ContextRequestIDKey, interceptor.ServerCtxRequestID(ctx)) //nolint
	userExampleMap, err := s.iDao.GetByIDs(ctx, req.Ids)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("ids", req.Ids), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	userExamples := []*serverNameExampleV1.UserExample{}
	for _, id := range req.Ids {
		if v, ok := userExampleMap[id]; ok {
			record, err := convertUserExample(v)
			if err != nil {
				logger.Warn("convertUserExample error", logger.Err(err), logger.Any("userExample", v), interceptor.ServerCtxRequestIDField(ctx))
				return nil, ecode.StatusInternalServerError.ToRPCErr()
			}
			userExamples = append(userExamples, record)
		}
	}

	return &serverNameExampleV1.ListUserExampleByIDsReply{UserExamples: userExamples}, nil
}

// List of records by query parameters
func (s *userExample) List(ctx context.Context, req *serverNameExampleV1.ListUserExampleRequest) (*serverNameExampleV1.ListUserExampleReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	params := &query.Params{}
	err = copier.Copy(params, req.Params)
	if err != nil {
		return nil, ecode.StatusListUserExample.Err()
	}
	params.Size = int(req.Params.Limit)

	ctx = context.WithValue(ctx, interceptor.ContextRequestIDKey, interceptor.ServerCtxRequestID(ctx)) //nolint
	records, total, err := s.iDao.GetByColumns(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "query params error:") {
			logger.Warn("GetByColumns error", logger.Err(err), logger.Any("params", params), interceptor.ServerCtxRequestIDField(ctx))
			return nil, ecode.StatusInvalidParams.Err()
		}
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("params", params), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	userExamples := []*serverNameExampleV1.UserExample{}
	for _, record := range records {
		data, err := convertUserExample(record)
		if err != nil {
			logger.Warn("convertUserExample error", logger.Err(err), logger.Any("id", record.ID), interceptor.ServerCtxRequestIDField(ctx))
			continue
		}
		userExamples = append(userExamples, data)
	}

	return &serverNameExampleV1.ListUserExampleReply{
		Total:        total,
		UserExamples: userExamples,
	}, nil
}

func convertUserExample(record *model.UserExample) (*serverNameExampleV1.UserExample, error) {
	value := &serverNameExampleV1.UserExample{}
	err := copier.Copy(value, record)
	if err != nil {
		return nil, err
	}
	value.Id = record.ID
	value.CreatedAt = record.CreatedAt.Unix()
	value.UpdatedAt = record.UpdatedAt.Unix()
	return value, nil
}
