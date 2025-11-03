package systems

import (
	"github.com/zodimo/go-gio-ark-todo/components"

	"github.com/google/uuid"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark/ecs"
)

var _ app.System = (*AddPendingTodos)(nil)

type AddPendingTodos struct {
	uiRes ecs.Resource[components.UI]
}

func NewAddPendingTodos() *AddPendingTodos {
	return &AddPendingTodos{}
}

func (s *AddPendingTodos) Initialize(w *ecs.World) {
	s.uiRes = ecs.NewResource[components.UI](w)
}

func (s *AddPendingTodos) Update(w *ecs.World) {
	uiState := s.uiRes.Get().UIState
	if len(uiState.PendingTodos) > 0 {
		for _, todoText := range uiState.PendingTodos {
			// Create new Todo entity
			builder := ecs.NewMap[components.Todo](w)
			uuid := uuid.New()
			builder.NewEntity(&components.Todo{ID: uuid.String(), Text: todoText})
		}

		uiState.PendingTodos = []string{}
	}
}
func (s *AddPendingTodos) Finalize(w *ecs.World) {
	// no finalization needed
}
