package facility_test

// import (
// 	"context"
// 	"fmt"
// 	"strconv"
// 	"testing"

// 	"github.com/google/uuid"
// 	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
// 	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
// 	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
// 	"github.com/segmentio/ksuid"
// )

// func TestUnit_CreateFacility(t *testing.T) {
// 	f := testFakeInfrastructureInteractor
// 	ctx := context.Background()
// 	name := "Kanairo One"
// 	code := ksuid.New().String()
// 	county := enums.CountyTypeNairobi
// 	description := "This is just for mocking"

// 	type args struct {
// 		ctx      context.Context
// 		facility dto.FacilityInput
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantNil bool
// 		wantErr bool
// 	}{
// 		{
// 			name: "happy case - valid payload",
// 			args: args{
// 				ctx: ctx,
// 				facility: dto.FacilityInput{
// 					Name:        name,
// 					Code:        code,
// 					Active:      true,
// 					County:      county,
// 					Description: description,
// 				},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "sad case - facility code not defined",
// 			args: args{
// 				ctx: ctx,
// 				facility: dto.FacilityInput{
// 					Name:        name,
// 					Active:      true,
// 					County:      county,
// 					Description: description,
// 				},
// 			},
// 			wantErr: true,
// 			wantNil: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "sad case - facility code not defined" {
// 				fakeCreate.GetOrCreateFacilityFn = func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
// 					return nil, fmt.Errorf("failed to create facility")
// 				}
// 			}
// 			if tt.name == "happy case - valid payload" {
// 				fakeQuery.RetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
// 					return nil, fmt.Errorf("failed query and retrieve facility by MFLCode")
// 				}

// 				fakeCreate.GetOrCreateFacilityFn = func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
// 					return &domain.Facility{
// 						Name:        facility.Name,
// 						Code:        facility.Code,
// 						Active:      facility.Active,
// 						County:      facility.County,
// 						Description: facility.Description,
// 					}, nil
// 				}
// 			}

// 			got, err := f.GetOrCreateFacility(tt.args.ctx, tt.args.facility)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("UseCaseFacilityImpl.GetOrCreateFacility() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("OnboardingDb.GetOrCreateFacility() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if tt.wantErr && got != nil {
// 				t.Errorf("expected facility to be nil for %v", tt.name)
// 				return
// 			}

// 			if !tt.wantErr && got == nil {
// 				t.Errorf("expected facility not to be nil for %v", tt.name)
// 				return
// 			}
// 		})
// 	}
// 	// TODO: add teardown
// }

// func TestUseCaseFacilityImpl_RetrieveFacility_Unittest(t *testing.T) {
// 	ctx := context.Background()

// 	f := testFakeInfrastructureInteractor

// 	ID := uuid.New().String()

// 	type args struct {
// 		ctx      context.Context
// 		id       *string
// 		isActive bool
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:      ctx,
// 				id:       &ID,
// 				isActive: true,
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:      ctx,
// 				id:       &ID,
// 				isActive: true,
// 			},
// 			wantErr: true,
// 		},

// 		{
// 			name: "Sad case - no id",
// 			args: args{
// 				ctx:      ctx,
// 				isActive: false,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeCreate.GetOrCreateFacilityFn = func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
// 					ID := uuid.New().String()
// 					return &domain.Facility{
// 						ID:          &ID,
// 						Name:        "facility.Name",
// 						Code:        "facility.Code",
// 						Active:      true,
// 						County:      "facility.County",
// 						Description: "facility.Description",
// 					}, nil
// 				}

// 				fakeQuery.RetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
// 					ID := uuid.New().String()
// 					return &domain.Facility{
// 						ID:          &ID,
// 						Name:        "facility.Name",
// 						Code:        "facility.Code",
// 						Active:      true,
// 						County:      "facility.County",
// 						Description: "facility.Description",
// 					}, nil
// 				}
// 			}

// 			if tt.name == "Sad case" {
// 				fakeQuery.RetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
// 					return nil, fmt.Errorf("an error occurred while retrieving facility")
// 				}
// 			}

// 			if tt.name == "Sad case - no id" {
// 				fakeQuery.RetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
// 					return nil, fmt.Errorf("an error occurred while retrieving facility")
// 				}
// 			}
// 			got, err := f.RetrieveFacility(tt.args.ctx, tt.args.id, tt.args.isActive)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("UseCaseFacilityImpl.RetrieveFacility() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if tt.wantErr && got != nil {
// 				t.Errorf("expected facilities to be nil for %v", tt.name)
// 				return
// 			}

// 			if !tt.wantErr && got == nil {
// 				t.Errorf("expected facilities not to be nil for %v", tt.name)
// 				return
// 			}
// 		})
// 	}
// }

// func TestUseCaseFacilityImpl_RetrieveFacilityByMFLCode_Unittest(t *testing.T) {
// 	ctx := context.Background()

// 	f := testFakeInfrastructureInteractor

// 	MFLCode := ksuid.New().String()

// 	type args struct {
// 		ctx      context.Context
// 		MFLCode  string
// 		isActive bool
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:      ctx,
// 				MFLCode:  MFLCode,
// 				isActive: true,
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:      ctx,
// 				MFLCode:  MFLCode,
// 				isActive: true,
// 			},
// 			wantErr: true,
// 		},

// 		{
// 			name: "Sad case#1",
// 			args: args{
// 				ctx:      ctx,
// 				MFLCode:  MFLCode,
// 				isActive: false,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeCreate.GetOrCreateFacilityFn = func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
// 					ID := uuid.New().String()
// 					return &domain.Facility{
// 						ID:          &ID,
// 						Name:        "facility.Name",
// 						Code:        "facility.Code",
// 						Active:      true,
// 						County:      "facility.County",
// 						Description: "facility.Description",
// 					}, nil
// 				}

// 				fakeQuery.RetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
// 					ID := uuid.New().String()
// 					return &domain.Facility{
// 						ID:          &ID,
// 						Name:        "facility.Name",
// 						Code:        "facility.Code",
// 						Active:      true,
// 						County:      "facility.County",
// 						Description: "facility.Description",
// 					}, nil
// 				}
// 			}

// 			if tt.name == "Sad case" {
// 				fakeCreate.GetOrCreateFacilityFn = func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
// 					ID := uuid.New().String()
// 					return &domain.Facility{
// 						ID:          &ID,
// 						Name:        "facility.Name",
// 						Code:        "facility.Code",
// 						Active:      true,
// 						County:      "facility.County",
// 						Description: "facility.Description",
// 					}, nil
// 				}

// 				fakeQuery.RetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
// 					return nil, fmt.Errorf("an error occurred while retrieving facility by MFLCode")
// 				}
// 			}

// 			if tt.name == "Sad case#1" {
// 				fakeQuery.RetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
// 					return nil, fmt.Errorf("an error occurred while retrieving facility by MFLCode")
// 				}
// 			}
// 			got, err := f.RetrieveFacilityByMFLCode(tt.args.ctx, tt.args.MFLCode, tt.args.isActive)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("UseCaseFacilityImpl.RetrieveFacilityByMFLCode() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if tt.wantErr && got != nil {
// 				t.Errorf("expected facilities to be nil for %v", tt.name)
// 				return
// 			}

// 			if !tt.wantErr && got == nil {
// 				t.Errorf("expected facilities not to be nil for %v", tt.name)
// 				return
// 			}
// 		})
// 	}
// }

// func TestUnit_ListFacilities(t *testing.T) {
// 	ctx := context.Background()

// 	f := testFakeInfrastructureInteractor
// 	code := ksuid.New().String()

// 	noSearchTerm := ""
// 	searchTerm := "term"

// 	noFilterInput := []*dto.FiltersInput{}

// 	formatBool := strconv.FormatBool(true)

// 	filterInput := []*dto.FiltersInput{
// 		{
// 			DataType: enums.FilterDataTypeName,
// 			Value:    "Kanairo One",
// 		},
// 		{
// 			DataType: enums.FilterDataTypeMFLCode,
// 			Value:    code,
// 		},
// 		{
// 			DataType: enums.FilterDataTypeActive,
// 			Value:    formatBool,
// 		},
// 		{
// 			DataType: enums.FilterDataTypeCounty,
// 			Value:    enums.CountyTypeNairobi.String(),
// 		},
// 	}

// 	filterEmptyName := []*dto.FiltersInput{
// 		{
// 			DataType: enums.FilterDataTypeName,
// 			Value:    "",
// 		},
// 		{
// 			DataType: enums.FilterDataTypeMFLCode,
// 			Value:    code,
// 		},
// 		{
// 			DataType: enums.FilterDataTypeActive,
// 			Value:    formatBool,
// 		},
// 		{
// 			DataType: enums.FilterDataTypeCounty,
// 			Value:    enums.CountyTypeNairobi.String(),
// 		},
// 	}
// 	filterEmptyMFLCode := []*dto.FiltersInput{
// 		{
// 			DataType: enums.FilterDataTypeName,
// 			Value:    "Kanairo One",
// 		},
// 		{
// 			DataType: enums.FilterDataTypeMFLCode,
// 			Value:    "",
// 		},
// 		{
// 			DataType: enums.FilterDataTypeActive,
// 			Value:    formatBool,
// 		},
// 		{
// 			DataType: enums.FilterDataTypeCounty,
// 			Value:    enums.CountyTypeNairobi.String(),
// 		},
// 	}

// 	filterInvalidBool := []*dto.FiltersInput{
// 		{
// 			DataType: enums.FilterDataTypeName,
// 			Value:    "Kanairo One",
// 		},
// 		{
// 			DataType: enums.FilterDataTypeMFLCode,
// 			Value:    code,
// 		},
// 		{
// 			DataType: enums.FilterDataTypeActive,
// 			Value:    "invalid",
// 		},
// 		{
// 			DataType: enums.FilterDataTypeCounty,
// 			Value:    enums.CountyTypeNairobi.String(),
// 		},
// 	}

// 	filterInvalidCounty := []*dto.FiltersInput{
// 		{
// 			DataType: enums.FilterDataTypeName,
// 			Value:    "Kanairo One",
// 		},
// 		{
// 			DataType: enums.FilterDataTypeMFLCode,
// 			Value:    code,
// 		},
// 		{
// 			DataType: enums.FilterDataTypeActive,
// 			Value:    formatBool,
// 		},
// 		{
// 			DataType: enums.FilterDataTypeCounty,
// 			Value:    "kanairo",
// 		},
// 	}
// 	paginationInput := dto.PaginationsInput{
// 		Limit:       1,
// 		CurrentPage: 1,
// 	}
// 	paginationInputNoCurrentPage := dto.PaginationsInput{
// 		Limit: 1,
// 	}

// 	type args struct {
// 		ctx              context.Context
// 		searchTerm       *string
// 		filterInput      []*dto.FiltersInput
// 		PaginationsInput dto.PaginationsInput
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:              ctx,
// 				searchTerm:       &noSearchTerm,
// 				filterInput:      noFilterInput,
// 				PaginationsInput: paginationInput,
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "valid: with valid filters",
// 			args: args{
// 				ctx:              ctx,
// 				searchTerm:       &noSearchTerm,
// 				filterInput:      filterInput,
// 				PaginationsInput: paginationInput,
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:              ctx,
// 				searchTerm:       &searchTerm,
// 				filterInput:      filterInput,
// 				PaginationsInput: paginationInput,
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "invalid: missing current page",
// 			args: args{
// 				ctx:              ctx,
// 				searchTerm:       &searchTerm,
// 				filterInput:      filterInput,
// 				PaginationsInput: paginationInputNoCurrentPage,
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "invalid: empty name passed",
// 			args: args{
// 				ctx:              ctx,
// 				searchTerm:       &searchTerm,
// 				filterInput:      filterEmptyName,
// 				PaginationsInput: paginationInput,
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "invalid: empty MFL code",
// 			args: args{
// 				ctx:              ctx,
// 				searchTerm:       &searchTerm,
// 				filterInput:      filterEmptyMFLCode,
// 				PaginationsInput: paginationInput,
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "invalid: invalid bool",
// 			args: args{
// 				ctx:              ctx,
// 				searchTerm:       &searchTerm,
// 				filterInput:      filterInvalidBool,
// 				PaginationsInput: paginationInput,
// 			},
// 			wantErr: true,
// 		},

// 		{
// 			name: "invalid: invalid county",
// 			args: args{
// 				ctx:              ctx,
// 				searchTerm:       &searchTerm,
// 				filterInput:      filterInvalidCounty,
// 				PaginationsInput: paginationInput,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Sad case" {
// 				fakeCreate.GetOrCreateFacilityFn = func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
// 					ID := uuid.New().String()
// 					return &domain.Facility{
// 						ID:          &ID,
// 						Name:        facility.Name,
// 						Code:        facility.Code,
// 						Active:      facility.Active,
// 						County:      facility.County,
// 						Description: facility.Description,
// 					}, nil
// 				}
// 				fakeQuery.ListFacilitiesFn = func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, PaginationsInput dto.PaginationsInput) (*domain.FacilityPage, error) {
// 					return nil, fmt.Errorf("failed to list facilities")
// 				}
// 			}
// 			if tt.name == "invalid: missing current page" {
// 				fakeQuery.ListFacilitiesFn = func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, PaginationsInput dto.PaginationsInput) (*domain.FacilityPage, error) {
// 					return nil, fmt.Errorf("failed to list facilities")
// 				}

// 			}
// 			if tt.name == "invalid: missing current page" {
// 				fakeQuery.ListFacilitiesFn = func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, PaginationsInput dto.PaginationsInput) (*domain.FacilityPage, error) {
// 					return nil, fmt.Errorf("failed to list facilities")
// 				}

// 			}
// 			if tt.name == "invalid: empty name passed" {
// 				fakeQuery.ListFacilitiesFn = func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, PaginationsInput dto.PaginationsInput) (*domain.FacilityPage, error) {
// 					return nil, fmt.Errorf("failed to list facilities")
// 				}

// 			}
// 			if tt.name == "invalid: empty MFL code" {
// 				fakeQuery.ListFacilitiesFn = func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, PaginationsInput dto.PaginationsInput) (*domain.FacilityPage, error) {
// 					return nil, fmt.Errorf("failed to list facilities")
// 				}

// 			}
// 			if tt.name == "invalid: invalid bool" {
// 				fakeQuery.ListFacilitiesFn = func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, PaginationsInput dto.PaginationsInput) (*domain.FacilityPage, error) {
// 					return nil, fmt.Errorf("failed to list facilities")
// 				}

// 			}
// 			if tt.name == "invalid: invalid county" {
// 				fakeQuery.ListFacilitiesFn = func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, PaginationsInput dto.PaginationsInput) (*domain.FacilityPage, error) {
// 					return nil, fmt.Errorf("failed to list facilities")
// 				}

// 			}

// 			got, err := f.ListFacilities(tt.args.ctx, tt.args.searchTerm, tt.args.filterInput, tt.args.PaginationsInput)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("OnboardingDb.ListFacilities() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if tt.wantErr && got != nil {
// 				t.Errorf("expected facilities to be nil for %v", tt.name)
// 				return
// 			}

// 			if !tt.wantErr && got == nil {
// 				t.Errorf("expected facilities not to be nil for %v", tt.name)
// 				return
// 			}
// 		})
// 	}
// }
