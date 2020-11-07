package cache

import (
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/graymeta/stow"
)

type item struct {
	overlay stow.Container
	base    stow.Item
}

// ID gets a unique string describing this Item.
func (i *item) ID() string {
	return i.base.ID()
}

// Name gets a human-readable name describing this Item.
func (i *item) Name() string {
	return i.base.Name()
}

// URL gets a URL for this item.
// For example:
// local: file:///path/to/something
// azure: azure://host:port/api/something
//    s3: s3://host:post/etc
func (i *item) URL() *url.URL {
	return i.base.URL()
}

// Size gets the size of the Item's contents in bytes.
func (i *item) Size() (int64, error) {
	return i.base.Size()
}

// Open opens the Item for reading.
// Calling code must close the io.ReadCloser.
func (i *item) Open() (io.ReadCloser, error) {
	name := i.base.Name()
	size, _ := i.base.Size()
	overlayItem, err := i.overlay.Item(name)

	if err == nil {
		fmt.Printf("Getting from overlay\n")
		return overlayItem.Open()
	}

	fmt.Printf("Getting from remote\n")

	r, err := i.base.Open()
	if err != nil {
		return nil, err
	}

	oi, err := i.overlay.Put(name, r, size, nil)
	if err != nil {
		return nil, err
	}

	/*
		pr, pw := io.Pipe()
		tr := newTeeReaderCloser(r, pw)

		go func() {
			defer pw.Close()
			fmt.Printf("Writing to overlay\n")
			oi, err := i.overlay.Put(name, pr, size, nil)
			if err != nil {
				fmt.Printf("Err??? %s", err)
			} else {
				s, _ := oi.Size()
				fmt.Printf("Name: %s Size: %d", oi.Name(), s)
			}
		}()
	*/

	return oi.Open()
}

// ETag is a string that is different when the Item is
// different, and the same when the item is the same.
// Usually this is the last modified datetime.
func (i *item) ETag() (string, error) {
	return i.base.ETag()
}

// LastMod returns the last modified date of the file.
func (i *item) LastMod() (time.Time, error) {
	return i.base.LastMod()
}

// Metadata gets a map of key/values that belong
// to this Item.
func (i *item) Metadata() (map[string]interface{}, error) {
	return i.base.Metadata()
}
