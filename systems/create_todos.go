package systems

import (
	"github.com/zodimo/go-gio-ark-todo/components"

	"github.com/google/uuid"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark/ecs"
)

var _ app.System = (*CreateTodos)(nil)

type CreateTodos struct {
	initialTodos []string
}

func NewCreateTodos(initialTodos []string) *CreateTodos {
	return &CreateTodos{
		initialTodos: initialTodos,
	}
}

func (s *CreateTodos) Initialize(w *ecs.World) {
	builder := ecs.NewMap[components.Todo](w)
	for _, todoText := range s.initialTodos {
		uuid := uuid.New()
		builder.NewEntity(&components.Todo{ID: uuid.String(), Text: todoText})
	}
}
func (s *CreateTodos) Update(w *ecs.World) {
	// no runtime logic needed
}
func (s *CreateTodos) Finalize(w *ecs.World) {
	// no finalization needed
}
