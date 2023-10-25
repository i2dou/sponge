// Package types define the structure of request parameters and respond results in this package
package types

import (
	"time"

	"github.com/i2dou/sponge/pkg/mysql/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// todo generate the request and response struct to here
// delete the templates code start

// CreateUserExampleRequest request params
type CreateUserExampleRequest struct {
	Name     string `json:"name" binding:"min=2"`         // username
	Email    string `json:"email" binding:"email"`        // email
	Password string `json:"password" binding:"md5"`       // password
	Phone    string `json:"phone" binding:"e164"`         // phone number, e164 rules, e.g. +8612345678901
	Avatar   string `json:"avatar" binding:"min=5"`       // avatar
	Age      int    `json:"age" binding:"gt=0,lt=120"`    // age
	Gender   int    `json:"gender" binding:"gte=0,lte=2"` // gender, 1:Male, 2:Female, other values:unknown
}

// UpdateUserExampleByIDRequest request params
type UpdateUserExampleByIDRequest struct {
	ID       uint64 `json:"id" binding:"-"`      // id
	Name     string `json:"name" binding:""`     // username
	Email    string `json:"email" binding:""`    // email
	Password string `json:"password" binding:""` // password
	Phone    string `json:"phone" binding:""`    // phone number
	Avatar   string `json:"avatar" binding:""`   // avatar
	Age      int    `json:"age" binding:""`      // age
	Gender   int    `json:"gender" binding:""`   // gender, 1:Male, 2:Female, other values:unknown
}

// UserExampleObjDetail detail
type UserExampleObjDetail struct {
	ID        string    `json:"id"`        // id
	Name      string    `json:"name"`      // username
	Email     string    `json:"email"`     // email
	Phone     string    `json:"phone"`     // phone number
	Avatar    string    `json:"avatar"`    // avatar
	Age       int       `json:"age"`       // age
	Gender    int       `json:"gender"`    // gender, 1:Male, 2:Female, other values:unknown
	Status    int       `json:"status"`    // account status, 1:inactive, 2:activated, 3:blocked
	LoginAt   int64     `json:"loginAt"`   // login timestamp
	CreatedAt time.Time `json:"createdAt"` // create time
	UpdatedAt time.Time `json:"updatedAt"` // update time
}

// delete the templates code end

// CreateUserExampleRespond only for api docs
type CreateUserExampleRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateUserExampleByIDRespond only for api docs
type UpdateUserExampleByIDRespond struct {
	Result
}

// GetUserExampleByIDRespond only for api docs
type GetUserExampleByIDRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		UserExample UserExampleObjDetail `json:"userExample"`
	} `json:"data"` // return data
}

// DeleteUserExampleByIDRespond only for api docs
type DeleteUserExampleByIDRespond struct {
	Result
}

// DeleteUserExamplesByIDsRequest request params
type DeleteUserExamplesByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// DeleteUserExamplesByIDsRespond only for api docs
type DeleteUserExamplesByIDsRespond struct {
	Result
}

// GetUserExampleByConditionRequest request params
type GetUserExampleByConditionRequest struct {
	query.Conditions
}

// GetUserExampleByConditionRespond only for api docs
type GetUserExampleByConditionRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		UserExample UserExampleObjDetail `json:"userExample"`
	} `json:"data"` // return data
}

// ListUserExamplesByIDsRequest request params
type ListUserExamplesByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListUserExamplesByIDsRespond only for api docs
type ListUserExamplesByIDsRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		UserExamples []UserExampleObjDetail `json:"userExamples"`
	} `json:"data"` // return data
}

// ListUserExamplesRequest request params
type ListUserExamplesRequest struct {
	query.Params
}

// ListUserExamplesRespond only for api docs
type ListUserExamplesRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		UserExamples []UserExampleObjDetail `json:"userExamples"`
	} `json:"data"` // return data
}
