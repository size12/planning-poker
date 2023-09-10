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
	id          uuid.UUID
	name        string
	accessLevel access.Access

	voted *voting.Vote
}

func NewPlayer() (*Player, error) {
	userID, err := uuid.NewUUID()

	if err != nil {
		return nil, err
	}

	name := nameGenerator.Generate()

	return &Player{
		id:          userID,
		name:        name,
		accessLevel: access.User,
	}, nil
}

func (p *Player) ID() uuid.UUID {
	return p.id
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) AccessLevel() access.Access {
	return p.accessLevel
}

func (p *Player) SetAccessLevel(level access.Access) error {
	p.accessLevel = level
	return nil
}

func (p *Player) Vote(vote *voting.Vote) error {
	p.voted = vote
	return nil
}

func (p *Player) ClearVote() error {
	p.voted = nil
	return nil
}
