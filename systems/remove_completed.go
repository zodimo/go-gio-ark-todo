package systems

import (
	"github.com/zodimo/go-gio-ark-todo/components"

	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark/ecs"
)

var _ app.System = (*RemoveCompleted)(nil)

type RemoveCompleted struct {
	filter *ecs.Filter2[components.Todo, components.TodoCompleted]
	UIRes  ecs.Resource[components.UI]
}

func NewRemoveCompleted() *RemoveCompleted {
	return &RemoveCompleted{}
}
func (s *RemoveCompleted) Initialize(w *ecs.World) {
	s.filter = s.filter.New(w)
	s.UIRes = ecs.NewResource[components.UI](w)
}
func (s *RemoveCompleted) Update(w *ecs.World) {
	ui := s.UIRes.Get()

	if ui.UIState.ClearCompletedClicked {
		query := s.filter.Query()
		entities := []ecs.Entity{}
		todoIds := []string{}
		for query.Next() {
			entities = append(entities, query.Entity())

			todo, _ := query.Get()
			todoIds = append(todoIds, todo.ID)
		}

		for _, entity := range entities {
			w.RemoveEntity(entity)
		}
		ui.UIState.ClearCompletedClicked = false

		for _, todoId := range todoIds {
			ui.DeleteClickableForTodo(todoId)
		}
	}
}
func (s *RemoveCompleted) Finalize(w *ecs.World) {
	// no finalization needed
}
