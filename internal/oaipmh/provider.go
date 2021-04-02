package oaipmh

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

const (
	xmlnsXsi          = "http://www.w3.org/2001/XMLSchema-instance"
	xsiSchemaLocation = "http://www.openarchives.org/OAI/2.0/ http://www.openarchives.org/OAI/2.0/OAI-PMH.xsd"
)

var (
	ErrVerbMissing              = Error{Code: "badVerb", Value: "Verb is missing"}
	ErrVerbRepeated             = Error{Code: "badVerb", Value: "Verb can't be repeated"}
	ErrNoSetHierarchy           = Error{Code: "noSetHierarchy", Value: "Sets are not supported"}
	ErrIDDoesNotExist           = Error{Code: "idDoesNotExist", Value: "Identifier is unknown or illegal"}
	ErrNoRecordsMatch           = Error{Code: "noRecordsMatch", Value: "No records match"}
	ErrResumptiontokenRepeated  = Error{Code: "badArgument", Value: "Argument 'resumptionToken' can't be repeated"}
	ErrResumptiontokenExclusive = Error{Code: "badArgument", Value: "resumptionToken cannot be combined with other attributes"}
	ErrMetadataPrefixMissing    = Error{Code: "badArgument", Value: "Argument 'metadataPrefix' is missing"}
	ErrIdentifierMissing        = Error{Code: "badArgument", Value: "Argument 'identifier' is missing"}

	OAIDC = MetadataFormat{
		MetadataPrefix:    "oai_dc",
		Schema:            "http://www.openarchives.org/OAI/2.0/oai_dc.xsd",
		MetadataNamespace: "http://www.openarchives.org/OAI/2.0/oai_dc/",
	}

	verbRe = regexp.MustCompile("^Identify|ListMetadataFormats|ListSets|ListIdentifiers|ListRecords|GetRecord$")
)

type Request struct {
	XMLName         xml.Name `xml:"request"`
	URL             string   `xml:",chardata"`
	Verb            string   `xml:"verb,attr,omitempty"`
	MetadataPrefix  string   `xml:"metadataPrefix,attr,omitempty"`
	Identifier      string   `xml:"identifier,attr,omitempty"`
	Set             string   `xml:"set,attr,omitempty"`
	From            string   `xml:"from,attr,omitempty"`
	Until           string   `xml:"until,attr,omitempty"`
	ResumptionToken string   `xml:"resumptionToken,attr,omitempty"`
}

type response struct {
	XMLName           xml.Name `xml:"http://www.openarchives.org/OAI/2.0/ OAI-PMH"`
	XmlnsXsi          string   `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation string   `xml:"xsi:schemaLocation,attr"`
	ResponseDate      string   `xml:"responseDate"`
	Request           Request
	Errors            []Error
	Body              interface{}
}

type Error struct {
	XMLName xml.Name `xml:"error"`
	Code    string   `xml:"code,attr"`
	Value   string   `xml:",chardata"`
}

type Identify struct {
	XMLName           xml.Name `xml:"Identify"`
	RepositoryName    string   `xml:"repositoryName"`
	BaseURL           string   `xml:"baseURL"`
	ProtocolVersion   string   `xml:"protocolVersion"`
	AdminEmail        []string `xml:"adminEmail"`
	Granularity       string   `xml:"granularity"`
	EarliestDatestamp string   `xml:"earliestDatestamp"`
	DeletedRecord     string   `xml:"deletedRecord"`
}

type ListMetadataFormats struct {
	XMLName         xml.Name         `xml:"ListMetadataFormats"`
	MetadataFormats []MetadataFormat `xml:"metadataFormat"`
}

type ListSets struct {
	XMLName xml.Name `xml:"ListSets"`
	Sets    []Set    `xml:"set"`
}

type GetRecord struct {
	XMLName xml.Name `xml:"GetRecord"`
	Record  *Record  `xml:"record"`
}

type ListRecords struct {
	XMLName         xml.Name         `xml:"ListRecords"`
	Records         []*Record        `xml:"records"`
	ResumptionToken *ResumptionToken `xml:"resumptionToken"`
}

type MetadataFormat struct {
	MetadataPrefix    string `xml:"metadataPrefix"`
	Schema            string `xml:"schema"`
	MetadataNamespace string `xml:"metadataNamespace"`
}

type Set struct {
	SetSpec string `xml:"setSpec"`
	SetName string `xml:"setName"`
}

type Header struct {
	Status     string   `xml:"status,attr,omitempty"`
	Identifier string   `xml:"identifier"`
	Datestamp  string   `xml:"datestamp"`
	SetSpec    []string `xml:"setSpec"`
}

type Metadata struct {
	XML []byte `xml:",innerxml"`
}

type Record struct {
	Header   Header   `xml:"header"`
	Metadata Metadata `xml:"metadata"`
}

type ResumptionToken struct {
	ExpirationDate   string `xml:"expirationDate,attr,omitempty"`
	CompleteListSize int    `xml:"completeListSize,attr,omitempty"`
	Cursor           int    `xml:"cursor,attr,omitempty"`
	Value            string `xml:",chardata"`
}

type provider struct {
	ProviderOptions
}

type ProviderOptions struct {
	RepositoryName    string
	BaseURL           string
	AdminEmail        []string
	Granularity       string
	EarliestDatestamp string
	DeletedRecord     string
	MetadataFormats   []MetadataFormat
	Sets              []Set
	GetRecord         func(string, string) *Record
	ListRecords       func(*Request) ([]*Record, *ResumptionToken)
}

func NewProvider(opts ProviderOptions) http.Handler {
	p := &provider{
		ProviderOptions: opts,
	}

	if p.Granularity == "" {
		p.Granularity = "YYYY-MM-DDThh:mm:ssZ"
	}
	if p.DeletedRecord == "" {
		p.DeletedRecord = "persistent"
	}

	return p
}

// TODO description, compression
func (p *provider) identify(r *response) {
	r.Body = &Identify{
		RepositoryName:    p.RepositoryName,
		BaseURL:           p.BaseURL,
		ProtocolVersion:   "2.0",
		AdminEmail:        p.AdminEmail,
		Granularity:       p.Granularity,
		EarliestDatestamp: p.EarliestDatestamp,
		DeletedRecord:     p.DeletedRecord,
	}
}

// TODO identifier, idDoesNotExist, noMetadataFormats
func (p *provider) listMetadataFormats(r *response) {
	r.Body = &ListMetadataFormats{
		MetadataFormats: p.MetadataFormats,
	}
}

// TODO resumptionToken, badResumptionToken
func (p *provider) listSets(r *response) {
	if len(p.Sets) == 0 {
		r.Errors = append(r.Errors, ErrNoSetHierarchy)
		return
	}
	r.Body = &ListSets{
		Sets: p.Sets,
	}
}

func (p *provider) listIdentifiers(r *response) {
}

// TODO badResumptionToken, cannotDisseminateFormat, noSetHierarchy
func (p *provider) listRecords(r *response) {
	recs, token := p.ListRecords(&r.Request)
	if len(recs) == 0 {
		r.Errors = append(r.Errors, ErrNoRecordsMatch)
		return
	}
	r.Body = &ListRecords{
		Records:         recs,
		ResumptionToken: token,
	}
}

// TODO cannotDisseminateFormat
func (p *provider) getRecord(r *response) {
	// TODO also return error
	rec := p.GetRecord(r.Request.Identifier, r.Request.MetadataPrefix)
	if rec == nil {
		r.Errors = append(r.Errors, ErrIDDoesNotExist)
		return
	}
	r.Body = &GetRecord{
		Record: rec,
	}
}

func (p *provider) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := &response{
		XmlnsXsi:          xmlnsXsi,
		XsiSchemaLocation: xsiSchemaLocation,
		ResponseDate:      time.Now().UTC().Format(time.RFC3339),
		Request: Request{
			URL: p.BaseURL + r.URL.Path,
		},
	}

	q := r.URL.Query()

	res.setVerb(q)

	switch res.Request.Verb {
	case "Identify":
		res.allowAttrs(q)
		if len(res.Errors) == 0 {
			p.identify(res)
		}
	case "ListMetadataFormats":
		res.allowAttrs(q, "metadataPrefix", "identifier")
		if len(res.Errors) == 0 {
			p.listMetadataFormats(res)
		}
	case "ListSets":
		res.allowResumptionToken(q)
		if len(res.Errors) == 0 {
			p.listSets(res)
		}
	case "ListIdentifiers":
		res.allowResumptionToken(q)
		res.allowAttrs(q, "metadataPrefix", "from", "until", "set")
		res.requireMetadataPrefix()
		if len(res.Errors) == 0 {
			p.listIdentifiers(res)
		}
	case "ListRecords":
		res.allowResumptionToken(q)
		res.allowAttrs(q, "metadataPrefix", "from", "until", "set")
		res.requireMetadataPrefix()
		if len(res.Errors) == 0 {
			p.listRecords(res)
		}
	case "GetRecord":
		res.allowAttrs(q, "metadataPrefix", "identifier")
		res.requireMetadataPrefix()
		res.requireIdentifier()
		if len(res.Errors) == 0 {
			p.getRecord(res)
		}
	}

	res.render(200, w)
}

func (r *response) render(status int, w http.ResponseWriter) {
	out, err := xml.MarshalIndent(r, "", " ")
	if err != nil {
		log.Panic(err)
	}
	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(status)
	w.Write([]byte(xml.Header))
	w.Write(out)
}

func (r *response) setVerb(q url.Values) {
	vals, found := q["verb"]

	if !found {
		r.Errors = append(r.Errors, ErrVerbMissing)
		return
	}

	if len(vals) > 1 {
		r.Errors = append(r.Errors, ErrVerbRepeated)
		return
	}

	if !verbRe.MatchString(vals[0]) {
		r.Errors = append(r.Errors, Error{Code: "badVerb", Value: fmt.Sprintf("Verb '%s' is illegal", vals[0])})
		return
	}

	r.Request.Verb = vals[0]
}

func (r *response) allowResumptionToken(q url.Values) {
	vals, found := q["resumptionToken"]
	if !found {
		return
	}
	if len(vals) > 1 {
		r.Errors = append(r.Errors, ErrResumptiontokenRepeated)
		return
	}
	r.Request.ResumptionToken = vals[0]
}

func (r *response) allowAttrs(q url.Values, attrs ...string) {
	// resumptionToken is exclusive
	if r.Request.ResumptionToken != "" && len(q) > 2 {
		r.Errors = append(r.Errors, ErrResumptiontokenExclusive)
		return
	}

	for attr, vals := range q {
		if attr == "verb" {
			continue
		}

		var validAttr bool
		for _, a := range attrs {
			if attr == a {
				validAttr = true
				break
			}
		}

		if !validAttr {
			err := Error{Code: "badArgument", Value: fmt.Sprintf("Argument '%s' is illegal", attr)}
			r.Errors = append(r.Errors, err)
			continue
		}
		if len(vals) > 1 {
			err := Error{Code: "badArgument", Value: fmt.Sprintf("Argument '%s' can't be repeated", attr)}
			r.Errors = append(r.Errors, err)
			continue
		}

		val := vals[0]

		switch attr { // TODO checks
		case "metadataPrefix":
			r.Request.MetadataPrefix = val
		case "identifier":
			r.Request.Identifier = val
		case "from":
			r.Request.From = val
		case "until":
			r.Request.Until = val
		case "set":
			r.Request.Set = val
		}
	}
}

func (r *response) requireMetadataPrefix() {
	if r.Request.ResumptionToken == "" && r.Request.MetadataPrefix == "" {
		r.Errors = append(r.Errors, ErrMetadataPrefixMissing)
	}
}

func (r *response) requireIdentifier() {
	if r.Request.ResumptionToken == "" && r.Request.Identifier == "" {
		r.Errors = append(r.Errors, ErrIdentifierMissing)
	}
}
