package main

import (
	"github.com/sjwiesman/statefun-go/pkg/flink/statefun"
	"net/http"
)

var (
	ServerType = statefun.FunctionType{
		Namespace: "digital-twin",
		Type:      "server",
	}

	RackType = statefun.FunctionType{
		Namespace: "digital-twin",
		Type:      "rack",
	}

	ModelType = statefun.FunctionType{
		Namespace: "model",
		Type:      "predictive-maintenance",
	}
)

func main() {
	registry := statefun.NewFunctionRegistry()
	registry.RegisterFunctionPointer(ServerType, Server)
	registry.RegisterFunctionPointer(RackType, Rack)

	http.Handle("/functions", registry)
	_ = http.ListenAndServe(":8001", nil)
}
