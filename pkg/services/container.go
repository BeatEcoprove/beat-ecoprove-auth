package services

type Services struct {
	Token ITokenService
	PG    IPGService
	Email IEmailService
}
