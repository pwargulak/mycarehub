package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// SetInProgressBy is the resolver for the setInProgressBy field.
func (r *mutationResolver) SetInProgressBy(ctx context.Context, serviceRequestID string, staffID string) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.ServiceRequest.SetInProgressBy(ctx, serviceRequestID, staffID)
}

// CreateServiceRequest is the resolver for the createServiceRequest field.
func (r *mutationResolver) CreateServiceRequest(ctx context.Context, input dto.ServiceRequestInput) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.ServiceRequest.CreateServiceRequest(ctx, &input)
}

// ResolveServiceRequest is the resolver for the resolveServiceRequest field.
func (r *mutationResolver) ResolveServiceRequest(ctx context.Context, staffID string, requestID string, action []string, comment *string) (bool, error) {
	return r.mycarehub.ServiceRequest.ResolveServiceRequest(ctx, &staffID, &requestID, action, comment)
}

// VerifyClientPinResetServiceRequest is the resolver for the verifyClientPinResetServiceRequest field.
func (r *mutationResolver) VerifyClientPinResetServiceRequest(ctx context.Context, clientID string, serviceRequestID string, cccNumber string, phoneNumber string, physicalIdentityVerified bool, state string) (bool, error) {
	return r.mycarehub.ServiceRequest.VerifyClientPinResetServiceRequest(ctx, clientID, serviceRequestID, cccNumber, phoneNumber, physicalIdentityVerified, state)
}

// VerifyStaffPinResetServiceRequest is the resolver for the verifyStaffPinResetServiceRequest field.
func (r *mutationResolver) VerifyStaffPinResetServiceRequest(ctx context.Context, phoneNumber string, serviceRequestID string, verificationStatus string) (bool, error) {
	return r.mycarehub.ServiceRequest.VerifyStaffPinResetServiceRequest(ctx, phoneNumber, serviceRequestID, verificationStatus)
}

// GetServiceRequests is the resolver for the getServiceRequests field.
func (r *queryResolver) GetServiceRequests(ctx context.Context, requestType *string, requestStatus *string, facilityID string, flavour feedlib.Flavour) ([]*domain.ServiceRequest, error) {
	return r.mycarehub.ServiceRequest.GetServiceRequests(ctx, requestType, requestStatus, facilityID, flavour)
}

// GetPendingServiceRequestsCount is the resolver for the getPendingServiceRequestsCount field.
func (r *queryResolver) GetPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCountResponse, error) {
	return r.mycarehub.ServiceRequest.GetPendingServiceRequestsCount(ctx, facilityID)
}

// SearchServiceRequests is the resolver for the searchServiceRequests field.
func (r *queryResolver) SearchServiceRequests(ctx context.Context, searchTerm string, flavour feedlib.Flavour, requestType string, facilityID string) ([]*domain.ServiceRequest, error) {
	return r.mycarehub.ServiceRequest.SearchServiceRequests(ctx, searchTerm, flavour, requestType, facilityID)
}
