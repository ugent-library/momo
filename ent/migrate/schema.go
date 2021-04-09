// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// RecsColumns holds the columns for the "recs" table.
	RecsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "collection", Type: field.TypeString},
		{Name: "type", Type: field.TypeString},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "metadata", Type: field.TypeJSON},
		{Name: "source", Type: field.TypeBytes, Nullable: true},
	}
	// RecsTable holds the schema information for the "recs" table.
	RecsTable = &schema.Table{
		Name:        "recs",
		Columns:     RecsColumns,
		PrimaryKey:  []*schema.Column{RecsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// RepresentationsColumns holds the columns for the "representations" table.
	RepresentationsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "name", Type: field.TypeString},
		{Name: "data", Type: field.TypeBytes},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "rec_representations", Type: field.TypeUUID, Nullable: true},
	}
	// RepresentationsTable holds the schema information for the "representations" table.
	RepresentationsTable = &schema.Table{
		Name:       "representations",
		Columns:    RepresentationsColumns,
		PrimaryKey: []*schema.Column{RepresentationsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "representations_recs_representations",
				Columns:    []*schema.Column{RepresentationsColumns[5]},
				RefColumns: []*schema.Column{RecsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "representation_name_rec_representations",
				Unique:  true,
				Columns: []*schema.Column{RepresentationsColumns[1], RepresentationsColumns[5]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		RecsTable,
		RepresentationsTable,
	}
)

func init() {
	RepresentationsTable.ForeignKeys[0].RefTable = RecsTable
}