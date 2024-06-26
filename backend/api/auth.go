package api

import (
	"backend/db"
	"backend/types"
	"backend/util"
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ErrUnauthorized     = "Session not found; you must be logged in" // Authentication failed
	ErrInvalidLogin     = "Invalid login"                            // Invalid login credentials
	ErrUsernameTaken    = "Username is taken"
	ErrEmailTaken       = "Email is taken"
	ErrSignupSave       = "Failed to save user account" // New account failed to save into DB
	ErrPasswordMismatch = "Passwords do not match"
	ErrBlankFields      = "One or more fields are blank!"
)

func arePasswordsSame(signupInfo types.User) bool {
	return signupInfo.Password == signupInfo.ConfirmPass
}

// Returns true if any of the signup form fields are blank, else returns false
func signupFieldsBlank(signupInfo types.User) bool {
	if signupInfo.Username == "" ||
		signupInfo.Email == "" ||
		signupInfo.Password == "" ||
		signupInfo.ConfirmPass == "" {
		return true
	}

	return false
}

/*
TODO: add confirm password check 2/20
*/
func (s *Server) HandleSignup(w http.ResponseWriter, r *http.Request) {
	// read json
	var signupInfo types.User
	err := util.ReadJSONReq[types.User](r, &signupInfo)
	if err != nil {
		http.Error(w, "Could not decode request body into JSON", http.StatusBadRequest)
		return
	}

	if signupFieldsBlank(signupInfo) {
		http.Error(w, ErrBlankFields, http.StatusBadRequest)
		return
	}

	if !arePasswordsSame(signupInfo) {
		http.Error(w, ErrPasswordMismatch, http.StatusBadRequest)
		return
	}

	usernameUnique := db.CheckFieldUniqueness("username", signupInfo.Username)
	if !usernameUnique { // username not unique
		http.Error(w, ErrUsernameTaken, http.StatusConflict)
		return
	}

	emailUnique := db.CheckFieldUniqueness("email", signupInfo.Email)
	if !emailUnique { // email not unique
		http.Error(w, ErrEmailTaken, http.StatusConflict)
		return
	}

	signupInfo.Subscribed = false

	//set balance to 0
	balance, err := primitive.ParseDecimal128("0")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	signupInfo.Balance = balance

	// create new session before saving info to database
	sessionManager := GetSessionManager()
	session, err := sessionManager.CreateSession(SessionTemplate{SessionID: ""})
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}
	signupInfo.SessionID = session.SessionID

	// delete session since we don't want to store the session in memory yet
	sessionManager.DeleteSession(session.SessionID)

	// hash password
	hashedPassword, err := util.HashPassword(signupInfo.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	signupInfo.Password = hashedPassword

	// insert signupInfo into DB
	_, err = db.InsertIntoUsersCollection(signupInfo)
	if err != nil {
		http.Error(w, ErrSignupSave, http.StatusInternalServerError)
		return
	}

	/*
		On the frontend, the client should be redirected to the login page
	*/
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

const CookieExpiration = 30 // 30 minutes

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// read json
	var loginInfo types.User
	err := util.ReadJSONReq[types.User](r, &loginInfo)
	if err != nil {
		http.Error(w, "Could not decode request body into JSON", http.StatusBadRequest)
		return
	}

	// find document by username, since each username is constrained to be unique
	var userResult types.User
	usersCollection := db.GetCollection("users")
	res := usersCollection.FindOne(context.Background(), bson.M{
		"username": loginInfo.Username,
	})

	if res.Err() != nil {
		http.Error(w, ErrInvalidLogin, http.StatusUnauthorized)
		return
	}

	if err = res.Decode(&userResult); err != nil {
		http.Error(w, ErrInvalidLogin, http.StatusUnauthorized)
		return
	}

	// compare passwords
	err = util.CheckPassword(loginInfo.Password, userResult.Password)
	if err != nil {
		http.Error(w, ErrInvalidLogin, http.StatusUnauthorized)
		return
	}

	// save appropriate session data to session store
	sessionManager := GetSessionManager()
	sessionExpired := sessionManager.IsSessionExpired(r)

	var session *Session

	if sessionExpired {
		// generate a new sessionID
		session, err = sessionManager.CreateSession(SessionTemplate{SessionID: ""})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = usersCollection.UpdateByID(
			context.Background(),
			userResult.UserID,
			bson.M{"$set": bson.M{"sessionid": session.SessionID}},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("New session created:", session.SessionID)

		// generate cookie
		expiration := time.Now().Add(CookieExpiration * time.Minute)
		cookie := http.Cookie{
			Name:     SESSIONID_COOKIE_NAME,
			Value:    session.SessionID,
			Path:     "http://127.0.0.1:1573",
			Expires:  expiration,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		}

		http.SetCookie(w, &cookie)
	} else {
		// Client has sessionID cookie already stored in DB and it's not expired
		session, err = sessionManager.CreateSession(SessionTemplate{
			SessionID: userResult.SessionID,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Existing session found:", userResult.SessionID)
	}

	session.Store["username"] = userResult.Username
	session.Store["userid"] = userResult.UserID

	w.Write([]byte("success"))
}

func (s *Server) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// delete session
	session := r.Context().Value(SessionKey).(*Session)

	sessionManager := GetSessionManager()
	sessionManager.DeleteSession(session.SessionID)

	cookie := http.Cookie{
		Name:     SESSIONID_COOKIE_NAME,
		Value:    "",
		Path:     "http://127.0.0.1:1573",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

type CtxSessionKey string

// Key name for the attached context session value from AuthMiddleware
const SessionKey CtxSessionKey = "session"

/*
This middleware checks if the client is authenticated by retrieving
their sessionID cookie and validating that their session exists.

If it does exist, the <next> handler will be called with the session
attached to the request context, with the key name of <SessionKey>.
Otherwise, a 401 status code will be returned
*/
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionManager := GetSessionManager()

		session, isLoggedIn := sessionManager.IsLoggedIn(r)
		if !isLoggedIn {
			http.Error(w, ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), SessionKey, session)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
