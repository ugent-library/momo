package pg

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/google/uuid"
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
	ID              string         `gorm:"type:uuid;not null;uniqueIndex"`
	Collection      string         `gorm:"not null;index"`
	Type            string         `gorm:"not null;index"`
	Metadata        datatypes.JSON `gorm:"not null"`
	Source          datatypes.JSON
	Representations []Representation `gorm:"foreignKey:RecFK"`
}

type Representation struct {
	Model
	RecFK int64
	Name  string `gorm:"not null;index"`
	Data  []byte `gorm:"not null"`
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
	s.db.AutoMigrate(&Rec{}, &Representation{})
	return s, nil
}

func (s *store) GetRec(id string) (*engine.Rec, error) {
	r := Rec{}
	res := s.db.Where("id = ?", id).First(&r)
	if res.Error != nil {
		return nil, res.Error
	}
	return reifyRec(&r), nil
}

func (s *store) EachRec(fn func(*engine.Rec) bool) error {
	rows, err := s.db.Model(&Rec{}).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		rec := Rec{}
		if err := s.db.ScanRows(rows, &rec); err != nil {
			return err
		}
		if ok := fn(reifyRec(&rec)); !ok {
			return nil
		}
	}
	return nil
}

func (s *store) AddRec(rec *engine.Rec) error {
	if rec.ID == "" {
		rec.ID = uuid.NewString()
	}
	r := Rec{
		ID:         rec.ID,
		Type:       rec.Type,
		Collection: rec.Collection,
		Metadata:   datatypes.JSON(rec.RawMetadata),
		Source:     datatypes.JSON(rec.RawSource),
	}

	res := s.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"type", "collection",
			"metadata", "source", "updated_at"}),
	}).Create(&r)

	rec.CreatedAt = r.CreatedAt
	rec.UpdatedAt = r.UpdatedAt

	return res.Error
}

func (s *store) AddRepresentation(recID, name string, data []byte) error {
	rec := Rec{}
	res := s.db.Where("id = ?", recID).First(&rec)
	if res.Error != nil {
		return res.Error
	}

	rep := Representation{
		RecFK: rec.PK,
		Name:  name,
		Data:  data,
	}

	return s.db.Model(&rec).Association("Representations").Append(&rep)
}

func (s *store) Reset() error {
	err := s.db.Migrator().DropTable(&Rec{}, &Representation{})
	if err != nil {
		return err
	}
	return s.db.AutoMigrate(&Rec{}, &Representation{})

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
