package controllers

import "github.com/task/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetUsers))).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetUser))).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Products routes
	s.Router.HandleFunc("/products", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetProducts))).Methods("GET")
	s.Router.HandleFunc("/products/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetProduct))).Methods("GET")
	s.Router.HandleFunc("/products", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthseller(s.CreateProduct))).Methods("POST")
	s.Router.HandleFunc("/products/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthseller(s.UpdateProduct))).Methods("PUT")
	s.Router.HandleFunc("/products/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteProduct)).Methods("DELETE")
	s.Router.HandleFunc("/buy", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthBuyer((s.BuyProduct)))).Methods("POST")

}
