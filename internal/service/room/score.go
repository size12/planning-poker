package room

import (
	"math"
)

func (r *Room) CalculateScore() {
	r.Lock()
	defer r.Unlock()

	score := 0.0
	cnt := 0

	for _, player := range r.Players {
		if player.Voted != nil {
			score += player.Voted.Value
			cnt += 1
		}
	}

	if cnt == 0 {
		r.Score = "?"
		return
	}

	r.Score = r.findCardByScore(score / float64(cnt))
}

func (r *Room) findCardByScore(score float64) string {
	result := "?"
	diff := math.Inf(1)

	for _, card := range r.Buttons.Votes {

		if math.Abs(card.Value-score) < diff {
			diff = math.Abs(card.Value - score)
			result = card.Name
		}
	}

	return result
}
