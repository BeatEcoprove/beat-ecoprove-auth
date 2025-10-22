package mappers

import (
	"github.com/BeatEcoprove/identityService/internal/domain"
	"github.com/BeatEcoprove/identityService/pkg/contracts"
	"github.com/BeatEcoprove/identityService/pkg/services"
)

func MapProfileIdsToString(profiles []domain.Profile) []string {
	var ids []string = make([]string, 0)

	for _, profile := range profiles {
		ids = append(ids, profile.ID)
	}

	return ids
}

func ToAuthResponse(
	identityUser *domain.IdentityUser,
	mainProfile *domain.Profile,
	profiles []domain.Profile,
	role string,
	accessToken,
	refreshToken *services.JwtToken,
) *contracts.AuthResponse {
	return &contracts.AuthResponse{
		Details: contracts.AccountResponse{
			UserID:     identityUser.ID,
			Email:      identityUser.Email,
			ProfileID:  mainProfile.ID,
			ProfileIds: MapProfileIdsToString(profiles),
			Role:       role,
		},
		AccessToken:            accessToken.Token,
		AccessTokenExpiration:  accessToken.ExpireAt,
		RefreshToken:           refreshToken.Token,
		RefreshTokenExpiration: refreshToken.ExpireAt,
	}
}
