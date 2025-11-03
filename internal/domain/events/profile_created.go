package events

type ProfileCreatedEvent struct {
	ProfileId string `json:"profile_id"`
	AuthId    string `json:"auth_id"`
	Role      string `json:"role"`
}

func (e *ProfileCreatedEvent) GetEventType() string {
	return "profile_created"
}
