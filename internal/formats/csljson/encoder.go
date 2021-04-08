package csljson

import (
	"encoding/json"
	"io"

	"github.com/ugent-library/momo/internal/engine"
)

type Item struct {
	ID             string   `json:"id"`
	Type           string   `json:"type,omitempty"`
	Title          string   `json:"title,omitempty"`
	Author         []Person `json:"author,omitempty"`
	Edition        string   `json:"edition,omitempty"`
	Issued         Issued   `json:"issued,omitempty"`
	Publisher      string   `json:"publisher,omitempty"`
	PublisherPlace string   `json:"publisher-place,omitempty"`
	DOI            string   `json:"DOI,omitempty"`
	ISBN           string   `json:"ISBN,omitempty"`
}

type Issued struct {
	Raw string `json:"raw,omitempty"`
}

type Person struct {
	Family string `json:"family,omitempty"`
	Given  string `json:"given,omitempty"`
}

type encoder struct {
	json *json.Encoder
}

func NewEncoder(w io.Writer) engine.RecEncoder {
	return &encoder{json.NewEncoder(w)}
}

func (e *encoder) Encode(rec *engine.Rec) error {
	item := Item{
		ID:             rec.ID,
		Title:          rec.GetString("title"),
		Edition:        rec.GetString("edition"),
		Publisher:      rec.GetString("publisher"),
		PublisherPlace: rec.GetString("placeOfPublication"),
	}
	switch rec.Type {
	case "Book":
		item.Type = "book"
	case "JournalArticle":
		item.Type = "article-journal"
	case "Chapter", "BookChapter":
		item.Type = "chapter"
	case "Thesis":
		item.Type = "thesis"
	}
	item.Issued.Raw = rec.GetString("publicationDate")
	for _, v := range rec.GetStringSlice("author[*].name") {
		item.Author = append(item.Author, Person{Family: v})
		break
	}
	for _, v := range rec.GetStringSlice("doi") {
		item.DOI = v
		break
	}
	for _, v := range rec.GetStringSlice("isbn") {
		item.ISBN = v
		break
	}

	return e.json.Encode(&item)
}
