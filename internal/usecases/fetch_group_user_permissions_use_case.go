package usecases

import (
	"github.com/BeatEcoprove/identityService/internal/repositories"
	"github.com/BeatEcoprove/identityService/pkg/contracts"
	fails "github.com/BeatEcoprove/identityService/pkg/errors"
)

type (
	FetchGroupUserPermissionsInput struct {
		// AuthId   string FIXME: right now it doens't need authentication, but later must authenticate the microservice
		GroupID  string
		MmeberID string
	}

	FetchGroupUserPermissionsUseCase struct {
		memberRepo repositories.IMemberChatRepository
	}
)

func NewFetchGroupUserPermissionsUseCase(
	memberRepo repositories.IMemberChatRepository,
) *FetchGroupUserPermissionsUseCase {
	return &FetchGroupUserPermissionsUseCase{
		memberRepo: memberRepo,
	}
}

func (apu *FetchGroupUserPermissionsUseCase) Handle(request FetchGroupUserPermissionsInput) (*contracts.GroupPermissionsResponse, error) {
	permissions, err := apu.memberRepo.GetPermissions(request.GroupID)

	if err != nil {
		return nil, fails.GROUP_NOT_FOUND
	}

	for _, permission := range permissions {
		if permission.MemberId == request.MmeberID {
			return &contracts.GroupPermissionsResponse{
				MemberID: permission.MemberId,
				Role:     permission.Role,
			}, nil
		}
	}

	return nil, fails.MEMBER_NOT_FOUND
}
