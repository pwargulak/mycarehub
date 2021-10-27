package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
)

//GetFacilities returns a slice of healthcare facilities in the platform.
func (d *OnboardingDb) GetFacilities(ctx context.Context) ([]*domain.Facility, error) {
	var facility []*domain.Facility
	facilities, err := d.query.GetFacilities(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get facilities: %v", err)
	}

	if len(facilities) == 0 {
		return facility, nil
	}
	for _, m := range facilities {
		active, err := strconv.ParseBool(m.Active)
		if err != nil {
			return nil, fmt.Errorf("failed to parse facility.Active to boolean")
		}
		singleFacility := domain.Facility{
			ID:          m.FacilityID,
			Name:        m.Name,
			Code:        m.Code,
			Active:      active,
			County:      m.County,
			Description: m.Description,
		}

		facility = append(facility, &singleFacility)
	}

	return facility, nil
}

// RetrieveFacility gets a facility by ID from the database
func (d *OnboardingDb) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	if id == nil {
		return nil, fmt.Errorf("facility ID should be defined")
	}
	facilitySession, err := d.query.RetrieveFacility(ctx, id, isActive)
	if err != nil {
		return nil, fmt.Errorf("failed query and retrieve one facility: %s", err)
	}

	return d.mapFacilityObjectToDomain(facilitySession), nil
}

// RetrieveByFacilityMFLCode gets a facility by ID from the database
func (d *OnboardingDb) RetrieveByFacilityMFLCode(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
	if MFLCode == "" {
		return nil, fmt.Errorf("facility ID should be defined")
	}
	facilitySession, err := d.query.RetrieveFacilityByMFLCode(ctx, MFLCode, isActive)
	if err != nil {
		return nil, fmt.Errorf("failed query and retrieve facility by MFLCode: %s", err)
	}

	return d.mapFacilityObjectToDomain(facilitySession), nil
}

// GetUserProfileByUserID fetches a user profile facility using the user ID
func (d *OnboardingDb) GetUserProfileByUserID(ctx context.Context, userID string, flavour feedlib.Flavour) (*domain.User, error) {
	user, err := d.query.GetUserProfileByUserID(ctx, userID, flavour)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user profile by user ID: %s", err)
	}

	return d.mapProfileObjectToDomain(user), nil
}

// GetUserPINByUserID fetches a user profile facility using the user ID
func (d *OnboardingDb) GetUserPINByUserID(ctx context.Context, userID string) (*domain.UserPIN, error) {
	pinData, err := d.query.GetUserPINByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed query and retrieve user PIN data by user ID: %s", err)
	}

	return d.mapPINObjectToDomain(pinData), nil
}

// GetClientProfileByClientID retrieves a client profile using the client ID
func (d *OnboardingDb) GetClientProfileByClientID(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
	client, err := d.query.GetClientProfileByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client profile by client ID %v: %v", clientID, err)
	}

	return d.mapClientObjectToDomain(client), err
}

// GetStaffProfileByStaffID retrieves a staff profile by staffProfileID
func (d *OnboardingDb) GetStaffProfileByStaffID(ctx context.Context, staffProfileID string) (*domain.StaffUserProfile, error) {
	staffUserObject, err := d.query.GetStaffProfileByStaffID(ctx, staffProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get StaffProfile by staffProfileID %v: %v", staffProfileID, err)
	}
	return d.mapStaffProfileUserObjectToDomain(staffUserObject), nil

}

// GetStaffProfileByStaffNumber retrieves a staff profile by staffNumber
func (d *OnboardingDb) GetStaffProfileByStaffNumber(ctx context.Context, staffNumber string) (*domain.StaffUserProfile, error) {
	staffUserObject, err := d.query.GetStaffProfileByStaffNumber(ctx, staffNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get StaffProfile by staffNumber %v: %v", staffNumber, err)
	}
	return d.mapStaffProfileUserObjectToDomain(staffUserObject), nil
}
