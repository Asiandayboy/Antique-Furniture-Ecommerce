package api

import (
	"net/http"
	"sync"

	"github.com/google/uuid"
)

type SessionStore map[string]any

type Session struct {
	SessionId uuid.UUID
	Store     SessionStore
}

type SessionManager struct {
	// a map of all currently running sessions
	Sessions map[string]*Session
}

var instance *SessionManager
var once sync.Once

const SESSIONID_COOKIE_NAME = "afpsid" // stands for antique furniture project session id

/*
This function creates a single session manager if there isn't one.
If there is one, it returns the existing one. Only one session
manager can exist at a time
*/
func GetSessionManager() *SessionManager {
	once.Do(func() {
		instance = &SessionManager{
			Sessions: make(map[string]*Session),
		}
	})

	return instance
}

/*
Creates a new session with a UUID and adds it to the map of currently running sessions.

This method will return an error if an error occurs when generating a new sessionId
*/
func (s *SessionManager) CreateSession() (*Session, error) {

	// generate sessionId
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	// create new session
	session := &Session{
		SessionId: id,
		Store:     make(map[string]any),
	}

	// insert session into session manager
	s.Sessions[id.String()] = session

	return session, nil
}

/*
Retrieves the session associated with the provided sessionID.

If a session can be found with the provided sessionID, the session and true will be
returned. If not, nil and false will be returned
*/
func (s *SessionManager) GetSession(sessionId string) (*Session, bool) {
	if session, sessionExists := s.Sessions[sessionId]; !sessionExists {
		return nil, false
	} else {
		return session, true
	}
}

/*
Deletes the session associated with the provided sessionID from the
map of currently running sessions
*/
func (s *SessionManager) DeleteSession(sessionID string) {
	delete(s.Sessions, sessionID)
}

/*
Checks if the client is logged in.

Reads sessionID from request cookie and returns the session and true if a session can be
retrieved with the provided ID. If a session can't be found, then it will return nil and false.
*/
func (s *SessionManager) IsLoggedIn(r *http.Request) (*Session, bool) {
	cookie, err := r.Cookie(SESSIONID_COOKIE_NAME)
	if err != nil {
		return nil, false
	}

	sessionID := cookie.Value
	SessionManager := GetSessionManager()

	return SessionManager.GetSession(sessionID)
}
