package main

import "github.com/4726/game/services/matchmaking/pb"

type Server struct {
	rankedQueue Queue
	unrankedQueue Queue
	matches map[uint64]Match
	queueTimes QueueTimes
}

func (s *Server) Join(in *pb.JoinQueueRequest, outStream pb.Queue_JoinServer) error {
	var queue Queue
	if in.GetQueueType() == QueueType_UNRANKED {
		queue = s.unrankedQueue
	} else {
		queue = s.rankedQueue
	}
	
	foundCh := make(chan uint64, 1)
	qd := QueueData{in.GetUserID(), in.GetRating(), foundCh, time.Now()}
	users, err := queue.FindAndDeleteWithinRatingRangeOrEnqueue(qd, 100, 10)
	if err != nil {
		//already in quee
		return
	}

	var matchID uint64
	if len(users) > 0 {
		matchID, _ = getMatchID()
		for _, v := range users {
			v.FoundCh <- matchID
		}
	} else {
		matchID = <- foundCh
		matches[matchID] == NewMatch(users, in.GetQueueType())
	}
	resp := &JoinQueueResponse {
		in.GetUserId(),
		in.GetQueueType(),
		matchID,
		true, 
		uint32(20),
	}

	return outStream.Send(resp)
}

func (s *Server) Leave(ctx context.Context, in *pb.LeaveQueueRequest) (*pb.LeaveQueueResponse, error) {
	var queue Queue
	if in.GetQueueType() == QueueType_UNRANKED {
		queue = s.unrankedQueue
	} else {
		queue = s.rankedQueue
	}

	queue.DeleteOne(in.GetUserId())

	return &pb.LeaveQueueResponse{in.GetUserId(), in.GetQueueType()}, nil
}

func (s *Server) Accept(in *pb.AcceptQueueRequest, outStream pb.Queue_AcceptServer) error {
	match := s.matches[in.GetMatchId()]
	ch := make(chan MatchResponse, 1)
	if err := match.Accept(in.GetUserId()); err != nil {
		return err
	}

	for {
		resp <- ch
		if err := outStream.Send(resp); err != nil {
			return err
		}
		if resp.GetCancelled() {
			return nil
		}
	}
}

func (s *Server) Decline(ctx context.Context, in *pb.DeclineQueueRequest) (*pb.DeclineQueueResponse, error) {
	match := s.matches[in.GetMatchId()]
	ch := make(chan MatchResponse, 1)
	if err := match.Decline(in.GetUserId()); err != nil {
		return nil, err
	}

	return &DeclineQueueResponse{in.GetUserId(), in.GetQueueType()}, nil
}

func (s *Server) Info(ctx context.Context, in *pb.QueueInfoRequest) (*pb.QueueInfoResponse, error) {
	estimatedWaitTime := s.queueTimes.EstimatedWaitTime(in.GetRating())
	
	return &QueueInfoResponse{
		uint32(estimatedWaitTime.Seconds())
	}, nil
}

func (s *Server) getMatchID() (uint64, error) {
	return 1, nil
}