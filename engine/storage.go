package engine

type Storage interface {
	GetRec(string) (*Rec, error)
	GetAllRecs(chan<- *Rec) error
	AddRec(*Rec) error
}

type SearchStorage interface {
	SearchRecs(SearchArgs) (*RecHits, error)
	CreateRecIndex() error
	DeleteRecIndex() error
	AddRecs(<-chan *Rec)
}
