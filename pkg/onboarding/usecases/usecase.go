package usecases

import (
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/facility"
	engagementSvc "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
)

// Interactor is an implementation of the usecases interface
type Interactor struct {
	*engagementSvc.ServiceEngagementImpl
	*facility.UseCaseFacilityImpl
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(infrastructure infrastructure.Interactor) Interactor {
	var engagement *engagementSvc.ServiceEngagementImpl
	facility := facility.NewFacilityUsecase(infrastructure)

	return Interactor{
		engagement,
		facility,
	}
}