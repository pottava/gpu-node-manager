package controllers

import (
	"fmt"
	"os"

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
	return c.RenderText(fmt.Sprintf("Revision: %s", revision))
}

func (c App) Login() revel.Result {
	return c.Render()
}

func (c App) Index() revel.Result {
	return c.Render()
}
