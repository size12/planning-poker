package voting

import "fmt"

type VoteButton struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type VotePack struct {
	Votes []VoteButton `json:"votes"`
}

func (v *VoteButton) String() string {
	return fmt.Sprintf("%s:%.1f", v.Name, v.Value)
}
