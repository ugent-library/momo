package oaidc

import (
	"encoding/xml"
	"io"

	"github.com/ugent-library/momo/internal/engine"
)

const startTag = `
<oai_dc:dc xmlns="http://www.openarchives.org/OAI/2.0/oai_dc/"
xmlns:oai_dc="http://www.openarchives.org/OAI/2.0/oai_dc/"
xmlns:dc="http://purl.org/dc/elements/1.1/"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/oai_dc/ http://www.openarchives.org/OAI/2.0/oai_dc.xsd">
`

const endTag = `</oai_dc:dc>`

type encoder struct {
	writer io.Writer
}

type visitor func(io.Writer, *engine.Rec) error

var visitors = []visitor{
	addID,
	addTitle,
	addAuthor,
	addAbstract,
	addPublisher,
	addDOI,
	addISBN,
}

func NewEncoder(w io.Writer) engine.RecEncoder {
	return &encoder{w}
}

func (e *encoder) Encode(rec *engine.Rec) (err error) {
	w := e.writer

	if _, err = io.WriteString(w, startTag); err != nil {
		return
	}

	for _, v := range visitors {
		if err = v(w, rec); err != nil {
			return
		}
	}

	if _, err = io.WriteString(w, endTag); err != nil {
		return
	}

	return
}

func addID(w io.Writer, rec *engine.Rec) error {
	return addTag(w, "identifier", rec.ID)
}

func addTitle(w io.Writer, rec *engine.Rec) error {
	if v := rec.Metadata().Title; v != "" {
		return addTag(w, "title", v)
	}
	return nil
}

func addAuthor(w io.Writer, rec *engine.Rec) error {
	for _, v := range rec.Metadata().Author {
		if err := addTag(w, "contributor", v.Name); err != nil {
			return err
		}
	}
	return nil
}

func addAbstract(w io.Writer, rec *engine.Rec) error {
	for _, v := range rec.Metadata().Abstract {
		if err := addTag(w, "description", v.Text); err != nil {
			return err
		}
	}
	return nil
}

func addPublisher(w io.Writer, rec *engine.Rec) error {
	if v := rec.Metadata().Publisher; v != "" {
		return addTag(w, "publisher", v)
	}
	return nil
}

func addDOI(w io.Writer, rec *engine.Rec) error {
	return addTag(w, "identifier", rec.Metadata().DOI...)
}

func addISBN(w io.Writer, rec *engine.Rec) error {
	for _, v := range rec.Metadata().ISBN {
		if err := addTag(w, "identifier", "ISBN: "+v); err != nil {
			return err
		}
	}
	return nil
}

func addTag(w io.Writer, tag string, vals ...string) error {
	for _, val := range vals {
		// TODO catch errors
		io.WriteString(w, `<dc:`)
		io.WriteString(w, tag)
		io.WriteString(w, `>`)
		xml.EscapeText(w, []byte(val))
		io.WriteString(w, `</dc:`)
		io.WriteString(w, tag)
		io.WriteString(w, ">\n")
	}
	return nil
}
