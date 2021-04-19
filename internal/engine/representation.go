package engine

import (
	"bytes"
	"time"
)

type RepresentationEngine interface {
	GetRepresentation(string, string) (*Representation, error)
	AddRepresentation(*Representation) error
}

type Representation struct {
	ID        string    `json:"id"`
	RecID     string    `json:"rec_id"`
	Format    string    `json:"format"`
	Data      []byte    `json:"data"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (e *engine) GetRepresentation(recID, format string) (*Representation, error) {
	switch format {
	case "mla":
		return e.store.GetRepresentation(recID, format)
	default:
		rec, err := e.GetRec(recID)
		if err != nil {
			return nil, err
		}
		var b bytes.Buffer
		enc := e.NewRecEncoder(&b, format)
		if enc == nil {
			panic("Unknown format " + format)
		}
		if err := enc.Encode(rec); err != nil {
			return nil, err
		}
		rep := &Representation{
			RecID:  recID,
			Format: format,
			Data:   b.Bytes(),
		}
		return rep, nil
	}
}

func (e *engine) AddRepresentation(rep *Representation) error {
	return e.store.AddRepresentation(rep)
}
