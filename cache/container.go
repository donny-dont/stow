package cache

import (
	"errors"
	"io"

	"github.com/graymeta/stow"
)

type container struct {
	overlay stow.Container
	base    stow.Container
}

// ID gets a unique string describing this Container.
func (c *container) ID() string {
	return c.base.ID()
}

// Name gets a human-readable name describing this Container.
func (c *container) Name() string {
	return c.base.Name()
}

// Item gets an item by its ID.
func (c *container) Item(id string) (stow.Item, error) {
	it, err := c.base.Item(id)
	if err != nil {
		return nil, err
	}

	return &item{overlay: c.overlay, base: it}, nil
}

// Items gets a page of items with the specified
// prefix for this Container.
// The specified cursor is a pointer to the start of
// the items to get. It it obtained from a previous
// call to this method, or should be CursorStart for the
// first page.
// count is the number of items to return per page.
// The returned cursor can be checked with IsCursorEnd to
// decide if there are any more items or not.
func (c *container) Items(prefix, cursor string, count int) ([]stow.Item, string, error) {
	its, startAfter, err := c.base.Items(prefix, cursor, count)
	if err != nil {
		return nil, "", err
	}

	wrapped := make([]stow.Item, len(its))
	for i, val := range its {
		wrapped[i] = &item{overlay: c.overlay, base: val}
	}

	return wrapped, startAfter, nil
}

// RemoveItem removes the Item with the specified ID.
func (c *container) RemoveItem(id string) error {
	return errors.New("readonly")
}

// Put creates a new Item with the specified name, and contents
// read from the reader.
func (c *container) Put(name string, r io.Reader, size int64, metadata map[string]interface{}) (stow.Item, error) {
	return nil, errors.New("readonly")
}
