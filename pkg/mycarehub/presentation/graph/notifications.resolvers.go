package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// SendFCMNotification is the resolver for the sendFCMNotification field.
func (r *mutationResolver) SendFCMNotification(ctx context.Context, registrationTokens []string, data map[string]interface{}, notification firebasetools.FirebaseSimpleNotificationInput) (bool, error) {
	return r.mycarehub.Notification.SendNotification(ctx, registrationTokens, data, &notification)
}

// ReadNotifications is the resolver for the readNotifications field.
func (r *mutationResolver) ReadNotifications(ctx context.Context, ids []string) (bool, error) {
	return r.mycarehub.Notification.ReadNotifications(ctx, ids)
}

// FetchNotifications is the resolver for the fetchNotifications field.
func (r *queryResolver) FetchNotifications(ctx context.Context, userID string, flavour feedlib.Flavour, paginationInput dto.PaginationsInput, filters *domain.NotificationFilters) (*domain.NotificationsPage, error) {
	return r.mycarehub.Notification.FetchNotifications(ctx, userID, flavour, paginationInput, filters)
}

// FetchNotificationTypeFilters is the resolver for the fetchNotificationTypeFilters field.
func (r *queryResolver) FetchNotificationTypeFilters(ctx context.Context, flavour feedlib.Flavour) ([]*domain.NotificationTypeFilter, error) {
	return r.mycarehub.Notification.FetchNotificationTypeFilters(ctx, flavour)
}
