package auth

type SessionStore struct{}

func NewSessionStore() *SessionStore {
	return &SessionStore{}
}
