package subscriptions

import (
	"github.com/google/uuid"
)

// CreateSubscription
// swagger:model CreateSubscriptionRequest
type CreateSubscription struct {
	// Name of the service being subscribed to
	ServiceName string `json:"service_name" validate:"required,min=2,max=100"`

	// Price
	Price float64 `json:"price" validate:"required,gt=0"`

	// Unique identifier of the user who owns the subscription
	UserID string `json:"user_id" validate:"required,uuid4"`

	// Date when the subscription start (MM-YYYY format)
	StartDate string `json:"start_date" validate:"required,monthyear"`

	// Date when the subscription end (MM-YYYY format)
	EndDate string `json:"end_date" validate:"required,monthyear"`
}

// UpdateSubscription
// swagger:model UpdateSubscription
type UpdateSubscription struct {
	// New service name
	ServiceName *string `json:"service_name,omitempty" validate:"omitempty,min=2,max=100"`

	// New monthly cost
	Price *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`

	// New start date (MM-YYYY format)
	StartDate *string `json:"start_date,omitempty" validate:"omitempty,monthyear"`

	// New end date (MM-YYYY format)
	EndDate *string `json:"end_date,omitempty" validate:"omitempty,monthyear"`
}

// ResSubscription
// swagger:model SubscriptionResponse
type ResSubscription struct {
	// Unique identifier of the subscription
	ID uuid.UUID `json:"id"`

	// Name of the subscribed service
	ServiceName string `json:"service_name"`

	// Monthly subscription cost in RUB
	Price float64 `json:"price"`

	// Unique identifier of the user
	UserID uuid.UUID `json:"user_id"`

	// Subscription start date (MM-YYYY format)
	StartDate string `json:"start_date"`

	// Optional subscription end date (MM-YYYY format)
	EndDate *string `json:"end_date,omitempty"`
}

// SubscriptionSummary
// swagger:model SubscriptionSummary
type SubscriptionSummary struct {
	// User ID to filter by
	UserID string `json:"user_id" validate:"omitempty,uuid4"`

	// Service name to filter by
	ServiceName string `json:"service_name" validate:"omitempty,min=2,max=100"`

	// Start of the period (MM-YYYY format)
	StartDate string `json:"start_date" validate:"required,monthyear"`

	// End of the period (MM-YYYY format)
	EndDate string `json:"end_date" validate:"required,monthyear"`
}

// ResSubscriptionSummary
// swagger:model ResSubscriptionSummary
type ResSubscriptionSummary struct {
	// Total cost for the period
	TotalPrice float64 `json:"total_price"`

	// Period start (MM-YYYY format)
	StartDate string `json:"start_date"`

	// Period end (MM-YYYY format)
	EndDate string `json:"end_date"`

	// Number of subscriptions matched
	Count int `json:"count"`
}

// ErrorResponse represents API error response
// swagger:model ErrorResponse
type ErrorResponse struct {
	// Error message
	Message string `json:"message"`

	// Optional list of detailed errors
	Errors []string `json:"errors,omitempty"`

	// HTTP status code
	StatusCode int `json:"status_code"`
}

// SubscriptionList contains filtering parameters
// swagger:parameters subscriptionList
type SubscriptionList struct {
	// User ID to filter by
	UserID string `json:"user_id"`

	// Maximum number of results to return
	Limit string `json:"limit"`

	// Offset for pagination
	Offset string `json:"offset"`
}
