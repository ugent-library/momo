// Code generated by entc, DO NOT EDIT.

package representation

import (
	"time"

	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the representation type in the database.
	Label = "representation"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldFormat holds the string denoting the format field in the database.
	FieldFormat = "format"
	// FieldData holds the string denoting the data field in the database.
	FieldData = "data"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// EdgeRec holds the string denoting the rec edge name in mutations.
	EdgeRec = "rec"
	// Table holds the table name of the representation in the database.
	Table = "representations"
	// RecTable is the table the holds the rec relation/edge.
	RecTable = "representations"
	// RecInverseTable is the table name for the Rec entity.
	// It exists in this package in order to avoid circular dependency with the "rec" package.
	RecInverseTable = "recs"
	// RecColumn is the table column denoting the rec relation/edge.
	RecColumn = "rec_id"
)

// Columns holds all SQL columns for representation fields.
var Columns = []string{
	FieldID,
	FieldFormat,
	FieldData,
	FieldCreatedAt,
	FieldUpdatedAt,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "representations"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"rec_id",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// FormatValidator is a validator for the "format" field. It is called by the builders before save.
	FormatValidator func(string) error
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)
