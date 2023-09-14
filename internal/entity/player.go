package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/goombaio/namegenerator"
	"github.com/size12/planning-poker/internal/entity/voting"
)

var nameGenerator namegenerator.Generator

func init() {
	seed := time.Now().UTC().UnixNano()
	nameGenerator = namegenerator.NewNameGenerator(seed)
}

type Player struct {
	ID    uuid.UUID    `json:"id"`
	Name  string       `json:"name"`
	Voted *voting.Vote `json:"voted"`
}

func NewPlayer(id uuid.UUID) (*Player, error) {
	name := nameGenerator.Generate()

	return &Player{
		ID:   id,
		Name: name,
	}, nil
}

func (p *Player) Vote(vote *voting.Vote) {
	p.Voted = vote
}

func (p *Player) ClearVote() error {
	p.Voted = nil
	return nil
}
