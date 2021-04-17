// Code generated by entc, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/ugent-library/momo/ent/rec"
	"github.com/ugent-library/momo/ent/representation"
	"github.com/ugent-library/momo/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	recFields := schema.Rec{}.Fields()
	_ = recFields
	// recDescCollection is the schema descriptor for collection field.
	recDescCollection := recFields[1].Descriptor()
	// rec.CollectionValidator is a validator for the "collection" field. It is called by the builders before save.
	rec.CollectionValidator = recDescCollection.Validators[0].(func(string) error)
	// recDescType is the schema descriptor for type field.
	recDescType := recFields[2].Descriptor()
	// rec.TypeValidator is a validator for the "type" field. It is called by the builders before save.
	rec.TypeValidator = recDescType.Validators[0].(func(string) error)
	// recDescCreatedAt is the schema descriptor for created_at field.
	recDescCreatedAt := recFields[3].Descriptor()
	// rec.DefaultCreatedAt holds the default value on creation for the created_at field.
	rec.DefaultCreatedAt = recDescCreatedAt.Default.(func() time.Time)
	// recDescUpdatedAt is the schema descriptor for updated_at field.
	recDescUpdatedAt := recFields[4].Descriptor()
	// rec.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	rec.DefaultUpdatedAt = recDescUpdatedAt.Default.(func() time.Time)
	// rec.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	rec.UpdateDefaultUpdatedAt = recDescUpdatedAt.UpdateDefault.(func() time.Time)
	// recDescID is the schema descriptor for id field.
	recDescID := recFields[0].Descriptor()
	// rec.DefaultID holds the default value on creation for the id field.
	rec.DefaultID = recDescID.Default.(func() uuid.UUID)
	representationFields := schema.Representation{}.Fields()
	_ = representationFields
	// representationDescFormat is the schema descriptor for format field.
	representationDescFormat := representationFields[1].Descriptor()
	// representation.FormatValidator is a validator for the "format" field. It is called by the builders before save.
	representation.FormatValidator = representationDescFormat.Validators[0].(func(string) error)
	// representationDescCreatedAt is the schema descriptor for created_at field.
	representationDescCreatedAt := representationFields[3].Descriptor()
	// representation.DefaultCreatedAt holds the default value on creation for the created_at field.
	representation.DefaultCreatedAt = representationDescCreatedAt.Default.(func() time.Time)
	// representationDescUpdatedAt is the schema descriptor for updated_at field.
	representationDescUpdatedAt := representationFields[4].Descriptor()
	// representation.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	representation.DefaultUpdatedAt = representationDescUpdatedAt.Default.(func() time.Time)
	// representation.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	representation.UpdateDefaultUpdatedAt = representationDescUpdatedAt.UpdateDefault.(func() time.Time)
	// representationDescID is the schema descriptor for id field.
	representationDescID := representationFields[0].Descriptor()
	// representation.DefaultID holds the default value on creation for the id field.
	representation.DefaultID = representationDescID.Default.(func() uuid.UUID)
}
