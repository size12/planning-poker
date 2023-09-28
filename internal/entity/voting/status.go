package voting

type Status string

const (
	RoomStatusWaiting  Status = "WAITING"
	RoomStatusVoting   Status = "VOTING"
	RoomStatusRevealed Status = "REVEALED"
)

const (
	PlayerVoting    Status = "PLAYER_VOTING"
	PlayerObserving Status = "PLAYER_OBSERVING"
)
