package engine

type Storage interface {
	GetRec(string) (*Rec, error)
	EachRec(func(*Rec) bool) error
	AddRec(*Rec) error
	AddRepresentation(string, string, []byte) error
	Reset() error
}

type SearchStorage interface {
	SearchRecs(SearchArgs) (*RecHits, error)
	SearchMoreRecs(string) (*RecHits, error)
	AddRecs(<-chan *Rec)
	CreateRecIndex() error
	DeleteRecIndex() error
	Reset() error
}
