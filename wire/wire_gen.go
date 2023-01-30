// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"context"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases"
)

// Injectors from wire.go:

// InitializeUseCases is an injector that initializes the use cases
func InitializeUseCases(ctx context.Context) (*usecases.MyCareHub, error) {
	myCareHub, err := ProviderUseCases()
	if err != nil {
		return nil, err
	}
	return myCareHub, nil
}