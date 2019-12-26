package controllers

import "github.com/jawee/case-api/api/middlewares"

func (s *Server) initializeRoutes() {

  // Home route
  s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.HOME)).Methods("GET")

  // Login route
  s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

  // Users routes
  s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
  s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
  s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
  s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
  s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

  // Cases routes
  s.Router.HandleFunc("/cases", middlewares.SetMiddlewareJSON(s.CreateCase)).Methods("POST")
  s.Router.HandleFunc("/cases", middlewares.SetMiddlewareJSON(s.GetCases)).Methods("GET")
  s.Router.HandleFunc("/cases/{id}", middlewares.SetMiddlewareJSON(s.GetCase)).Methods("GET")
  s.Router.HandleFunc("/cases/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateCase))).Methods("PUT")
  s.Router.HandleFunc("/cases/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

  // CasePriorities routes
  s.Router.HandleFunc("/case_priorities", middlewares.SetMiddlewareJSON(s.CreateCasePriority)).Methods("POST")
  s.Router.HandleFunc("/case_priorities", middlewares.SetMiddlewareJSON(s.GetCasePriorities)).Methods("GET")
  s.Router.HandleFunc("/case_priorities/{id}", middlewares.SetMiddlewareJSON(s.GetCasePriority).Methods("GET")
  s.Router.HandleFunc("/case_priorities/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateCasePriority))).Methods("PUT")
  s.Router.HandleFunc("/case_priorities/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteCasePriority)).Methods("DELETE")

  // CaseStatuses routes
  s.Router.HandleFunc("/case_statuses", middlewares.SetMiddlewareJSON(s.CreateCaseStatus)).Methods("POST")
  s.Router.HandleFunc("/case_statuses", middlewares.SetMiddlewareJSON(s.GetCaseStatuses)).Methods("GET")
  s.Router.HandleFunc("/case_statuses/{id}", middlewares.SetMiddlewareJSON(s.GetCaseStatus).Methods("GET")
  s.Router.HandleFunc("/case_statuses/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateCaseStatus))).Methods("PUT")
  s.Router.HandleFunc("/case_statuses/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteCaseStatus)).Methods("DELETE")
