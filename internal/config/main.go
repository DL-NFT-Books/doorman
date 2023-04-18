package config

import (
	networker "github.com/dl-nft-books/network-svc/connector"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
)

type Config interface {
	comfig.Logger
	types.Copuser
	comfig.Listenerer
	ServiceConfiger
	networker.NetworkConfigurator
}

type config struct {
	comfig.Logger
	types.Copuser
	comfig.Listenerer
	getter kv.Getter
	ServiceConfiger
	networker.NetworkConfigurator
}

func New(getter kv.Getter) Config {
	return &config{
		getter:              getter,
		Copuser:             copus.NewCopuser(getter),
		Listenerer:          comfig.NewListenerer(getter),
		Logger:              comfig.NewLogger(getter, comfig.LoggerOpts{}),
		ServiceConfiger:     NewServiceConfiger(getter),
		NetworkConfigurator: networker.NewNetworkConfigurator(getter),
	}
}
