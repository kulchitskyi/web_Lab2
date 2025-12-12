package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	er "deu/internal/errors"
	"deu/internal/models"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func validateRequest(s interface{}) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errorsMap := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errorsMap[err.Field()] = fmt.Sprintf("Field '%s' failed validation: tag '%s'.",
			err.Field(),
			err.Tag(),
			err.Param(),
		)
	}

	return errorsMap
}

type Handler struct {
	Service *UserService
}

// validateAndGetID checks if the ID string is a valid UUID and extracts it.
func validateAndGetID(w http.ResponseWriter, parts []string) (string, bool) {
	if len(parts) < 3 || parts[2] == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "User ID missing in path"})
		return "", false
	}
	id := parts[2]
	
	if err := validate.Var(id, "required,uuid"); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "User ID must be a valid UUID"})
		return "", false
	}
	return id, true
}


// GET /users
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.Service.GetAll(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// GET /users/{id}
func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	
	id, ok := validateAndGetID(w, parts)
	if !ok {
		return
	}

	u, err := h.Service.GetById(r.Context(), id)
	if err != nil {
		if err == er.ErrUserNotFound {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, u)
}

// POST /users
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var u models.UserCreateRequest
	
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
		return
	}

	if errorsMap := validateRequest(u); errorsMap != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"validation_errors": errorsMap})
		return
	}
	
	user, err := h.Service.Create(r.Context(), &u)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, *user)
}

// PATCH /users/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	
	id, ok := validateAndGetID(w, parts)
	if !ok {
		return
	}

	var u models.UserUpdateRequest
	
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
		return
	}

	if errorsMap := validateRequest(u); errorsMap != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"validation_errors": errorsMap})
		return
	}

	err := h.Service.Update(r.Context(), id, &u)
	if err != nil {
		if err == er.ErrUserNotFound {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// DELETE /users/{id}
func (h *Handler) DeleteById(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	
	id, ok := validateAndGetID(w, parts)
	if !ok {
		return
	}

	err := h.Service.DeleteById(r.Context(), id)
	if err != nil {
		if err == er.ErrUserNotFound {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// POST /users/{id}/places/{place_id}
func (h *Handler) AddVisitedPlace(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	
	if len(parts) < 5 || parts[2] == "" || parts[4] == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing user ID or place ID"})
		return
	}
	userID := parts[2]
	placeID := parts[4]
	
	if err := validate.Var(userID, "required,uuid"); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "User ID must be a valid UUID"})
		return
	}

	if err := validate.Var(placeID, "required,uuid"); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Place ID must be a valid UUID"})
		return
	}

	err := h.Service.AddVisitedPlace(r.Context(), userID, placeID)

	if err != nil {
		switch err {
		case er.ErrUserNotFound:
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		case er.ErrPlaceNotFound:
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Place not found"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "Place added to user's visited list"})
}

// GET /users/{id}/places/{place_id}
func (h *Handler) CheckIfVisited(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 5 || parts[2] == "" || parts[4] == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing user ID or place ID"})
		return
	}
	userID := parts[2]
	placeID := parts[4]

	if err := validate.Var(userID, "required,uuid"); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "User ID must be a valid UUID"})
		return
	}

	if err := validate.Var(placeID, "required,uuid"); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Place ID must be a valid UUID"})
		return
	}

	visited, err := h.Service.HasVisitedPlace(r.Context(), userID, placeID)

	if err != nil {
		if err == er.ErrUserNotFound || err == er.ErrPlaceNotFound {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if visited {
		writeJSON(w, http.StatusOK, map[string]interface{}{"visited": true, "message": "User has visited this place"})
	} else {
		writeJSON(w, http.StatusOK, map[string]interface{}{"visited": false, "message": "User has NOT visited this place"})
	}
}

// DELETE /users/{id}/places/{place_id}
func (h *Handler) RemoveVisitedPlace(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	
	if len(parts) < 5 || parts[2] == "" || parts[4] == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing user ID or place ID"})
		return
	}
	userID := parts[2]
	placeID := parts[4]
	
	if err := validate.Var(userID, "required,uuid"); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "User ID must be a valid UUID"})
		return
	}

	if err := validate.Var(placeID, "required,uuid"); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Place ID must be a valid UUID"})
		return
	}

	err := h.Service.RemoveVisitedPlace(r.Context(), userID, placeID)

	if err != nil {
		switch err {
		case er.ErrUserNotFound:
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		case er.ErrPlaceNotFound:
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Place not found"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "Place removed from user's visited list"})
}

// DELETE /users
func (h *Handler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	err := h.Service.DeleteAll(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "all users deleted"})
}