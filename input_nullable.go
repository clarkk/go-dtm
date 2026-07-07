package dtm

import (
	"bytes"
	"encoding/json/v2"
)

var bytes_null = []byte("null")

type Nullable[T any] struct {
	Value **T
}

func NewNullable[T any](p *T) Nullable[T] {
	return Nullable[T]{
		Value: &p,
	}
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	n.Value = new(*T)
	
	if bytes.Equal(data, bytes_null) {
		return nil
	}
	
	var inner T
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}
	
	*n.Value = &inner
	return nil
}

func (n Nullable[T]) Assign(dest *T) bool {
	if n.Value == nil || *n.Value == nil || dest == nil {
		return false
	}
	*dest = **n.Value
	return true
}