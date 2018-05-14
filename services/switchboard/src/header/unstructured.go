package header

import (
	"encoding/hex"
)

type UnstructuredConf struct {
	Size uint
}

type Unstructured []byte

func NewUnstructured(conf UnstructuredConf) *Unstructured {
	slice := make(Unstructured, conf.Size)
	return &slice
}

func (header *Unstructured) Id() []byte {
	return *header
}

func (header *Unstructured) String() string {
	return hex.EncodeToString(header.Id())
}
