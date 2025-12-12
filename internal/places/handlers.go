package places

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
		errorsMap[err.Field()] = fmt.Sprintf("Field '%s' failed validation: tag '%s'. Required value: %s",
			err.Field(),
			err.Tag(),
			err.Param(),
		)
	}

	return errorsMap
}

func validateAndGetID(w http.ResponseWriter, parts []string) (string, bool) {
	if len(parts) < 3 || parts[2] == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Place ID missing in path"})
		return "", false
	}
	id := parts[2]
	
	if err := validate.Var(id, "required,uuid"); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Place ID must be a valid UUID"})
		return "", false
	}
	return id, true
}

type Handler struct {
	Service       *PlaceService
	AllowDeletion bool
}

// GET /places
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.Service.GetAll(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// GET /places/{id}
func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	
	id, ok := validateAndGetID(w, parts)
	if !ok {
		return
	}

	p, err := h.Service.GetById(r.Context(), id)
	if err != nil {
		if err == er.ErrPlaceNotFound {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Place not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, p)
}

// POST /places
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var p models.PlaceCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
		return
	}
	
	if errorsMap := validateRequest(p); errorsMap != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"validation_errors": errorsMap})
		return
	}

	place, err := h.Service.Create(r.Context(), &p)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, *place)
}

// PATCH /places/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	
	id, ok := validateAndGetID(w, parts)
	if !ok {
		return
	}
	
	var p models.PlaceUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
		return
	}
	
	if errorsMap := validateRequest(p); errorsMap != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"validation_errors": errorsMap})
		return
	}

	err := h.Service.Update(r.Context(), id, &p)
	if err != nil {
		if err == er.ErrPlaceNotFound {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Place not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// DELETE /places/{id}
func (h *Handler) DeleteById(w http.ResponseWriter, r *http.Request) {
	if !h.AllowDeletion {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Place deletion is disabled by system configuration"})
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	
	id, ok := validateAndGetID(w, parts)
	if !ok {
		return
	}

	err := h.Service.DeleteById(r.Context(), id)
	if err != nil {
		if err == er.ErrPlaceNotFound {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Place not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// DELETE /places
func (h *Handler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	if !h.AllowDeletion {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Place deletion is disabled by system configuration"})
		return
	}

	err := h.Service.DeleteAll(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "all deleted"})
}