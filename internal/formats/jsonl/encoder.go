package jsonl

import (
	"encoding/json"
	"io"

	"github.com/ugent-library/momo/internal/engine"
)

type encoder struct {
	json *json.Encoder
}

func NewEncoder(w io.Writer) engine.RecEncoder {
	return &encoder{json.NewEncoder(w)}
}

func (e *encoder) Encode(rec *engine.Rec) error {
	return e.json.Encode(rec)
}
