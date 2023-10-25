package handler

import (
	"net/http"
	"testing"
	"time"

	serverNameExampleV1 "github.com/i2dou/sponge/api/serverNameExample/v1"
	"github.com/i2dou/sponge/internal/cache"
	"github.com/i2dou/sponge/internal/dao"
	"github.com/i2dou/sponge/internal/ecode"
	"github.com/i2dou/sponge/internal/model"

	"github.com/i2dou/sponge/pkg/gin/response"
	"github.com/i2dou/sponge/pkg/gohttp"
	"github.com/i2dou/sponge/pkg/gotest"
	"github.com/i2dou/sponge/pkg/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/i2dou/sponge/api/types"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
)

func newUserExamplePbHandler() *gotest.Handler {
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

	// init mock handler
	h := gotest.NewHandler(d, testData)
	h.IHandler = &userExamplePbHandler{userExampleDao: d.IDao.(dao.UserExampleDao)}
	iHandler := h.IHandler.(serverNameExampleV1.UserExampleLogicer)

	testFns := []gotest.RouterInfo{
		{
			FuncName: "Create",
			Method:   http.MethodPost,
			Path:     "/userExample",
			HandlerFunc: func(c *gin.Context) {
				req := &serverNameExampleV1.CreateUserExampleRequest{}
				_ = c.ShouldBindJSON(req)
				_, err := iHandler.Create(c, req)
				if err != nil {
					response.Error(c, ecode.ErrCreateUserExample)
					return
				}
				response.Success(c)
			},
		},
		{
			FuncName: "DeleteByID",
			Method:   http.MethodDelete,
			Path:     "/userExample/:id",
			HandlerFunc: func(c *gin.Context) {
				req := &serverNameExampleV1.DeleteUserExampleByIDRequest{
					Id: utils.StrToUint64(c.Param("id")),
				}
				_, err := iHandler.DeleteByID(c, req)
				if err != nil {
					response.Error(c, ecode.ErrDeleteByIDUserExample)
					return
				}
				response.Success(c)
			},
		},
		{
			FuncName: "DeleteByIDs",
			Method:   http.MethodPost,
			Path:     "/userExample/delete/ids",
			HandlerFunc: func(c *gin.Context) {
				req := &serverNameExampleV1.DeleteUserExampleByIDsRequest{}
				_ = c.ShouldBindJSON(req)
				_, err := iHandler.DeleteByIDs(c, req)
				if err != nil {
					response.Error(c, ecode.ErrDeleteByIDsUserExample)
					return
				}
				response.Success(c)
			},
		},
		{
			FuncName: "UpdateByID",
			Method:   http.MethodPut,
			Path:     "/userExample/:id",
			HandlerFunc: func(c *gin.Context) {
				req := &serverNameExampleV1.UpdateUserExampleByIDRequest{}
				_ = c.ShouldBindJSON(req)
				req.Id = utils.StrToUint64(c.Param("id"))
				_, err := iHandler.UpdateByID(c, req)
				if err != nil {
					response.Error(c, ecode.ErrUpdateByIDUserExample)
					return
				}
				response.Success(c)
			},
		},
		{
			FuncName: "GetByID",
			Method:   http.MethodGet,
			Path:     "/userExample/:id",
			HandlerFunc: func(c *gin.Context) {
				req := &serverNameExampleV1.GetUserExampleByIDRequest{
					Id: utils.StrToUint64(c.Param("id")),
				}
				_, err := iHandler.GetByID(c, req)
				if err != nil {
					response.Error(c, ecode.ErrGetByIDUserExample)
					return
				}
				response.Success(c)
			},
		},
		{
			FuncName: "GetByCondition",
			Method:   http.MethodPost,
			Path:     "/userExample/condition",
			HandlerFunc: func(c *gin.Context) {
				req := &serverNameExampleV1.GetUserExampleByConditionRequest{}
				_ = c.ShouldBindJSON(req)
				_, err := iHandler.GetByCondition(c, req)
				if err != nil {
					response.Error(c, ecode.ErrGetByConditionUserExample)
					return
				}
				response.Success(c)
			},
		},
		{
			FuncName: "ListByIDs",
			Method:   http.MethodPost,
			Path:     "/userExample/list/ids",
			HandlerFunc: func(c *gin.Context) {
				req := &serverNameExampleV1.ListUserExampleByIDsRequest{}
				_ = c.ShouldBindJSON(req)
				_, err := iHandler.ListByIDs(c, req)
				if err != nil {
					response.Error(c, ecode.ErrListByIDsUserExample)
					return
				}
				response.Success(c)
			},
		},
		{
			FuncName: "List",
			Method:   http.MethodPost,
			Path:     "/userExample/list",
			HandlerFunc: func(c *gin.Context) {
				req := &serverNameExampleV1.ListUserExampleRequest{}
				_ = c.ShouldBindJSON(req)
				_, err := iHandler.List(c, req)
				if err != nil {
					response.Error(c, ecode.ErrListUserExample)
					return
				}
				response.Success(c)
			},
		},
	}

	h.GoRunHTTPServer(testFns)

	time.Sleep(time.Millisecond * 200)
	return h
}

func Test_userExamplePbHandler_Create(t *testing.T) {
	h := newUserExamplePbHandler()
	defer h.Close()
	testData := &serverNameExampleV1.CreateUserExampleRequest{}
	_ = copier.Copy(testData, h.TestData.(*model.UserExample))

	h.MockDao.SQLMock.ExpectBegin()
	args := h.MockDao.GetAnyArgs(h.TestData)
	h.MockDao.SQLMock.ExpectExec("INSERT INTO .*").
		WithArgs(args[:len(args)-1]...). // adjusted for the amount of test data
		WillReturnResult(sqlmock.NewResult(1, 1))
	h.MockDao.SQLMock.ExpectCommit()

	result := &gohttp.StdResult{}
	err := gohttp.Post(result, h.GetRequestURL("Create"), testData)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", result)
	// delete the templates code start
	result = &gohttp.StdResult{}
	testData = &serverNameExampleV1.CreateUserExampleRequest{
		Name:     "foo",
		Password: "f447b20a7fcbf53a5d5be013ea0b15af",
		Email:    "foo@bar.com",
		Phone:    "16000000001",
		Avatar:   "http://foo/1.jpg",
		Age:      10,
		Gender:   1,
	}
	err = gohttp.Post(result, h.GetRequestURL("Create"), testData)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", result)

	h.MockDao.SQLMock.ExpectBegin()
	h.MockDao.SQLMock.ExpectCommit()
	// create error test
	result = &gohttp.StdResult{}
	err = gohttp.Post(result, h.GetRequestURL("Create"), testData)
	assert.NoError(t, err)
	// delete the templates code end
}

func Test_userExamplePbHandler_DeleteByID(t *testing.T) {
	h := newUserExamplePbHandler()
	defer h.Close()
	testData := h.TestData.(*model.UserExample)

	h.MockDao.SQLMock.ExpectBegin()
	h.MockDao.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(h.MockDao.AnyTime, testData.ID). // adjusted for the amount of test data
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	h.MockDao.SQLMock.ExpectCommit()

	result := &gohttp.StdResult{}
	err := gohttp.Delete(result, h.GetRequestURL("DeleteByID", testData.ID))
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}

	// zero id error test
	err = gohttp.Delete(result, h.GetRequestURL("DeleteByID", 0))
	assert.NoError(t, err)

	// delete error test
	err = gohttp.Delete(result, h.GetRequestURL("DeleteByID", 111))
	assert.NoError(t, err)
}

func Test_userExamplePbHandler_DeleteByIDs(t *testing.T) {
	h := newUserExamplePbHandler()
	defer h.Close()
	testData := h.TestData.(*model.UserExample)

	h.MockDao.SQLMock.ExpectBegin()
	h.MockDao.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(h.MockDao.AnyTime, testData.ID). // adjusted for the amount of test data
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	h.MockDao.SQLMock.ExpectCommit()

	result := &gohttp.StdResult{}
	err := gohttp.Post(result, h.GetRequestURL("DeleteByIDs"), &serverNameExampleV1.DeleteUserExampleByIDsRequest{Ids: []uint64{testData.ID}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}

	// zero id error test
	err = gohttp.Post(result, h.GetRequestURL("DeleteByIDs"), &serverNameExampleV1.DeleteUserExampleByIDsRequest{})
	assert.NoError(t, err)

	// get error test
	err = gohttp.Post(result, h.GetRequestURL("DeleteByIDs"), &serverNameExampleV1.DeleteUserExampleByIDsRequest{Ids: []uint64{111}})
	assert.NoError(t, err)
}

func Test_userExamplePbHandler_UpdateByID(t *testing.T) {
	h := newUserExamplePbHandler()
	defer h.Close()
	testData := &serverNameExampleV1.UpdateUserExampleByIDRequest{}
	_ = copier.Copy(testData, h.TestData.(*model.UserExample))
	testData.Id = h.TestData.(*model.UserExample).ID

	h.MockDao.SQLMock.ExpectBegin()
	h.MockDao.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(h.MockDao.AnyTime, testData.Id). // adjusted for the amount of test data
		WillReturnResult(sqlmock.NewResult(int64(testData.Id), 1))
	h.MockDao.SQLMock.ExpectCommit()

	result := &gohttp.StdResult{}
	err := gohttp.Put(result, h.GetRequestURL("UpdateByID", testData.Id), testData)
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}

	// zero id error test
	err = gohttp.Put(result, h.GetRequestURL("UpdateByID", 0), testData)
	assert.NoError(t, err)

	// update error test
	err = gohttp.Put(result, h.GetRequestURL("UpdateByID", 111), testData)
	assert.NoError(t, err)
}

func Test_userExamplePbHandler_GetByID(t *testing.T) {
	h := newUserExamplePbHandler()
	defer h.Close()
	testData := h.TestData.(*model.UserExample)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	h.MockDao.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	result := &gohttp.StdResult{}
	err := gohttp.Get(result, h.GetRequestURL("GetByID", testData.ID))
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}

	// zero id error test
	err = gohttp.Get(result, h.GetRequestURL("GetByID", 0))
	assert.NoError(t, err)

	// get error test
	err = gohttp.Get(result, h.GetRequestURL("GetByID", 111))
	assert.NoError(t, err)
}

func Test_userExamplePbHandler_GetByCondition(t *testing.T) {
	h := newUserExamplePbHandler()
	defer h.Close()
	testData := h.TestData.(*model.UserExample)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	h.MockDao.SQLMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	result := &gohttp.StdResult{}
	err := gohttp.Post(result, h.GetRequestURL("GetByCondition"), &serverNameExampleV1.GetUserExampleByConditionRequest{
		Conditions: &types.Conditions{
			Columns: []*types.Column{
				{
					Name:  "id",
					Value: utils.Uint64ToStr(testData.ID),
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}

	// zero error test
	err = gohttp.Post(result, h.GetRequestURL("GetByCondition"), nil)
	assert.NoError(t, err)

	// valid error test
	err = gohttp.Post(result, h.GetRequestURL("GetByCondition"), &serverNameExampleV1.GetUserExampleByConditionRequest{
		Conditions: &types.Conditions{
			Columns: []*types.Column{
				{
					Name:  "id",
					Value: "111",
					Exp:   "unknown",
				},
			},
		},
	})

	// get error test
	err = gohttp.Post(result, h.GetRequestURL("GetByCondition"), &serverNameExampleV1.GetUserExampleByConditionRequest{
		Conditions: &types.Conditions{
			Columns: []*types.Column{
				{
					Name:  "id",
					Value: "111",
				},
			},
		},
	})
	assert.NoError(t, err)
}

func Test_userExamplePbHandler_ListByIDs(t *testing.T) {
	h := newUserExamplePbHandler()
	defer h.Close()
	testData := h.TestData.(*model.UserExample)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	h.MockDao.SQLMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	result := &gohttp.StdResult{}
	err := gohttp.Post(result, h.GetRequestURL("ListByIDs"), &serverNameExampleV1.ListUserExampleByIDsRequest{Ids: []uint64{testData.ID}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}

	// zero id error test
	err = gohttp.Post(result, h.GetRequestURL("ListByIDs"), &serverNameExampleV1.ListUserExampleByIDsRequest{})
	assert.NoError(t, err)

	// get error test
	err = gohttp.Post(result, h.GetRequestURL("ListByIDs"), &serverNameExampleV1.ListUserExampleByIDsRequest{Ids: []uint64{111}})
	assert.NoError(t, err)
}

func Test_userExamplePbHandler_List(t *testing.T) {
	h := newUserExamplePbHandler()
	defer h.Close()
	testData := h.TestData.(*model.UserExample)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	h.MockDao.SQLMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	result := &gohttp.StdResult{}
	err := gohttp.Post(result, h.GetRequestURL("List"), &serverNameExampleV1.ListUserExampleRequest{
		Params: &types.Params{
			Page:  0,
			Limit: 10,
			Sort:  "ignore count", // ignore test count
		}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}

	// nil params error test
	err = gohttp.Post(result, h.GetRequestURL("List"), &serverNameExampleV1.ListUserExampleRequest{})
	assert.NoError(t, err)

	// get error test
	err = gohttp.Post(result, h.GetRequestURL("List"), &serverNameExampleV1.ListUserExampleRequest{Params: &types.Params{
		Page:  0,
		Limit: 10,
	}})
	assert.NoError(t, err)
}

func TestNewUserExamplePbHandler(t *testing.T) {
	defer func() {
		recover()
	}()
	_ = NewUserExamplePbHandler()
}
