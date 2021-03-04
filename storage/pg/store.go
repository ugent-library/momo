package pg

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/lib/pq"
	"github.com/ugent-library/momo/records"
)

type Model struct {
	ID        int64          `gorm:"primaryKey"`
	CreatedAt time.Time      `gorm:"autoCreateTime:milli;type:timestamp;index"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:milli;type:timestamp;index"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp;index"`
}

type Rec struct {
	Model
	RecID      string         `gorm:"not null;uniqueIndex"`
	Type       string         `gorm:"not null;index"`
	Collection pq.StringArray `gorm:"type:text[];not null;index"`
	Title      string         `gorm:"not null"`
	Metadata   datatypes.JSON `gorm:"not null"`
}

type Store struct {
	db *gorm.DB
}

func New(dsn string) (*Store, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	s := &Store{db: db}
	s.db.AutoMigrate(&Rec{})
	return s, nil
}

func (s *Store) GetRec(id string) (*records.Rec, error) {
	r := Rec{}
	res := s.db.Where("rec_id = ?", id).First(&r)
	if res.Error != nil {
		return nil, res.Error
	}
	rec := records.Rec{
		ID:         r.RecID,
		Type:       r.Type,
		Collection: r.Collection,
		Title:      r.Title,
		Metadata:   json.RawMessage(r.Metadata),
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
	return &rec, nil
}

func (s *Store) AddRec(rec *records.Rec) error {
	r := Rec{
		RecID:      rec.ID,
		Type:       rec.Type,
		Collection: pq.StringArray(rec.Collection),
		Title:      rec.Title,
		Metadata:   datatypes.JSON(rec.Metadata),
	}
	// TODO upsert
	res := s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&r)
	return res.Error
}
