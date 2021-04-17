package sql

// edit schema in './ent/schema'
// generate client with 'ent generate --idtype string ./ent/schema/'

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/ugent-library/momo/ent"
	"github.com/ugent-library/momo/ent/migrate"
	entrec "github.com/ugent-library/momo/ent/rec"
	entrep "github.com/ugent-library/momo/ent/representation"
	"github.com/ugent-library/momo/internal/engine"
)

type store struct {
	db     *sql.DB
	client *ent.Client
}

func New(dialect, dsn string) (*store, error) {
	driver, err := entsql.Open(dialect, dsn)
	if err != nil {
		return nil, err
	}
	db := driver.DB()
	client := ent.NewClient(ent.Driver(driver))
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, err
	}
	s := &store{client: client, db: db}
	return s, nil
}

func (s *store) GetRec(id string) (*engine.Rec, error) {
	r, err := s.client.Rec.
		Query().
		Where(entrec.ID(uuid.MustParse(id))).
		Only(context.Background())
	if err != nil {
		return nil, err
	}
	rec := &engine.Rec{
		ID:         r.ID.String(),
		Collection: r.Collection,
		Type:       r.Type,
		Metadata:   r.Metadata,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
		RawSource:  r.Source,
	}
	return rec, nil
}

// ent: https://github.com/ent/ent/issues/215
func (s *store) EachRec(fn func(*engine.Rec) bool) error {
	rows, err := s.db.Query("SELECT id, collection, type, metadata, created_at, updated_at, source FROM recs")

	defer rows.Close()

	for rows.Next() {
		var rec engine.Rec
		var rawMetadata json.RawMessage
		err = rows.Scan(
			&rec.ID,
			&rec.Collection,
			&rec.Type,
			&rawMetadata,
			&rec.CreatedAt,
			&rec.UpdatedAt,
			&rec.RawSource,
		)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(rawMetadata, &rec.Metadata); err != nil {
			return err
		}
		if ok := fn(&rec); !ok {
			return nil
		}
	}

	return rows.Err()
}

// ent: support upsert
func (s *store) AddRec(rec *engine.Rec) error {
	r, err := s.client.Rec.
		Create().
		SetID(uuid.MustParse(rec.ID)).
		SetCollection(rec.Collection).
		SetType(rec.Type).
		SetMetadata(rec.Metadata).
		SetSource(rec.RawSource).
		Save(context.Background())
	if err != nil {
		return err
	}

	rec.CreatedAt = r.CreatedAt
	rec.UpdatedAt = r.UpdatedAt

	return nil
}

func (s *store) GetRepresentation(recID, name string) (*engine.Representation, error) {
	r, err := s.client.Representation.
		Query().
		Where(
			entrep.HasRecWith(entrec.ID(uuid.MustParse(recID))),
			entrep.Format(name)).
		Only(context.Background())
	if err != nil {
		return nil, err
	}
	rep := &engine.Representation{
		ID:        r.ID.String(),
		Format:    r.Format,
		Data:      r.Data,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
	return rep, nil
}

func (s *store) AddRepresentation(recID, name string, data []byte) error {
	_, err := s.client.Representation.
		Create().
		SetFormat(name).
		SetData(data).
		SetRecID(uuid.MustParse(recID)).
		Save(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// ent: add migrate.WithDropTable(true)
func (s *store) Reset() error {
	for _, tbl := range migrate.Tables {
		stmt := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tbl.Name)
		if _, err := s.db.Exec(stmt); err != nil {
			return err
		}
	}
	return s.client.Schema.Create(context.Background())
}
