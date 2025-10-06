package handlers

import (
	"fmt"

	"github.com/BeatEcoprove/identityService/internal/domain"
	"github.com/BeatEcoprove/identityService/internal/domain/events"
	"github.com/BeatEcoprove/identityService/internal/repositories"
)

type InviteAcceptedHandler struct {
	memberChatRepository repositories.IMemberChatRepository
	authRepository       repositories.IAuthRepository
}

func NewInviteAcceptedHandler(
	memberRepository repositories.IMemberChatRepository,
	authRepository repositories.IAuthRepository,
) *InviteAcceptedHandler {
	return &InviteAcceptedHandler{
		memberChatRepository: memberRepository,
		authRepository:       authRepository,
	}
}

func (handler *InviteAcceptedHandler) Call(payload any) error {
	event, ok := payload.(*events.InviteAcceptedEvent)

	if !ok {
		return fmt.Errorf("failed to cast, payload %+v", payload)
	}

	groupEntry, err := handler.memberChatRepository.GetByGroupId(event.GroupId)

	if err != nil {
		return fmt.Errorf("failed to find group entry")
	}

	if !handler.authRepository.ExistsUserWithId(event.InviteeId) {
		return fmt.Errorf("failed to find user with id %s", event.InviteeId)
	}

	if groupEntry.IsMember(event.InviteeId) {
		return fmt.Errorf("failed, because is already member of the group")
	}

	memberChat := domain.NewMemberChat(event.GroupId, event.InviteeId, domain.GetChatRoleByInt(event.Role))

	if err := handler.memberChatRepository.Create(memberChat); err != nil {
		return fmt.Errorf("failed, because something went wrong")
	}

	return nil
}
