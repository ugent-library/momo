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

// RepresentationCreate is the builder for creating a Representation entity.
type RepresentationCreate struct {
	config
	mutation *RepresentationMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (rc *RepresentationCreate) SetName(s string) *RepresentationCreate {
	rc.mutation.SetName(s)
	return rc
}

// SetData sets the "data" field.
func (rc *RepresentationCreate) SetData(b []byte) *RepresentationCreate {
	rc.mutation.SetData(b)
	return rc
}

// SetCreatedAt sets the "created_at" field.
func (rc *RepresentationCreate) SetCreatedAt(t time.Time) *RepresentationCreate {
	rc.mutation.SetCreatedAt(t)
	return rc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (rc *RepresentationCreate) SetNillableCreatedAt(t *time.Time) *RepresentationCreate {
	if t != nil {
		rc.SetCreatedAt(*t)
	}
	return rc
}

// SetUpdatedAt sets the "updated_at" field.
func (rc *RepresentationCreate) SetUpdatedAt(t time.Time) *RepresentationCreate {
	rc.mutation.SetUpdatedAt(t)
	return rc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (rc *RepresentationCreate) SetNillableUpdatedAt(t *time.Time) *RepresentationCreate {
	if t != nil {
		rc.SetUpdatedAt(*t)
	}
	return rc
}

// SetID sets the "id" field.
func (rc *RepresentationCreate) SetID(u uuid.UUID) *RepresentationCreate {
	rc.mutation.SetID(u)
	return rc
}

// SetRecID sets the "rec" edge to the Rec entity by ID.
func (rc *RepresentationCreate) SetRecID(id uuid.UUID) *RepresentationCreate {
	rc.mutation.SetRecID(id)
	return rc
}

// SetNillableRecID sets the "rec" edge to the Rec entity by ID if the given value is not nil.
func (rc *RepresentationCreate) SetNillableRecID(id *uuid.UUID) *RepresentationCreate {
	if id != nil {
		rc = rc.SetRecID(*id)
	}
	return rc
}

// SetRec sets the "rec" edge to the Rec entity.
func (rc *RepresentationCreate) SetRec(r *Rec) *RepresentationCreate {
	return rc.SetRecID(r.ID)
}

// Mutation returns the RepresentationMutation object of the builder.
func (rc *RepresentationCreate) Mutation() *RepresentationMutation {
	return rc.mutation
}

// Save creates the Representation in the database.
func (rc *RepresentationCreate) Save(ctx context.Context) (*Representation, error) {
	var (
		err  error
		node *Representation
	)
	rc.defaults()
	if len(rc.hooks) == 0 {
		if err = rc.check(); err != nil {
			return nil, err
		}
		node, err = rc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*RepresentationMutation)
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
func (rc *RepresentationCreate) SaveX(ctx context.Context) *Representation {
	v, err := rc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// defaults sets the default values of the builder before save.
func (rc *RepresentationCreate) defaults() {
	if _, ok := rc.mutation.CreatedAt(); !ok {
		v := representation.DefaultCreatedAt()
		rc.mutation.SetCreatedAt(v)
	}
	if _, ok := rc.mutation.UpdatedAt(); !ok {
		v := representation.DefaultUpdatedAt()
		rc.mutation.SetUpdatedAt(v)
	}
	if _, ok := rc.mutation.ID(); !ok {
		v := representation.DefaultID()
		rc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (rc *RepresentationCreate) check() error {
	if _, ok := rc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New("ent: missing required field \"name\"")}
	}
	if v, ok := rc.mutation.Name(); ok {
		if err := representation.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf("ent: validator failed for field \"name\": %w", err)}
		}
	}
	if _, ok := rc.mutation.Data(); !ok {
		return &ValidationError{Name: "data", err: errors.New("ent: missing required field \"data\"")}
	}
	if _, ok := rc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New("ent: missing required field \"created_at\"")}
	}
	if _, ok := rc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New("ent: missing required field \"updated_at\"")}
	}
	return nil
}

func (rc *RepresentationCreate) sqlSave(ctx context.Context) (*Representation, error) {
	_node, _spec := rc.createSpec()
	if err := sqlgraph.CreateNode(ctx, rc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return _node, nil
}

func (rc *RepresentationCreate) createSpec() (*Representation, *sqlgraph.CreateSpec) {
	var (
		_node = &Representation{config: rc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: representation.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: representation.FieldID,
			},
		}
	)
	if id, ok := rc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := rc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: representation.FieldName,
		})
		_node.Name = value
	}
	if value, ok := rc.mutation.Data(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBytes,
			Value:  value,
			Column: representation.FieldData,
		})
		_node.Data = value
	}
	if value, ok := rc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: representation.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := rc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: representation.FieldUpdatedAt,
		})
		_node.UpdatedAt = value
	}
	if nodes := rc.mutation.RecIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   representation.RecTable,
			Columns: []string{representation.RecColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: rec.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.rec_representations = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// RepresentationCreateBulk is the builder for creating many Representation entities in bulk.
type RepresentationCreateBulk struct {
	config
	builders []*RepresentationCreate
}

// Save creates the Representation entities in the database.
func (rcb *RepresentationCreateBulk) Save(ctx context.Context) ([]*Representation, error) {
	specs := make([]*sqlgraph.CreateSpec, len(rcb.builders))
	nodes := make([]*Representation, len(rcb.builders))
	mutators := make([]Mutator, len(rcb.builders))
	for i := range rcb.builders {
		func(i int, root context.Context) {
			builder := rcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*RepresentationMutation)
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
func (rcb *RepresentationCreateBulk) SaveX(ctx context.Context) []*Representation {
	v, err := rcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}