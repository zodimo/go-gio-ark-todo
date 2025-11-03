package components

import (
	"gioui.org/app"
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
