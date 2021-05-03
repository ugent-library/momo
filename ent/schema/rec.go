package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Rec holds the schema definition for the Rec entity.
type Rec struct {
	ent.Schema
}

func (Rec) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("collection").
			NotEmpty(),
		field.String("type").
			NotEmpty(),
		field.JSON("metadata", map[string]interface{}{}),
		field.String("source_id").
			NotEmpty().
			Unique(),
		field.String("source_format").
			Optional(),
		field.Bytes("source_metadata").
			Optional(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Immutable(),
	}
}

func (Rec) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("representations", Representation.Type).
			StorageKey(edge.Column("rec_id")),
	}
}

func (Rec) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("collection"),
		index.Fields("type"),
		index.Fields("created_at"),
		index.Fields("updated_at"),
	}
}
