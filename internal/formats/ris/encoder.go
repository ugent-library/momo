package ris

import (
	"fmt"
	"io"

	"github.com/ugent-library/momo/internal/engine"
)

type encoder struct {
	writer io.Writer
}

type visitor func(io.Writer, *engine.Rec) error

var visitors = []visitor{
	addID,
	addTitle,
	addAuthor,
	addAbstract,
	addEdition,
	addPublisher,
	addPlaceOfPublication,
	addPublicationDate,
	addDOI,
	addISBN,
	addEndOfRecord,
}

func NewEncoder(w io.Writer) engine.RecEncoder {
	return &encoder{w}
}

func (e *encoder) Encode(rec *engine.Rec) (err error) {
	w := e.writer

	for _, v := range visitors {
		if err = v(w, rec); err != nil {
			return
		}
	}
	return
}

func addID(w io.Writer, rec *engine.Rec) error {
	return addTag(w, "ID", rec.ID)
}

func addEndOfRecord(w io.Writer, rec *engine.Rec) error {
	_, err := io.WriteString(w, "ER  - \r\n")
	return err
}

func addTitle(w io.Writer, rec *engine.Rec) error {
	if v := rec.Metadata().Title; v != "" {
		return addTag(w, "TI", v)
	}
	return nil
}

func addAuthor(w io.Writer, rec *engine.Rec) error {
	for _, v := range rec.Metadata().Author {
		return addTag(w, "AU", v.Name)
	}
	return nil
}

func addAbstract(w io.Writer, rec *engine.Rec) error {
	for _, v := range rec.Metadata().Abstract {
		return addTag(w, "AB", v.Text)
	}
	return nil
}

func addEdition(w io.Writer, rec *engine.Rec) error {
	if v := rec.Metadata().Edition; v != "" {
		return addTag(w, "ET", v)
	}
	return nil
}

func addPublisher(w io.Writer, rec *engine.Rec) error {
	if v := rec.Metadata().Publisher; v != "" {
		return addTag(w, "PB", v)
	}
	return nil
}

func addPlaceOfPublication(w io.Writer, rec *engine.Rec) error {
	if v := rec.Metadata().PlaceOfPublication; v != "" {
		return addTag(w, "CY", v)
	}
	return nil
}

func addPublicationDate(w io.Writer, rec *engine.Rec) error {
	if v := rec.Metadata().PublicationDate; v != "" {
		// TODO
	}
	return nil
}

func addDOI(w io.Writer, rec *engine.Rec) error {
	return addTag(w, "DO", rec.Metadata().DOI...)
}

func addISBN(w io.Writer, rec *engine.Rec) error {
	return addTag(w, "SN", rec.Metadata().ISBN...)
}

func addNote(w io.Writer, rec *engine.Rec) error {
	for _, v := range rec.Metadata().Note {
		return addTag(w, "N1", v.Text)
	}
	return nil
}

func addTag(w io.Writer, tag string, vals ...string) error {
	for _, val := range vals {
		if _, err := fmt.Fprintf(w, "%s  - %s\r\n", tag, val); err != nil {
			return err
		}
	}
	return nil
}
