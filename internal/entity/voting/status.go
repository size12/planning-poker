package voting

type Status string

const (
	RoomStatusWaiting  Status = "WAITING"
	RoomStatusVoting   Status = "VOTING"
	RoomStatusRevealed Status = "REVEALED"
)
