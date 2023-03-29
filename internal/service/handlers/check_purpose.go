package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"github.com/dl-nft-books/doorman/internal/service/helpers"
	"github.com/dl-nft-books/doorman/resources"
)

func CheckPurpose(w http.ResponseWriter, r *http.Request) {
	logger := helpers.Log(r)

	purpose, _, err := helpers.GetAccessTokenInfo(r)
	if err != nil {
		logger.WithError(err).Debug("failed to retrieve access token")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	err = helpers.ValidatePurposes(purpose)
	if err != nil {
		logger.WithError(err).Debug("unknown purpose")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	result := resources.Purpose{
		Type: purpose,
	}

	ape.Render(w, result)
}
