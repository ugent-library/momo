package oaipmh

import (
	"encoding/xml"
	"log"
	"net/http"
)

const header = `
<OAI-PMH xmlns="http://www.openarchives.org/OAI/2.0/"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/ http://www.openarchives.org/OAI/2.0/OAI-PMH.xsd">
`
const footer = `
</OAI-PMH>
`

var (
	errBadVerb = err{Code: "badVerb", Msg: "illegal OAI verb"}
)

type err struct {
	XMLName xml.Name `xml:"error"`
	Code    string   `xml:"code,attr"`
	Msg     string   `xml:",chardata"`
}

type provider struct{}

func NewProvider() http.Handler {
	return &provider{}
}

func (p *provider) identify(w http.ResponseWriter, r *http.Request) {
	p.render(200, errBadVerb, w, r)
}

func (p *provider) listMetadataFormats(w http.ResponseWriter, r *http.Request) {
	p.render(200, errBadVerb, w, r)
}

func (p *provider) listSets(w http.ResponseWriter, r *http.Request) {
	p.render(200, errBadVerb, w, r)
}

func (p *provider) listIdentifiers(w http.ResponseWriter, r *http.Request) {
	p.render(200, errBadVerb, w, r)
}

func (p *provider) listRecords(w http.ResponseWriter, r *http.Request) {
	p.render(200, errBadVerb, w, r)
}

func (p *provider) getRecord(w http.ResponseWriter, r *http.Request) {
	p.render(200, errBadVerb, w, r)
}

func (p *provider) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	verb := r.URL.Query().Get("verb")

	switch verb {
	case "Identify":
		p.identify(w, r)
	case "ListMetadataFormats":
		p.listMetadataFormats(w, r)
	case "ListSets":
		p.listSets(w, r)
	case "ListIdentifiers":
		p.listIdentifiers(w, r)
	case "ListRecords":
		p.listRecords(w, r)
	case "GetRecord":
		p.getRecord(w, r)
	default:
		p.render(200, errBadVerb, w, r)
	}
}

func (p *provider) render(status int, res interface{}, w http.ResponseWriter, r *http.Request) {
	body, err := xml.Marshal(res)
	if err != nil {
		log.Panic(err)
	}

	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(status)
	w.Write([]byte(xml.Header))
	w.Write([]byte(header))
	w.Write(body)
	w.Write([]byte(footer))
}
