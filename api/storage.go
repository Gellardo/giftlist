package api

import (
	"errors"
)

var store Storage

// Storage describes the used interface for storing and retrieving lists for the API to use.
type Storage interface {
	GetList(id string) (*List, error)
	StoreList(l *List) error
}

// easyStore implements the Storage interface and is mainly intended for mocking.
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
	s.lists[l.ID] = l
	return nil
}

func (s *easyStore) empty() {
	for id, _ := range s.lists {
		delete(s.lists, id)
	}
}
