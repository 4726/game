package main

import "github.com/4726/game/services/matchmaking/pb"

type QueueService struct {
	queues map[pb.QueueType]Queue
	matches map[uint64]Match
	queueTimes QueueTimes
	opts ServerOptions
}

type QueueServiceOptions struct {
	//rating range where players can get matched with each other
	//ex) RatingRange of 100 allows a player with 1000 rating to match with a player with 1100 rating
	RatingRange int
	//number of players in a single match
	PlayerCount int
}

func NewQueueService(opts QueueServiceOptions) *QueueService {
	queues := map[pb.QueueType]Queue{}
	queues[pb.QueueType_UNRANKED] = NewQueue()
	queues[pb.QueueType_RANKED] = NewQueue()
	return &QueueService{
		queues,
		map[uint64]Match{},
		NewQueueTimes(1000),
		opts,
	}
}

func (s *QueueService) Join(in *pb.JoinQueueRequest, outStream pb.Queue_JoinServer) error {
	queue := s.queues[in.GetQueueType()]
	
	foundCh := make(chan uint64, 1)
	startTime := time.Now()
	qd := QueueData{in.GetUserID(), in.GetRating(), foundCh, startTime. false}
	users, err := queue.FindAndMarkMatchFoundWithinRatingRangeOrEnqueue(qd, s.opts.RatingRange, s.opts.PlayerCount)
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

func (s *QueueService) Leave(ctx context.Context, in *pb.LeaveQueueRequest) (*pb.LeaveQueueResponse, error) {
	queue := s.queues[in.GetQueueType()]

	queue.DeleteOne(in.GetUserId())

	return &pb.LeaveQueueResponse{in.GetUserId(), in.GetQueueType()}, nil
}

func (s *QueueService) Accept(in *pb.AcceptQueueRequest, outStream pb.Queue_AcceptServer) error {
	queue := s.queues[in.GetQueueType()]

	match := s.matches[in.GetMatchId()]
	ch := make(chan MatchStatus, 1)
	if err := match.Accept(in.GetUserId(), ch); err != nil {
		return err
	}

	for {
		status <- ch
		resp := &AcceptQueueResponse{
			TotalAccepted: status.TotalAccepted,
			TotalNeeded: stauts.TotalNeeded,
			QueueType: in.GetQueueType(),
			Cancelled: status.Cancelled,
		}
		if err := outStream.Send(resp); err != nil {
			return err
		}
		if resp.GetCancelled() {
			//someone declined, adds user back into queue
			queue.MarkMatchFound(in.GetUserId(), false)
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

func (s *QueueService) Decline(ctx context.Context, in *pb.DeclineQueueRequest) (*pb.DeclineQueueResponse, error) {
	match := s.matches[in.GetMatchId()]
	ch := make(chan MatchResponse, 1)
	if err := match.Decline(in.GetUserId()); err != nil {
		return nil, err
	}

	return &DeclineQueueResponse{in.GetUserId(), in.GetQueueType()}, nil
}

func (s *QueueService) Info(ctx context.Context, in *pb.QueueInfoRequest) (*pb.QueueInfoResponse, error) {
	estimatedWaitTime := s.queueTimes.EstimatedWaitTime(in.GetRating(), 100)
	
	return &QueueInfoResponse{
		uint32(estimatedWaitTime.Seconds()),
	}, nil
}

func (s *QueueService) getMatchID() (uint64, error) {
	return 1, nil
}