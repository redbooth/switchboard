package header

type EmptyConf struct{}

type Empty struct{}

func NewEmpty(conf EmptyConf) *Empty {
	return &Empty{}
}

func (header *Empty) Id() []byte {
	return []byte{}
}

func (header *Empty) String() string {
	return ""
}
