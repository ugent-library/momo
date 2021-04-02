package oaipmh

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	xmlnsXsi          = "http://www.w3.org/2001/XMLSchema-instance"
	xsiSchemaLocation = "http://www.openarchives.org/OAI/2.0/ http://www.openarchives.org/OAI/2.0/OAI-PMH.xsd"
)

var (
	OAIDC = MetadataFormat{
		MetadataPrefix:    "oai_dc",
		Schema:            "http://www.openarchives.org/OAI/2.0/oai_dc.xsd",
		MetadataNamespace: "http://www.openarchives.org/OAI/2.0/oai_dc/",
	}

	yes                      struct{}
	verbs                    = map[string]struct{}{"Identify": yes, "ListMetadataFormats": yes, "ListSets": yes, "ListIdentifiers": yes, "ListRecords": yes, "GetRecord": yes}
	identifyAttrs            = map[string]struct{}{"verb": yes}
	listMetadataFormatsAttrs = map[string]struct{}{"verb": yes, "identifier": yes}
	listSetsAttrs            = map[string]struct{}{"verb": yes, "resumptionToken": yes}
	listRecordsAttrs         = map[string]struct{}{"verb": yes, "resumptionToken": yes, "metadataPrefix": yes, "set": yes, "from": yes, "until": yes}
	getRecordAttrs           = map[string]struct{}{"verb": yes, "metadataPrefix": yes, "identifier": yes}

	errVerbMissing              = Error{Code: "badVerb", Value: "Verb is missing"}
	errVerbRepeated             = Error{Code: "badVerb", Value: "Verb can't be repeated"}
	errNoSetHierarchy           = Error{Code: "noSetHierarchy", Value: "Sets are not supported"}
	errIDDoesNotExist           = Error{Code: "idDoesNotExist", Value: "Identifier is unknown or illegal"}
	errNoRecordsMatch           = Error{Code: "noRecordsMatch", Value: "No records match"}
	errResumptiontokenExclusive = Error{Code: "badArgument", Value: "resumptionToken cannot be combined with other attributes"}
	errMetadataPrefixMissing    = Error{Code: "badArgument", Value: "Argument 'metadataPrefix' is missing"}
	errIdentifierMissing        = Error{Code: "badArgument", Value: "Argument 'identifier' is missing"}
	errFromMissing              = Error{Code: "badArgument", Value: "Argument 'from' is missing"}
	errUntilMissing             = Error{Code: "badArgument", Value: "Argument 'until' is missing"}
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
	provider          *provider
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

type ListIdentifiers struct {
	XMLName         xml.Name         `xml:"ListIdentifiers"`
	Headers         []*Header        `xml:"header"`
	ResumptionToken *ResumptionToken `xml:"resumptionToken"`
}

type ListRecords struct {
	XMLName         xml.Name         `xml:"ListRecords"`
	Records         []*Record        `xml:"record"`
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
	GetRecord         func(*Request) *Record
	ListIdentifiers   func(*Request) ([]*Header, *ResumptionToken)
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
		r.Errors = append(r.Errors, errNoSetHierarchy)
		return
	}
	r.Body = &ListSets{
		Sets: p.Sets,
	}
}

// TODO badResumptionToken, cannotDisseminateFormat
func (p *provider) listIdentifiers(r *response) {
	headers, token := p.ListIdentifiers(&r.Request)
	if len(headers) == 0 {
		r.Errors = append(r.Errors, errNoRecordsMatch)
		return
	}
	r.Body = &ListIdentifiers{
		Headers:         headers,
		ResumptionToken: token,
	}
}

// TODO badResumptionToken, cannotDisseminateFormat
func (p *provider) listRecords(r *response) {
	recs, token := p.ListRecords(&r.Request)
	if len(recs) == 0 {
		r.Errors = append(r.Errors, errNoRecordsMatch)
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
	rec := p.GetRecord(&r.Request)
	if rec == nil {
		r.Errors = append(r.Errors, errIDDoesNotExist)
		return
	}
	r.Body = &GetRecord{
		Record: rec,
	}
}

func (p *provider) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := &response{
		provider:          p,
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
		if len(res.Errors) == 0 {
			p.identify(res)
		}
	case "ListMetadataFormats":
		res.validateAttrs(q, listMetadataFormatsAttrs)
		res.Request.Identifier = res.getAttr(q, "identifier")

		if len(res.Errors) == 0 {
			p.listMetadataFormats(res)
		}
	case "ListSets":
		res.validateAttrs(q, listSetsAttrs)
		res.Request.ResumptionToken = res.getAttr(q, "resumptionToken")

		if len(res.Errors) == 0 {
			p.listSets(res)
		}
	case "ListIdentifiers":
		res.validateAttrs(q, listRecordsAttrs)
		res.Request.ResumptionToken = res.getAttr(q, "resumptionToken")
		res.setMetadataPrefix(q)
		res.setSet(q)
		res.setFromUntil(q)

		if len(res.Errors) == 0 {
			p.listIdentifiers(res)
		}
	case "ListRecords":
		res.validateAttrs(q, listRecordsAttrs)
		res.Request.ResumptionToken = res.getAttr(q, "resumptionToken")
		res.setMetadataPrefix(q)
		res.setSet(q)
		res.setFromUntil(q)

		if len(res.Errors) == 0 {
			p.listRecords(res)
		}
	case "GetRecord":
		res.validateAttrs(q, getRecordAttrs)
		res.setMetadataPrefix(q)
		res.Request.Identifier = res.getAttr(q, "identifier")
		if res.Request.Identifier == "" {
			res.Errors = append(res.Errors, errIdentifierMissing)
		}

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

func (r *response) validateAttrs(q url.Values, attrs map[string]struct{}) {
	for attr := range q {
		if _, ok := attrs[attr]; !ok {
			r.Errors = append(r.Errors, Error{Code: "badArgument", Value: fmt.Sprintf("Attribute '%s' is illegal", attr)})
		}
	}
}

func (r *response) setVerb(q url.Values) {
	vals, ok := q["verb"]

	if !ok {
		r.Errors = append(r.Errors, errVerbMissing)
		return
	}
	if len(vals) > 1 {
		r.Errors = append(r.Errors, errVerbRepeated)
		return
	}
	if _, ok := verbs[vals[0]]; !ok {
		r.Errors = append(r.Errors, Error{Code: "badVerb", Value: fmt.Sprintf("Verb '%s' is illegal", vals[0])})
		return
	}

	r.Request.Verb = vals[0]
}

func (r *response) setMetadataPrefix(q url.Values) {
	val := r.getAttr(q, "metadataPrefix")

	if val != "" && r.Request.ResumptionToken != "" {
		r.Errors = append(r.Errors, errResumptiontokenExclusive)
		return
	}

	if val == "" && r.Request.ResumptionToken == "" {
		r.Errors = append(r.Errors, errMetadataPrefixMissing)
		return
	}

	var valid bool
	for _, fmt := range r.provider.MetadataFormats {
		if val == fmt.MetadataPrefix {
			valid = true
			break
		}
	}
	if !valid {
		err := Error{Code: "cannotDisseminateFormat", Value: fmt.Sprintf("Metadata format '%s' is not supported", val)}
		r.Errors = append(r.Errors, err)
		return
	}

	r.Request.MetadataPrefix = val
}

func (r *response) setSet(q url.Values) {
	val := r.getAttr(q, "set")

	if val != "" && r.Request.ResumptionToken != "" {
		r.Errors = append(r.Errors, errResumptiontokenExclusive)
		return
	}

	if val != "" && len(r.provider.Sets) == 0 {
		r.Errors = append(r.Errors, errNoSetHierarchy)
		return
	}

	// TODO check set exists
	r.Request.Set = val
}

func (r *response) setFromUntil(q url.Values) {
	f := r.getAttr(q, "from")
	u := r.getAttr(q, "until")
	if f != "" && u != "" {
		// TODO validate format
		r.Request.From = f
		r.Request.Until = u
	} else if f == "" && u != "" {
		r.Errors = append(r.Errors, errFromMissing)
	} else if f != "" && u == "" {
		r.Errors = append(r.Errors, errUntilMissing)
	}
}

func (r *response) getAttr(q url.Values, attr string) string {
	vals, ok := q[attr]
	if !ok {
		return ""
	}

	if len(vals) > 1 {
		err := Error{Code: "badArgument", Value: fmt.Sprintf("Argument '%s' can't be repeated", attr)}
		r.Errors = append(r.Errors, err)
		return ""
	}

	if vals[0] == "" {
		err := Error{Code: "badArgument", Value: fmt.Sprintf("Argument '%s' can't be empty", attr)}
		r.Errors = append(r.Errors, err)
		return ""
	}

	return vals[0]
}
