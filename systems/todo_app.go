package systems

import (
	"context"

	"github.com/zodimo/go-gio-ark-todo/components"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	arkApp "github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark/ecs"
)

type C = layout.Context
type D = layout.Dimensions

type TodoApp struct {
	App    *arkApp.App
	window *app.Window
	ctx    context.Context

	settings *components.Settings
}

func NewTodoApp(window *app.Window) *TodoApp {

	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	a := &TodoApp{
		window: window,
		ctx:    context.Background(),
	}

	a.App = arkApp.New()

	a.settings = &components.Settings{
		ScreenWidth:  640,
		ScreenHeight: 480,
		Title:        "ARK + Gioui Todo App",
	}

	ecs.AddResource(&a.App.World, a.settings)
	ecs.AddResource(&a.App.World, components.NewUI(a.window, components.NewUIWidgets(th), components.NewUIState()))

	a.App.AddSystem(NewAddPendingTodos())
	a.App.AddSystem(NewCreateTodos([]string{
		"Learn ARK ECS",
		"Build gioui app",
		"Integrate ARK with gioui",
	}))
	a.App.AddSystem(NewTogglePendingTodo())
	a.App.AddSystem(NewRemoveCompleted())
	a.App.AddSystem(NewTodoStats())
	a.App.AddSystem(NewUpdateUI())
	a.App.AddUISystem(NewDrawUI())

	a.App.Initialize()

	return a
}

func (ta *TodoApp) Run() error {
	gtxRes := ecs.NewResource[layout.Context](&ta.App.World)

	var ops op.Ops
	for {
		switch e := ta.window.Event().(type) {
		case app.DestroyEvent:
			ta.App.Finalize()
			return nil
		case app.FrameEvent:
			if gtxRes.Has() {
				gtxRes.Remove()
			}
			gtx := app.NewContext(&ops, e)
			gtxRes.Add(&gtx)

			// inv := op.InvalidateCmd{At: gtx.Now}
			// gtx.Execute(inv)

			// ta.window.Invalidate()

			ta.App.Update()
			ta.App.UpdateUI()

			e.Frame(gtx.Ops)
		}
	}
}

func (ta *TodoApp) GetSettings() *components.Settings {
	return ta.settings
}
