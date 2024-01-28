package api

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	Port string
	Mux  *http.ServeMux
}

func NewServer(port string) *Server {
	return &Server{
		Port: port,
		Mux:  http.NewServeMux(),
	}
}

func (s *Server) HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "AntiqueFurniture Website")
}

/*
Starts the server to begin listening for requests
*/
func (s *Server) Start() {
	s.Mux.HandleFunc("/", s.HandleRoot)
	s.Mux.HandleFunc("/login", s.HandleLogin)
	s.Mux.HandleFunc("/signup", s.HandleSignup)
	s.Mux.HandleFunc("/logout", s.HandleLogout)
	s.Mux.HandleFunc("/list_furniture", s.HandleListFurniture)
	s.Mux.HandleFunc("/get_furnitures", s.HandleGetFurnitures)
	s.Mux.HandleFunc("/get_furniture", s.HandleGetFurniture)

	s.Use("/checkout", s.HandleCheckout, AuthMiddleware)
	s.Use("/account", s.HandleAccount, AuthMiddleware)

	// initialize SessionManager
	GetSessionManager()

	log.Printf("\x1b[34mListening on port %s\x1b[0m\n", s.Port)
	err := http.ListenAndServe(s.Port, s.Mux)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

/*
This method takes an endpoint and its handler, and then applies
the middleware to the handler in the order they were provided
*/
func (s *Server) Use(
	pattern string,
	handler http.HandlerFunc,
	middlewares ...func(http.HandlerFunc) http.HandlerFunc,
) {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	s.Mux.HandleFunc(pattern, handler)
}
