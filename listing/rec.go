package listing

type RecHits struct {
	Total int       `json:"total"`
	Hits  []*RecHit `json:"hits"`
}

type RecHit struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}
