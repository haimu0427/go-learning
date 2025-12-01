//go:build wireinject
// +build wireinject

package main

import "github.com/google/wire"

func InitializeEvent() (Event, func(), error) {
	wire.Build(
		NewEvent,
		NewGreeter,
		NewSimpleMessage,
		wire.Bind(new(IMessage), new(*SimpleMessage)),
	)
	return Event{}, nil, nil
}
