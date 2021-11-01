package infrastructure

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// Create represents a contract that contains all `create` ops to the database
//
// All the  contracts for create operations are assembled here
type Create interface {
	GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
	RegisterClient(
		ctx context.Context,
		userInput *dto.UserInput,
		clientInput *dto.ClientProfileInput,
	) (*domain.ClientUserProfile, error)
	SavePin(ctx context.Context, pinInput *domain.UserPIN) (bool, error)
}

// Delete represents all the deletion action interfaces
type Delete interface {
	DeleteFacility(ctx context.Context, id string) (bool, error)
}

// Query contains all query methods
type Query interface {
	RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	GetFacilities(ctx context.Context) ([]*domain.Facility, error)
	RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error)
	GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error)
	ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, PaginationsInput dto.PaginationsInput) (*domain.FacilityPage, error)
	GetUserPINByUserID(ctx context.Context, userID string) (*domain.UserPIN, error)
}
