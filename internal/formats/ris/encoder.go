package ris

import (
	"fmt"
	"io"

	"github.com/ugent-library/momo/internal/engine"
	"github.com/ugent-library/momo/internal/metadata"
)

type encoder struct {
	writer io.Writer
}

type visitor func(io.Writer, metadata.Rec) error

var visitors = []visitor{
	addID,
	addTitle,
	addAuthor,
	addAbstract,
	addEdition,
	addPublisher,
	addPlaceOfPublication,
	addDOI,
	addISBN,
}

func NewEncoder(w io.Writer) engine.RecEncoder {
	return &encoder{w}
}

func (e *encoder) Encode(rec *engine.Rec) (err error) {
	w := e.writer
	r := metadata.WrapRec(rec)

	for _, v := range visitors {
		if err = v(w, r); err != nil {
			return
		}
	}
	_, err = io.WriteString(w, "ER  - \r\n")

	return
}

func addID(w io.Writer, rec metadata.Rec) error {
	return addTag(w, "ID", rec.ID)
}

func addTitle(w io.Writer, rec metadata.Rec) error {
	return addTag(w, "TI", rec.Title())
}

func addAuthor(w io.Writer, rec metadata.Rec) error {
	for _, val := range rec.Author() {
		if err := addTag(w, "AU", val.Name); err != nil {
			return err
		}
	}
	return nil
}

func addAbstract(w io.Writer, rec metadata.Rec) error {
	for _, val := range rec.Abstract() {
		if err := addTag(w, "AB", val.Text); err != nil {
			return err
		}
	}
	return nil
}

func addEdition(w io.Writer, rec metadata.Rec) error {
	return addTag(w, "ET", rec.Edition())
}

func addPublisher(w io.Writer, rec metadata.Rec) error {
	return addTag(w, "PB", rec.Publisher())
}

func addPlaceOfPublication(w io.Writer, rec metadata.Rec) error {
	return addTag(w, "CY", rec.PlaceOfPublication())
}

func addDOI(w io.Writer, rec metadata.Rec) error {
	return addTag(w, "DO", rec.DOI()...)
}

func addISBN(w io.Writer, rec metadata.Rec) error {
	return addTag(w, "SN", rec.ISBN()...)
}

func addTag(w io.Writer, tag string, vals ...string) error {
	for _, val := range vals {
		if val == "" {
			continue
		}
		if _, err := fmt.Fprintf(w, "%s  - %s\r\n", tag, val); err != nil {
			return err
		}
	}
	return nil
}
