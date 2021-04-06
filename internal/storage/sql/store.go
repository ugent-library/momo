package sql

import (
	"context"

	"github.com/ugent-library/momo/internal/engine"
	"github.com/ugent-library/momo/internal/storage/sql/ent"

	_ "github.com/lib/pq"
)

type store struct {
	client *ent.Client
}

func New(driver, dsn string) (*store, error) {
	client, err := ent.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, err
	}
	s := &store{client}
	return s, nil
}

func (s *store) AddRec(rec *engine.Rec) error {
	return nil
}
