package components

import (
	"fmt"

	"gioui.org/app"
	"gioui.org/widget"
)

type UI struct {
	Window    *app.Window
	UIWidgets *UIWidgets
	UIState   *UIState
}

func NewUI(window *app.Window, uiWidgets *UIWidgets, uiState *UIState) *UI {
	return &UI{
		Window:    window,
		UIWidgets: uiWidgets,
		UIState:   uiState,
	}
}

func (d *UI) GetToggleClickableForTodo(todoId string) *widget.Clickable {
	if d.UIWidgets.TodoToggleButtons == nil {
		fmt.Println("Initialize toggle buttons")
		d.UIWidgets.TodoToggleButtons = make(map[string]*widget.Clickable)
	}
	if btn, exists := d.UIWidgets.TodoToggleButtons[todoId]; exists {
		return btn
	}
	fmt.Printf("creating toggle button for todo %s\n", todoId)
	newBtn := &widget.Clickable{}
	d.UIWidgets.TodoToggleButtons[todoId] = newBtn
	return newBtn
}

func (d *UI) GetRemoveClickableForTodo(todoId string) *widget.Clickable {
	if d.UIWidgets.TodoRemoveButtons == nil {
		fmt.Println("Initialize remove buttons")
		d.UIWidgets.TodoRemoveButtons = make(map[string]*widget.Clickable)
	}
	if btn, exists := d.UIWidgets.TodoRemoveButtons[todoId]; exists {
		return btn
	}
	fmt.Printf("creating remove button for todo %s\n", todoId)
	newBtn := &widget.Clickable{}
	d.UIWidgets.TodoRemoveButtons[todoId] = newBtn
	return newBtn
}

func (d *UI) DeleteClickableForTodo(todoId string) {
	delete(d.UIWidgets.TodoRemoveButtons, todoId)
	delete(d.UIWidgets.TodoToggleButtons, todoId)
}
