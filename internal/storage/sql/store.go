package sql

// edit schema in './ent/schema'
// generate client with 'ent generate --idtype string ./ent/schema'

import (
	"context"
	"database/sql"
	"fmt"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/ugent-library/momo/ent"
	"github.com/ugent-library/momo/ent/migrate"
	entrec "github.com/ugent-library/momo/ent/rec"
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
	}
	return rec, nil
}

// ent: cursor support
// TODO scan directly in type?
func (s *store) EachRec(fn func(*engine.Rec) bool) error {
	offset := 0
	limit := 200

	ctx := context.Background()
	recs, err := s.client.Rec.Query().Limit(limit).All(ctx)
	if err != nil {
		return err
	}

	for len(recs) > 0 {
		for _, r := range recs {
			rec := engine.Rec{
				ID:         r.ID.String(),
				Collection: r.Collection,
				Type:       r.Type,
				Metadata:   r.Metadata,
				CreatedAt:  r.CreatedAt,
				UpdatedAt:  r.UpdatedAt,
			}
			if ok := fn(&rec); !ok {
				return nil
			}
		}

		offset += limit
		if recs, err = s.client.Rec.Query().Limit(limit).Offset(offset).All(ctx); err != nil {
			return err
		}
	}

	return nil
}

// ent: support upsert
func (s *store) AddRec(rec *engine.Rec) error {
	r, err := s.client.Rec.
		Create().
		SetID(uuid.MustParse(rec.ID)).
		SetCollection(rec.Collection).
		SetType(rec.Type).
		SetMetadata(rec.Metadata).
		Save(context.Background())
	if err != nil {
		return err
	}

	rec.CreatedAt = r.CreatedAt
	rec.UpdatedAt = r.UpdatedAt

	return nil
}

func (s *store) AddRepresentation(recID, name string, data []byte) error {
	_, err := s.client.Representation.
		Create().
		SetName(name).
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
