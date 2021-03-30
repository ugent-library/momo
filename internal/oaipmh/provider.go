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
	errBadVerb        = oaiError{Code: "badVerb", Msg: "Illegal OAI verb"}
	errNoSetHierarchy = oaiError{Code: "noSetHierarchy", Msg: "Sets are not supported"}

	OAIDC = MetadataFormat{
		MetadataPrefix:    "oai_dc",
		Schema:            "http://www.openarchives.org/OAI/2.0/oai_dc.xsd",
		MetadataNamespace: "http://www.openarchives.org/OAI/2.0/oai_dc/",
	}
)

type oaiRequest struct {
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

type oaiResponse struct {
	XMLName           xml.Name `xml:"http://www.openarchives.org/OAI/2.0/ OAI-PMH"`
	XmlnsXsi          string   `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation string   `xml:"xsi:schemaLocation,attr"`
	ResponseDate      string   `xml:"responseDate"`
	Request           oaiRequest
	Body              interface{}
}

type oaiError struct {
	XMLName xml.Name `xml:"error"`
	Code    string   `xml:"code,attr"`
	Msg     string   `xml:",chardata"`
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

type MetadataFormat struct {
	MetadataPrefix    string `xml:"metadataPrefix"`
	Schema            string `xml:"schema"`
	MetadataNamespace string `xml:"metadataNamespace"`
}

type Set struct {
	SetSpec string `xml:"setSpec"`
	SetName string `xml:"setName"`
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
func (p *provider) identify(r *http.Request) *Identify {
	return &Identify{
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
func (p *provider) listMetadataFormats(r *http.Request) interface{} {
	return &ListMetadataFormats{
		MetadataFormats: p.MetadataFormats,
	}
}

// TODO resumptionToken, badArgument, badResumptionToken
func (p *provider) listSets(r *http.Request) interface{} {
	if len(p.Sets) == 0 {
		return errNoSetHierarchy
	}
	return &ListSets{
		Sets: p.Sets,
	}
}

func (p *provider) listIdentifiers(r *http.Request) interface{} {
	return errBadVerb
}

func (p *provider) listRecords(r *http.Request) interface{} {
	return errBadVerb
}

func (p *provider) getRecord(r *http.Request) interface{} {
	return errBadVerb
}

func (p *provider) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := oaiRequest{URL: p.BaseURL + r.URL.Path}
	err := p.formDecoder.Decode(&req, r.URL.Query())
	if err != nil {
		log.Panic(err)
	}

	res := oaiResponse{
		XmlnsXsi:          xmlnsXsi,
		XsiSchemaLocation: xsiSchemaLocation,
		ResponseDate:      time.Now().UTC().Format(time.RFC3339),
		Request:           req,
	}

	switch req.Verb {
	case "Identify":
		res.Body = p.identify(r)
	case "ListMetadataFormats":
		res.Body = p.listMetadataFormats(r)
	case "ListSets":
		res.Body = p.listSets(r)
	case "ListIdentifiers":
		res.Body = p.listIdentifiers(r)
	case "ListRecords":
		res.Body = p.listRecords(r)
	case "GetRecord":
		res.Body = p.getRecord(r)
	default:
		res.Body = errBadVerb
	}

	res.render(200, w)
}

func (r oaiResponse) render(status int, w http.ResponseWriter) {
	out, err := xml.MarshalIndent(r, "", " ")
	if err != nil {
		log.Panic(err)
	}
	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(status)
	w.Write([]byte(xml.Header))
	w.Write(out)
}
