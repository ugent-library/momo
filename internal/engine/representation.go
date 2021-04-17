package engine

import "time"

type Representation struct {
	ID        string    `json:"id"`
	Format    string    `json:"format"`
	Data      []byte    `json:"data"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
