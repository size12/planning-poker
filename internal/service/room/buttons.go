package room

import "github.com/size12/planning-poker/internal/entity/voting"

func (r *Room) SetButtonsPack(pack voting.VotePack) {
	r.Lock()
	defer r.Unlock()

	r.Buttons = pack
}
