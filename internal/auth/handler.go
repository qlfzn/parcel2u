package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,min=3,max=72"`
	Role     string `json:"role" validate:"omitempty,oneof=student dispatcher admin"`
}

type LoginUserPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token,omitempty"`
}

type Handler struct {
	userService UserService
}

func NewHandler(userService UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

var validate = validator.New()

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("decoder error: %s", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(payload); err != nil {
		http.Error(w, "validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	user := &User{
		Username: payload.Username,
		Role:     payload.Role,
	}

	if err := user.Password.SetHash(payload.Password); err != nil {
		http.Error(w, "failed to process password", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	if err := h.userService.Create(ctx, user); err != nil {
		http.Error(w, "failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := CreateToken(user.Username)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		User:  user,
		Token: token,
	}

	h.writeJSONResponse(w, response, http.StatusCreated)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var payload LoginUserPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("LoginUser: invalid request body: %v", err)
		h.writeErrorResponse(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(payload); err != nil {
		log.Printf("LoginUser: validation failed: %v", err)
		h.writeErrorResponse(w, "validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user, err := h.userService.GetByUsername(ctx, payload.Username)
	if err != nil {
		log.Printf("LoginUser: user not found or db error for username '%s': %v", payload.Username, err)
		h.writeErrorResponse(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := user.Password.CompareHash(payload.Password); err != nil {
		log.Printf("LoginUser: password mismatch for username '%s': %v", payload.Username, err)
		h.writeErrorResponse(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := CreateToken(user.Username)
	if err != nil {
		log.Printf("LoginUser: failed to generate token for username '%s': %v", payload.Username, err)
		h.writeErrorResponse(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		User:  user,
		Token: token,
	}

	log.Printf("LoginUser: successful login for username '%s'", payload.Username)
	h.writeJSONResponse(w, response, http.StatusOK)
}

func (h *Handler) writeJSONResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeErrorResponse(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
