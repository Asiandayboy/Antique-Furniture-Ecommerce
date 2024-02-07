package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/stripe/stripe-go/v76"
)

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

const (
	ErrPostMethod   = "Request must be a POST request"
	ErrGetMethod    = "Request must be a GET request"
	ErrPutMethod    = "Request must be a PUT request"
	ErrDeleteMethod = "Request must be a DELETE request"
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
	s.Post("/login", s.HandleLogin)
	s.Post("/signup", s.HandleSignup)
	s.Post("/logout", s.HandleLogout, AuthMiddleware)

	s.Post("/list_furniture", s.HandleListFurniture, AuthMiddleware)
	s.Get("/get_furnitures", s.HandleGetFurnitures)
	s.Get("/get_furniture", s.HandleGetFurniture)

	s.Get("/account", s.HandleAccountGET, AuthMiddleware)
	s.Put("/account", s.HandleAccountPUT, AuthMiddleware)
	s.Get("/account/address", s.HandleAddressGET, AuthMiddleware)
	s.Post("/account/address", s.HandleAddressPOST, AuthMiddleware)
	s.Put("/account/address", s.HandleAddressPUT, AuthMiddleware)
	s.Delete("/account/address", s.HandleAddressDELETE, AuthMiddleware)

	s.Post("/checkout", s.HandleCheckout, AuthMiddleware)

	// handle auth in the handler bc cookies aren't sent when Stripe sends the webhook
	s.Post("/checkout_webhook", s.HandleStripeWebhook)

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
	err = s.httpServer.ListenAndServe()
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

// Middleware to verify request method
func requestMethodMiddleware(method string) MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				var errMsg string
				switch method {
				case "GET":
					errMsg = ErrGetMethod
				case "POST":
					errMsg = ErrPostMethod
				case "PUT":
					errMsg = ErrPutMethod
				case "DELETE":
					errMsg = ErrDeleteMethod
				default:
					errMsg = fmt.Sprintf("Unsupported method: %s", method)
				}
				http.Error(w, errMsg, http.StatusMethodNotAllowed)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

/*
Registers the given handler for the
given pattern for POST request methods only. The
Post middleware is appended last to <middlewares>
argument
*/
func (s *Server) Post(
	pattern string,
	handler http.HandlerFunc,
	middlewares ...MiddlewareFunc,
) {
	middlewares = append(middlewares, requestMethodMiddleware("POST"))
	s.Use(pattern, handler, middlewares...)
}

/*
Registers the given handler for the
given pattern for GET request methods only. The
Get middleware is appended last to <middlewares>
argument
*/
func (s *Server) Get(
	pattern string,
	handler http.HandlerFunc,
	middlewares ...MiddlewareFunc,
) {
	middlewares = append(middlewares, requestMethodMiddleware("GET"))
	s.Use(pattern, handler, middlewares...)
}

/*
Registers the given handler for the
given pattern for PUT request methods only. The
Put middleware is appended last to <middlewares>
argument
*/
func (s *Server) Put(
	pattern string,
	handler http.HandlerFunc,
	middlewares ...MiddlewareFunc,
) {
	middlewares = append(middlewares, requestMethodMiddleware("PUT"))
	s.Use(pattern, handler, middlewares...)
}

/*
Registers the given handler for the
given pattern for DELETE request methods only. The
Delete middleware is appended last to <middlewares>
argument
*/
func (s *Server) Delete(
	pattern string,
	handler http.HandlerFunc,
	middlewares ...MiddlewareFunc,
) {
	middlewares = append(middlewares, requestMethodMiddleware("DELETE"))
	s.Use(pattern, handler, middlewares...)
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
