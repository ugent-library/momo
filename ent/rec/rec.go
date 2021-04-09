// Code generated by entc, DO NOT EDIT.

package rec

import (
	"time"

	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the rec type in the database.
	Label = "rec"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCollection holds the string denoting the collection field in the database.
	FieldCollection = "collection"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldMetadata holds the string denoting the metadata field in the database.
	FieldMetadata = "metadata"
	// EdgeRepresentations holds the string denoting the representations edge name in mutations.
	EdgeRepresentations = "representations"
	// Table holds the table name of the rec in the database.
	Table = "recs"
	// RepresentationsTable is the table the holds the representations relation/edge.
	RepresentationsTable = "representations"
	// RepresentationsInverseTable is the table name for the Representation entity.
	// It exists in this package in order to avoid circular dependency with the "representation" package.
	RepresentationsInverseTable = "representations"
	// RepresentationsColumn is the table column denoting the representations relation/edge.
	RepresentationsColumn = "rec_representations"
)

// Columns holds all SQL columns for rec fields.
var Columns = []string{
	FieldID,
	FieldCollection,
	FieldType,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldMetadata,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// CollectionValidator is a validator for the "collection" field. It is called by the builders before save.
	CollectionValidator func(string) error
	// TypeValidator is a validator for the "type" field. It is called by the builders before save.
	TypeValidator func(string) error
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)
