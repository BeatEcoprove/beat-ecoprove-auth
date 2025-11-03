package events

type UserCreatedEvent struct {
	AuthID    string `json:"auth_id"`
	ProfileID string `json:"profile_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

func (e *UserCreatedEvent) GetEventType() string {
	return "user_created"
}
