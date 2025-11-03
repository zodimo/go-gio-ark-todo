package components

import (
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type todoITemState int

// done , not done, deleted, not deleted

const (
	TodoItemStateActive todoITemState = iota
	TodoItemStateCompleted
	TodoItemStateCompletedDeleted
)

type ViewState int

const (
	ViewAll ViewState = iota
	ViewActive
	ViewCompleted
)

// UIState component - holds UI-related state
type UIState struct {
	InputText string

	CurrentView                ViewState
	AddTodoSubmitButtonClicked bool
	NewTodoSubmitted           bool

	ViewAllTodoClicked       bool
	ViewActiveTodoClicked    bool
	ViewCompletedTodoClicked bool

	DeleteTodoClicked bool

	ClearCompletedClicked bool

	PendingTodos []string

	PendingToggleTodo *PendingToggleTodo

	IsDirty bool
}

func NewUIState() *UIState {
	return &UIState{
		InputText: "",

		CurrentView:                ViewAll,
		AddTodoSubmitButtonClicked: false,
		NewTodoSubmitted:           false,

		ViewAllTodoClicked:       false,
		ViewActiveTodoClicked:    false,
		ViewCompletedTodoClicked: false,

		DeleteTodoClicked: false,

		ClearCompletedClicked: false,

		PendingTodos: []string{},

		PendingToggleTodo: nil,
		IsDirty:           true,
	}
}

type TodoFilter int

const (
	FilterAll TodoFilter = iota
	FilterActive
	FilterCompleted
)

type TodoCompleted struct{}

type TodoDeleted struct{}

type Todo struct {
	ID   string
	Text string
}

type TodoState struct {
	State todoITemState
}

type PendingToggleTodo struct {
	TodoID      string
	IsCompleted bool
}

type UIWidgets struct {
	Theme *material.Theme

	TodoToggleButtons map[string]*widget.Clickable

	Editor               *widget.Editor
	AddButton            *widget.Clickable
	ClearCompletedButton *widget.Clickable

	AllFilter       *widget.Clickable
	ActiveFilter    *widget.Clickable
	CompletedFilter *widget.Clickable
}

func NewUIWidgets(theme *material.Theme) *UIWidgets {
	editor := &widget.Editor{}
	editor.SingleLine = true
	editor.Submit = true

	return &UIWidgets{
		Theme:                theme,
		TodoToggleButtons:    make(map[string]*widget.Clickable),
		Editor:               editor,
		AddButton:            &widget.Clickable{},
		ClearCompletedButton: &widget.Clickable{},
		AllFilter:            &widget.Clickable{},
		ActiveFilter:         &widget.Clickable{},
		CompletedFilter:      &widget.Clickable{},
	}
}

type TodoItem struct {
	Todo        Todo
	IsCompleted bool
}
