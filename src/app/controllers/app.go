package controllers

import (
	"fmt"
	"os"

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
	if len(util.Date) == 0 {
		return c.RenderText(fmt.Sprintf("Revision: %s\n", revision))
	}
	return c.RenderText(fmt.Sprintf(
		"Revision: %s (built at %s)\n", revision, util.Date))
}

func (c App) Login() revel.Result {
	return c.Render()
}

func (c App) Index() revel.Result {
	return c.Render()
}
