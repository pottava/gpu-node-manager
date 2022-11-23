package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	gc "github.com/pottava/gpu-node-manager/src/app/googlecloud"
	"github.com/revel/revel"
)

type Notebooks struct {
	*revel.Controller
}

func (c Notebooks) Index() revel.Result {
	return c.Render()
}

func (c Notebooks) Create() revel.Result {
	return c.Render()
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
		Menu     string `json:"menu"`
		Runtime  string `json:"runtime"`
		ProxyUri string `json:"proxyUri"`
		State    string `json:"state"`
	}
	results := []*Result{}
	for _, note := range notebooks {
		runtime, err := gc.DescribeManagedNotebook(c.Request.Context(), note.Runtime)
		if err == nil {
			results = append(results, &Result{
				Menu:     note.Menu,
				Runtime:  note.Runtime,
				ProxyUri: runtime.AccessConfig.ProxyUri,
				State:    runtime.State.String(),
			})
		}
	}
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
		time.Now().Format("20060102150405"),
	)
	ctx := c.Request.Context()

	params := struct {
		Menu string `json:"menu"`
	}{}
	c.Params.BindJSON(&params)

	if err = gc.CreateManagedNotebook(ctx, name, email, params.Menu); err != nil {
		c.Log.Errorf("Failed to create a notebook: %v", err)
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderError(errors.New("内部エラー"))
	}
	if err = gc.SaveNotebook(ctx, name, email, params.Menu); err != nil {
		c.Log.Errorf("Failed to save a notebook: %v", err)
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderError(errors.New("内部エラー"))
	}
	return c.RenderJSON(nil)
}
