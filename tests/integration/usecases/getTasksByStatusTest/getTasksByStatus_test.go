package getTasksByStatusTest

import (
	"LOTestTask/internal/di"
	"LOTestTask/pkg"
	"LOTestTask/pkg/createTaskPkg"
	"LOTestTask/pkg/getTasksByStatusPkg"
	"LOTestTask/tests/integration/usecases"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGetTasksByStatus(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	initInfr(ctx, cancel)

	testCases := initTestCases()

	prepare()

	for _, tc := range testCases {
		t.Run(tc.TestName, func(t *testing.T) {
			path := strings.Builder{}
			path.WriteString("http://localhost:9191/tasks")
			for i, status := range tc.TestData.Statuses {
				if i == 0 {
					path.WriteString(fmt.Sprintf("?statuses=%s", status))
				}

				path.WriteString(fmt.Sprintf("&statuses=%s", status))
			}

			resp, err := http.Get(path.String())
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

				data := getTasksByStatusPkg.GetTasksByStatusV0Response{}

				if err := json.Unmarshal(body, &data); err != nil {
					t.Errorf("can not unmarshal json - %v", err)
				}

				if len(data.Tasks) != len(tc.ExpectedOutput.Tasks) {
					t.Errorf("wrong length of tasks, expected - %d, actual - %d", len(tc.ExpectedOutput.Tasks), len(data.Tasks))
				}

				for _, task := range data.Tasks {
					for _, outTask := range tc.ExpectedOutput.Tasks {
						if task.ID == outTask.ID {
							if task.Discription != outTask.Discription {
								t.Errorf("wrong output discription, expected - %s, actual - %s", outTask, task.Discription)
							}

							if task.AuthorID != outTask.AuthorID {
								t.Errorf("wrong output author id, expected - %s, actual - %s", outTask.AuthorID, task.AuthorID)
							}

							if string(task.Status) != string(outTask.Status) {
								t.Errorf("wrong output status, expected - %s, actual - %s", string(outTask.Status), string(task.Status))
							}

							if !(task.CreatedAt.Sub(outTask.CreatedAt) < 3*time.Second) || !(outTask.CreatedAt.Sub(task.CreatedAt) < 3*time.Second) {
								t.Errorf("wrong output created at, expected - %s, actual - %s", outTask.CreatedAt, task.CreatedAt)
							}

							if !(task.UpdatedAt.Sub(outTask.UpdatedAt) < 3*time.Second) || !(outTask.UpdatedAt.Sub(task.UpdatedAt) < 3*time.Second) {
								t.Errorf("wrong output updated at, expected - %s, actual - %s", outTask.UpdatedAt, task.UpdatedAt)
							}
						}
					}
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

func prepare() error {
	data1 := createTaskPkg.CreateTaskV0Request{
		ID:          "1",
		Discription: "test",
		AuthorID:    "11",
		Status:      "open",
	}

	data2 := createTaskPkg.CreateTaskV0Request{
		ID:          "2",
		Discription: "test",
		AuthorID:    "11",
		Status:      "done",
	}

	jsonData1, err := json.Marshal(data1)
	if err != nil {
		return err
	}

	jsonData2, err := json.Marshal(data2)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:9191/tasks", "application/json", bytes.NewBuffer(jsonData1))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("status code not 200")
	}

	resp2, err := http.Post("http://localhost:9191/tasks", "application/json", bytes.NewBuffer(jsonData2))
	if err != nil {
		return err
	}

	if resp2.StatusCode != 200 {
		return errors.New("status code not 200")
	}

	return nil
}

func initTestCases() []usecases.GetTasksByStatusTestCase {
	cases := []usecases.GetTasksByStatusTestCase{
		{
			TestName: "[TestGetTasksByStatus] successfull get",
			TestData: getTasksByStatusPkg.GetTasksByStatusV0Request{
				Statuses: []pkg.Status{"open", "done"},
			},
			ExpectedError: nil,
			ExpectedOutput: getTasksByStatusPkg.GetTasksByStatusV0Response{
				Tasks: []getTasksByStatusPkg.Task{
					{
						ID:          "1",
						Discription: "test",
						AuthorID:    "11",
						Status:      "open",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
					{
						ID:          "2",
						Discription: "test",
						AuthorID:    "11",
						Status:      "done",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				},
			},
			ExpectedStatusCode: 200,
		},
		{
			TestName: "[TestGetTasksByStatus] not found",
			TestData: getTasksByStatusPkg.GetTasksByStatusV0Request{
				Statuses: []pkg.Status{"in%20test"},
			},
			ExpectedError: nil,
			ExpectedOutput: getTasksByStatusPkg.GetTasksByStatusV0Response{
				Tasks: nil,
			},
			ExpectedStatusCode: 200,
		},
		{
			TestName: "[TestGetTasksByStatus] successfull get",
			TestData: getTasksByStatusPkg.GetTasksByStatusV0Request{
				Statuses: []pkg.Status{},
			},
			ExpectedError: nil,
			ExpectedOutput: getTasksByStatusPkg.GetTasksByStatusV0Response{
				Tasks: []getTasksByStatusPkg.Task{
					{
						ID:          "1",
						Discription: "test",
						AuthorID:    "11",
						Status:      "open",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
					{
						ID:          "2",
						Discription: "test",
						AuthorID:    "11",
						Status:      "done",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				},
			},
			ExpectedStatusCode: 200,
		},
		{
			TestName: "[TestGetTasksByStatus] wrong status",
			TestData: getTasksByStatusPkg.GetTasksByStatusV0Request{
				Statuses: []pkg.Status{"test"},
			},
			ExpectedError:      nil,
			ExpectedStatusCode: 400,
		},
	}

	return cases
}
