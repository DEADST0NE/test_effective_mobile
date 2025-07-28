package subscriptions

import (
	"context"
	"fmt"
	"time"

	entities "effective_mobile/src/_entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepo struct {
	db *gorm.DB
}

func NewSubscriptionRepo(db *gorm.DB) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) Create(ctx context.Context, sub *entities.Subscriptions) error {
	return r.db.WithContext(ctx).Create(sub).Error
}

func (r *SubscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*entities.Subscriptions, error) {
	var sub entities.Subscriptions
	err := r.db.WithContext(ctx).First(&sub, "id = ?", id).Error
	return &sub, err
}

func (r *SubscriptionRepo) Update(ctx context.Context, sub *entities.Subscriptions) error {
	return r.db.WithContext(ctx).Save(sub).Error
}

func (r *SubscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Subscriptions{}, "id = ?", id).Error
}

func (r *SubscriptionRepo) List(ctx context.Context, filter ListOptions) ([]entities.Subscriptions, error) {
	var subs []entities.Subscriptions

	query := r.db.WithContext(ctx).Model(&entities.Subscriptions{})

	if filter.UserID != nil {
		query = query.Where("user_id = ?", filter.UserID)
	}

	if filter.Limit != nil {
		query = query.Limit(*filter.Limit)
	}

	if filter.Offset != nil {
		query = query.Offset(*filter.Offset)
	}

	if err := query.Find(&subs).Error; err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	return subs, nil
}

func (r *SubscriptionRepo) GetSubscriptionSummary(
	ctx context.Context,
	userID *uuid.UUID,
	serviceName string,
	startDate, endDate time.Time,
) (float64, int, error) {
	var result struct {
		Total float64
		Count int
	}

	query := r.db.WithContext(ctx).
		Model(&entities.Subscriptions{}).
		Select("COALESCE(SUM(price), 0) as total, COUNT(*) as count").
		Where("start_date <= ? AND (end_date >= ? OR end_date IS NULL)", endDate, startDate)

	if userID != nil {
		query = query.Where("user_id = ?", userID)
	}

	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	if err := query.Scan(&result).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to calculate summary: %w", err)
	}

	return result.Total, result.Count, nil
}

type ListOptions struct {
	UserID *uuid.UUID
	Limit  *int
	Offset *int
}
