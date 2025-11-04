package systems

import (
	"github.com/zodimo/go-gio-ark-todo/components"

	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark/ecs"
)

var _ app.System = (*RemoveTodo)(nil)

type RemoveTodo struct {
	filter *ecs.Filter1[components.Todo]
	uiRes  ecs.Resource[components.UI]
}

func NewRemoveTodo() *RemoveTodo {
	return &RemoveTodo{}
}

func (s *RemoveTodo) Initialize(w *ecs.World) {
	s.filter = ecs.NewFilter1[components.Todo](w)
	s.uiRes = ecs.NewResource[components.UI](w)
}

func (s *RemoveTodo) Update(w *ecs.World) {

	ui := s.uiRes.Get()
	uiState := ui.UIState
	if uiState.PendingRemoveTodo == nil {
		return
	}

	query := s.filter.Query()
	entity := ecs.Entity{}
	for query.Next() {
		todo := query.Get()
		if todo.ID == uiState.PendingRemoveTodo.TodoID {
			entity = query.Entity()
			break
		}
	}
	query.Close()

	if entity.IsZero() {
		return
	}

	w.RemoveEntity(entity)
	ui.DeleteClickableForTodo(uiState.PendingRemoveTodo.TodoID)
	uiState.PendingRemoveTodo = nil
	uiState.IsDirty = true

}

func (s *RemoveTodo) Finalize(w *ecs.World) {
	// no finalization needed
}
