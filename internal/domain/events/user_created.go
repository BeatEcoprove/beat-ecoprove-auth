package events

type UserCreatedEvent struct {
	AuthId    string `json:"auth_id"`
	ProfileId string `json:"profile_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

func (e *UserCreatedEvent) GetEventType() string {
	return "user_created"
}
