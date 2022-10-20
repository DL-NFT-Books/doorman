package service

import (
	gosdk "gitlab.com/tokene/go-sdk"

	"gitlab.com/tokend/nft-books/doorman/internal/config"
	"gitlab.com/tokend/nft-books/doorman/internal/service/handlers"
	"gitlab.com/tokend/nft-books/doorman/internal/service/helpers"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router(cfg config.Config) chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			helpers.CtxLog(s.log),
			helpers.CtxServiceConfig(cfg.ServiceConfig()),
			//TODO change when admin's contracts added
			helpers.CtxNodeAdmins(gosdk.NewNodeAdminsMock(cfg.AdminsConfig().Admins...)),
		),
	)

	r.Route("/integrations/doorman", func(r chi.Router) {
		r.Get("/validate-token", handlers.ValidateJWT)
		r.Get("/refresh-token", handlers.RefreshJwt)
		r.Get("/token-pair", handlers.GenerateJwtPair)
		r.Get("/check-permission/{owner}", handlers.CheckResourcePermission)
		r.Get("/check-purpose", handlers.CheckPurpose)
	})

	return r
}
