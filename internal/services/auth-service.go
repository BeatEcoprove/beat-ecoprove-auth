package services

import (
	"errors"

	"github.com/BeatEcoprove/identityService/internal/domain"
	"github.com/BeatEcoprove/identityService/internal/repositories"
	"github.com/BeatEcoprove/identityService/pkg/contracts"
)

type SignUpInput struct {
	Email    string
	Password string
	Role     int
}

type AuthService struct {
	authRepo    *repositories.AuthRepository
	profileRepo *repositories.ProfileRepository
}

func NewAuthService(
	authRepo *repositories.AuthRepository,
	profileRepo *repositories.ProfileRepository,
) *AuthService {
	return &AuthService{
		authRepo:    authRepo,
		profileRepo: profileRepo,
	}
}

func (as *AuthService) SignUp(input SignUpInput) (*contracts.SignUpResponse, error) {
	// Setps to create an user
	/// 1. Exists already this user registered?
	if ok := as.authRepo.ExistsUserWithEmail(input.Email); ok {
		return nil, errors.New("ups this user already exits")
	}

	/// 2. If not then register the user
	identityUser := domain.NewIdentityUser(
		input.Email,
		input.Password,
		domain.Role(input.Role),
	)

	if err := as.authRepo.Create(identityUser); err != nil {
		return nil, err
	}

	/// 3. Generate an ProfileID -> Which will be the Main Profile
	//// Generate the profile info and the type
	profile := domain.NewProfile(identityUser.ID, domain.Main)

	if err := as.profileRepo.Create(profile); err != nil {
		return nil, err
	}

	/// 4. Generate the authorization tokens -> Hydra

	return &contracts.SignUpResponse{
		ProfileID: profile.ID,
	}, nil
}
