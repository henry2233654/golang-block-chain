package services

import (
	"fmt"
	"golang-block-chain/repositories"
)

const (
	UniqueKeyIsDuplicate = "101" // 某個唯一欄位值重複

	NotExist      = "201" // 資源不存在
	BeingUsed     = "202" // 資源正被使用中
	NotAllowed    = "203" // 不允許的行為
	IncorrectData = "204" // 不正確的資料
)

type ServiceError interface {
	error
	Details() []*ServiceErrorDetail
}

type ServiceErrorDetail struct {
	Index   *int        `json:"index"`
	Field   *string     `json:"field"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Err     interface{} `json:"err"`
}

func (err *ServiceErrorDetail) Error() string {
	return err.Message
}

func NewServiceErrorDetail(index *int, field *string, code string, message string, err interface{}) *ServiceErrorDetail {
	return &ServiceErrorDetail{
		Index:   index,
		Field:   field,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

type ValidateError struct {
	ServiceError
	Errors []*ServiceErrorDetail
}

func (err *ValidateError) Error() string {
	return "invalid data"
}

func (err *ValidateError) Details() []*ServiceErrorDetail {
	return err.Errors
}

func NewValidateError(errors []*ServiceErrorDetail) *ValidateError {
	return &ValidateError{Errors: errors}
}

// type InvalidDataError struct {
// 	error
// 	Invalids map[string]interface{}
// }

// func (err *InvalidDataError) Error() string {
// 	return "invalid data"
// }

type NotExistError struct {
	ServiceError
	Index        *int
	ResourceName string
	Expected     interface{}
}

func (err *NotExistError) Error() string {
	return fmt.Sprintf("no [%s] was found matching the [%v]", err.ResourceName, err.Expected)
}

func (err *NotExistError) Details() []*ServiceErrorDetail {
	return []*ServiceErrorDetail{
		{
			Index:   err.Index,
			Field:   nil,
			Code:    NotExist,
			Message: fmt.Sprintf("no [%s] was found matching the [%v]", err.ResourceName, err.Expected),
		},
	}
}

func NewNotExistError(index *int, resourceName string, expected interface{}) *NotExistError {
	return &NotExistError{Index: index, ResourceName: resourceName, Expected: expected}
}

type DuplicateError struct {
	ServiceError
	Index *int
	Err   *repositories.UniqueConstrainError
}

func (e *DuplicateError) Error() string {
	return "some unique key is duplicated"
}

func (err *DuplicateError) Details() []*ServiceErrorDetail {
	return []*ServiceErrorDetail{
		{
			Index:   err.Index,
			Field:   nil,
			Code:    UniqueKeyIsDuplicate,
			Message: "some unique key is duplicated",
			Err:     err.Err,
		},
	}
}

func NewDuplicateError(index *int, err *repositories.UniqueConstrainError) *DuplicateError {
	return &DuplicateError{
		Index: index,
		Err:   err,
	}
}
