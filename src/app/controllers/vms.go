package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	gc "github.com/pottava/gpu-node-manager/src/app/googlecloud"
	"github.com/pottava/gpu-node-manager/src/app/util"
	"github.com/revel/revel"
)

type VMs struct {
	*revel.Controller
}

func (c VMs) Index() revel.Result {
	return c.Render()
}

func (c VMs) Create() revel.Result {
	isNotProduction := (util.AppStage() != "prod")
	return c.Render(isNotProduction)
}

func (c VMs) ListAPI() revel.Result {
	token, err := gc.Auth(c.Request)
	if err != nil {
		c.Log.Errorf("Failed to get vms: %v", err)
		c.Response.SetStatus(http.StatusUnauthorized)
		return c.RenderError(errors.New("認証エラー"))
	}
	email := gc.VerifiedEmail(token)

	vms, err := gc.GetVMs(c.Request.Context(), email)
	if err != nil {
		c.Log.Errorf("Failed to get vms: %v", err)
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderError(errors.New("内部エラー"))
	}

	type Result struct {
		Menu      string `json:"menu"`
		Name      string `json:"name"`
		State     string `json:"state"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
	results := []*Result{}
	for _, vm := range vms {
		results = append(results, &Result{
			Menu:      vm.Menu,
			Name:      vm.Name,
			State:     strconv.FormatBool(vm.Active),
			CreatedAt: util.DateToStr(vm.CreatedAt),
			UpdatedAt: util.DateToStr(vm.UpdatedAt),
		})
	}
	return c.RenderJSON(results)
}

func (c VMs) CreateAPI() revel.Result {
	token, err := gc.Auth(c.Request)
	if err != nil {
		c.Log.Errorf("Failed to get vms: %v", err)
		c.Response.SetStatus(http.StatusUnauthorized)
		return c.RenderError(errors.New("認証エラー"))
	}
	email := gc.VerifiedEmail(token)
	name := fmt.Sprintf("vm-%s-%s",
		strings.ToLower(re.ReplaceAllString(strings.Split(email, "@")[0], "")),
		time.Now().Format("20060102-150405"),
	)
	ctx := c.Request.Context()

	params := struct {
		Menu string `json:"menu"`
	}{}
	c.Params.BindJSON(&params)

	// Create a managed notebook
	if err = gc.CreateVM(ctx, name, email, params.Menu); err != nil {
		c.Log.Errorf("Failed to create a notebook: %v", err)
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderError(errors.New("内部エラー"))
	}
	// Save a record to Firestore
	if err = gc.SaveVM(ctx, name, email, params.Menu); err != nil {
		c.Log.Errorf("Failed to save a notebook: %v", err)
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderError(errors.New("内部エラー"))
	}
	return c.RenderJSON(nil)
}

func (c VMs) UpdateAPI() revel.Result {
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
	_, err = gc.GetVM(ctx, email, params.ID)
	if err != nil {
		c.Log.Errorf("Failed to find the VM: %v (name: %s)", err, params.ID)
		c.Response.SetStatus(http.StatusNotFound)
		return c.RenderError(errors.New("扱える仮想マシンがありません"))
	}
	// Stop and start the VM
	switch params.Action {
	case "start":
		if err = gc.StartVM(ctx, params.ID); err != nil {
			c.Log.Errorf("Failed to start the VM instance: %v (name: %s)", err, params.ID)
			c.Response.SetStatus(http.StatusInternalServerError)
			return c.RenderError(errors.New("内部エラー"))
		}
	case "stop":
		if err = gc.StopVM(ctx, params.ID); err != nil {
			c.Log.Errorf("Failed to stop the VM instance: %v (name: %s)", err, params.ID)
			c.Response.SetStatus(http.StatusInternalServerError)
			return c.RenderError(errors.New("内部エラー"))
		}
	}
	return c.RenderJSON(nil)
}

func (c VMs) DeleteAPI() revel.Result {
	token, err := gc.Auth(c.Request)
	if err != nil {
		c.Log.Errorf("Failed to delete a VM: %v", err)
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
	vm, err := gc.GetVM(ctx, email, params.ID)
	if err != nil {
		c.Log.Errorf("Failed to find the VM: %v (name: %s)", err, params.ID)
		c.Response.SetStatus(http.StatusNotFound)
		return c.RenderError(errors.New("扱える仮想マシンがありません"))
	}
	// Delete the VM
	if err := gc.DeleteVM(ctx, params.ID); err != nil {
		c.Log.Errorf("Failed to delete the VM instance: %v (name: %s)", err, params.ID)
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderError(errors.New("内部エラー"))
	}
	// Update the record
	deactivate := map[string]interface{}{
		"active": false,
	}
	if err := gc.UpdateVM(ctx, vm.FirestoreID, deactivate); err != nil {
		c.Log.Errorf("Failed to update the VM record: %v (name: %s)", err, params.ID)
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderError(errors.New("内部エラー"))
	}
	return c.RenderJSON(nil)
}
