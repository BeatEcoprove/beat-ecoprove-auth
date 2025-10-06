package events

type UserCreatedEvent struct {
	PublicId string `json:"public_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func (e *UserCreatedEvent) GetEventType() string {
	return "user_created"
}
