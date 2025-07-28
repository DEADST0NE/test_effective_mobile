package subscriptions

import (
	"encoding/json"
	"net/http"

	"effective_mobile/src/_core/validator"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type SubscriptionController struct {
	service *SubscriptionService
}

func NewSubscriptionController(service *SubscriptionService) *SubscriptionController {
	return &SubscriptionController{service: service}
}

func (c *SubscriptionController) RegisterRoutes(r *mux.Router) {
	validator.Init()

	r.HandleFunc("/subscriptions/summary", c.GetSubscriptionSummary).Methods("GET")
	r.HandleFunc("/subscriptions", c.Create).Methods("POST")
	r.HandleFunc("/subscriptions/{id}", c.GetByID).Methods("GET")
	r.HandleFunc("/subscriptions/{id}", c.Update).Methods("PUT")
	r.HandleFunc("/subscriptions/{id}", c.Delete).Methods("DELETE")
	r.HandleFunc("/subscriptions", c.List).Methods("GET")
}

// Create godoc
// @Summary Create a new subscription
// @Description Creates a new service subscription
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param request body CreateSubscription true "Subscription data"
// @Success 201 {object} ResSubscription
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions [post]
func (c *SubscriptionController) Create(w http.ResponseWriter, r *http.Request) {
	var data CreateSubscription
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validator.Validate.Struct(data); err != nil {
		errorResponse(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	resp, err := c.service.Create(r.Context(), data)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to create subscription", err.Error())
		return
	}

	responseWriter(w, http.StatusCreated, resp)
}

// GetByID godoc
// @Summary Get subscription by ID
// @Description Returns a single subscription by its ID
// @Tags Subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} ResSubscription
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/{id} [get]
func (c *SubscriptionController) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	subscription, err := c.service.GetByID(r.Context(), id)
	if err != nil {
		errorResponse(w, http.StatusNotFound, "Subscription not found")
		return
	}

	responseWriter(w, http.StatusOK, subscription)
}

// Update godoc
// @Summary Update subscription
// @Description Updates an existing subscription
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param request body UpdateSubscription true "Subscription update data"
// @Success 200 {object} ResSubscription
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/{id} [put]
func (c *SubscriptionController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	var data UpdateSubscription
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validator.Validate.Struct(data); err != nil {
		errorResponse(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	updated, err := c.service.Update(r.Context(), id, data)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to update subscription", err.Error())
		return
	}

	responseWriter(w, http.StatusOK, updated)
}

// Delete godoc
// @Summary Delete subscription
// @Description Deletes a subscription by ID
// @Tags Subscriptions
// @Param id path string true "Subscription ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/{id} [delete]
func (c *SubscriptionController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	if err := c.service.Delete(r.Context(), id); err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to delete subscription", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List godoc
// @Summary List subscriptions
// @Description Returns a list of subscriptions with optional filtering
// @Tags Subscriptions
// @Produce json
// @Param request query SubscriptionList true "Summary list subctiptions"
// @Success 200 {array} ResSubscription
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions [get]
func (c *SubscriptionController) List(w http.ResponseWriter, r *http.Request) {
	filter := SubscriptionList{
		UserID: r.URL.Query().Get("user_id"),
		Limit:  r.URL.Query().Get("limit"),
		Offset: r.URL.Query().Get("offset"),
	}

	subscriptions, err := c.service.List(r.Context(), filter)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to list subscriptions", err.Error())
		return
	}

	responseWriter(w, http.StatusOK, subscriptions)
}

func responseWriter(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func errorResponse(w http.ResponseWriter, statusCode int, message string, details ...string) {
	response := ErrorResponse{
		Message:    message,
		Errors:     details,
		StatusCode: statusCode,
	}
	responseWriter(w, statusCode, response)
}

// GetSubscriptionSummary godoc
// @Summary Get subscription summary
// @Description Calculate total cost of subscriptions for selected period with optional filters
// @Tags Subscriptions
// @Produce json
// @Param request query SubscriptionSummary true "Summary request parameters"
// @Success 200 {object} ResSubscriptionSummary
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/summary [get]
func (c *SubscriptionController) GetSubscriptionSummary(w http.ResponseWriter, r *http.Request) {
	req := SubscriptionSummary{
		UserID:      r.URL.Query().Get("user_id"),
		ServiceName: r.URL.Query().Get("service_name"),
		StartDate:   r.URL.Query().Get("start_date"),
		EndDate:     r.URL.Query().Get("end_date"),
	}

	if err := validator.Validate.Struct(req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	var userID *uuid.UUID
	if req.UserID != "" {
		parsedID, err := uuid.Parse(req.UserID)
		if err != nil {
			errorResponse(w, http.StatusBadRequest, "Invalid user ID format", err.Error())
			return
		}
		userID = &parsedID
	}

	summary, err := c.service.GetSubscriptionSummary(
		r.Context(),
		userID,
		req.ServiceName,
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to calculate summary", err.Error())
		return
	}

	responseWriter(w, http.StatusOK, summary)
}
