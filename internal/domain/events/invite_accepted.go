package events

type InviteAcceptedEvent struct {
	InviteId  string `json:"invite_id"`
	GroupId   string `json:"group_id"`
	InviteeId string `json:"invitee_id"`
	Role      int    `json:"role"`
}

func (e *InviteAcceptedEvent) GetEventType() string {
	return "invite_accepted"
}
