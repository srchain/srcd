package srcdb


type Putter interface {
	Put(key []byte, value []byte) error
}

type Deleter interface {
	Delete(key []byte) error
}

type Batch interface {
	Putter
	Deleter
	ValueSize() int
	Write() error
	Reset()

}

type Database interface {
	 Putter
	 Get(key []byte) ([]byte,error)
	 Has(key []byte) (bool,error)
	 Close()
	 NewBatch() Batch
}


