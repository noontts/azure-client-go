package internal

import (
	"azureclient/internal/controller"
	"azureclient/internal/middleware"
	"net/http"
)

const (
   RouteMembers    = "/members"
   RouteMemberByID = "/members/"
   RouteSendEmail  = "/send-email"
)

// SetupRoutes registers all HTTP routes and handlers
func SetupRoutes(mux *http.ServeMux, memberController *controller.MemberController, emailController *controller.EmailController) {
   // POST /send-email
   mux.HandleFunc(RouteSendEmail, methodHandler(http.MethodPost, middleware.ErrorHandler(emailController.SendEmail)))
	// GET /members
	mux.HandleFunc(RouteMembers, methodHandler(http.MethodGet, middleware.ErrorHandler(memberController.GetMembers)))
	// POST /members
	mux.HandleFunc(RouteMembers, methodHandler(http.MethodPost, middleware.ErrorHandler(memberController.CreateMember)))

	// GET /members/{id}
	mux.HandleFunc(RouteMemberByID, methodHandler(http.MethodGet, middleware.ErrorHandler(memberController.GetMemberByID)))
	// PUT /members/{id}
	mux.HandleFunc(RouteMemberByID, methodHandler(http.MethodPut, middleware.ErrorHandler(memberController.UpdateMember)))
	// DELETE /members/{id}
	mux.HandleFunc(RouteMemberByID, methodHandler(http.MethodDelete, middleware.ErrorHandler(memberController.DeleteMember)))
}

// methodHandler returns a handler that only responds to the given method, otherwise 405
func methodHandler(method string, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handlerFunc(w, r)
	}
}
