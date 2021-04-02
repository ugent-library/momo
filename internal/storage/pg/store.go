package pg

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ugent-library/momo/internal/engine"
)

type Model struct {
	PK        int64          `gorm:"primaryKey"`
	CreatedAt time.Time      `gorm:"not null;autoCreateTime:milli;type:timestamp;index"`
	UpdatedAt time.Time      `gorm:"not null;autoUpdateTime:milli;type:timestamp;index"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp;index"`
}

type Rec struct {
	Model
	ID         string         `gorm:"not null;uniqueIndex:idx_recs_id_collection"`
	Collection string         `gorm:"not null;uniqueIndex:idx_recs_id_collection"`
	Type       string         `gorm:"not null;index"`
	Metadata   datatypes.JSON `gorm:"not null"`
	Source     datatypes.JSON
}

type store struct {
	db *gorm.DB
}

func New(dsn string) (*store, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	s := &store{db: db}
	s.db.AutoMigrate(&Rec{})
	return s, nil
}

func (s *store) GetRec(coll string, id string) (*engine.Rec, error) {
	r := Rec{}
	res := s.db.Where("collection = ?", coll).Where("id = ?", id).First(&r)
	if res.Error != nil {
		return nil, res.Error
	}
	return reifyRec(&r), nil
}

func (s *store) AllRecs() engine.RecCursor {
	rows, err := s.db.Model(&Rec{}).Rows()
	return &recCursor{db: s.db, rows: rows, err: err}
}

func (s *store) AddRec(rec *engine.Rec) error {
	r := Rec{
		ID:         rec.ID,
		Type:       rec.Type,
		Collection: rec.Collection,
		Metadata:   datatypes.JSON(rec.RawMetadata),
		Source:     datatypes.JSON(rec.RawSource),
	}

	res := s.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}, {Name: "collection"}},
		DoUpdates: clause.AssignmentColumns([]string{"type", "collection",
			"metadata", "source", "updated_at"}),
	}).Create(&r)

	return res.Error
}

func (s *store) Reset() error {
	err := s.db.Migrator().DropTable(&Rec{})
	if err != nil {
		return err
	}
	return s.db.AutoMigrate(&Rec{})

}

func reifyRec(r *Rec) *engine.Rec {
	return &engine.Rec{
		ID:          r.ID,
		Type:        r.Type,
		Collection:  r.Collection,
		RawMetadata: json.RawMessage(r.Metadata),
		RawSource:   json.RawMessage(r.Source),
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
