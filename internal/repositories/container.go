package repositories

type Repositories struct {
	Auth       IAuthRepository
	Profile    IProfileRepository
	MemberChat IMemberChatRepository
}
