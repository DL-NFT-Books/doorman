package service

import (
	"github.com/dl-nft-books/doorman/internal/config"
	"github.com/dl-nft-books/doorman/internal/service/handlers"
	"github.com/dl-nft-books/doorman/internal/service/helpers"
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
			helpers.CtxNetworkConnector(*cfg.NetworkConnector()),
		),
	)

	r.Route("/integrations/doorman", func(r chi.Router) {
		r.Route("/validate-token", func(r chi.Router) {
			r.Get("/", handlers.ValidateJWT)
			r.Get("/admin", handlers.ValidateJWT)
		})
		r.Get("/refresh-token", handlers.RefreshJwt)
		r.Get("/token-pair", handlers.GenerateJwtPair)
		r.Get("/check-permission/{owner}", handlers.CheckResourcePermission)
		r.Get("/check-purpose", handlers.CheckPurpose)
	})

	return r
}
