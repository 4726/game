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
	startTime := time.Now()
	qd := QueueData{in.GetUserID(), in.GetRating(), foundCh, startTime. false}
	users, err := queue.FindAndMarkMatchFoundWithinRatingRangeOrEnqueue(qd, 100, 10)
	if err != nil {
		//already in queue
		return
	}

	matchFoundFn := func(status QueueStatus) (bool, error) {
		status = <- foundCh
		if status.MatchStarted {
			return true, nil
		}
		s.queueTimes.Add(QueueDuration{in.GetRating(), time.Since(startTime)})

		resp := &JoinQueueResponse {
			in.GetUserId(),
			in.GetQueueType(),
			matchID,
			true, 
			uint32(20),
		}
	
		err := outStream.Send(resp)
		return false, err
	}()

	var status QueueStatus
	if len(users) > 0 {
		status, _ = getMatchID()
		matches[matchID] == NewMatch(users, in.GetQueueType())
		for _, v := range users {
			v.FoundCh <- QueueStatus{matchID, false}
		}
		s.queueTimes.Add(QueueDuration{in.GetRating(), time.Since(startTime)})

		resp := &JoinQueueResponse {
			in.GetUserId(),
			in.GetQueueType(),
			matchID,
			true, 
			uint32(20),
		}
	
		if err := outStream.Send(resp); err != nil {
			return err
		}
	}
	for {
		matchStarted, err := matchFoundFn(status)
		if err != nil {
			return err
		}
		if matchStated {
			return nil
		}
	}
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
	var queue Queue
	if in.GetQueueType() == QueueType_UNRANKED {
		queue = s.unrankedQueue
	} else {
		queue = s.rankedQueue
	}

	match := s.matches[in.GetMatchId()]
	ch := make(chan *pb.AcceptQueueResponse, 1)
	if err := match.Accept(in.GetUserId(), ch); err != nil {
		return err
	}

	for {
		resp <- ch
		if err := outStream.Send(resp); err != nil {
			return err
		}
		if resp.GetCancelled() {
			queue.MarkMatchFound(in.GetUserId(), false) //readd user to queue
			return nil
		}
		if resp.GetTotalAccepted() == resp.GetTotalNeeded() {
			//everyone accepted, remove user from queue
			//also removes match
			queue.SetMatchStartedAndDelete(in.GetUserId(), in.GetMatchId(), true)
			delete(s.matches, in.GetMatchId())
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