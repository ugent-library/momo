package engine

type Storage interface {
	GetRec(string, string) (*Rec, error)
	GetAllRecs() RecCursor
	AddRec(*Rec) error
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
