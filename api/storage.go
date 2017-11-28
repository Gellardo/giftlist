package api

import (
	"errors"
)

var store Storage

type Storage interface {
	GetList(id string) (*List, error)
	StoreList(l *List) error
}

type easyStore struct {
	lists map[string]*List
}

func (s *easyStore) GetList(id string) (*List, error) {
	l := s.lists[id]
	if l == nil {
		return nil, errors.New("not found")
	}
	return l, nil
}
func (s *easyStore) StoreList(l *List) error {
	s.lists[l.Id] = l
	return nil
}
