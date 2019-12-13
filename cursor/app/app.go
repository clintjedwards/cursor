package app

import (
	"log"

	"github.com/clintjedwards/cursor/config"
	"github.com/clintjedwards/cursor/master"
)

// StartServices initializes a GRPC-web compatible webserver and a GPRC service
func StartServices() {
	config, err := config.FromEnv()
	if err != nil {
		log.Fatal(err)
	}

	api := master.NewCursorMaster(config)
	grpcServer := master.CreateGRPCServer(api)

	master.InitGRPCService(config, grpcServer)
}
