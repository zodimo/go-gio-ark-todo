package systems

import (
	"github.com/zodimo/go-gio-ark-todo/components"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark/ecs"
)

var _ app.System = (*UpdateUI)(nil)

type UpdateUI struct {
	uiRes  ecs.Resource[components.UI]
	gtxRes ecs.Resource[layout.Context]

	allTodosfilter      *ecs.Filter1[components.Todo]
	completedTdosfilter *ecs.Filter2[components.Todo, components.TodoCompleted]
}

func NewUpdateUI() *UpdateUI {
	return &UpdateUI{}
}

func (s *UpdateUI) Initialize(w *ecs.World) {
	s.uiRes = ecs.NewResource[components.UI](w)
	s.gtxRes = ecs.NewResource[layout.Context](w)
	s.allTodosfilter = ecs.NewFilter1[components.Todo](w)
	s.completedTdosfilter = ecs.NewFilter2[components.Todo, components.TodoCompleted](w)
}

func (s *UpdateUI) Update(w *ecs.World) {
	ui := s.uiRes.Get()
	gtx := *(s.gtxRes.Get())
	for {
		event, ok := ui.UIWidgets.Editor.Update(gtx)
		if !ok {
			break
		}
		if _, ok := event.(widget.SubmitEvent); ok {
			s.appendPendingTodo(w, ui.UIWidgets.Editor.Text())
			s.clearEditor()
		}
	}

	if ui.UIWidgets.AddButton.Clicked(gtx) {
		s.appendPendingTodo(w, ui.UIWidgets.Editor.Text())
		s.clearEditor()

	}

	if ui.UIWidgets.ActiveFilter.Clicked(gtx) {
		s.setCurrentView(w, components.ViewActive)
	}
	if ui.UIWidgets.CompletedFilter.Clicked(gtx) {
		s.setCurrentView(w, components.ViewCompleted)
	}

	if ui.UIWidgets.AllFilter.Clicked(gtx) {
		s.setCurrentView(w, components.ViewAll)
	}

	if ui.UIWidgets.ClearCompletedButton.Clicked(gtx) {
		s.clearCompleted(w)
	}

	// Handle todo toggle clicks
	query := s.allTodosfilter.Query()

	todoIdButtonClicked := ""
	for query.Next() {
		todo := query.Get()
		toggleBtn := s.getToggleClickableForTodo(todo.ID)
		if toggleBtn.Clicked(gtx) {
			todoIdButtonClicked = todo.ID

		}
	}
	if todoIdButtonClicked != "" {
		// r.toggleTodo(w, todoIdButtonClicked)
		s.toggleTodo(w, todoIdButtonClicked)
	}
}

func (s *UpdateUI) Finalize(w *ecs.World) {
	// no finalization needed
}

func (s *UpdateUI) clearCompleted(w *ecs.World) {
	uiState := s.uiRes.Get().UIState
	uiState.ClearCompletedClicked = true
	uiState.IsDirty = true
}

func (s *UpdateUI) clearEditor() {
	ui := s.uiRes.Get()
	ui.UIWidgets.Editor.SetText("")
	ui.UIState.IsDirty = true
}

func (s *UpdateUI) appendPendingTodo(w *ecs.World, text string) {
	if text == "" {
		return
	}
	uiState := s.uiRes.Get().UIState
	uiState.PendingTodos = append(uiState.PendingTodos, text)
	uiState.IsDirty = true
}

func (s *UpdateUI) setCurrentView(w *ecs.World, view components.ViewState) {
	uiState := s.uiRes.Get().UIState
	uiState.CurrentView = view
	uiState.IsDirty = true
}

func (r *UpdateUI) getToggleClickableForTodo(todoId string) *widget.Clickable {
	ui := r.uiRes.Get().UIWidgets
	if ui.TodoToggleButtons == nil {
		ui.TodoToggleButtons = make(map[string]*widget.Clickable)
	}
	if btn, exists := ui.TodoToggleButtons[todoId]; exists {
		return btn
	}
	newBtn := &widget.Clickable{}
	ui.TodoToggleButtons[todoId] = newBtn
	return newBtn
}

func (r *UpdateUI) toggleTodo(w *ecs.World, todoId string) {

	// Remove completed component
	query := r.allTodosfilter.Query()

	var pendingToggle *components.PendingToggleTodo

	for query.Next() {
		todo := query.Get()
		if todo.ID == todoId {
			isCompleted := r.isTodoCompleted(todo.ID)
			pendingToggle = &components.PendingToggleTodo{
				TodoID:      todoId,
				IsCompleted: isCompleted,
			}
		}
	}

	if pendingToggle == nil {
		return
	}

	uiState := r.uiRes.Get().UIState
	localCopy := *pendingToggle
	uiState.PendingToggleTodo = &localCopy
	uiState.IsDirty = true

}

func (d *UpdateUI) isTodoCompleted(TodoId string) bool {
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
