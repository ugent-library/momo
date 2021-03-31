package oaipmh

import (
	"encoding/xml"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/form/v4"
)

const (
	xmlnsXsi          = "http://www.w3.org/2001/XMLSchema-instance"
	xsiSchemaLocation = "http://www.openarchives.org/OAI/2.0/ http://www.openarchives.org/OAI/2.0/OAI-PMH.xsd"
)

var (
	ErrBadVerb        = Error{Code: "badVerb", Value: "Illegal OAI verb"}
	ErrNoSetHierarchy = Error{Code: "noSetHierarchy", Value: "Sets are not supported"}
	ErrIDDoesNotExist = Error{Code: "idDoesNotExist", Value: "Identifier is unknown or illegal"}
	ErrNoRecordsMatch = Error{Code: "noRecordsMatch", Value: "No records match"}

	OAIDC = MetadataFormat{
		MetadataPrefix:    "oai_dc",
		Schema:            "http://www.openarchives.org/OAI/2.0/oai_dc.xsd",
		MetadataNamespace: "http://www.openarchives.org/OAI/2.0/oai_dc/",
	}
)

type Request struct {
	XMLName         xml.Name `xml:"request"`
	URL             string   `xml:",chardata"`
	Verb            string   `xml:"verb,attr,omitempty" form:"verb"`
	MetadataPrefix  string   `xml:"metadataPrefix,attr,omitempty" form:"metadataPrefix"`
	Identifier      string   `xml:"identifier,attr,omitempty" form:"identifier"`
	Set             string   `xml:"set,attr,omitempty" form:"set"`
	From            string   `xml:"from,attr,omitempty" form:"from"`
	Until           string   `xml:"until,attr,omitempty" form:"until"`
	ResumptionToken string   `xml:"resumptionToken,attr,omitempty" form:"resumptionToken"`
}

type response struct {
	XMLName           xml.Name `xml:"http://www.openarchives.org/OAI/2.0/ OAI-PMH"`
	XmlnsXsi          string   `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation string   `xml:"xsi:schemaLocation,attr"`
	ResponseDate      string   `xml:"responseDate"`
	Request           Request
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
	formDecoder *form.Decoder
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
		formDecoder:     form.NewDecoder(),
	}

	if p.Granularity == "" {
		p.Granularity = "YYYY-MM-DDThh:mm:ssZ"
	}
	if p.DeletedRecord == "" {
		p.DeletedRecord = "persistent"
	}

	return p
}

// TODO badArgument, description, compression
func (p *provider) identify(res *response) {
	res.Body = &Identify{
		RepositoryName:    p.RepositoryName,
		BaseURL:           p.BaseURL,
		ProtocolVersion:   "2.0",
		AdminEmail:        p.AdminEmail,
		Granularity:       p.Granularity,
		EarliestDatestamp: p.EarliestDatestamp,
		DeletedRecord:     p.DeletedRecord,
	}
}

// TODO identifier, badArgument, idDoesNotExist, noMetadataFormats
func (p *provider) listMetadataFormats(res *response) {
	res.Body = &ListMetadataFormats{
		MetadataFormats: p.MetadataFormats,
	}
}

// TODO resumptionToken, badArgument, badResumptionToken
func (p *provider) listSets(res *response) {
	if len(p.Sets) == 0 {
		res.Body = ErrNoSetHierarchy
		return
	}
	res.Body = &ListSets{
		Sets: p.Sets,
	}
}

func (p *provider) listIdentifiers(res *response) {
	res.Body = ErrBadVerb
}

// TODO badArgument, badResumptionToken, cannotDisseminateFormat, noRecordsMatch, noSetHierarchy
func (p *provider) listRecords(res *response) {
	recs, token := p.ListRecords(&res.Request)
	if len(recs) == 0 {
		res.Body = ErrNoRecordsMatch
		return
	}
	res.Body = &ListRecords{
		Records:         recs,
		ResumptionToken: token,
	}
}

// TODO badArgument, cannotDisseminateFormat
func (p *provider) getRecord(res *response) {
	// TODO also return error
	rec := p.GetRecord(res.Request.Identifier, res.Request.MetadataPrefix)
	if rec == nil {
		res.Body = ErrIDDoesNotExist
		return
	}
	res.Body = &GetRecord{
		Record: rec,
	}
}

func (p *provider) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := Request{URL: p.BaseURL + r.URL.Path}
	err := p.formDecoder.Decode(&req, r.URL.Query())
	if err != nil {
		log.Panic(err)
	}

	res := &response{
		XmlnsXsi:          xmlnsXsi,
		XsiSchemaLocation: xsiSchemaLocation,
		ResponseDate:      time.Now().UTC().Format(time.RFC3339),
		Request:           req,
	}

	switch req.Verb {
	case "Identify":
		p.identify(res)
	case "ListMetadataFormats":
		p.listMetadataFormats(res)
	case "ListSets":
		p.listSets(res)
	case "ListIdentifiers":
		p.listIdentifiers(res)
	case "ListRecords":
		p.listRecords(res)
	case "GetRecord":
		p.getRecord(res)
	default:
		res.Body = ErrBadVerb
	}

	res.render(200, w)
}

func (r response) render(status int, w http.ResponseWriter) {
	out, err := xml.MarshalIndent(r, "", " ")
	if err != nil {
		log.Panic(err)
	}
	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(status)
	w.Write([]byte(xml.Header))
	w.Write(out)
}
