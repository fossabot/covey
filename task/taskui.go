package task

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/task/types"
	"github.com/chabad360/covey/ui"
	"github.com/go-playground/pure/v5"
	json "github.com/json-iterator/go"
)

func uiTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []types.Task
	err := storage.DB.QueryRow(context.Background(),
		"SELECT jsonb_agg(to_jsonb(tasks)) FROM tasks").Scan(&tasks)
	if err != nil {
		common.ErrorWriter(w, err)
	}
	p := &ui.Page{
		Title:   "Tasks",
		URL:     strings.Split(r.URL.Path, "/"),
		Details: struct{ Tasks []types.Task }{Tasks: tasks},
	}
	t := ui.GetTemplate("tasksAll")
	err = t.ExecuteTemplate(w, "base", p)
	if err != nil {
		common.ErrorWriter(w, err)
	}
}

func uiTaskSingle(w http.ResponseWriter, r *http.Request) {
	vars := pure.RequestVars(r)
	task, ok := GetTask(vars.URLParam("task"))
	if !ok {
		common.ErrorWriter404(w, vars.URLParam("task"))
		return
	}

	p := &ui.Page{
		Title:   fmt.Sprintf("Task %s", vars.URLParam("task")),
		URL:     strings.Split(r.URL.Path, "/"),
		Details: struct{ Task types.ITask }{Task: task},
	}

	t := ui.GetTemplate("tasksSingle")
	err := t.ExecuteTemplate(w, "base", p)
	if err != nil {
		common.ErrorWriter(w, err)
	}
}

func UITaskNew(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var nodes []string
		storage.DB.QueryRow(context.Background(), "SELECT jsonb_agg(name) FROM nodes;").Scan(&nodes)

		p := &ui.Page{
			Title: fmt.Sprintf("New Task"),
			URL:   strings.Split(r.URL.Path, "/"),
			Details: struct {
				Plugins []string
				Nodes   []string
			}{[]string{"shell"}, nodes},
		}

		t := ui.GetTemplate("tasksNew")
		err := t.ExecuteTemplate(w, "base", p)
		if err != nil {
			common.ErrorWriter(w, err)
		}
	} else {
		u, err := url.QueryUnescape(r.FormValue("command"))
		if err != nil {
			common.ErrorWriter(w, err)
			return
		}
		task, err := json.Marshal(map[string]interface{}{
			"node":    r.FormValue("node"),
			"plugin":  r.FormValue("plugin"),
			"details": map[string]string{"command": u}})
		if err != nil {
			common.ErrorWriter(w, err)
			return
		}
		t, err := NewTask(task)
		if err != nil {
			common.ErrorWriter(w, err)
			return
		}

		http.Redirect(w, r, "/tasks/"+t.GetID(), http.StatusFound)
	}
}

// RegisterUIHandlers registers the HTTP handlers for the task UI.
func RegisterUIHandlers(r pure.IRouteGroup) {
	r.Get("", uiTasks)
	r.Get("/:task", uiTaskSingle)
}
