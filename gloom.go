package gloom

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
)

const FILTER_SIZE = 32

var (
	ErrNotFound = errors.New("nothing found for that key")
)

type Gloom struct {
	filter  big.Int
	hashers []func([]byte) int
}

// Add a hasher
func (g *Gloom) Add(f func([]byte) []byte) {
	g.hashers = append(g.hashers, func(b []byte) int {
		h := f(b)
		n := binary.BigEndian.Uint64(h)
		i := n % FILTER_SIZE
		return int(i)
	})
}

// Put an element by (always) writing to disk and updating the filter
func (g *Gloom) Put(key []byte, value []byte) error {
	f, err := os.Create(string(key))
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(value); err != nil {
		return fmt.Errorf("error writing value to disk: %w", err)
	}

	for _, h := range g.hashers {
		g.filter.SetBit(&g.filter, h(key), 1)
	}

	return nil
}

// Get an element by checking the filter and (maybe) reading from disk
func (g *Gloom) Get(key []byte) ([]byte, error) {
	for _, h := range g.hashers {
		if g.filter.Bit(h(key)) == 0 {
			return nil, ErrNotFound
		}
	}

	f, err := os.Open(string(key))
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	value, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error reading value from disk: %w", err)
	}

	return value, nil
}
