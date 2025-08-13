package createTaskTest

import (
	"LOTestTask/internal/di"
	"LOTestTask/pkg/createTaskPkg"
	"LOTestTask/tests/integration/usecases"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestCreateTask(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	initInfr(ctx, cancel)

	testCases := initTestCases()

	for _, tc := range testCases {
		t.Run(tc.TestName, func(t *testing.T) {
			jsonData, err := json.Marshal(tc.TestData)
			if err != nil {
				t.Errorf("can not marshal json - %v", err)
			}

			resp, err := http.Post("http://localhost:9191/tasks", "application/json", bytes.NewBuffer(jsonData))
			if err != tc.ExpectedError {
				t.Errorf("wrong error, expected - %v, actual - %v", tc.ExpectedError, err)
			}

			if resp.StatusCode != tc.ExpectedStatusCode {
				t.Errorf("wrong status code, expected - %d, actual - %d", tc.ExpectedStatusCode, resp.StatusCode)
			}

			if resp.StatusCode == 200 {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("can not read body - %v", err)
				}

				data := createTaskPkg.CreateTaskV0Response{}

				if err := json.Unmarshal(body, &data); err != nil {
					t.Errorf("can not unmarshal json - %v", err)
				}

				if data != tc.ExpectedOutput {
					t.Errorf("wrong output, expected - %v, actual - %v", tc.ExpectedOutput, data)
				}
			}

			resp.Body.Close()
		})
	}

	cancel()
}

func initInfr(ctx context.Context, cancel context.CancelFunc) {
	os.Setenv("CACHE_SPACE", "30")
	os.Setenv("SERVER_ADDR", "localhost")
	os.Setenv("SERVER_PORT", "9191")
	os.Setenv("LOGGER_LOG_LEVEL", "info")

	di := di.DI{}
	if err := di.Init(); err != nil {
		log.Fatalf("can not init service - %v", err)
	}

	go func() {
		if err := di.Start(); err != nil {
			log.Printf("error ocured while starting server - %v", err)
			cancel()
		}
	}()

	go func() {
		di.Stop(ctx)
	}()
}

func initTestCases() []usecases.CreateTaskTestCase {
	cases := []usecases.CreateTaskTestCase{
		{
			TestName: "[TestCreateTask] successfull creation",
			TestData: createTaskPkg.CreateTaskV0Request{
				ID:          "1",
				Discription: "test",
				AuthorID:    "11",
				Status:      "open",
			},
			ExpectedError: nil,
			ExpectedOutput: createTaskPkg.CreateTaskV0Response{
				ID: "1",
			},
			ExpectedStatusCode: 200,
		},
		{
			TestName: "[TestCreateTask] wrong status (bad request)",
			TestData: createTaskPkg.CreateTaskV0Request{
				ID:          "2",
				Discription: "test",
				AuthorID:    "11",
				Status:      "test",
			},
			ExpectedError:      nil,
			ExpectedStatusCode: 400,
		},
		{
			TestName: "[TestCreateTask] duplicate id (bad request)",
			TestData: createTaskPkg.CreateTaskV0Request{
				ID:          "1",
				Discription: "test",
				AuthorID:    "11",
				Status:      "done",
			},
			ExpectedError:      nil,
			ExpectedStatusCode: 400,
		},
		{
			TestName: "[TestCreateTask] without authorID (bad request)",
			TestData: createTaskPkg.CreateTaskV0Request{
				ID:          "2",
				Discription: "test",
				Status:      "done",
			},
			ExpectedError:      nil,
			ExpectedStatusCode: 400,
		},
		{
			TestName: "[TestCreateTask] without id (bad request)",
			TestData: createTaskPkg.CreateTaskV0Request{
				Discription: "test",
				AuthorID:    "11",
				Status:      "done",
			},
			ExpectedError:      nil,
			ExpectedStatusCode: 400,
		},
		{
			TestName: "[TestCreateTask] without status (bad request)",
			TestData: createTaskPkg.CreateTaskV0Request{
				ID:          "2",
				Discription: "test",
				AuthorID:    "11",
			},
			ExpectedError:      nil,
			ExpectedStatusCode: 400,
		},
		{
			TestName: "[TestCreateTask] successfull creation without discription",
			TestData: createTaskPkg.CreateTaskV0Request{
				ID:       "2",
				AuthorID: "11",
				Status:   "open",
			},
			ExpectedError: nil,
			ExpectedOutput: createTaskPkg.CreateTaskV0Response{
				ID: "2",
			},
			ExpectedStatusCode: 200,
		},
	}

	return cases
}
