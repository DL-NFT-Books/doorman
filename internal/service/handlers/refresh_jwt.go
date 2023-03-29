package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"github.com/dl-nft-books/doorman/internal/service/helpers"
	"github.com/dl-nft-books/doorman/resources"
)

func RefreshJwt(w http.ResponseWriter, r *http.Request) {
	logger := helpers.Log(r)

	address, err := helpers.GetTokenInfo(r)
	if err != nil {
		logger.WithError(err).Debug("failed to retrieve refresh token")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// success logic
	claims := resources.JwtClaimsAttributes{
		EthAddress: address,
		Purpose:    resources.Purpose{Type: "session"},
	}

	accessToken, sessioExp, err := helpers.GenerateJWT(claims.EthAddress, claims.Purpose.Type, helpers.ServiceConfig(r))
	if err != nil {
		logger.WithError(err).Debug(err)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	refreshToken, refreshExp, err := helpers.GenerateRefreshToken(claims.EthAddress, helpers.ServiceConfig(r))
	if err != nil {
		logger.WithError(err).Debug(err)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := NewJwtPairResponseModel(
		NewJwtModel(accessToken, string(resources.ACCESS_JWT), sessioExp),
		NewJwtModel(refreshToken, string(resources.REFRESH_JWT), refreshExp),
	)

	ape.Render(w, response)
}
