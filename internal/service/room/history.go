package room

import (
	"time"

	"github.com/size12/planning-poker/internal/entity"
)

func (r *Room) SaveToHistory() error {
	r.Lock()
	defer r.Unlock()
	r.History = append(r.History, entity.Score{
		Time:  time.Now().Format("2006-01-02 15:04:05"),
		Score: r.Score,
	})

	return nil
}
