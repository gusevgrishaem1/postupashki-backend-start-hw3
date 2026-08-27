package domain

// Session represents an authenticated task-service session.
type Session struct {
	UserID    string
	SessionID string
}
