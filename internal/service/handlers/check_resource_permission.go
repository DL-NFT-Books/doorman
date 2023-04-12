package handlers

import (
	"github.com/dl-nft-books/doorman/solidity/generated/contractsregistry"
	"github.com/dl-nft-books/doorman/solidity/generated/rolemanager"
	"net/http"
	"strings"

	"github.com/dl-nft-books/doorman/internal/service/helpers"
	"github.com/dl-nft-books/doorman/internal/service/requests"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
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
	networker := helpers.NetworkConnector(r)
	networks, err := networker.GetNetworksDetailed()
	if err != nil {
		logger.WithError(err).Debug("failed to get networks")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	isAdmin := true
	for _, network := range networks.Data {
		contractsRegistry, err := contractsregistry.NewContractsregistry(common.HexToAddress(network.FactoryAddress), network.RpcUrl)
		if err != nil {
			logger.WithError(err).Debug("failed to create contract registry")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		roleManagerContract, err := contractsRegistry.GetRoleManagerContract(nil)
		roleManager, err := rolemanager.NewRolemanager(roleManagerContract, network.RpcUrl)
		if err != nil {
			logger.WithError(err).Debug("failed to create role manager")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		isAdmin, err = roleManager.RolemanagerCaller.HasAnyRole(nil, common.HexToAddress(address))
		if err != nil {
			logger.WithError(err).Debug("failed to check is admin")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		// check if user is admin at least in one network
		if !isAdmin {
			break
		}
	}
	if address != strings.ToLower(owner) && !isAdmin {
		logger.WithError(err).Debug("user has no rights to get resource")
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
