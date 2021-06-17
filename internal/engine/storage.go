package engine

type Storage interface {
	GetRec(string) (*Rec, error)
	EachRec(func(*Rec) bool) error
	AddRecBySourceID(*Rec) error
	UpdateRecMetadata(string, map[string]interface{}) (*Rec, error)
	GetRepresentation(string, string) (*Representation, error)
	AddRepresentation(*Representation) error
	Reset() error
}

type SearchStorage interface {
	SearchRecs(SearchArgs) (*RecHits, error)
	SearchMoreRecs(string) (*RecHits, error)
	AddRec(*Rec) error
	AddRecs(<-chan *Rec)
	CreateRecIndex() error
	DeleteRecIndex() error
	Reset() error
}
