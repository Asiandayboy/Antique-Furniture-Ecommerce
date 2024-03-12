package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/rs/cors"
	"github.com/stripe/stripe-go/v76"
)

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

const (
	ErrMethodNotAllowed = "Method Not Allowed"
)

type Server struct {
	Port       string
	Mux        *http.ServeMux
	httpServer *http.Server
}

func NewServer(port string) *Server {
	m := http.NewServeMux()
	s := &http.Server{
		Addr:    port,
		Handler: m,
	}
	return &Server{
		Port:       port,
		Mux:        m,
		httpServer: s,
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
	s.Use("POST /login", s.HandleLogin, logEndpointHit)
	s.Use("POST /signup", s.HandleSignup, logEndpointHit)
	s.Use("POST /logout", s.HandleLogout, AuthMiddleware, logEndpointHit)

	s.Use("POST /list_furniture", s.HandleListFurniture, AuthMiddleware, logEndpointHit)
	s.Use("GET /get_furnitures", s.HandleGetFurnitures, logEndpointHit)
	s.Use("GET /get_furniture", s.HandleGetFurniture, logEndpointHit)

	s.Use("GET /account", s.HandleAccountGET, AuthMiddleware, logEndpointHit)
	s.Use("PUT /account", s.HandleAccountPUT, AuthMiddleware, logEndpointHit)
	s.Use("GET /account/address", s.HandleAddressGET, AuthMiddleware, logEndpointHit)
	s.Use("POST /account/address", s.HandleAddressPOST, AuthMiddleware, logEndpointHit)
	s.Use("PUT /account/address", s.HandleAddressPUT, AuthMiddleware, logEndpointHit)
	s.Use("DELETE /account/address/{addressID}", s.HandleAddressDELETE, AuthMiddleware, logEndpointHit)
	s.Use("GET /account/purchase_history", s.HandlePurchaseHistory, AuthMiddleware, logEndpointHit)

	s.Use("POST /checkout", s.HandleCheckout, AuthMiddleware, logEndpointHit)

	// handle auth in the handler bc cookies aren't sent when Stripe sends the webhook
	s.Use("POST /checkout_webhook", s.HandleStripeWebhook, logEndpointHit)

	/*----------STRIPE-----------*/

	stripe.Key = os.Getenv("STRIPE_TEST_KEY")
	webhookURL := fmt.Sprintf("http://localhost%s/checkout_webhook", s.Port)

	/*
		Since this project is not hosted on the internet, we don't have a public
		URL. So, we need to add a local listener to watch for webhook requests

		You need to have the stripe.exe directory added to your PATH env variables,
		which you can do by installing the Stripe CLI, in order to execute this command
	*/
	command := exec.Command("stripe", "listen", "--forward-to", webhookURL)
	err := command.Start()
	if err != nil {
		log.Fatal("Failed to execute stripe listen command")
	}
	log.Println("\x1b[34mConnected local webhook listener\x1b[0m")

	/*---------------------------*/

	// initialize SessionManager
	GetSessionManager()

	log.Printf("\x1b[34mListening on port %s\x1b[0m\n", s.Port)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:5173"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	err = http.ListenAndServe(s.Port, c.Handler(s.Mux))
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

}

func (s *Server) Shutdown() {
	err := s.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatal("Failed to shutdown server")
	}
}

/*
This method takes an endpoint and its handler, and then applies
the middleware to the handler in the order they were provided,
where the last middleware provided is the one that gets executed
first in the chain
*/
func (s *Server) Use(
	pattern string,
	handler http.HandlerFunc,
	middlewares ...MiddlewareFunc,
) {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	s.Mux.HandleFunc(pattern, handler)
}

/*
A middleware which logs when an endpoint was hit
*/
func logEndpointHit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		method := r.Method

		// format example: ENDPOINT HIT -> GET /account
		log.Printf("\x1b[35mENDPOINT HIT -> %s %s\x1b[0m\n", method, url)

		next.ServeHTTP(w, r)
	}
}
