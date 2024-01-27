package api

import (
	"backend/db"
	"backend/types"
	"backend/util"
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *Server) HandleSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Request must be a POST request", http.StatusBadRequest)
		return
	}

	// read json
	var signupInfo types.User
	err := util.ReadJSONReq[types.User](r, &signupInfo)
	if err != nil {
		http.Error(w, "Could not decode request body into JSON", http.StatusBadRequest)
		return
	}

	usernameUnique := db.CheckFieldUniqueness("username", signupInfo.Username)
	if !usernameUnique { // username not unique
		http.Error(w, "Username is taken", http.StatusConflict)
		return
	}

	emailUnique := db.CheckFieldUniqueness("email", signupInfo.Email)
	if !emailUnique { // email not unique
		http.Error(w, "Email is taken", http.StatusConflict)
		return
	}

	// create new session before saving info to database
	sessionManager := GetSessionManager()
	session, err := sessionManager.CreateSession(SessionTemplate{SessionID: ""})
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}
	signupInfo.SessionId = session.SessionId

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
		http.Error(w, "Failed to save user account", http.StatusInternalServerError)
		return
	}

	/*
		On the frontend, the client should be redirected to the login page
	*/
	w.Write([]byte("success"))
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Request must be a POST request", http.StatusBadRequest)
		return
	}

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
		http.Error(w, "Invalid login", http.StatusUnauthorized)
		return
	}

	if err = res.Decode(&userResult); err != nil {
		http.Error(w, "Invalid login", http.StatusUnauthorized)
		return
	}

	// compare passwords
	err = util.CheckPassword(loginInfo.Password, userResult.Password)
	if err != nil {
		http.Error(w, "Invalid login", http.StatusUnauthorized)
		return
	}

	// save appropriate session data to session store
	sessionManager := GetSessionManager()
	session, exists := sessionManager.GetSession(userResult.SessionId)
	if !exists {
		// create new session
		session, err = sessionManager.CreateSession(SessionTemplate{SessionID: ""})
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}
	}
	session.Store["username"] = userResult.Username
	session.Store["userid"] = userResult.UserId

	// send back response with cookie
	cookie := http.Cookie{
		Name:  SESSIONID_COOKIE_NAME,
		Value: session.SessionId,
		Path:  "/",
	}

	http.SetCookie(w, &cookie)

	w.Write([]byte("success"))
}

func (s *Server) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Request must be a POST request", http.StatusMethodNotAllowed)
		return
	}

	// delete session
	sessionManager := GetSessionManager()
	sessionID, err := sessionManager.GetSessionID(r)
	if err == http.ErrNoCookie {
		http.Error(w, "SessionID cookie not found", http.StatusBadRequest)
		return
	}

	_, sessioExists := sessionManager.GetSession(sessionID)
	if !sessioExists {
		http.Error(w, "Session does not exist", http.StatusUnauthorized)
		return
	}

	sessionManager.DeleteSession(sessionID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}
