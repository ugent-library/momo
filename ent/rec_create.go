// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/ugent-library/momo/ent/rec"
	"github.com/ugent-library/momo/ent/representation"
)

// RecCreate is the builder for creating a Rec entity.
type RecCreate struct {
	config
	mutation *RecMutation
	hooks    []Hook
}

// SetCollection sets the "collection" field.
func (rc *RecCreate) SetCollection(s string) *RecCreate {
	rc.mutation.SetCollection(s)
	return rc
}

// SetType sets the "type" field.
func (rc *RecCreate) SetType(s string) *RecCreate {
	rc.mutation.SetType(s)
	return rc
}

// SetCreatedAt sets the "created_at" field.
func (rc *RecCreate) SetCreatedAt(t time.Time) *RecCreate {
	rc.mutation.SetCreatedAt(t)
	return rc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (rc *RecCreate) SetNillableCreatedAt(t *time.Time) *RecCreate {
	if t != nil {
		rc.SetCreatedAt(*t)
	}
	return rc
}

// SetUpdatedAt sets the "updated_at" field.
func (rc *RecCreate) SetUpdatedAt(t time.Time) *RecCreate {
	rc.mutation.SetUpdatedAt(t)
	return rc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (rc *RecCreate) SetNillableUpdatedAt(t *time.Time) *RecCreate {
	if t != nil {
		rc.SetUpdatedAt(*t)
	}
	return rc
}

// SetMetadata sets the "metadata" field.
func (rc *RecCreate) SetMetadata(m map[string]interface{}) *RecCreate {
	rc.mutation.SetMetadata(m)
	return rc
}

// SetSource sets the "source" field.
func (rc *RecCreate) SetSource(b []byte) *RecCreate {
	rc.mutation.SetSource(b)
	return rc
}

// SetID sets the "id" field.
func (rc *RecCreate) SetID(u uuid.UUID) *RecCreate {
	rc.mutation.SetID(u)
	return rc
}

// AddRepresentationIDs adds the "representations" edge to the Representation entity by IDs.
func (rc *RecCreate) AddRepresentationIDs(ids ...uuid.UUID) *RecCreate {
	rc.mutation.AddRepresentationIDs(ids...)
	return rc
}

// AddRepresentations adds the "representations" edges to the Representation entity.
func (rc *RecCreate) AddRepresentations(r ...*Representation) *RecCreate {
	ids := make([]uuid.UUID, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return rc.AddRepresentationIDs(ids...)
}

// Mutation returns the RecMutation object of the builder.
func (rc *RecCreate) Mutation() *RecMutation {
	return rc.mutation
}

// Save creates the Rec in the database.
func (rc *RecCreate) Save(ctx context.Context) (*Rec, error) {
	var (
		err  error
		node *Rec
	)
	rc.defaults()
	if len(rc.hooks) == 0 {
		if err = rc.check(); err != nil {
			return nil, err
		}
		node, err = rc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*RecMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = rc.check(); err != nil {
				return nil, err
			}
			rc.mutation = mutation
			node, err = rc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(rc.hooks) - 1; i >= 0; i-- {
			mut = rc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, rc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (rc *RecCreate) SaveX(ctx context.Context) *Rec {
	v, err := rc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// defaults sets the default values of the builder before save.
func (rc *RecCreate) defaults() {
	if _, ok := rc.mutation.CreatedAt(); !ok {
		v := rec.DefaultCreatedAt()
		rc.mutation.SetCreatedAt(v)
	}
	if _, ok := rc.mutation.UpdatedAt(); !ok {
		v := rec.DefaultUpdatedAt()
		rc.mutation.SetUpdatedAt(v)
	}
	if _, ok := rc.mutation.ID(); !ok {
		v := rec.DefaultID()
		rc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (rc *RecCreate) check() error {
	if _, ok := rc.mutation.Collection(); !ok {
		return &ValidationError{Name: "collection", err: errors.New("ent: missing required field \"collection\"")}
	}
	if v, ok := rc.mutation.Collection(); ok {
		if err := rec.CollectionValidator(v); err != nil {
			return &ValidationError{Name: "collection", err: fmt.Errorf("ent: validator failed for field \"collection\": %w", err)}
		}
	}
	if _, ok := rc.mutation.GetType(); !ok {
		return &ValidationError{Name: "type", err: errors.New("ent: missing required field \"type\"")}
	}
	if v, ok := rc.mutation.GetType(); ok {
		if err := rec.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf("ent: validator failed for field \"type\": %w", err)}
		}
	}
	if _, ok := rc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New("ent: missing required field \"created_at\"")}
	}
	if _, ok := rc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New("ent: missing required field \"updated_at\"")}
	}
	if _, ok := rc.mutation.Metadata(); !ok {
		return &ValidationError{Name: "metadata", err: errors.New("ent: missing required field \"metadata\"")}
	}
	return nil
}

func (rc *RecCreate) sqlSave(ctx context.Context) (*Rec, error) {
	_node, _spec := rc.createSpec()
	if err := sqlgraph.CreateNode(ctx, rc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return _node, nil
}

func (rc *RecCreate) createSpec() (*Rec, *sqlgraph.CreateSpec) {
	var (
		_node = &Rec{config: rc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: rec.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: rec.FieldID,
			},
		}
	)
	if id, ok := rc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := rc.mutation.Collection(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: rec.FieldCollection,
		})
		_node.Collection = value
	}
	if value, ok := rc.mutation.GetType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: rec.FieldType,
		})
		_node.Type = value
	}
	if value, ok := rc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: rec.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := rc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: rec.FieldUpdatedAt,
		})
		_node.UpdatedAt = value
	}
	if value, ok := rc.mutation.Metadata(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: rec.FieldMetadata,
		})
		_node.Metadata = value
	}
	if value, ok := rc.mutation.Source(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBytes,
			Value:  value,
			Column: rec.FieldSource,
		})
		_node.Source = value
	}
	if nodes := rc.mutation.RepresentationsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   rec.RepresentationsTable,
			Columns: []string{rec.RepresentationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: representation.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// RecCreateBulk is the builder for creating many Rec entities in bulk.
type RecCreateBulk struct {
	config
	builders []*RecCreate
}

// Save creates the Rec entities in the database.
func (rcb *RecCreateBulk) Save(ctx context.Context) ([]*Rec, error) {
	specs := make([]*sqlgraph.CreateSpec, len(rcb.builders))
	nodes := make([]*Rec, len(rcb.builders))
	mutators := make([]Mutator, len(rcb.builders))
	for i := range rcb.builders {
		func(i int, root context.Context) {
			builder := rcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*RecMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, rcb.builders[i+1].mutation)
				} else {
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, rcb.driver, &sqlgraph.BatchCreateSpec{Nodes: specs}); err != nil {
						if cerr, ok := isSQLConstraintError(err); ok {
							err = cerr
						}
					}
				}
				mutation.done = true
				if err != nil {
					return nil, err
				}
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, rcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (rcb *RecCreateBulk) SaveX(ctx context.Context) []*Rec {
	v, err := rcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}
