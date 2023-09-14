package room

import "github.com/google/uuid"

func (r *Room) SetAdmin(ID uuid.UUID) {
	r.Lock()
	defer r.Unlock()

	r.adminID = ID
}

// IsAdmin validates if player with given ID is admin.
func (r *Room) IsAdmin(ID uuid.UUID) bool {
	r.RLock()
	defer r.RUnlock()

	return r.adminID == ID
}
