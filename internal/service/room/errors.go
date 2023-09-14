package room

import "errors"

var ErrPlayerNotFound = errors.New("player with such ID does not exists")
var ErrEmptyRoomName = errors.New("room must have non-empty name")
var ErrAccessDenied = errors.New("access denied")
