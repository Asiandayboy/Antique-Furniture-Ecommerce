package api

import (
	"github.com/google/uuid"
	"sync"
)

type SessionStore map[string]any

type Session struct {
	SessionId uuid.UUID
	Store     SessionStore
}

type SessionManager struct {
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
Creates a new session and adds it to the map of stored Sessions.
This method will return an error if an error occurs when generating
a new sessionId
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
Creates a session if the session with the sessionId cannot be found, else
returns the existing session.

Returns an error if an error occured when creating a new session
*/
func (s *SessionManager) GetSession(sessionId string) (*Session, error) {
	if session, sessionExists := s.Sessions[sessionId]; !sessionExists {
		return s.CreateSession()
	} else {
		return session, nil
	}
}
