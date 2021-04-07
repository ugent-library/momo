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

func addEndOfRecord(w io.Writer, rec *engine.Rec) error {
	_, err := io.WriteString(w, "ER  - \r\n")
	return err
}

func addID(w io.Writer, rec *engine.Rec) error {
	return addTag(w, "ID", rec.ID)
}

func addTitle(w io.Writer, rec *engine.Rec) error {
	return addTag(w, "TI", rec.GetString("title"))
}

func addAuthor(w io.Writer, rec *engine.Rec) error {
	vals := rec.GetStringSlice("author[*].name")
	return addTag(w, "AU", vals...)
}

func addAbstract(w io.Writer, rec *engine.Rec) error {
	vals := rec.GetStringSlice("abstract[*].text")
	return addTag(w, "AB", vals...)
}

func addEdition(w io.Writer, rec *engine.Rec) error {
	return addTag(w, "ET", rec.GetString("edition"))
}

func addPublisher(w io.Writer, rec *engine.Rec) error {
	return addTag(w, "PB", rec.GetString("publisher"))
}

func addPlaceOfPublication(w io.Writer, rec *engine.Rec) error {
	return addTag(w, "CY", rec.GetString("placeOfPublication"))
}

func addDOI(w io.Writer, rec *engine.Rec) error {
	vals := rec.GetStringSlice("doi")
	return addTag(w, "DO", vals...)
}

func addISBN(w io.Writer, rec *engine.Rec) error {
	vals := rec.GetStringSlice("isbn")
	return addTag(w, "SN", vals...)
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
