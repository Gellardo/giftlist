package main

import (
	"errors"
)

type Storage interface {
	GetList(id string) (*list, error)
	StoreList(l *list) error
}

type easyStore struct {
	lists map[string]*list
}

func (s *easyStore) GetList(id string) (*list, error) {
	l := s.lists[id]
	if l == nil {
		return nil, errors.New("not found")
	}
	return l, nil
}
func (s *easyStore) StoreList(l *list) error {
	s.lists[l.Id] = l
	return nil
}
