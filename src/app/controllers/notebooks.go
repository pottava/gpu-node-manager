package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	gc "github.com/pottava/gpu-node-manager/src/app/googlecloud"
	"github.com/pottava/gpu-node-manager/src/app/util"
	"github.com/revel/revel"
)

type Notebooks struct {
	*revel.Controller
}

func (c Notebooks) Index() revel.Result {
	return c.Render()
}

func (c Notebooks) Create() revel.Result {
	isNotProduction := (util.AppStage() != "prod")
	return c.Render(isNotProduction)
}

func (c Notebooks) ListAPI() revel.Result {
	token, err := gc.Auth(c.Request)
	if err != nil {
		c.Log.Errorf("Failed to get notebooks: %v", err)
		c.Response.SetStatus(http.StatusUnauthorized)
		return c.RenderError(errors.New("認証エラー"))
	}
	email := gc.VerifiedEmail(token)

	notebooks, err := gc.GetNotebooks(c.Request.Context(), email)
	if err != nil {
		c.Log.Errorf("Failed to get notebooks: %v", err)
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderError(errors.New("内部エラー"))
	}

	type Result struct {
		Menu      string `json:"menu"`
		Runtime   string `json:"runtime"`
		ProxyUri  string `json:"proxyUri"`
		State     string `json:"state"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
	results := []*Result{}
	wg := &sync.WaitGroup{}
	wg.Add(len(notebooks))
	for _, note := range notebooks {
		go func(note *gc.Notebook) {
			runtime, err := gc.DescribeManagedNotebook(c.Request.Context(), note.Runtime)
			if err == nil {
				results = append(results, &Result{
					Menu:      note.Menu,
					Runtime:   note.Runtime,
					ProxyUri:  runtime.AccessConfig.ProxyUri,
					State:     runtime.State.String(),
					CreatedAt: util.DateToStr(note.CreatedAt),
					UpdatedAt: util.DateToStr(note.UpdatedAt),
				})
			} else {
				if note.Active {
					results = append(results, &Result{
						Menu:      note.Menu,
						Runtime:   note.Runtime,
						ProxyUri:  "",
						State:     err.Error(),
						CreatedAt: util.DateToStr(note.CreatedAt),
						UpdatedAt: util.DateToStr(note.UpdatedAt),
					})
				} else {
					results = append(results, &Result{
						Menu:      note.Menu,
						Runtime:   note.Runtime,
						ProxyUri:  "",
						State:     "DELETED",
						CreatedAt: util.DateToStr(note.CreatedAt),
						UpdatedAt: util.DateToStr(note.UpdatedAt),
					})
				}
			}
			wg.Done()
		}(note)
	}
	wg.Wait()
	return c.RenderJSON(results)
}

var (
	re = regexp.MustCompile(`[^a-zA-Z0-9]+`)
)

func (c Notebooks) CreateAPI() revel.Result {
	token, err := gc.Auth(c.Request)
	if err != nil {
		c.Log.Errorf("Failed to get notebooks: %v", err)
		c.Response.SetStatus(http.StatusUnauthorized)
		return c.RenderError(errors.New("認証エラー"))
	}
	email := gc.VerifiedEmail(token)
	name := fmt.Sprintf("note-%s-%s",
		strings.ToLower(re.ReplaceAllString(strings.Split(email, "@")[0], "")),
		time.Now().Format("20060102-150405"),
	)
	ctx := c.Request.Context()

	params := struct {
		Menu string `json:"menu"`
	}{}
	c.Params.BindJSON(&params)

	// Create a managed notebook
	if err = gc.CreateManagedNotebook(ctx, name, email, params.Menu); err != nil {
		c.Log.Errorf("Failed to create a notebook: %v", err)
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderError(errors.New("内部エラー"))
	}
	if err = gc.BindRole(c.Request, "user:"+email, "roles/notebooks.runner"); err != nil {
		c.Log.Errorf("Failed to bind a role to the user: %v", err)
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderError(errors.New("内部エラー"))
	}
	// Create a storage bucket
	name = fmt.Sprintf("bucket-users-%s",
		strings.ToLower(re.ReplaceAllString(strings.Split(email, "@")[0], "")),
	)
	if err = gc.MakeBucket(ctx, name); err != nil {
		c.Log.Errorf("Failed to create a bucket: %v", err)
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderError(errors.New("内部エラー"))
	}
	if err = gc.AddRoleToBucket(ctx, name, "user:"+email, "storage.objectAdmin"); err != nil {
		c.Log.Errorf("Failed to bind a role to the user: %v", err)
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderError(errors.New("内部エラー"))
	}
	// Save a record to Firestore
	if err = gc.SaveNotebook(ctx, name, email, params.Menu); err != nil {
		c.Log.Errorf("Failed to save a notebook: %v", err)
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderError(errors.New("内部エラー"))
	}
	return c.RenderJSON(nil)
}

func (c Notebooks) UpdateAPI() revel.Result {
	token, err := gc.Auth(c.Request)
	if err != nil {
		c.Log.Errorf("Failed to update a notebook: %v", err)
		c.Response.SetStatus(http.StatusUnauthorized)
		return c.RenderError(errors.New("認証エラー"))
	}
	email := gc.VerifiedEmail(token)
	ctx := c.Request.Context()

	params := struct {
		Action string `json:"action"`
		ID     string `json:"id"`
	}{}
	c.Params.BindJSON(&params)

	// Check its owner
	if _, err = gc.GetNotebook(ctx, email, params.ID); err != nil {
		c.Log.Errorf("Failed to find the notebook: %v (name: %s)", err, params.ID)
		c.Response.SetStatus(http.StatusNotFound)
		return c.RenderError(errors.New("扱えるノートがありません"))
	}
	// Stop and start the notebook
	switch params.Action {
	case "start":
		if err = gc.StartManagedNotebook(ctx, params.ID); err != nil {
			c.Log.Errorf("Failed to start the notebook instance: %v (name: %s)", err, params.ID)
			c.Response.SetStatus(http.StatusInternalServerError)
			return c.RenderError(errors.New("内部エラー"))
		}
	case "stop":
		if err = gc.StopManagedNotebook(ctx, params.ID); err != nil {
			c.Log.Errorf("Failed to stop the notebook instance: %v (name: %s)", err, params.ID)
			c.Response.SetStatus(http.StatusInternalServerError)
			return c.RenderError(errors.New("内部エラー"))
		}
	}
	return c.RenderJSON(nil)
}

func (c Notebooks) DeleteAPI() revel.Result {
	token, err := gc.Auth(c.Request)
	if err != nil {
		c.Log.Errorf("Failed to delete a notebook: %v", err)
		c.Response.SetStatus(http.StatusUnauthorized)
		return c.RenderError(errors.New("認証エラー"))
	}
	email := gc.VerifiedEmail(token)
	ctx := c.Request.Context()

	params := struct {
		ID string `json:"id"`
	}{}
	c.Params.BindJSON(&params)

	// Check its owner
	note, err := gc.GetNotebook(ctx, email, params.ID)
	if err != nil {
		c.Log.Errorf("Failed to find the notebook: %v (name: %s)", err, params.ID)
		c.Response.SetStatus(http.StatusNotFound)
		return c.RenderError(errors.New("扱えるノートがありません"))
	}
	// Delete the notebook
	if err = gc.DeleteManagedNotebook(ctx, params.ID); err != nil {
		c.Log.Errorf("Failed to delete the notebook instance: %v (name: %s)", err, params.ID)
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderError(errors.New("内部エラー"))
	}
	// Update the record
	deactivate := map[string]interface{}{
		"active": false,
	}
	if err := gc.UpdateNotebook(ctx, note.FirestoreID, deactivate); err != nil {
		c.Log.Errorf("Failed to update the notebook record: %v (name: %s)", err, params.ID)
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderError(errors.New("内部エラー"))
	}
	return c.RenderJSON(nil)
}
