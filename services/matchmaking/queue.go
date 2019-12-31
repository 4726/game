package main

import "github.com/4726/game/services/matchmaking/pb"

func Join(*pb.JoinQueueRequest, pb.Queue_JoinServer) error {

}

func Leave(ctx context.Context, *pb.LeaveQueueRequest) (*pb.LeaveQueueResponse, error) {

}

func Accept(*pb.AcceptQueueRequest, pb.Queue_AcceptServer) error {

}

func Decline(ctx context.Context, *pb.DeclineQueueRequest) (*pb.DeclineQueueResponse, error) {


}

func Info(ctx context.Context, *pb.QueueInfoRequest) (*pb.QueueInfoResponse, error) {

}