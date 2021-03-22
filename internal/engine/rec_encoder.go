package engine

import "io"

type RecEncoderEngine interface {
	NewRecEncoder(io.Writer, string) RecEncoder
}

type RecEncoderFactory func(io.Writer) RecEncoder

type RecEncoder interface {
	Encode(*Rec) error
}

func (e *engine) NewRecEncoder(w io.Writer, format string) RecEncoder {
	e.recEncodersMu.Lock()
	defer e.recEncodersMu.Unlock()
	factory, found := e.recEncoders[format]
	if !found {
		return nil
	}
	return factory(w)
}
