package mock

import (
	"context"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// CommunityUsecaseMock contains the community usecase mock methods
type CommunityUsecaseMock struct {
	MockListMembersFn        func(ctx context.Context, input *stream.QueryOption) ([]*domain.Member, error)
	MockCreateCommunityFn    func(ctx context.Context, input dto.CommunityInput) (*domain.Community, error)
	MockListCommunityMembers func(ctx context.Context, communityID string) ([]*domain.CommunityMember, error)
}

// NewCommunityUsecaseMock initializes a new instance of the Community usecase happy cases
func NewCommunityUsecaseMock() *CommunityUsecaseMock {
	return &CommunityUsecaseMock{
		MockListMembersFn: func(ctx context.Context, input *stream.QueryOption) ([]*domain.Member, error) {
			return []*domain.Member{
				{
					ID:   uuid.New().String(),
					Role: "user",
				},
			}, nil
		},
		MockCreateCommunityFn: func(ctx context.Context, input dto.CommunityInput) (*domain.Community, error) {
			return &domain.Community{
				ID:          uuid.New().String(),
				Name:        "test",
				Description: "test",
				AgeRange: &domain.AgeRange{
					LowerBound: 1,
					UpperBound: 3,
				},
				Gender:     []enumutils.Gender{enumutils.AllGender[0]},
				ClientType: []enums.ClientType{enums.AllClientType[0]},
				InviteOnly: true,
			}, nil
		},
		MockListCommunityMembers: func(ctx context.Context, communityID string) ([]*domain.CommunityMember, error) {
			return []*domain.CommunityMember{
				{
					UserID: uuid.New().String(),
					User: domain.Member{
						ID: uuid.New().String(),
					},
				},
			}, nil
		},
	}
}

// ListMembers mocks the implementation for listing getstream users
func (c CommunityUsecaseMock) ListMembers(ctx context.Context, input *stream.QueryOption) ([]*domain.Member, error) {
	return c.MockListMembersFn(ctx, input)
}

// CreateCommunity mocks the implementation of creating communities
func (c CommunityUsecaseMock) CreateCommunity(ctx context.Context, input dto.CommunityInput) (*domain.Community, error) {
	return c.MockCreateCommunityFn(ctx, input)
}

// ListCommunityMembers mocks the implementation of listing members
func (c CommunityUsecaseMock) ListCommunityMembers(ctx context.Context, communityID string) ([]*domain.CommunityMember, error) {
	return c.MockListCommunityMembers(ctx, communityID)
}
