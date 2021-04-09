package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Representation holds the schema definition for the Representation entity.
type Representation struct {
	ent.Schema
}

// Fields of the Representation.
func (Representation) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("name").
			NotEmpty(),
		// ent: why no NotEmpty?
		field.Bytes("data"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Immutable(),
	}
}

// Edges of the Representation.
func (Representation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("rec", Rec.Type).
			Ref("representations").
			Unique(),
	}
}

// Indexes of the Representation.
func (Representation) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Edges("rec").
			Unique(),
	}
}
