package facility_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	pubsubMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility/mock"
)

func TestUseCaseFacilityImpl_RetrieveFacility_Unittest(t *testing.T) {
	ctx := context.Background()

	ID := uuid.New().String()

	type args struct {
		ctx      context.Context
		id       *string
		isActive bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: retrieve facility",
			args: args{
				ctx:      ctx,
				id:       &ID,
				isActive: true,
			},
			wantErr: false,
		},

		{
			name: "Sad case - no id",
			args: args{
				ctx:      ctx,
				isActive: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeFacility := mock.NewFacilityUsecaseMock()

			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt)

			if tt.name == "Sad case - no id" {
				fakeFacility.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("an error occurred while retrieving facility")
				}
			}
			got, err := f.RetrieveFacility(tt.args.ctx, tt.args.id, tt.args.isActive)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.RetrieveFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCaseFacilityImpl_RetrieveFacilityByIdentifier_Unittest(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx        context.Context
		identifier dto.FacilityIdentifierInput
		isActive   bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: retrieve facility by identifier",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
				isActive: true,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to retrieve facility by identifier",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
				isActive: false,
			},
			wantErr: true,
		},
		{
			name: "Sad case: missing input",
			args: args{
				ctx:      ctx,
				isActive: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()
			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt)

			if tt.name == "Sad case: unable to retrieve facility by identifier" {
				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("an error occurred while retrieving facility by Identifier")
				}
			}

			got, err := f.RetrieveFacilityByIdentifier(tt.args.ctx, &tt.args.identifier, tt.args.isActive)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.RetrieveFacilityByIdentifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected facilities to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUnit_ListProgramFacilities(t *testing.T) {
	ctx := context.Background()
	searchTerm := "term"

	paginationInput := dto.PaginationsInput{
		Limit:       1,
		CurrentPage: 1,
	}

	type args struct {
		ctx              context.Context
		searchTerm       *string
		filterInput      []*dto.FiltersInput
		paginationsInput *dto.PaginationsInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: list facilities",
			args: args{
				ctx:        ctx,
				searchTerm: &searchTerm,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeName,
						Value:    "value",
					},
				},
				paginationsInput: &paginationInput,
			},
			wantErr: false,
		},
		{
			name: "Happy case- empty search input",
			args: args{
				ctx:              ctx,
				filterInput:      nil,
				paginationsInput: &paginationInput,
			},
			wantErr: false,
		},
		{
			name: "Sad case- failed to get logged in user",
			args: args{
				ctx:        ctx,
				searchTerm: &searchTerm,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeName,
						Value:    "value",
					},
				},
				paginationsInput: &paginationInput,
			},
			wantErr: true,
		},
		{
			name: "Sad case- failed to get user profile",
			args: args{
				ctx:        ctx,
				searchTerm: &searchTerm,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeName,
						Value:    "value",
					},
				},
				paginationsInput: &paginationInput,
			},
			wantErr: true,
		},
		{
			name: "Sad case- failed to list facilities",
			args: args{
				ctx:        ctx,
				searchTerm: &searchTerm,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeName,
						Value:    "value",
					},
				},
				paginationsInput: &paginationInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt)

			if tt.name == "Sad case- failed to get logged in user" {
				fakeExt.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case- failed to get user profile" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case- failed to list facilities" {
				fakeDB.MockListProgramFacilitiesFn = func(ctx context.Context, programID, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("failed to list facilities")
				}
			}

			got, err := f.ListProgramFacilities(tt.args.ctx, tt.args.searchTerm, tt.args.filterInput, tt.args.paginationsInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.ListProgramFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected facilities to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCaseFacilityImpl_Inactivate_Unittest(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx        context.Context
		identifier dto.FacilityIdentifierInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - empty mflCode",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - invalid mflCode",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - very bad mflCode",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt)

			if tt.name == "Sad Case - empty mflCode" {
				fakeDB.MockInactivateFacilityFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - invalid mflCode" {
				fakeDB.MockInactivateFacilityFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - very bad mflCode" {
				fakeDB.MockInactivateFacilityFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := f.InactivateFacility(tt.args.ctx, &tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.Inactivate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("PGInstance.InactivateFacility() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestUseCaseFacilityImpl_Reactivate_Unittest(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx        context.Context
		identifier dto.FacilityIdentifierInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - empty mflCode",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - invalid mflCode",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - very bad mflCode",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt)

			if tt.name == "Sad Case - empty mflCode" {
				fakeDB.MockReactivateFacilityFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - invalid mflCode" {
				fakeDB.MockReactivateFacilityFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - very bad mflCode" {
				fakeDB.MockReactivateFacilityFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := f.ReactivateFacility(tt.args.ctx, &tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.ReactivateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("PGInstance.ReactivateFacility() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestUseCaseFacilityImpl_DeleteFacility(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx        context.Context
		identifier dto.FacilityIdentifierInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully delete facility",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - unable delete facility",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeFacility := mock.NewFacilityUsecaseMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt)

			if tt.name == "Happy Case - Successfully delete facility" {
				fakeFacility.MockDeleteFacilityFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad Case - unable delete facility" {
				fakeFacility.MockDeleteFacilityFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := f.DeleteFacility(tt.args.ctx, &tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.DeleteFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCaseFacilityImpl.DeleteFacility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCaseFacilityImpl_ListFacilities(t *testing.T) {
	ctx := context.Background()
	searchTerm := "term"

	paginationInput := dto.PaginationsInput{
		Limit:       1,
		CurrentPage: 1,
	}

	type args struct {
		ctx              context.Context
		searchTerm       *string
		filterInput      []*dto.FiltersInput
		paginationsInput *dto.PaginationsInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: list facilities",
			args: args{
				ctx:        ctx,
				searchTerm: &searchTerm,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeName,
						Value:    "value",
					},
				},
				paginationsInput: &paginationInput,
			},
			wantErr: false,
		},
		{
			name: "Happy case- empty search input",
			args: args{
				ctx:              ctx,
				filterInput:      nil,
				paginationsInput: &paginationInput,
			},
			wantErr: false,
		},

		{
			name: "Sad case- failed to list facilities",
			args: args{
				ctx:        ctx,
				searchTerm: &searchTerm,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeName,
						Value:    "value",
					},
				},
				paginationsInput: &paginationInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt)

			if tt.name == "Sad case- failed to list facilities" {
				fakeDB.MockListProgramFacilitiesFn = func(ctx context.Context, programID, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("failed to list facilities")
				}
			}

			got, err := f.ListProgramFacilities(tt.args.ctx, tt.args.searchTerm, tt.args.filterInput, tt.args.paginationsInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.ListProgramFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected facilities to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCaseFacilityImpl_SyncFacilities(t *testing.T) {
	ctx := context.Background()
	fakeDB := pgMock.NewPostgresMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeExt := extensionMock.NewFakeExtension()

	f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx: ctx,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - unable to notify pubsub",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - unable to get facilities without FHIROrganisationID",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad Case - unable to notify pubsub" {
				fakePubsub.MockNotifyCreateOrganizationFn = func(ctx context.Context, facility *domain.Facility) error {
					return fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad Case - unable to get facilities without FHIROrganisationID" {
				fakeDB.MockGetFacilitiesWithoutFHIRIDFn = func(ctx context.Context) ([]*domain.Facility, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			err := f.SyncFacilities(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.SyncFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCaseFacilityImpl_UpdateFacility(t *testing.T) {
	ctx := context.Background()
	fakeDB := pgMock.NewPostgresMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeExt := extensionMock.NewFakeExtension()

	f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt)

	UUID := uuid.New().String()

	type args struct {
		ctx                context.Context
		updateFacilityData *dto.UpdateFacilityPayload
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx: ctx,
				updateFacilityData: &dto.UpdateFacilityPayload{
					FacilityID:         UUID,
					FHIROrganisationID: UUID,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case",
			args: args{
				ctx: ctx,
				updateFacilityData: &dto.UpdateFacilityPayload{
					FacilityID:         UUID,
					FHIROrganisationID: UUID,
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Sad Case" {
				fakeDB.MockUpdateFacilityFn = func(ctx context.Context, facility *domain.Facility, updateData map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}
			err := f.UpdateFacility(tt.args.ctx, tt.args.updateFacilityData)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.UpdateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCaseFacilityImpl_AddFacilityContact(t *testing.T) {

	type args struct {
		ctx        context.Context
		facilityID string
		contact    string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: success adding facility contact",
			args: args{
				ctx:        context.Background(),
				facilityID: gofakeit.UUID(),
				contact:    interserviceclient.TestUserPhoneNumber,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: fail to normalize phone number",
			args: args{
				ctx:        context.Background(),
				facilityID: gofakeit.UUID(),
				contact:    "072897",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: fail to update facility",
			args: args{
				ctx:        context.Background(),
				facilityID: gofakeit.UUID(),
				contact:    interserviceclient.TestUserPhoneNumber,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt)

			if tt.name == "sad case: fail to update facility" {
				fakeDB.MockUpdateFacilityFn = func(ctx context.Context, facility *domain.Facility, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update facility")
				}
			}

			got, err := f.AddFacilityContact(tt.args.ctx, tt.args.facilityID, tt.args.contact)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.AddFacilityContact() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCaseFacilityImpl.AddFacilityContact() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCaseFacilityImpl_CreateFacilities(t *testing.T) {
	type args struct {
		ctx        context.Context
		facilities []*domain.Facility
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create facilities",
			args: args{
				ctx: context.Background(),
				facilities: []*domain.Facility{
					{
						Name:        gofakeit.BS(),
						Phone:       "0777777777",
						Active:      true,
						Country:     "Kenya",
						Description: gofakeit.BS(),
						Identifier: domain.FacilityIdentifier{
							Active: true,
							Type:   enums.FacilityIdentifierTypeMFLCode,
							Value:  "392893828",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to create facilities",
			args: args{
				ctx: context.Background(),
				facilities: []*domain.Facility{
					{
						Name:        gofakeit.BS(),
						Phone:       "0999999999",
						Active:      true,
						Country:     "Kenya",
						Description: gofakeit.BS(),
						Identifier: domain.FacilityIdentifier{
							Active: true,
							Type:   enums.FacilityIdentifierTypeMFLCode,
							Value:  "09090908",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to notify create facility in clinical service",
			args: args{
				ctx: context.Background(),
				facilities: []*domain.Facility{
					{
						Name:        gofakeit.BS(),
						Phone:       "0999999999",
						Active:      true,
						Country:     "Kenya",
						Description: gofakeit.BS(),
						Identifier: domain.FacilityIdentifier{
							Active: true,
							Type:   enums.FacilityIdentifierTypeMFLCode,
							Value:  "09090908",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to notify create facility in cms",
			args: args{
				ctx: context.Background(),
				facilities: []*domain.Facility{
					{
						Name:        gofakeit.BS(),
						Phone:       "0999999999",
						Active:      true,
						Country:     "Kenya",
						Description: gofakeit.BS(),
						Identifier: domain.FacilityIdentifier{
							Active: true,
							Type:   enums.FacilityIdentifierTypeMFLCode,
							Value:  "09090908",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()
			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt)

			if tt.name == "Sad case: failed to create facilities" {
				fakeDB.MockCreateFacilitiesFn = func(ctx context.Context, facilities []*domain.Facility) ([]*domain.Facility, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to notify create facility in clinical service" {
				fakePubsub.MockNotifyCreateOrganizationFn = func(ctx context.Context, facility *domain.Facility) error {
					return fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: failed to notify create facility in cms" {
				fakePubsub.MockNotifyCreateCMSFacilityFn = func(ctx context.Context, facility *dto.CreateCMSFacilityPayload) error {
					return fmt.Errorf("an error occurred")
				}
			}

			got, err := f.CreateFacilities(tt.args.ctx, tt.args.facilities)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.CreateFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a a value to be returned, got: %v", got)
			}
		})
	}
}

func TestUseCaseFacilityImpl_PublishFacilitiesToCMS(t *testing.T) {
	id := gofakeit.UUID()
	id2 := gofakeit.UUID()
	type args struct {
		ctx        context.Context
		facilities []*domain.Facility
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: publish facilities to cms",
			args: args{
				ctx: context.Background(),
				facilities: []*domain.Facility{
					{
						ID:   &id,
						Name: gofakeit.BS(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed publish facilities to cms",
			args: args{
				ctx: context.Background(),
				facilities: []*domain.Facility{
					{
						ID:   &id2,
						Name: gofakeit.BS(),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt)

			if tt.name == "Sad case: failed publish facilities to cms" {
				fakePubsub.MockNotifyCreateCMSFacilityFn = func(ctx context.Context, facility *dto.CreateCMSFacilityPayload) error {
					return fmt.Errorf("an error occurred")
				}
			}

			if err := f.PublishFacilitiesToCMS(tt.args.ctx, tt.args.facilities); (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.PublishFacilitiesToCMS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCaseFacilityImpl_CmdAddFacilityToProgram(t *testing.T) {

	type args struct {
		ctx         context.Context
		facilityIDs []string
		programID   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: add facilities to program",
			args: args{
				ctx:         context.Background(),
				facilityIDs: []string{gofakeit.UUID()},
				programID:   gofakeit.UUID(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: fail to add facility to program",
			args: args{
				ctx:         context.Background(),
				facilityIDs: []string{gofakeit.UUID()},
				programID:   gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: failed to notify cms",
			args: args{
				ctx:         context.Background(),
				facilityIDs: []string{gofakeit.UUID()},
				programID:   gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		fakeDB := pgMock.NewPostgresMock()
		fakePubsub := pubsubMock.NewPubsubServiceMock()
		fakeExt := extensionMock.NewFakeExtension()
		f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt)
		if tt.name == "sad case: fail to add facility to program" {
			fakeDB.MockAddFacilityToProgramFn = func(ctx context.Context, programID string, facilityIDs []string) ([]*domain.Facility, error) {
				return nil, fmt.Errorf("failed to add facility to program")
			}
		}
		if tt.name == "sad case: failed to notify cms" {
			fakePubsub.MockNotifyCMSAddFacilityToProgramFn = func(ctx context.Context, payload *dto.CMSLinkFacilityToProgramPayload) error {
				return fmt.Errorf("failed to notify cms")
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.AddFacilityToProgram(tt.args.ctx, tt.args.facilityIDs, tt.args.programID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.CmdAddFacilityToProgram() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCaseFacilityImpl.CmdAddFacilityToProgram() = %v, want %v", got, tt.want)
			}
		})
	}
}
