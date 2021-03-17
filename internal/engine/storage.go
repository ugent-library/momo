package engine

type Storage interface {
	GetRec(string) (*Rec, error)
	GetAllRecs(chan<- *Rec) error
	AddRec(*Rec) error
	Reset() error
}

type SearchStorage interface {
	SearchRecs(SearchArgs) (*RecHits, error)
	AddRecs(<-chan *Rec)
	CreateRecIndex() error
	DeleteRecIndex() error
	Reset() error
}
