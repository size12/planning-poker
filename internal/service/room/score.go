package room

import "fmt"

func (r *Room) CalculateScore() string {
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
		r.Score = "-"
		return r.Score
	}

	r.Score = fmt.Sprintf("%.1f", score/float64(cnt))

	return r.Score
}
