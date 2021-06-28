// Code generated by entc, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/ugent-library/momo/ent/rec"
)

// Rec is the model entity for the Rec schema.
type Rec struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// Collection holds the value of the "collection" field.
	Collection string `json:"collection,omitempty"`
	// Type holds the value of the "type" field.
	Type string `json:"type,omitempty"`
	// Metadata holds the value of the "metadata" field.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// SourceID holds the value of the "source_id" field.
	SourceID string `json:"source_id,omitempty"`
	// SourceFormat holds the value of the "source_format" field.
	SourceFormat string `json:"source_format,omitempty"`
	// SourceMetadata holds the value of the "source_metadata" field.
	SourceMetadata []byte `json:"source_metadata,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the RecQuery when eager-loading is set.
	Edges RecEdges `json:"edges"`
}

// RecEdges holds the relations/edges for other nodes in the graph.
type RecEdges struct {
	// Representations holds the value of the representations edge.
	Representations []*Representation `json:"representations,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}
	
// RepresentationsOrErr returns the Representations value or an error if the edge
// was not loaded in eager-loading.
func (e RecEdges) RepresentationsOrErr() ([]*Representation, error) {
	if e.loadedTypes[0] {
		return e.Representations, nil
	}
	return nil, &NotLoadedError{edge: "representations"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Rec) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case rec.FieldMetadata, rec.FieldSourceMetadata:
			values[i] = new([]byte)
		case rec.FieldCollection, rec.FieldType, rec.FieldSourceID, rec.FieldSourceFormat:
			values[i] = new(sql.NullString)
		case rec.FieldCreatedAt, rec.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case rec.FieldID:
			values[i] = new(uuid.UUID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Rec", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Rec fields.
func (r *Rec) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case rec.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				r.ID = *value
			}
		case rec.FieldCollection:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field collection", values[i])
			} else if value.Valid {
				r.Collection = value.String
			}
		case rec.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				r.Type = value.String
			}
		case rec.FieldMetadata:

			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field metadata", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &r.Metadata); err != nil {
					return fmt.Errorf("unmarshal field metadata: %w", err)
				}
			}
		case rec.FieldSourceID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field source_id", values[i])
			} else if value.Valid {
				r.SourceID = value.String
			}
		case rec.FieldSourceFormat:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field source_format", values[i])
			} else if value.Valid {
				r.SourceFormat = value.String
			}
		case rec.FieldSourceMetadata:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field source_metadata", values[i])
			} else if value != nil {
				r.SourceMetadata = *value
			}
		case rec.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				r.CreatedAt = value.Time
			}
		case rec.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				r.UpdatedAt = value.Time
			}
		}
	}
	return nil
}

// QueryRepresentations queries the "representations" edge of the Rec entity.
func (r *Rec) QueryRepresentations() *RepresentationQuery {
	return (&RecClient{config: r.config}).QueryRepresentations(r)
}

// Update returns a builder for updating this Rec.
// Note that you need to call Rec.Unwrap() before calling this method if this Rec
// was returned from a transaction, and the transaction was committed or rolled back.
func (r *Rec) Update() *RecUpdateOne {
	return (&RecClient{config: r.config}).UpdateOne(r)
}

// Unwrap unwraps the Rec entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (r *Rec) Unwrap() *Rec {
	tx, ok := r.config.driver.(*txDriver)
	if !ok {
		panic("ent: Rec is not a transactional entity")
	}
	r.config.driver = tx.drv
	return r
}

// String implements the fmt.Stringer.
func (r *Rec) String() string {
	var builder strings.Builder
	builder.WriteString("Rec(")
	builder.WriteString(fmt.Sprintf("id=%v", r.ID))
	builder.WriteString(", collection=")
	builder.WriteString(r.Collection)
	builder.WriteString(", type=")
	builder.WriteString(r.Type)
	builder.WriteString(", metadata=")
	builder.WriteString(fmt.Sprintf("%v", r.Metadata))
	builder.WriteString(", source_id=")
	builder.WriteString(r.SourceID)
	builder.WriteString(", source_format=")
	builder.WriteString(r.SourceFormat)
	builder.WriteString(", source_metadata=")
	builder.WriteString(fmt.Sprintf("%v", r.SourceMetadata))
	builder.WriteString(", created_at=")
	builder.WriteString(r.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(r.UpdatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// Recs is a parsable slice of Rec.
type Recs []*Rec

func (r Recs) config(cfg config) {
	for _i := range r {
		r[_i].config = cfg
	}
}
