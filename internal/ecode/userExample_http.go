package ecode

import (
	"github.com/i2dou/sponge/pkg/errcode"
)

// userExample business-level http error codes.
// the userExampleNO value range is 1~100, if the same number appears, it will cause a failure to start the service.
var (
	userExampleNO       = 1
	userExampleName     = "userExample"
	userExampleBaseCode = errcode.HCode(userExampleNO)

	ErrCreateUserExample         = errcode.NewError(userExampleBaseCode+1, "failed to create "+userExampleName)
	ErrDeleteByIDUserExample     = errcode.NewError(userExampleBaseCode+2, "failed to delete "+userExampleName)
	ErrDeleteByIDsUserExample    = errcode.NewError(userExampleBaseCode+3, "failed to delete by batch ids "+userExampleName)
	ErrUpdateByIDUserExample     = errcode.NewError(userExampleBaseCode+4, "failed to update "+userExampleName)
	ErrGetByIDUserExample        = errcode.NewError(userExampleBaseCode+5, "failed to get "+userExampleName+" details")
	ErrGetByConditionUserExample = errcode.NewError(userExampleBaseCode+6, "failed to get "+userExampleName+" details by conditions")
	ErrListByIDsUserExample      = errcode.NewError(userExampleBaseCode+7, "failed to list by batch ids "+userExampleName)
	ErrListUserExample           = errcode.NewError(userExampleBaseCode+8, "failed to list of "+userExampleName)
	// error codes are globally unique, adding 1 to the previous error code
)
