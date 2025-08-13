package usecases

import (
	"LOTestTask/pkg/createTaskPkg"
	"LOTestTask/pkg/getTaskByIDPkg"
	"LOTestTask/pkg/getTasksByStatusPkg"
)

type CreateTaskTestCase struct {
	TestName           string
	TestData           createTaskPkg.CreateTaskV0Request
	ExpectedError      error
	ExpectedOutput     createTaskPkg.CreateTaskV0Response
	ExpectedStatusCode int
}

type GetTaskByIDTestCase struct {
	TestName           string
	TestData           getTaskByIDPkg.GetTaskByIDV0Request
	ExpectedError      error
	ExpectedOutput     getTaskByIDPkg.GetTaskByIDV0Response
	ExpectedStatusCode int
}

type GetTasksByStatusTestCase struct {
	TestName           string
	TestData           getTasksByStatusPkg.GetTasksByStatusV0Request
	ExpectedError      error
	ExpectedOutput     getTasksByStatusPkg.GetTasksByStatusV0Response
	ExpectedStatusCode int
}
