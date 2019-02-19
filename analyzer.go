package main

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/src-d/lookout-sdk.v0/pb"

	"github.com/sanity-io/litter"
	"google.golang.org/grpc"
	log "gopkg.in/src-d/go-log.v1"
)

type analyzer struct{}

var portToListen = 9930
var dataSrvAddr, _ = pb.ToGoGrpcAddress("ipv4://localhost:10301")
var version = "0.0.1"

func (*analyzer) NotifyReviewEvent(ctx context.Context, review *pb.ReviewEvent) (*pb.EventResponse, error) {
	comment := pb.Comment{
		Text: fmt.Sprintf(
			"Analyzer called with settings:\n```golang\n%s\n```",
			litter.Sdump(review.Configuration))}

	return &pb.EventResponse{AnalyzerVersion: version, Comments: []*pb.Comment{&comment}}, nil
}

func (*analyzer) NotifyPushEvent(context.Context, *pb.PushEvent) (*pb.EventResponse, error) {
	return &pb.EventResponse{}, nil
}

func main() {
	lis, err := pb.Listen(fmt.Sprintf("ipv4://0.0.0.0:%d", portToListen))
	if err != nil {
		log.Errorf(err, "failed to listen on port: %d", portToListen)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAnalyzerServer(grpcServer, &analyzer{})
	log.Infof("starting gRPC Analyzer server at port %d", portToListen)
	grpcServer.Serve(lis)
}
