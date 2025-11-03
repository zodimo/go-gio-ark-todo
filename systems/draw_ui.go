package systems

import (
	"fmt"
	"image/color"

	"github.com/zodimo/go-gio-ark-todo/components"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark/ecs"
)

var _ app.UISystem = (*DrawUI)(nil)

type DrawUI struct {
	uiRes  ecs.Resource[components.UI]
	gtxRes ecs.Resource[layout.Context]

	allTodosfilter      *ecs.Filter1[components.Todo]
	completedTdosfilter *ecs.Filter2[components.Todo, components.TodoCompleted]
}

func NewDrawUI() *DrawUI {
	return &DrawUI{}
}

func (d *DrawUI) InitializeUI(w *ecs.World) {
	d.uiRes = ecs.NewResource[components.UI](w)
	d.gtxRes = ecs.NewResource[layout.Context](w)

	d.allTodosfilter = ecs.NewFilter1[components.Todo](w)
	d.completedTdosfilter = ecs.NewFilter2[components.Todo, components.TodoCompleted](w)
}

func (d *DrawUI) UpdateUI(w *ecs.World) {
	gtx := *(d.gtxRes.Get())
	d.Layout(w, gtx)
}
func (d *DrawUI) PostUpdateUI(w *ecs.World) {
	// no post update needed
}
func (d *DrawUI) FinalizeUI(w *ecs.World) {
	// no finalization needed
}

func (d *DrawUI) Layout(w *ecs.World, gtx C) D {
	todos := []components.TodoItem{}

	ui := d.uiRes.Get()

	uiWidgets := ui.UIWidgets
	uiState := ui.UIState

	if uiState.IsDirty {
		inv := op.InvalidateCmd{At: gtx.Now}
		gtx.Execute(inv)
		uiState.IsDirty = false

	}

	switch uiState.CurrentView {
	case components.ViewAll:
		// Render all todos
		query := d.allTodosfilter.Query()
		for query.Next() {
			todo := query.Get()
			if d.isTodoCompleted(todo.ID) {
				todos = append(todos, components.TodoItem{Todo: *todo, IsCompleted: true})
			} else {
				todos = append(todos, components.TodoItem{Todo: *todo})
			}
		}

	case components.ViewActive:
		// Render active todos
		query := d.allTodosfilter.Query()
		for query.Next() {
			todo := query.Get()
			if !d.isTodoCompleted(todo.ID) {
				todos = append(todos, components.TodoItem{Todo: *todo})
			}

		}
	case components.ViewCompleted:
		// Render completed todos
		query := d.completedTdosfilter.Query()
		for query.Next() {
			todo, _ := query.Get()
			todos = append(todos, components.TodoItem{Todo: *todo, IsCompleted: true})
		}
	default:
		// Default case
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// Header
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top: unit.Dp(20), Bottom: unit.Dp(20),
				Left: unit.Dp(20), Right: unit.Dp(20),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				title := material.H3(uiWidgets.Theme, "ARK + Gioui Todo App")
				title.Alignment = text.Middle
				return title.Layout(gtx)
			})
		}),
		// Input section
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(10),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						editor := material.Editor(uiWidgets.Theme, uiWidgets.Editor, "Add a new todo...")
						return editor.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {

						btn := material.Button(uiWidgets.Theme, uiWidgets.AddButton, "Add")
						return layout.Inset{Left: unit.Dp(10)}.Layout(gtx, btn.Layout)
					}),
				)
			})
		}),

		// Filter buttons
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(10),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {

						btn := material.Button(uiWidgets.Theme, uiWidgets.AllFilter, "All")
						if uiState.CurrentView == components.ViewAll {
							btn.Background = color.NRGBA{R: 100, G: 150, B: 200, A: 255}
						}
						return btn.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {

						btn := material.Button(uiWidgets.Theme, uiWidgets.ActiveFilter, "Active")
						if uiState.CurrentView == components.ViewActive {
							btn.Background = color.NRGBA{R: 100, G: 150, B: 200, A: 255}
						}
						return layout.Inset{Left: unit.Dp(10)}.Layout(gtx, btn.Layout)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {

						btn := material.Button(uiWidgets.Theme, uiWidgets.CompletedFilter, "Completed")
						if uiState.CurrentView == components.ViewCompleted {
							btn.Background = color.NRGBA{R: 100, G: 150, B: 200, A: 255}
						}
						return layout.Inset{Left: unit.Dp(10)}.Layout(gtx, btn.Layout)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{}
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(uiWidgets.Theme, uiWidgets.ClearCompletedButton, "Clear Completed")
						if uiState.ClearCompletedClicked {
							btn.Background = color.NRGBA{R: 100, G: 150, B: 200, A: 255}
						}
						return layout.Inset{Left: unit.Dp(10)}.Layout(gtx, btn.Layout)
					}),
				)
			})
		}),

		// Todo list
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {

			return layout.Inset{
				Left: unit.Dp(20), Right: unit.Dp(20),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				list := layout.List{
					Axis: layout.Vertical,
				}
				return list.Layout(gtx, len(todos), func(gtx layout.Context, i int) layout.Dimensions {
					todoItem := todos[i]

					return layout.Inset{Bottom: unit.Dp(5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

						return layout.Flex{
							Alignment: layout.Middle,
						}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {

								checkbox := "☐"
								if todoItem.IsCompleted {
									checkbox = "☑"
								}

								label := material.Body1(uiWidgets.Theme, checkbox)
								label.Color = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
								return layout.Inset{Right: unit.Dp(10)}.Layout(gtx, label.Layout)
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								label := material.Body1(uiWidgets.Theme, todoItem.Todo.Text)
								if todoItem.IsCompleted {
									label.Color = color.NRGBA{R: 150, G: 150, B: 150, A: 255}
								}

								// Make the entire row clickable
								// return ts.todoButtons[i].Layout(gtx, label.Layout)
								toggleBtn := d.getToggleClickableForTodo(todoItem.Todo.ID)
								return toggleBtn.Layout(gtx, label.Layout)
								// return label.Layout(gtx)
							}),
						)
					})
				})
			})

		}),
	)

}

func (d *DrawUI) getToggleClickableForTodo(todoId string) *widget.Clickable {
	uiWidgets := d.uiRes.Get().UIWidgets
	if uiWidgets.TodoToggleButtons == nil {
		fmt.Println("Initialize toggle buttons")
		uiWidgets.TodoToggleButtons = make(map[string]*widget.Clickable)
	}
	if btn, exists := uiWidgets.TodoToggleButtons[todoId]; exists {
		return btn
	}
	fmt.Printf("createing toggle button for todo %s\n", todoId)
	newBtn := &widget.Clickable{}
	uiWidgets.TodoToggleButtons[todoId] = newBtn
	return newBtn
}

func (d *DrawUI) isTodoCompleted(TodoId string) bool {
	query := d.completedTdosfilter.Query()
	var result bool
	for query.Next() {
		todo, _ := query.Get()
		if todo.ID == TodoId {
			result = true
			break
		}

	}
	if result {
		// close the world lock because we break early
		query.Close()
	}
	return result
}
