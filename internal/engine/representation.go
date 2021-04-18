package engine

import "time"

type RepresentationEngine interface {
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

func (e *engine) AddRepresentation(rep *Representation) error {
	return e.store.AddRepresentation(rep)
}
