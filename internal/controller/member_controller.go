package controller

import (
	"azureclient/internal/middleware"
	"azureclient/internal/model"
	"azureclient/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type MemberController struct {
	Service service.MemberService
}

func NewMemberController(s service.MemberService) *MemberController {
	return &MemberController{Service: s}
}

func (c *MemberController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/members", middleware.ErrorHandler(c.CreateMember)).Methods("POST")
	r.HandleFunc("/members/{id}", middleware.ErrorHandler(c.GetMemberByID)).Methods("GET"))
	r.HandleFunc("/members/{id}", middleware.ErrorHandler(c.UpdateMember)).Methods("PUT")
	r.HandleFunc("/members/{id}", middleware.ErrorHandler(c.DeleteMember)).Methods("DELETE")
	r.HandleFunc("/members", middleware.ErrorHandler(c.ListMembers)).Methods("GET")
}

func (c *MemberController) CreateMember(w http.ResponseWriter, r *http.Request) error {
	var member model.Member
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	if err := c.Service.CreateMember(r.Context(), &member); err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(member)
	return nil
}

func (c *MemberController) GetMemberByID(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return nil
	}
	member, err := c.Service.GetMemberByID(r.Context(), uint(id))
	if err != nil {
		return err
	}
	json.NewEncoder(w).Encode(member)
	return nil
}

func (c *MemberController) UpdateMember(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return nil
	}
	var member model.Member
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	member.ID = uint(id)
	if err := c.Service.UpdateMember(r.Context(), &member); err != nil {
		return err
	}
	json.NewEncoder(w).Encode(member)
	return nil
}

func (c *MemberController) DeleteMember(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return nil
	}
	if err := c.Service.DeleteMember(r.Context(), uint(id)); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (c *MemberController) ListMembers(w http.ResponseWriter, r *http.Request) error {
	members, err := c.Service.ListMembers(r.Context())
	if err != nil {
		return err
	}
	json.NewEncoder(w).Encode(members)
	return nil
}
