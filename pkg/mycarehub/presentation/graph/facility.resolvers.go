package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.26

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// DeleteFacility is the resolver for the deleteFacility field.
func (r *mutationResolver) DeleteFacility(ctx context.Context, identifier dto.FacilityIdentifierInput) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.DeleteFacility(ctx, &identifier)
}

// ReactivateFacility is the resolver for the reactivateFacility field.
func (r *mutationResolver) ReactivateFacility(ctx context.Context, identifier dto.FacilityIdentifierInput) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.ReactivateFacility(ctx, &identifier)
}

// InactivateFacility is the resolver for the inactivateFacility field.
func (r *mutationResolver) InactivateFacility(ctx context.Context, identifier dto.FacilityIdentifierInput) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.InactivateFacility(ctx, &identifier)
}

// AddFacilityContact is the resolver for the addFacilityContact field.
func (r *mutationResolver) AddFacilityContact(ctx context.Context, facilityID string, contact string) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.AddFacilityContact(ctx, facilityID, contact)
}

// AddFacilityToProgram is the resolver for the addFacilityToProgram field.
func (r *mutationResolver) AddFacilityToProgram(ctx context.Context, facilityIDs []string, programID string) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.AddFacilityToProgram(ctx, facilityIDs, programID)
}

// ListFacilities is the resolver for the listFacilities field.
func (r *queryResolver) ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationInput dto.PaginationsInput) (*domain.FacilityPage, error) {
	return r.mycarehub.Facility.ListFacilities(ctx, searchTerm, filterInput, &paginationInput)
}

// RetrieveFacility is the resolver for the retrieveFacility field.
func (r *queryResolver) RetrieveFacility(ctx context.Context, id string, active bool) (*domain.Facility, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.RetrieveFacility(ctx, &id, active)
}

// RetrieveFacilityByIdentifier is the resolver for the retrieveFacilityByIdentifier field.
func (r *queryResolver) RetrieveFacilityByIdentifier(ctx context.Context, identifier dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.RetrieveFacilityByIdentifier(ctx, &identifier, isActive)
}

// ListProgramFacilities is the resolver for the listProgramFacilities field.
func (r *queryResolver) ListProgramFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationInput dto.PaginationsInput) (*domain.FacilityPage, error) {
	return r.mycarehub.Facility.ListProgramFacilities(ctx, searchTerm, filterInput, &paginationInput)
}
