package handlers

import (
	"fmt"
	"log"

	"github.com/BeatEcoprove/identityService/internal/domain"
	"github.com/BeatEcoprove/identityService/internal/domain/events"
	"github.com/BeatEcoprove/identityService/internal/repositories"
)

type GroupCreatedHandler struct {
	memberChatRepository repositories.IMemberChatRepository
	authRepository       repositories.IAuthRepository
}

func NewGroupCreatedHandler(
	memberRepository repositories.IMemberChatRepository,
	authRepository repositories.IAuthRepository,
) *GroupCreatedHandler {
	return &GroupCreatedHandler{
		memberChatRepository: memberRepository,
		authRepository:       authRepository,
	}
}

func (handler *GroupCreatedHandler) Call(payload any) error {
	event, ok := payload.(*events.GroupCreatedEvent)

	if !ok {
		return fmt.Errorf("failed to cast, payload %+v", payload)
	}

	if handler.memberChatRepository.ExistsGroupById(event.GroupId) {
		return fmt.Errorf("failed to create permission stack, because group already assigned")
	}

	if !handler.authRepository.ExistsUserWithId(event.CreatorId) {
		return fmt.Errorf("failed to create permission stack, because user doens't exist")
	}

	memberChat := &domain.MemberChatPermission{
		GroupId:  event.GroupId,
		MemberId: event.CreatorId,
		Role:     string(domain.ChatAdmin),
	}

	if err := handler.memberChatRepository.Create(memberChat); err != nil {
		return fmt.Errorf("failed to create permission stack")
	}

	log.Printf("created permission stack, for %s", event.CreatorId)
	return nil
}
