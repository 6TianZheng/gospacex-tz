package init

import (
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gospacex-tz/bff/basic/config"
	"gospacex-tz/srv/basic/inits"
	__ "gospacex-tz/srv/basic/proto"
	"log"
)

func init() {
	GrpcInit()
	inits.MysqlInit()
}

func GrpcInit() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient("127.0.0.1:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	config.ProductClient = __.NewProductClient(conn)
}
