package sql

// edit schema in './ent/schema'
// generate client with 'ent generate --idtype string ./ent/schema/'

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

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
		ID:             r.ID.String(),
		Collection:     r.Collection,
		Type:           r.Type,
		Metadata:       r.Metadata,
		SourceID:       r.SourceID,
		SourceFormat:   r.SourceFormat,
		SourceMetadata: r.SourceMetadata,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
	return rec, nil
}

// ent: https://github.com/ent/ent/issues/215
func (s *store) EachRec(fn func(*engine.Rec) bool) error {
	rows, err := s.db.Query(`SELECT id, collection, type, metadata, source_id, source_format, source_metadata, created_at, updated_at
	FROM recs`)

	defer rows.Close()

	for rows.Next() {
		var rec engine.Rec
		var rawMetadata json.RawMessage

		err = rows.Scan(&rec.ID, &rec.Collection, &rec.Type, &rawMetadata, &rec.SourceID, &rec.SourceFormat,
			&rec.SourceMetadata, &rec.CreatedAt, &rec.UpdatedAt)
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
func (s *store) AddRecBySourceID(rec *engine.Rec) error {
	stmt := `INSERT INTO recs (id, collection, type, metadata, source_id, source_format, source_metadata, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (source_id) DO
   	UPDATE SET collection = $2, type = $3, metadata = $4, source_format = $6, source_metadata = $7, updated_at = $9
	RETURNING id, created_at, updated_at`

	metadata, err := json.Marshal(rec.Metadata)
	if err != nil {
		return err
	}

	var (
		now       = time.Now()
		id        string
		createdAt time.Time
		updatedAt time.Time
	)

	err = s.db.QueryRow(stmt, rec.ID, rec.Collection, rec.Type, metadata, rec.SourceID, rec.SourceFormat,
		rec.SourceMetadata, now, now).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return err
	}

	rec.ID = id
	rec.CreatedAt = createdAt
	rec.UpdatedAt = updatedAt

	return nil
}

func (s *store) UpdateRecMetadata(id string, m map[string]interface{}) (*engine.Rec, error) {
	r, err := s.client.Rec.
		UpdateOneID(uuid.MustParse(id)).
		SetMetadata(m).
		Save(context.Background())
	if err != nil {
		return nil, err
	}
	rec := &engine.Rec{
		ID:             r.ID.String(),
		Collection:     r.Collection,
		Type:           r.Type,
		Metadata:       r.Metadata,
		SourceID:       r.SourceID,
		SourceFormat:   r.SourceFormat,
		SourceMetadata: r.SourceMetadata,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
	return rec, nil
}

func (s *store) GetRepresentation(recID, format string) (*engine.Representation, error) {
	r, err := s.client.Representation.
		Query().
		Where(
			entrep.HasRecWith(entrec.ID(uuid.MustParse(recID))),
			entrep.Format(format)).
		Only(context.Background())
	if err != nil {
		return nil, err
	}
	rep := &engine.Representation{
		ID:        r.ID.String(),
		RecID:     recID,
		Format:    r.Format,
		Data:      r.Data,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
	return rep, nil
}

func (s *store) AddRepresentation(rep *engine.Representation) error {
	_, err := s.client.Representation.
		Create().
		SetFormat(rep.Format).
		SetData(rep.Data).
		SetRecID(uuid.MustParse(rep.RecID)).
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
