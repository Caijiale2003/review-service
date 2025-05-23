// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package v1

import (
	fmt "fmt"
	errors "github.com/go-kratos/kratos/v2/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

// 为某个枚举单独设置错误码
func IsNeedLogin(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_NEED_LOGIN.String() && e.Code == 401
}

// 为某个枚举单独设置错误码
func ErrorNeedLogin(format string, args ...interface{}) *errors.Error {
	return errors.New(401, ErrorReason_NEED_LOGIN.String(), fmt.Sprintf(format, args...))
}

func IsDbError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_DB_ERROR.String() && e.Code == 402
}

func ErrorDbError(format string, args ...interface{}) *errors.Error {
	return errors.New(402, ErrorReason_DB_ERROR.String(), fmt.Sprintf(format, args...))
}

func IsOrderReviewed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_ORDER_REVIEWED.String() && e.Code == 400
}

func ErrorOrderReviewed(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_ORDER_REVIEWED.String(), fmt.Sprintf(format, args...))
}
