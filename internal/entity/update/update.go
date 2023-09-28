package update

import (
	"github.com/google/uuid"
	"github.com/size12/planning-poker/internal/entity/voting"
)

type Type string

const (
	Connected Type = "CONNECTED"

	Vote Type = "VOTE"

	ChangeStatus       Type = "CHANGE_STATUS"
	ChangePlayerStatus Type = "CHANGE_PLAYER_STATUS"
	ChangePlayerName   Type = "CHANGE_PLAYER_NAME"
)

type Update struct {
	Type     Type               `json:"type"`
	PlayerID uuid.UUID          `json:"player_id"`
	Vote     *voting.VoteButton `json:"vote"`
	Message  string             `json:"message"`
}
