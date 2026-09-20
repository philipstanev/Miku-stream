package sessions

type Session interface {
	GetInfo(sessionID string) (SessionInfo, error)
	PutUser(info SessionInfo) (string, error)
	RemoveUser(sessionID string) error
}

type SessionInfo struct {
	UserID string
	Role   string
}
