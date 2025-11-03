package handlers

type EventHandlers struct {
	GroupCreated   *GroupCreatedHandler
	InviteAccepted *InviteAcceptedHandler
	ProfileCreated *ProfileCreatedHandler
}
