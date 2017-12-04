package api

import (
	"errors"
)

// List stores the information of a giftlist.
type List struct {
	ID    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Items []Item `json:"items,omitempty"`
}

// Item represents a single entry in a giftlist.
type Item struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Link     string `json:"link,omitempty"`
	Assigned bool   `json:"assigned,omitempty"`
}

func (l *List) addItem(new Item) error {
	for _, i := range l.Items {
		if new.ID == i.ID {
			return errors.New("Item exists")
		}
	}
	l.Items = append(l.Items, new)
	return nil
}

func (l *List) getItem(itemid string) (*Item, error) {
	for _, i := range l.Items {
		if itemid == i.ID {
			return &i, nil
		}
	}
	return nil, errors.New("Item not found")
}
func (l *List) replaceItem(new *Item) error {
	for index, i := range l.Items {
		if new.ID == i.ID {
			l.Items[index] = *new
			return nil
		}
	}
	return errors.New("No item to replace")
}

func (l *List) deleteItem(itemid string) error {
	for index, i := range l.Items {
		if itemid == i.ID {
			l.Items = append(l.Items[:index], l.Items[index+1:]...)
			return nil
		}
	}
	return errors.New("No item to delete")
}

func (item *Item) merge(new Item) error {
	var def Item
	if item.ID != new.ID && new.ID != def.ID {
		return errors.New("Itemmerge: IDs do not match")
	}
	if item.Name != new.Name && new.Name != def.Name {
		item.Name = new.Name
	}
	if item.Link != new.Link && new.Link != def.Link {
		item.Link = new.Link
	}
	if item.Assigned != new.Assigned && new.Assigned != def.Assigned {
		item.Assigned = new.Assigned
	}
	return nil
}
