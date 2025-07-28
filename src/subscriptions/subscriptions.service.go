package subscriptions

import (
	"context"
	entities "effective_mobile/src/_entities"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo *SubscriptionRepo
}

func NewSubscriptionService(repo *SubscriptionRepo) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) Create(ctx context.Context, data CreateSubscription) (*ResSubscription, error) {
	userID, err := uuid.Parse(data.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	startDate, err := parseMonthYear(data.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %v", err)
	}

	var endDate *time.Time
	if data.EndDate != "" {
		ed, err := parseMonthYear(data.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date: %v", err)
		}
		endDate = &ed
	}

	sub := entities.Subscriptions{
		ServiceName: data.ServiceName,
		Price:       data.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := s.repo.Create(ctx, &sub); err != nil {
		return nil, err
	}

	response := &ResSubscription{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   formatMonthYear(sub.StartDate),
	}

	if sub.EndDate != nil {
		endDateStr := formatMonthYear(*sub.EndDate)
		response.EndDate = &endDateStr
	}

	return response, nil
}

func (s *SubscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*ResSubscription, error) {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return convertToResponse(sub), nil
}

func (s *SubscriptionService) Update(ctx context.Context, id uuid.UUID, data UpdateSubscription) (*ResSubscription, error) {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if data.ServiceName != nil {
		sub.ServiceName = *data.ServiceName
	}
	if data.Price != nil {
		sub.Price = *data.Price
	}
	if data.StartDate != nil {
		startDate, err := parseMonthYear(*data.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date: %v", err)
		}
		sub.StartDate = startDate
	}
	if data.EndDate != nil {
		endDate, err := parseMonthYear(*data.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date: %v", err)
		}
		sub.EndDate = &endDate
	}

	if err := s.repo.Update(ctx, sub); err != nil {
		return nil, err
	}

	return convertToResponse(sub), nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *SubscriptionService) List(ctx context.Context, filter SubscriptionList) ([]ResSubscription, error) {
	options, err := s.parseListOptions(filter)
	if err != nil {
		return nil, err
	}

	subs, err := s.repo.List(ctx, options)
	if err != nil {
		return nil, err
	}

	result := make([]ResSubscription, len(subs))
	for i, sub := range subs {
		result[i] = *convertToResponse(&sub)
	}

	return result, nil
}

func convertToResponse(sub *entities.Subscriptions) *ResSubscription {
	response := &ResSubscription{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   formatMonthYear(sub.StartDate),
	}

	if sub.EndDate != nil {
		endDateStr := formatMonthYear(*sub.EndDate)
		response.EndDate = &endDateStr
	}

	return response
}

func (s *SubscriptionService) parseListOptions(filter SubscriptionList) (ListOptions, error) {
	var options ListOptions

	if filter.UserID != "" {
		userID, err := uuid.Parse(filter.UserID)
		if err != nil {
			return options, fmt.Errorf("invalid user ID")
		}
		options.UserID = &userID
	}

	if filter.Limit != "" {
		limit, err := strconv.Atoi(filter.Limit)
		if err != nil || limit < 1 {
			return options, fmt.Errorf("limit must be greater than 0")
		}
		options.Limit = &limit
	} else {
		defaultLimit := 10
		options.Limit = &defaultLimit
	}

	if filter.Offset != "" {
		offset, err := strconv.Atoi(filter.Offset)
		if err != nil || offset < 0 {
			return options, fmt.Errorf("offset must be greater than or equal to 0")
		}
		options.Offset = &offset
	}

	return options, nil
}

func (s *SubscriptionService) GetSubscriptionSummary(
	ctx context.Context,
	userID *uuid.UUID,
	serviceName string,
	startDateStr, endDateStr string,
) (*ResSubscriptionSummary, error) {
	startDate, err := parseMonthYear(startDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %v", err)
	}

	endDate, err := parseMonthYear(endDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %v", err)
	}

	total, count, err := s.repo.GetSubscriptionSummary(
		ctx,
		userID,
		serviceName,
		startDate,
		endDate,
	)
	if err != nil {
		return nil, err
	}

	return &ResSubscriptionSummary{
		TotalPrice: total,
		StartDate:  startDateStr,
		EndDate:    endDateStr,
		Count:      count,
	}, nil
}

func parseMonthYear(monthYear string) (time.Time, error) {
	return time.Parse("01-2006", monthYear)
}

func formatMonthYear(t time.Time) string {
	return t.Format("01-2006")
}
