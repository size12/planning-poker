package handlers

import (
	"sync"
)

type Handlers struct {
	pool *sync.Map
}

func NewHandlers() (*Handlers, error) {
	pool := &sync.Map{}

	return &Handlers{pool: pool}, nil
}
