package update

import (
	"github.com/google/uuid"
	"github.com/size12/planning-poker/internal/entity/voting"
)

type Type string

const (
	Connected Type = "CONNECTED"

	Vote           Type = "VOTE"
	Reveal         Type = "REVEAL"
	StartNewVoting Type = "NEW_VOTING"

	AlreadyConnected Type = "ALREADY_CONNECTED"
)

type Update struct {
	Type     Type         `json:"type"`
	PlayerID uuid.UUID    `json:"player_id"`
	Vote     *voting.Vote `json:"vote"`
}
