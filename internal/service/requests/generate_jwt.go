package requests

import (
	"net/http"

	"gitlab.com/tokend/nft-books/doorman/internal/service/helpers"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/urlval"
)

type GenerateJwtRequest struct {
	EthAddress string `url:"eth_address"`
	Purpose    string `url:"purpose"`
	Type       string `url:"type"`
}

func NewGenerateJwt(r *http.Request) (GenerateJwtRequest, error) {
	var request GenerateJwtRequest

	if err := urlval.Decode(r.URL.Query(), &request); err != nil {
		return request, err
	}

	return request, request.Validate()
}

func (r *GenerateJwtRequest) Validate() error {
	return validation.Errors{
		"eth_address=": validation.Validate(&r.EthAddress, validation.Required, validation.Match(helpers.AddressRegexp)),
		"purpose=":     validation.Validate(r.Purpose, validation.Required, validation.By(helpers.ValidatePurposes)),
		"type=":        validation.Validate(&r.Type, validation.Required),
	}.Filter()
}
