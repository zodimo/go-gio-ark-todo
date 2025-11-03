package systems

import (
	"github.com/zodimo/go-gio-ark-todo/components"

	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark/ecs"
)

var _ app.System = (*TogglePendingTodo)(nil)

type TogglePendingTodo struct {
	allTodosfilter *ecs.Filter1[components.Todo]
	uiRes          ecs.Resource[components.UI]
}

func NewTogglePendingTodo() *TogglePendingTodo {
	return &TogglePendingTodo{}
}

func (s *TogglePendingTodo) Initialize(w *ecs.World) {
	s.allTodosfilter = ecs.NewFilter1[components.Todo](w)
	s.uiRes = ecs.NewResource[components.UI](w)
}

func (s *TogglePendingTodo) Update(w *ecs.World) {

	var targetEntity ecs.Entity
	var shouldComplete *bool

	uiState := s.uiRes.Get().UIState

	if uiState.PendingToggleTodo == nil {
		return
	}

	query := s.allTodosfilter.Query()
	for query.Next() {
		todo := query.Get()
		if todo.ID == uiState.PendingToggleTodo.TodoID {
			targetEntity = query.Entity()

			// Determine the action to take after the query is done
			isComp := uiState.PendingToggleTodo.IsCompleted
			shouldComplete = &isComp
			break
		}
	}
	// close the world lock because we break early
	query.Close()

	// 2. Perform the Structural Change (The Mutation Phase)
	if !targetEntity.IsZero() && shouldComplete != nil {

		if *shouldComplete {
			completedMapper := ecs.NewMap1[components.TodoCompleted](w)
			completedMapper.Remove(targetEntity)

		} else {
			completedMapper := ecs.NewMap1[components.TodoCompleted](w)
			completedMapper.Add(targetEntity, &components.TodoCompleted{})
		}

		uiState.PendingToggleTodo = nil
	}
}

func (s *TogglePendingTodo) Finalize(w *ecs.World) {
	// no finalization needed
}
