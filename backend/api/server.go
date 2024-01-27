package api

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	Port string
}

func NewServer(port string) *Server {
	return &Server{
		Port: port,
	}
}

/*
Starts the server to begin listening for requests
*/
func (s *Server) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.HandleRoot)
	mux.HandleFunc("/login", s.HandleLogin)
	mux.HandleFunc("/signup", s.HandleSignup)
	mux.HandleFunc("/logout", s.HandleLogout)
	mux.HandleFunc("/checkout", s.HandleCheckout)
	mux.HandleFunc("/list_furniture", s.HandleListFurniture)
	mux.HandleFunc("/get_furnitures", s.HandleGetFurnitures)
	mux.HandleFunc("/get_furniture", s.HandleGetFurniture)

	// initialize SessionManager
	GetSessionManager()

	log.Printf("\x1b[34mListening on port %s\x1b[0m\n", s.Port)
	err := http.ListenAndServe(s.Port, mux)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func (s *Server) HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "AntiqueFurniture Website")
}
