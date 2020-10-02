package main

import (
	"context"
	entities "datacenter/protobuf"
	"github.com/golang/protobuf/ptypes"
	"github.com/sjwiesman/statefun-go/pkg/flink/statefun"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	servers = "servers"
)

func Rack(ctx context.Context, runtime statefun.StatefulFunctionRuntime, msg *anypb.Any) error {
	update := entities.HealthUpdate{}
	if err := ptypes.UnmarshalAny(msg, &update); err != nil {
		return updateHealth(ctx, runtime, &update)
	}

	return nil
}

func updateHealth(ctx context.Context, runtime statefun.StatefulFunctionRuntime, msg *entities.HealthUpdate) error {
	serversPerRack := entities.ServersPerRack{}
	if _, err := runtime.Get(servers, &serversPerRack); err != nil {
		return err
	}

	serverId := statefun.Caller(ctx).Id
	serversPerRack.StatusByServer[serverId] = msg.Status
	return runtime.Set(servers, &serversPerRack)
}
