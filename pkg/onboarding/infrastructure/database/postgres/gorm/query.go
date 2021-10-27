package gorm

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/savannahghi/feedlib"
)

// Query contains all the db query methods
type Query interface {
	RetrieveFacility(ctx context.Context, id *string, isActive bool) (*Facility, error)
	RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*Facility, error)
	GetFacilities(ctx context.Context) ([]Facility, error)
	GetUserProfileByUserID(ctx context.Context, userID string, flavour feedlib.Flavour) (*User, error)
	GetUserPINByUserID(ctx context.Context, userID string) (*PINData, error)
	GetClientProfileByClientID(ctx context.Context, clientID string) (*ClientProfile, error)
	GetStaffProfileByStaffID(ctx context.Context, staffProfileID string) (*StaffUserProfile, error)
	GetStaffProfileByStaffNumber(ctx context.Context, staffNumber string) (*StaffUserProfile, error)
}

// RetrieveFacility fetches a single facility
func (db *PGInstance) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*Facility, error) {
	var facility Facility
	active := strconv.FormatBool(isActive)
	err := db.DB.Where(&Facility{FacilityID: id, Active: active}).Find(&facility).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get facility by ID %v: %v", id, err)
	}
	return &facility, nil
}

// RetrieveFacilityByMFLCode fetches a single facility using MFL Code
func (db *PGInstance) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*Facility, error) {
	var facility Facility
	active := strconv.FormatBool(isActive)
	if err := db.DB.Where(&Facility{Code: MFLCode, Active: active}).First(&facility).Error; err != nil {
		return nil, fmt.Errorf("failed to get facility by MFL Code %v and status %v: %v", MFLCode, active, err)
	}
	return &facility, nil
}

// GetUserProfileByUserID fetches a user profile facility using the user ID
func (db *PGInstance) GetUserProfileByUserID(ctx context.Context, userID string, flavour feedlib.Flavour) (*User, error) {
	var user User
	if err := db.DB.Where(&User{UserID: &userID, Flavour: flavour}).Preload("Contacts").First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user by userID %v: %v", userID, err)
	}
	return &user, nil
}

// GetUserPINByUserID fetches a user profile facility using the user ID
func (db *PGInstance) GetUserPINByUserID(ctx context.Context, userID string) (*PINData, error) {
	var pin PINData
	if err := db.DB.Where(&PINData{UserID: userID}).First(&pin).Error; err != nil {
		return nil, fmt.Errorf("failed to get facility by MFL Code %v: %v", userID, err)
	}
	return &pin, nil
}

// GetFacilities fetches all the healthcare facilities in the platform.
func (db *PGInstance) GetFacilities(ctx context.Context) ([]Facility, error) {
	var facility []Facility
	facilities := db.DB.Find(&facility).Error
	log.Printf("these are the facilities %v", facilities)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to query all facilities %v", err)
	// }
	log.Printf("these are the facilities %v", facility)
	return facility, nil
}

// GetClientProfileByClientID retrieves a client profile by ID
func (db *PGInstance) GetClientProfileByClientID(ctx context.Context, clientID string) (*ClientProfile, error) {
	var client ClientProfile
	if err := db.DB.Where(&ClientProfile{ID: &clientID}).First(&client).Error; err != nil {
		return nil, fmt.Errorf("failed to get client profile by client ID %v: %v", clientID, err)
	}

	return &client, nil
}

// GetStaffProfileByStaffID retrieves a staff profile by staffProfileID
func (db *PGInstance) GetStaffProfileByStaffID(ctx context.Context, staffProfileID string) (*StaffUserProfile, error) {
	var staff StaffProfile
	var user User
	if err := db.DB.Where(&StaffProfile{StaffProfileID: &staffProfileID}).Preload("Addresses").First(&staff).Error; err != nil {
		return nil, fmt.Errorf("failed to get a staff profile by ID %v : %v", &staffProfileID, err)
	}
	if err := db.DB.Where(&User{UserID: staff.UserID}).Preload("Contacts").First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get a staff user: %v", err)
	}
	return &StaffUserProfile{
		User:  &user,
		Staff: &staff,
	}, nil

}

// GetStaffProfileByStaffNumber retrieves a staff profile by staffNumber
func (db *PGInstance) GetStaffProfileByStaffNumber(ctx context.Context, staffNumber string) (*StaffUserProfile, error) {
	var staff StaffProfile
	var user User
	if err := db.DB.Where(&StaffProfile{StaffNumber: staffNumber}).Preload("Addresses").First(&staff).Error; err != nil {
		return nil, fmt.Errorf("failed to get a staff profile by ID %v : %v", staffNumber, err)
	}
	if err := db.DB.Where(&User{UserID: staff.UserID}).Preload("Contacts").First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get a staff user: %v", err)
	}
	return &StaffUserProfile{
		User:  &user,
		Staff: &staff,
	}, nil
}
