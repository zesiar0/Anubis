package proxy

import "Anubis/pkg/model"

type ConnectionOptions struct {
	Protocol          string
	user              string
	predefinedAccount *model.Account
}

type ConnectionOption func(options *ConnectionOptions)

func ConnectProtocol(protocol string) ConnectionOption {
	return func(opts *ConnectionOptions) {
		opts.Protocol = protocol
	}
}

func ConnectUser(user string) ConnectionOption {
	return func(opts *ConnectionOptions) {
		opts.user = user
	}
}

func ConnectAccount(account *model.Account) ConnectionOption {
	return func(opts *ConnectionOptions) {
		opts.predefinedAccount = account
	}
}
