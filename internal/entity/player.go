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
	ID     uuid.UUID     `json:"id"`
	Name   string        `json:"name"`
	Voted  *voting.Vote  `json:"voted"`
	Status voting.Status `json:"status"`
}

func NewPlayer(id uuid.UUID) (*Player, error) {
	name := nameGenerator.Generate()

	return &Player{
		ID:     id,
		Name:   name,
		Status: voting.PlayerObserving,
	}, nil
}

func (p *Player) SetVotingStatus(status voting.Status) {
	p.Status = status
}

func (p *Player) Vote(vote *voting.Vote) {
	p.Voted = vote
}

func (p *Player) ClearVote() error {
	p.Voted = nil
	return nil
}

func (p *Player) SetName(name string) error {
	p.Name = name
	return nil
}
