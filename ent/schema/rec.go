package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Rec holds the schema definition for the Rec entity.
type Rec struct {
	ent.Schema
}

// Fields of the Rec.
func (Rec) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("collection").
			NotEmpty(),
		field.String("type").
			NotEmpty(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Immutable(),
		field.JSON("metadata", map[string]interface{}{}),
		field.Bytes("source").
			Optional(),
	}
}

// Edges of the Rec.
func (Rec) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("representations", Representation.Type),
	}
}

// Indexes of the Rec.
func (Rec) Indexes() []ent.Index {
	return nil
}
