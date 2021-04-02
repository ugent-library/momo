package engine

type Storage interface {
	GetRec(string, string) (*Rec, error)
	AllRecs() RecCursor
	AddRec(*Rec) error
	Reset() error
}

type SearchStorage interface {
	SearchRecs(SearchArgs) (*RecHits, error)
	SearchAllRecs(SearchArgs) RecCursor
	AddRecs(<-chan *Rec)
	CreateRecIndex() error
	DeleteRecIndex() error
	Reset() error
}
