package handlers

import (
	"net/http"

	"gitlab.com/tokend/nft-books/doorman/internal/service/helpers"
	"gitlab.com/tokend/nft-books/doorman/internal/service/requests"
	"gitlab.com/tokend/nft-books/doorman/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func NewJwtModel(token string, tokenType string, expTime int64) resources.Jwt {
	model := resources.Jwt{
		Key:        resources.Key{ID: token, Type: resources.ResourceType(tokenType)},
		Attributes: resources.JwtAttributes{ExpiresIn: expTime},
	}
	return model
}

func NewJwtPairResponseModel(accessToken resources.Jwt, refreshToken resources.Jwt) resources.JwtPairResponse {
	included := resources.Included{}
	included.Add(&accessToken, &refreshToken)

	model := resources.JwtPairResponse{
		Data: resources.JwtPair{
			Key: resources.Key{Type: resources.JWT_PAIR},
			Relationships: resources.JwtPairRelationships{
				AccessToken: resources.Relation{
					Data: &resources.Key{
						ID:   accessToken.ID,
						Type: resources.ACCESS_JWT,
					},
				},
				RefreshToken: resources.Relation{
					Data: &resources.Key{
						ID:   refreshToken.ID,
						Type: resources.REFRESH_JWT,
					},
				},
			},
		},
		Included: included,
	}
	return model
}

func GenerateJwtPair(w http.ResponseWriter, r *http.Request) {
	logger := helpers.Log(r)

	request, err := requests.NewGenerateJwt(r)
	if err != nil {
		logger.WithError(err).Debug("bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	accessToken, sessioExp, err := helpers.GenerateJWT(request.EthAddress, request.Purpose, helpers.ServiceConfig(r))
	if err != nil {
		logger.WithError(err).Debug("cannot generate session token")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	refreshToken, refreshExp, err := helpers.GenerateRefreshToken(request.EthAddress, helpers.ServiceConfig(r))
	if err != nil {
		logger.WithError(err).Debug("cannot generate refresh token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := NewJwtPairResponseModel(
		NewJwtModel(accessToken, string(resources.ACCESS_JWT), sessioExp),
		NewJwtModel(refreshToken, string(resources.REFRESH_JWT), refreshExp),
	)
	ape.Render(w, response)
}
