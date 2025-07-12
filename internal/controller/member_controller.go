package controller

import (
	"azureclient/internal/model"
	"azureclient/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type MemberController struct {
	Service service.MemberService
}

func NewMemberController(s service.MemberService) *MemberController {
	return &MemberController{Service: s}
}

// For net/http, use ListMembers for /members and MemberHandler for /members/{id}

// Handles /members and /members/ (GET, POST)

// GET /members
func (c *MemberController) GetMembers(w http.ResponseWriter, r *http.Request) error {
	members, err := c.Service.ListMembers(r.Context())
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(members)
}

// POST /members
func (c *MemberController) CreateMember(w http.ResponseWriter, r *http.Request) error {
	var member model.Member
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		return err
	}
	if err := c.Service.CreateMember(r.Context(), &member); err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(member)
}

// GET /members/{id}
func (c *MemberController) GetMemberByID(w http.ResponseWriter, r *http.Request) error {
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		return err
	}
	member, err := c.Service.GetMemberByID(r.Context(), uint(id))
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(member)
}

// PUT /members/{id}
func (c *MemberController) UpdateMember(w http.ResponseWriter, r *http.Request) error {
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		return err
	}
	var member model.Member
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		return err
	}
	member.ID = uint(id)
	if err := c.Service.UpdateMember(r.Context(), &member); err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(member)
}

// DELETE /members/{id}
func (c *MemberController) DeleteMember(w http.ResponseWriter, r *http.Request) error {
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		return err
	}
	if err := c.Service.DeleteMember(r.Context(), uint(id)); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Helper to extract ID from /members/{id}
func extractIDFromPath(path string) (int, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return 0, http.ErrMissingFile // Not Found
	}
	idStr := parts[1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return id, nil
}
