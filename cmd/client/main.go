package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "rusprofile-wrapper/internal/rpc_server"
	"strconv"
	"time"
)

var	addr = flag.String("addr", "localhost:50051", "the address to connect to")


func runClient(INN uint64) {
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCompanyInfoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	r, err := c.FetchCompanyInfo(ctx, &pb.CompanyRequest{INN: INN})
	if err != nil {
		log.Fatalf("failed to fetch company info with following INN: %v; %v", err)
	}
	respINN := strconv.FormatUint(r.GetINN(), 10)
	KPP := strconv.FormatUint(uint64(r.GetKPP()), 10)
	log.Printf("INN: %s; KPP: %s; Company name: %s; Director full name: %s",
		respINN, KPP, r.GetCompanyName(), r.GetDirectorFullName())
}



func main() {
	INN := uint64(7840005720)
	runClient(INN)
}