package main

import (
	"flag"
	"fmt"
	"gospacex-tz/srv/basic/config"
	"gospacex-tz/srv/basic/inits"
	"gospacex-tz/srv/handler/service"
	"log"
	"net"

	"google.golang.org/grpc"
	_ "gospacex-tz/srv/basic/inits"
	__ "gospacex-tz/srv/basic/proto"
)

var (
	port = flag.Int("port", 8081, "The server port")
)

func main() {
	log.Println("Consul初始化成功")
	services, err := inits.GetServiceWithLoadBalancer(config.GlobalConfig.Consul.ServiceName)
	if err != nil {
		log.Printf("获取用户服务失败: %v", err)
	} else {
		log.Printf("获取到用户服务: %s, 地址: %s:%d", services.Service, services.Address, services.Port)
	}
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	__.RegisterProductServer(s, &service.Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	err = inits.ConsulShutdown()
	if err != nil {
		return
	}
	fmt.Println("服务已退出")
}
