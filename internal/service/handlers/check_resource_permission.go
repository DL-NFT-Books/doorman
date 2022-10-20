package handlers

import (
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/nft-books/doorman/internal/service/helpers"
	"gitlab.com/tokend/nft-books/doorman/internal/service/requests"
)

func CheckResourcePermission(w http.ResponseWriter, r *http.Request) {
	logger := helpers.Log(r)
	owner, err := requests.NewCheckResourcePermission(r)
	if err != nil {
		logger.WithError(err).Debug(err)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	_, address, err := helpers.GetAccessTokenInfo(r)
	if err != nil {
		logger.WithError(err).Debug("failed to retrieve access token")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	if address != strings.ToLower(owner) && !helpers.NodeAdmins(r).CheckAdmin(common.HexToAddress(address)) {
		logger.WithError(err).Debug("user has no rights to get resource")
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
