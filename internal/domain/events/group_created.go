package events

type GroupCreatedEvent struct {
	GroupId   string `json:"group_id"`
	CreatorId string `json:"creator_id"`
}

func (e *GroupCreatedEvent) GetEventType() string {
	return "group_created"
}
