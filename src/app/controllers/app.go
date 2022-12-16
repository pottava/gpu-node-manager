package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	gc "github.com/pottava/gpu-node-manager/src/app/googlecloud"
	"github.com/pottava/gpu-node-manager/src/app/util"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Health() revel.Result {
	revision := "local"
	if candidate, found := os.LookupEnv("K_REVISION"); found {
		revision = candidate
	}
	if len(util.BuildDate) == 0 {
		return c.RenderText(fmt.Sprintf("Revision: %s\n", revision))
	}
	return c.RenderText(fmt.Sprintf(
		"Revision: %s (built at %s)\n", revision, util.BuildDate))
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Login() revel.Result {
	return c.Render()
}

func (c App) ContextAPI() revel.Result {
	token, err := gc.Auth(c.Request)
	if err != nil {
		c.Log.Errorf("Failed to read its context: %v", err)
		c.Response.SetStatus(http.StatusUnauthorized)
		return c.RenderError(errors.New("認証エラー"))
	}
	email := gc.VerifiedEmail(token)

	type Result struct {
		Project  string `json:"prj"`
		Stage    string `json:"stg"`
		Revision string `json:"rev"`
		User     string `json:"usr"`
	}
	results := &Result{
		Project:  util.ProjectID(),
		Stage:    util.AppStage(),
		Revision: util.RunRevision(),
		User:     email,
	}
	return c.RenderJSON(results)
}

func (c App) PasswordResetAPI() revel.Result {
	token, err := gc.Auth(c.Request)
	if err != nil {
		c.Log.Errorf("Failed to read its context: %v", err)
		c.Response.SetStatus(http.StatusUnauthorized)
		return c.RenderError(errors.New("認証エラー"))
	}
	link, err := gc.PasswordResetLink(c.Request, gc.VerifiedEmail(token))
	if err != nil {
		c.Log.Errorf("Failed to generate a reset link: %v", err)
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderError(errors.New("内部エラー"))
	}
	return c.RenderJSON(map[string]string{"link": link})
}
