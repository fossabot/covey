package task

import (
	"context"
	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/test"
	"log"
	"os"
	"testing"
)

var task = &models.Task{
	ID:       "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
	State:    models.StateRunning,
	Plugin:   "test",
	Details:  map[string]string{"test": "test"},
	ExitCode: 0,
}

func TestAddTask(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want string
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
			`{"id": "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", "log": null, "node": "test", "time": "2000-01-01T01:01:01.000000001Z", "state": 2, "plugin": "test", "details": {"test": "test"}, "exit_code": 0}`},
		{"3", ""},
	}
	//revive:enable:line-length-limit

	testError := addTask(task)

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			var got []byte
			if db.Where("id = ?", tt.id).First(&got); string(got) != tt.want {
				t.Errorf("addTask() = %v, want %v, error: %v", string(got), tt.want, testError)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want string
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
			`{"id": "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", "log": ["hello", "world"], "node": "test", "time": "2000-01-01T01:01:01.000000001Z", "state": 2, "plugin": "test", "details": {"test": "test"}, "exit_code": 0}`},
		{"3", ""},
	}
	//revive:enable:line-length-limit

	tu := task
	tu.Log = []string{"hello", "world"}
	saveTask(tu)

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			var got []byte
			if db.QueryRow(context.Background(), "SELECT to_jsonb(tasks) - 'id_short' FROM tasks WHERE id = $1;",
				tt.id).Scan(&got); string(got) != tt.want {
				t.Errorf("updateTask() = %v, want %v, error: %v", string(got), tt.want, testError)
			}
		})
	}
}

func TestGetTaskJSON(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want string
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
			`{"id": "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", "log": ["hello", "world"], "node": "test", "time": "2000-01-01T01:01:01.000000001Z", "state": 2, "plugin": "test", "details": {"test": "test"}, "exit_code": 0}`},
		{"3", ""},
	}
	//revive:enable:line-length-limit

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			if got, err := getTask(tt.id); string(got) != tt.want {
				t.Errorf("getTaskJSON() = %v, want %v, error: %v", string(got), tt.want, err)
			}
		})
	}
}

func TestMain(m *testing.M) {
	pool, resource, pdb, err := test.Boilerplate()
	db = pdb
	storage.DB = pdb
	if err != nil {
		log.Fatalf("Could not setup DB connection: %s", err)
	}

	db.Create(task)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
