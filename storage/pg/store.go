package pg

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ugent-library/momo/records"
)

type Base struct {
	ID        int64          `gorm:"primaryKey"`
	CreatedAt time.Time      `gorm:"index"`
	UpdatedAt time.Time      `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Rec struct {
	Base
	RecID string `gorm:"uniqueIndex"`
	// Collection []string
	Type     string
	Title    string
	Metadata datatypes.JSON
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
	return nil, nil
}

func (s *Store) AddRec(rec *records.Rec) error {
	r := Rec{
		RecID:    rec.ID,
		Type:     rec.Type,
		Title:    rec.Title,
		Metadata: datatypes.JSON(rec.Metadata),
	}
	res := s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&r)
	return res.Error
}

func (s *Store) AddRecs(c <-chan *records.Rec) {
}
