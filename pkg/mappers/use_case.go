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
	accessToken,
	refreshToken *services.JwtToken,
) *contracts.AuthResponse {
	return &contracts.AuthResponse{
		TokenType:    "Bearer",
		AccessToken:  accessToken.Token,
		ExpiresIn:    accessToken.ExpireAt,
		RefreshToken: refreshToken.Token,
		Scope:        domain.GetPermissions(*identityUser),
	}
}
