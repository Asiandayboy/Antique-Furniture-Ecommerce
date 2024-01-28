package api

import (
	"net/http"
	"sync"

	"github.com/google/uuid"
)

type SessionStore map[string]any

const ErrSessionAlreadyExists string = "Session already exists"

type Session struct {
	SessionID string
	Store     SessionStore
}

type SessionManager struct {
	// a map of all currently running sessions
	Sessions map[string]*Session
}

// used to create a new session with CreateSession()
type SessionTemplate struct {
	SessionID string
}

type CreateSessionError struct {
	message string
}

func (c CreateSessionError) Error() string {
	return c.message
}

// a reference to the SessionManager
var instance *SessionManager
var once sync.Once

// name of the cookie for sessionIDs
const SESSIONID_COOKIE_NAME = "afpsid"

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
Creates a new session and adds it to the map of currently running sessions.

- Provide a template with an empty sessionID "" to generate a random sessionID with UUID

- Provide a template with a string to create a session with that string as the sessionID.
This will error if another session with that ID already exists

This method will return an error if an error occurs when generating a new sessionId.
*/
func (s *SessionManager) CreateSession(template SessionTemplate) (*Session, error) {
	var session *Session
	var id string
	if template.SessionID == "" {
		// generate sessionId
		uuidID, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}

		// create new session
		session = &Session{
			SessionID: uuidID.String(),
			Store:     make(map[string]any),
		}
		id = uuidID.String()
	} else {
		if _, exists := s.Sessions[template.SessionID]; exists {
			return nil, CreateSessionError{
				message: ErrSessionAlreadyExists,
			}
		}

		session = &Session{
			SessionID: template.SessionID,
			Store:     make(map[string]any),
		}
		id = template.SessionID
	}

	// insert session into session manager
	s.Sessions[id] = session

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
Returns the sessionID string from the cookie in the request or ErrNoCookie
*/
func (s *SessionManager) GetSessionID(r *http.Request) (string, error) {
	cookie, err := r.Cookie(SESSIONID_COOKIE_NAME)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

/*
Checks if the client is logged in.

Reads sessionID from request cookie and returns the session and true if a session can be
retrieved with the provided ID. If a session can't be found, then it will return nil and false.
*/
func (s *SessionManager) IsLoggedIn(r *http.Request) (*Session, bool) {
	sessionID, err := s.GetSessionID(r)
	if err != nil {
		return nil, false
	}

	return s.GetSession(sessionID)
}
