package service

import (
	"testing"
	"time"

	serverNameExampleV1 "github.com/i2dou/sponge/api/serverNameExample/v1"
	"github.com/i2dou/sponge/api/types"
	"github.com/i2dou/sponge/internal/cache"
	"github.com/i2dou/sponge/internal/dao"
	"github.com/i2dou/sponge/internal/model"

	"github.com/i2dou/sponge/pkg/gotest"
	"github.com/i2dou/sponge/pkg/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
)

func newUserExampleService() *gotest.Service {
	// todo additional test field information
	testData := &model.UserExample{}
	testData.ID = 1
	testData.CreatedAt = time.Now()
	testData.UpdatedAt = testData.CreatedAt

	// init mock cache
	c := gotest.NewCache(map[string]interface{}{utils.Uint64ToStr(testData.ID): testData})
	c.ICache = cache.NewUserExampleCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})

	// init mock dao
	d := gotest.NewDao(c, testData)
	d.IDao = dao.NewUserExampleDao(d.DB, c.ICache.(cache.UserExampleCache))

	// init mock service
	s := gotest.NewService(d, testData)
	serverNameExampleV1.RegisterUserExampleServer(s.Server, &userExample{
		UnimplementedUserExampleServer: serverNameExampleV1.UnimplementedUserExampleServer{},
		iDao:                           d.IDao.(dao.UserExampleDao),
	})

	// start up rpc server
	s.GoGrpcServer()
	time.Sleep(time.Millisecond * 100)

	// grpc client
	s.IServiceClient = serverNameExampleV1.NewUserExampleClient(s.GetClientConn())

	return s
}

func Test_userExampleService_Create(t *testing.T) {
	s := newUserExampleService()
	defer s.Close()
	testData := &serverNameExampleV1.CreateUserExampleRequest{}
	_ = copier.Copy(testData, s.TestData.(*model.UserExample))

	s.MockDao.SQLMock.ExpectBegin()
	args := s.MockDao.GetAnyArgs(s.TestData)
	s.MockDao.SQLMock.ExpectExec("INSERT INTO .*").
		WithArgs(args[:len(args)-1]...). // Modified according to the actual number of parameters
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.MockDao.SQLMock.ExpectCommit()

	reply, err := s.IServiceClient.(serverNameExampleV1.UserExampleClient).Create(s.Ctx, testData)
	t.Log(err, reply.String())

	// delete the templates code start
	testData = &serverNameExampleV1.CreateUserExampleRequest{
		Name:     "foo",
		Password: "f447b20a7fcbf53a5d5be013ea0b15af",
		Email:    "foo@bar.com",
		Phone:    "16000000001",
		Avatar:   "http://foo/1.jpg",
		Age:      10,
		Gender:   1,
	}
	reply, err = s.IServiceClient.(serverNameExampleV1.UserExampleClient).Create(s.Ctx, testData)
	t.Log(err, reply.String())

	s.MockDao.SQLMock.ExpectBegin()
	s.MockDao.SQLMock.ExpectCommit()
	reply, err = s.IServiceClient.(serverNameExampleV1.UserExampleClient).Create(s.Ctx, testData)
	assert.Error(t, err)
	// delete the templates code end
}

func Test_userExampleService_DeleteByID(t *testing.T) {
	s := newUserExampleService()
	defer s.Close()
	testData := &serverNameExampleV1.DeleteUserExampleByIDRequest{
		Id: s.TestData.(*model.UserExample).ID,
	}

	s.MockDao.SQLMock.ExpectBegin()
	s.MockDao.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(s.MockDao.AnyTime, testData.Id). // Modified according to the actual number of parameters
		WillReturnResult(sqlmock.NewResult(int64(testData.Id), 1))
	s.MockDao.SQLMock.ExpectCommit()

	reply, err := s.IServiceClient.(serverNameExampleV1.UserExampleClient).DeleteByID(s.Ctx, testData)
	assert.NoError(t, err)
	t.Log(reply.String())

	// zero id error test
	testData.Id = 0
	reply, err = s.IServiceClient.(serverNameExampleV1.UserExampleClient).DeleteByID(s.Ctx, testData)
	assert.Error(t, err)

	// delete error test
	testData.Id = 111
	reply, err = s.IServiceClient.(serverNameExampleV1.UserExampleClient).DeleteByID(s.Ctx, testData)
	assert.Error(t, err)
}

func Test_userExampleService_DeleteByIDs(t *testing.T) {
	s := newUserExampleService()
	defer s.Close()
	data := s.TestData.(*model.UserExample)
	testData := &serverNameExampleV1.DeleteUserExampleByIDsRequest{
		Ids: []uint64{data.ID},
	}

	s.MockDao.SQLMock.ExpectBegin()
	s.MockDao.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(s.MockDao.AnyTime, data.ID). // Modified according to the actual number of parameters
		WillReturnResult(sqlmock.NewResult(int64(data.ID), 1))
	s.MockDao.SQLMock.ExpectCommit()

	reply, err := s.IServiceClient.(serverNameExampleV1.UserExampleClient).DeleteByIDs(s.Ctx, testData)
	assert.NoError(t, err)
	t.Log(reply.String())

	// zero id error test
	testData.Ids = []uint64{}
	reply, err = s.IServiceClient.(serverNameExampleV1.UserExampleClient).DeleteByIDs(s.Ctx, testData)
	assert.Error(t, err)

	// delete error test
	testData.Ids = []uint64{111}
	reply, err = s.IServiceClient.(serverNameExampleV1.UserExampleClient).DeleteByIDs(s.Ctx, testData)
	assert.Error(t, err)
}

func Test_userExampleService_UpdateByID(t *testing.T) {
	s := newUserExampleService()
	defer s.Close()
	data := s.TestData.(*model.UserExample)
	testData := &serverNameExampleV1.UpdateUserExampleByIDRequest{}
	_ = copier.Copy(testData, s.TestData.(*model.UserExample))
	testData.Id = data.ID

	s.MockDao.SQLMock.ExpectBegin()
	s.MockDao.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(s.MockDao.AnyTime, testData.Id). // Modified according to the actual number of parameters
		WillReturnResult(sqlmock.NewResult(int64(testData.Id), 1))
	s.MockDao.SQLMock.ExpectCommit()

	reply, err := s.IServiceClient.(serverNameExampleV1.UserExampleClient).UpdateByID(s.Ctx, testData)
	assert.NoError(t, err)
	t.Log(reply.String())

	// zero id error test
	testData.Id = 0
	reply, err = s.IServiceClient.(serverNameExampleV1.UserExampleClient).UpdateByID(s.Ctx, testData)
	assert.Error(t, err)

	// upate error test
	testData.Id = 111
	reply, err = s.IServiceClient.(serverNameExampleV1.UserExampleClient).UpdateByID(s.Ctx, testData)
	assert.Error(t, err)
}

func Test_userExampleService_GetByID(t *testing.T) {
	s := newUserExampleService()
	defer s.Close()
	data := s.TestData.(*model.UserExample)
	testData := &serverNameExampleV1.GetUserExampleByIDRequest{
		Id: data.ID,
	}

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(data.ID, data.CreatedAt, data.UpdatedAt)

	s.MockDao.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.Id).
		WillReturnRows(rows)

	reply, err := s.IServiceClient.(serverNameExampleV1.UserExampleClient).GetByID(s.Ctx, testData)
	assert.NoError(t, err)
	t.Log(reply.String())

	// zero id error test
	testData.Id = 0
	reply, err = s.IServiceClient.(serverNameExampleV1.UserExampleClient).GetByID(s.Ctx, testData)
	assert.Error(t, err)

	// get error test
	testData.Id = 111
	reply, err = s.IServiceClient.(serverNameExampleV1.UserExampleClient).GetByID(s.Ctx, testData)
	assert.Error(t, err)
}

func Test_userExampleService_GetByCondition(t *testing.T) {
	s := newUserExampleService()
	defer s.Close()
	testData := s.TestData.(*model.UserExample)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	s.MockDao.SQLMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	reply, err := s.IServiceClient.(serverNameExampleV1.UserExampleClient).GetByCondition(s.Ctx, &serverNameExampleV1.GetUserExampleByConditionRequest{
		Conditions: &types.Conditions{
			Columns: []*types.Column{{
				Name:  "id",
				Value: "1",
			}},
		},
	})
	assert.NoError(t, err)
	t.Log(reply.String())

	// get error test
	reply, err = s.IServiceClient.(serverNameExampleV1.UserExampleClient).GetByCondition(s.Ctx, &serverNameExampleV1.GetUserExampleByConditionRequest{
		Conditions: &types.Conditions{},
	})
	assert.Error(t, err)
}

func Test_userExampleService_ListByIDs(t *testing.T) {
	s := newUserExampleService()
	defer s.Close()
	data := s.TestData.(*model.UserExample)
	testData := &serverNameExampleV1.ListUserExampleByIDsRequest{
		Ids: []uint64{data.ID},
	}

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(data.ID, data.CreatedAt, data.UpdatedAt)

	s.MockDao.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(data.ID).
		WillReturnRows(rows)

	reply, err := s.IServiceClient.(serverNameExampleV1.UserExampleClient).ListByIDs(s.Ctx, testData)
	assert.NoError(t, err)
	t.Log(reply.String())

	// get error test
	testData.Ids = []uint64{111}
	reply, err = s.IServiceClient.(serverNameExampleV1.UserExampleClient).ListByIDs(s.Ctx, testData)
	assert.Error(t, err)
}

func Test_userExampleService_List(t *testing.T) {
	s := newUserExampleService()
	defer s.Close()
	testData := s.TestData.(*model.UserExample)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	s.MockDao.SQLMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	reply, err := s.IServiceClient.(serverNameExampleV1.UserExampleClient).List(s.Ctx, &serverNameExampleV1.ListUserExampleRequest{
		Params: &types.Params{
			Page:  0,
			Limit: 10,
			Sort:  "ignore count", // ignore test count
		},
	})
	assert.NoError(t, err)
	t.Log(reply.String())

	// get error test
	reply, err = s.IServiceClient.(serverNameExampleV1.UserExampleClient).List(s.Ctx, &serverNameExampleV1.ListUserExampleRequest{
		Params: &types.Params{
			Page:  0,
			Limit: 10,
		},
	})
	assert.Error(t, err)
}

func Test_convertUserExample(t *testing.T) {
	testData := &model.UserExample{}
	testData.ID = 1
	testData.CreatedAt = time.Now()
	testData.UpdatedAt = testData.CreatedAt

	data, err := convertUserExample(testData)
	assert.NoError(t, err)

	t.Logf("%+v", data)
}
