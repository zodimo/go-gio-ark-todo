package systems

import (
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark/ecs"
	"github.com/zodimo/go-gio-ark-todo/components"
)

var _ app.System = (*TodoStats)(nil)

type TodoStats struct {
	allTodosfilter      *ecs.Filter1[components.Todo]
	completedTdosfilter *ecs.Filter2[components.Todo, components.TodoCompleted]
	uiRes               ecs.Resource[components.UI]
}

func NewTodoStats() *TodoStats {
	return &TodoStats{}
}

func (s *TodoStats) Initialize(w *ecs.World) {
	s.allTodosfilter = ecs.NewFilter1[components.Todo](w)
	s.completedTdosfilter = ecs.NewFilter2[components.Todo, components.TodoCompleted](w)
	s.uiRes = ecs.NewResource[components.UI](w)

}
func (s *TodoStats) Update(w *ecs.World) {
	ui := s.uiRes.Get()
	allTodosQuery := s.allTodosfilter.Query()
	ui.UIState.TotalTodos = allTodosQuery.Count()
	allTodosQuery.Close()

	completedTodosQuery := s.completedTdosfilter.Query()
	ui.UIState.TotalCompletedTodos = completedTodosQuery.Count()
	completedTodosQuery.Close()
}
func (s *TodoStats) Finalize(w *ecs.World) {
	// no finalization needed
}
