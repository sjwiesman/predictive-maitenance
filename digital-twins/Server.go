package main

import (
	"context"
	entities "datacenter/protobuf"
	"github.com/sjwiesman/statefun-go/pkg/flink/statefun"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
	"time"
)

const (
	history = "history"
	health  = "server-health-state"

	timeout = 60 * time.Second
)

func Server(ctx context.Context, runtime statefun.StatefulFunctionRuntime, msg *anypb.Any) error {
	report := entities.ServerMetricReport{}
	if err := msg.UnmarshalTo(&report); err == nil {
		err = onServerMetricReport(ctx, runtime, &report)
		if err != nil {
			return err
		}

		return registerServer(runtime, &report)
	}

	timer := entities.Heartbeat{}
	if err := msg.UnmarshalTo(&timer); err == nil {
		return onHeartBeat(runtime)
	}

	log.Printf("Unknown Type %v", msg)
	return nil
}

func registerServer(runtime statefun.StatefulFunctionRuntime, report *entities.ServerMetricReport) error {
	serverState := entities.ServerState{}

	exists, err := runtime.Get(health, &serverState)
	if err != nil {
		return err
	}

	if !exists {
		rack := statefun.Address{
			FunctionType: RackType,
			Id:           serverState.RackId,
		}

		err = runtime.Send(&rack, &entities.HealthUpdate{
			Status: entities.Status_HEALTHY,
		})

		if err != nil {
			return err
		}
	}

	serverState.ServerId = report.ServerId
	serverState.RackId = report.RackId
	serverState.LastReported = time.Now().UnixNano()
	serverState.Status = entities.Status_HEALTHY

	return runtime.Set(health, &serverState)
}

func onServerMetricReport(ctx context.Context, runtime statefun.StatefulFunctionRuntime, report *entities.ServerMetricReport) error {
	serverHistory := entities.ServerMetricHistory{}
	if _, err := runtime.Get(history, &serverHistory); err != nil {
		return err
	}

	serverHistory.Reports = append(serverHistory.Reports, report)
	if len(serverHistory.Reports) == 6 {
		serverHistory.Reports = serverHistory.Reports[1:]
	}

	if err := runtime.Set(history, &serverHistory); err != nil {
		return err
	}

	if err := runtime.SendAfter(statefun.Self(ctx), timeout, &entities.Heartbeat{}); err != nil {
		return err
	}

	model := &statefun.Address{
		FunctionType: ModelType,
		Id:           statefun.Self(ctx).Id,
	}

	return runtime.Send(model, &serverHistory)
}

func onHeartBeat(runtime statefun.StatefulFunctionRuntime) error {
	serverState := entities.ServerState{}

	if _, err := runtime.Get(health, &serverState); err != nil {
		return err
	}

	if serverState.LastReported-time.Now().UnixNano() <= timeout.Nanoseconds() {
		return nil
	}

	serverState.Status = entities.Status_FAILED
	if err := runtime.Set(health, &serverState); err != nil {
		return err
	}

	rack := statefun.Address{
		FunctionType: RackType,
		Id:           serverState.RackId,
	}

	return runtime.Send(&rack, &entities.HealthUpdate{
		Status: entities.Status_FAILED,
	})
}
