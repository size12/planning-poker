package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/goombaio/namegenerator"
	"github.com/size12/planning-poker/internal/entity/access"
	"github.com/size12/planning-poker/internal/entity/voting"
)

var nameGenerator namegenerator.Generator

func init() {
	seed := time.Now().UTC().UnixNano()
	nameGenerator = namegenerator.NewNameGenerator(seed)
}

type Player struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	AccessLevel access.Access `json:"access_level"`
	Voted       *voting.Vote  `json:"voted"`
}

func NewPlayer() (*Player, error) {
	userID, err := uuid.NewUUID()

	if err != nil {
		return nil, err
	}

	name := nameGenerator.Generate()

	return &Player{
		ID:          userID,
		Name:        name,
		AccessLevel: access.User,
	}, nil
}

func (p *Player) SetAccessLevel(level access.Access) error {
	p.AccessLevel = level
	return nil
}

func (p *Player) Vote(vote *voting.Vote) error {
	p.Voted = vote
	return nil
}

func (p *Player) ClearVote() error {
	p.Voted = nil
	return nil
}
