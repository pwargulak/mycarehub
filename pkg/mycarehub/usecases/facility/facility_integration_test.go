package facility_test

import (
	"context"
	"log"
	"strconv"
	"testing"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/segmentio/ksuid"
	"github.com/tj/assert"
)

func TestUseCaseFacilityImpl_CreateFacility(t *testing.T) {
	f := testInteractor
	ctx := context.Background()
	name := "Kanairo One"
	code := ksuid.New().String()
	county := enums.CountyTypeNairobi
	description := "This is just for mocking"

	type args struct {
		ctx      context.Context
		facility dto.FacilityInput
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "happy case - valid payload",
			args: args{
				ctx: ctx,
				facility: dto.FacilityInput{
					Name:        name,
					Code:        code,
					Active:      true,
					County:      county,
					Description: description,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case - facility code not defined",
			args: args{
				ctx: ctx,
				facility: dto.FacilityInput{
					Name:        name,
					Active:      true,
					County:      county,
					Description: description,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := f.FacilityUsecase.GetOrCreateFacility(tt.args.ctx, &tt.args.facility)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.GetOrCreateFacility() error = %v, wantErr %v", err, tt.wantErr)
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
	// Teardown
	_, err := f.FacilityUsecase.DeleteFacility(ctx, code)
	if err != nil {
		log.Printf("failed to delete facility: %v", err)
	}

}

func TestUseCaseFacilityImpl_RetrieveFacility_Integration(t *testing.T) {
	ctx := context.Background()

	f := testInteractor

	Code := "KN002"

	facilityInput := &dto.FacilityInput{
		Name:        "Kanairo One",
		Code:        Code,
		Active:      true,
		County:      enums.CountyTypeNairobi,
		Description: "This is just for mocking",
	}

	// Setup, create a facility
	facility, err := f.FacilityUsecase.GetOrCreateFacility(ctx, facilityInput)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
	}

	ID := facility.ID
	active := facility.Active

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
			name: "Happy case",
			args: args{
				ctx:      ctx,
				id:       ID,
				isActive: active,
			},
			wantErr: false,
		},

		{
			name: "Sad case - nil ID",
			args: args{
				ctx:      ctx,
				id:       nil,
				isActive: false,
			},
			wantErr: true,
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
			got, err := f.FacilityUsecase.RetrieveFacility(tt.args.ctx, tt.args.id, tt.args.isActive)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.RetrieveFacility() error = %v, wantErr %v", err, tt.wantErr)
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

	_, err = f.FacilityUsecase.DeleteFacility(ctx, string(facility.Code))
	if err != nil {
		t.Errorf("unable to delete facility")
		return
	}
}

func TestUseCaseFacilityImpl_DeleteFacility_Integrationtest(t *testing.T) {
	ctx := context.Background()

	u := testInteractor

	//Create facility
	facilityInput := &dto.FacilityInput{
		Name:        "Kanairo One",
		Code:        ksuid.New().String(),
		County:      enums.CountyTypeNairobi,
		Active:      true,
		Description: "This is just for integration testing",
	}

	// create a facility
	facility, err := u.FacilityUsecase.GetOrCreateFacility(ctx, facilityInput)
	assert.Nil(t, err)
	assert.NotNil(t, facility)

	// retrieve the facility
	facility1, err := u.FacilityUsecase.RetrieveFacility(ctx, facility.ID, true)
	assert.Nil(t, err)
	assert.NotNil(t, facility1)

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				id:  *facility1.ID,
			},
			wantErr: false,
		},

		{
			name: "Sad case - bad id",
			args: args{
				ctx: ctx,
				id:  ksuid.New().String(),
			},
			wantErr: false,
		},

		{
			name: "Sad case - empty id",
			args: args{
				ctx: ctx,
				id:  "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := u.FacilityUsecase.DeleteFacility(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.DeleteFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCaseFacilityImpl_RetrieveFacilityByMFLCode_Integration(t *testing.T) {
	ctx := context.Background()

	f := testInteractor

	facilityInput := &dto.FacilityInput{
		Name:        "Kanairo One",
		Code:        "KN001",
		Active:      true,
		County:      enums.CountyTypeNairobi,
		Description: "This is just for mocking",
	}

	// Setup, create a facility
	facility, err := f.FacilityUsecase.GetOrCreateFacility(ctx, facilityInput)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
	}

	mflCode := facility.Code

	invalidMFLCode := ksuid.New().String()

	type args struct {
		ctx      context.Context
		MFLCode  string
		isActive bool
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Facility
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:      ctx,
				MFLCode:  mflCode,
				isActive: true,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:      ctx,
				MFLCode:  invalidMFLCode,
				isActive: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.FacilityUsecase.RetrieveFacilityByMFLCode(tt.args.ctx, tt.args.MFLCode, tt.args.isActive)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.RetrieveFacilityByMFLCode() error = %v, wantErr %v", err, tt.wantErr)
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
	// Teardown
	_, err = f.FacilityUsecase.DeleteFacility(ctx, string(facility.Code))
	if err != nil {
		t.Errorf("unable to delete facility")
		return
	}
}

func TestUseCaseFacilityImpl_ListFacilities(t *testing.T) {
	ctx := context.Background()

	f := testInteractor

	code := ksuid.New().String()
	code2 := ksuid.New().String()

	facilityInput := &dto.FacilityInput{
		Name:        "Kanairo One",
		Code:        code,
		Active:      true,
		County:      enums.CountyTypeNairobi,
		Description: "This is just for mocking",
	}

	facilityInput2 := &dto.FacilityInput{
		Name:        "Baringo 2",
		Code:        code2,
		Active:      true,
		County:      enums.CountyTypeBaringo,
		Description: "This is just for mocking",
	}

	noSearchTerm := ""
	searchTerm := "term"

	noFilterInput := []*dto.FiltersInput{}

	formatBool := strconv.FormatBool(true)

	filterInput := []*dto.FiltersInput{
		{
			DataType: enums.FilterDataTypeName,
			Value:    "Kanairo One",
		},
		{
			DataType: enums.FilterDataTypeMFLCode,
			Value:    code,
		},
		{
			DataType: enums.FilterDataTypeActive,
			Value:    formatBool,
		},
		{
			DataType: enums.FilterDataTypeCounty,
			Value:    enums.CountyTypeNairobi.String(),
		},
	}

	filterEmptyName := []*dto.FiltersInput{
		{
			DataType: enums.FilterDataTypeName,
			Value:    "",
		},
		{
			DataType: enums.FilterDataTypeMFLCode,
			Value:    code,
		},
		{
			DataType: enums.FilterDataTypeActive,
			Value:    formatBool,
		},
		{
			DataType: enums.FilterDataTypeCounty,
			Value:    enums.CountyTypeNairobi.String(),
		},
	}
	filterEmptyMFLCode := []*dto.FiltersInput{
		{
			DataType: enums.FilterDataTypeName,
			Value:    "Kanairo One",
		},
		{
			DataType: enums.FilterDataTypeMFLCode,
			Value:    "",
		},
		{
			DataType: enums.FilterDataTypeActive,
			Value:    formatBool,
		},
		{
			DataType: enums.FilterDataTypeCounty,
			Value:    enums.CountyTypeNairobi.String(),
		},
	}

	filterInvalidBool := []*dto.FiltersInput{
		{
			DataType: enums.FilterDataTypeName,
			Value:    "Kanairo One",
		},
		{
			DataType: enums.FilterDataTypeMFLCode,
			Value:    code,
		},
		{
			DataType: enums.FilterDataTypeActive,
			Value:    "invalid",
		},
		{
			DataType: enums.FilterDataTypeCounty,
			Value:    enums.CountyTypeNairobi.String(),
		},
	}

	filterInvalidCounty := []*dto.FiltersInput{
		{
			DataType: enums.FilterDataTypeName,
			Value:    "Kanairo One",
		},
		{
			DataType: enums.FilterDataTypeMFLCode,
			Value:    code,
		},
		{
			DataType: enums.FilterDataTypeActive,
			Value:    formatBool,
		},
		{
			DataType: enums.FilterDataTypeCounty,
			Value:    "kanairo",
		},
	}

	paginationInput := dto.PaginationsInput{
		Limit:       1,
		CurrentPage: 1,
	}
	paginationInputNoCurrentPage := dto.PaginationsInput{
		Limit: 1,
	}

	// Setup
	// create a facility
	facility, err := f.FacilityUsecase.GetOrCreateFacility(ctx, facilityInput)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
	}
	// Create another Facility
	facility2, err := f.FacilityUsecase.GetOrCreateFacility(ctx, facilityInput2)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
	}

	type args struct {
		ctx              context.Context
		searchTerm       *string
		filterInput      []*dto.FiltersInput
		PaginationsInput dto.PaginationsInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FacilityPage
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:              ctx,
				searchTerm:       &noSearchTerm,
				filterInput:      noFilterInput,
				PaginationsInput: paginationInput,
			},
			wantErr: false,
		},

		{
			name: "valid: with valid filters",
			args: args{
				ctx:              ctx,
				searchTerm:       &noSearchTerm,
				filterInput:      filterInput,
				PaginationsInput: paginationInput,
			},
			wantErr: false,
		},

		{
			name: "invalid: no params passed",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing current page",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterInput,
				PaginationsInput: paginationInputNoCurrentPage,
			},
			wantErr: true,
		},
		{
			name: "invalid: empty name passed",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterEmptyName,
				PaginationsInput: paginationInput,
			},
			wantErr: true,
		},
		{
			name: "invalid: empty MFL code",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterEmptyMFLCode,
				PaginationsInput: paginationInput,
			},
			wantErr: true,
		},
		{
			name: "invalid: invalid bool",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterInvalidBool,
				PaginationsInput: paginationInput,
			},
			wantErr: true,
		},

		{
			name: "invalid: invalid county",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterInvalidCounty,
				PaginationsInput: paginationInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := f.FacilityUsecase.ListFacilities(tt.args.ctx, tt.args.searchTerm, tt.args.filterInput, tt.args.PaginationsInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.ListFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// assert we get at least one value in the filter
			if tt.name == "valid: with valid filters" {
				assert.GreaterOrEqual(t, len(got.Facilities), 1)
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected facilities to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, got.Pagination.Limit, tt.args.PaginationsInput.Limit)
				assert.Equal(t, got.Pagination.Limit, len(got.Facilities))
			}

		})
	}
	// Teardown
	_, err = f.FacilityUsecase.DeleteFacility(ctx, string(facility.Code))
	if err != nil {
		t.Errorf("unable to delete facility")
		return
	}
	_, err = f.FacilityUsecase.DeleteFacility(ctx, string(facility2.Code))
	if err != nil {
		t.Errorf("unable to delete facility")
		return
	}
}
