package database

import "go.uber.org/fx"

const MODEL_GROUP_NAME = "db-entity"

type Entity interface{}
type EntityRegistry struct {
	fx.In
	Models []Entity `group:"db-entity"`
}

// AsModel registers a model with FX group
func AsModel(models ...Entity) fx.Option {

	entities := make([]interface{}, len(models))
	for i, model := range models {
		entities[i] = fx.Annotate(
			func() Entity { return model },
			fx.ResultTags(`group:"`+MODEL_GROUP_NAME+`"`),
		)
	}

	return fx.Provide(entities...)
}
