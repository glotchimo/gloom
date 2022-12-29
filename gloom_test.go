package gloom

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"io"
	"os"
	"testing"
)

func TestGloomPut(t *testing.T) {
	// Set up Gloom w/ hashers
	var g Gloom
	g.Add(sha1.New().Sum)
	g.Add(sha256.New().Sum)
	g.Add(md5.New().Sum)

	// Store some data
	if err := g.Put([]byte("foo"), []byte("foo")); err != nil {
		t.Errorf("error putting data: %s", err.Error())
	}

	// Make sure the data exists on disk
	f, err := os.Open("foo")
	if err != nil {
		t.Errorf("error opening file: %s", err.Error())
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		t.Errorf("error reading data: %s", err.Error())
	}

	if string(data) != "foo" {
		t.Errorf("invalid data: %s", data)
	}

	// Clean up
	if err := os.Remove("foo"); err != nil {
		t.Errorf("error removing file: %s", err.Error())
	}
}

func TestGloomGet(t *testing.T) {
	// Set up Gloom w/ hashers
	var g Gloom
	g.Add(sha1.New().Sum)
	g.Add(sha256.New().Sum)
	g.Add(md5.New().Sum)

	// Store some data
	if err := g.Put([]byte("bd9JVE4Z"), []byte("bd9JVE4Z")); err != nil {
		t.Errorf("error putting data: %s", err.Error())
	}

	// Do valid check
	data, err := g.Get([]byte("bd9JVE4Z"))
	if err != nil {
		t.Errorf("error getting data: %s", err.Error())
	}

	if string(data) != "bd9JVE4Z" {
		t.Errorf("invalid data: %s", data)
	}

	// Do an invalid check
	if _, err := g.Get([]byte("dCbWTITq")); err != nil {
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("error getting data: %s", err.Error())
		}
	}

	// Clean up
	if err := os.Remove("bd9JVE4Z"); err != nil {
		t.Errorf("error removing file: %s", err.Error())
	}
}
