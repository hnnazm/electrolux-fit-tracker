package service

import (
	"fit-tracker/database"
	"time"

	"golang.org/x/net/context"
)

type (
	Config func(*userService) *userService

	userService struct {
		ingestorRepository database.IngestorRepository
	}

	GetUserDataInput struct {
		UserID string
		Date   time.Time
		Weight float64
	}

	GetUserDataResult struct {
		Steps            int64
		Distance         float64
		AverageHeartBeat float64
		KcalBurned       float64
	}
)

func New(configs ...Config) *userService {
	us := &userService{}

	for _, config := range configs {
		config(us)
	}

	return us
}

func (s userService) GetUserData(ctx context.Context, input *GetUserDataInput) (*GetUserDataResult, error) {
	var totalSteps, totalHeartBeat, totalMinutes int64
	var totalMET float64
	var res = new(GetUserDataResult)

	result, err := s.ingestorRepository.GetTraces(ctx, &database.GetTracesInput{
		UserID:    input.UserID,
		CreatedAt: input.Date,
	})
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return res, nil
	}

	for _, entry := range result {
		totalSteps += entry.Steps
		totalHeartBeat += entry.HeartBeat
		totalMET += entry.MET
		totalMinutes++ // Each entry represents 1 minute
	}

	averageHeartBeat := float64(totalHeartBeat) / float64(totalMinutes)
	averageMET := totalMET / float64(totalMinutes)

	kcalBurned := float64(totalMinutes) * (averageMET * 3.5 * input.Weight / 200)
	distance := float64(totalSteps) / 1000 * 0.7

	res.Steps = totalSteps
	res.Distance = distance
	res.AverageHeartBeat = averageHeartBeat
	res.KcalBurned = kcalBurned

	return res, nil
}

func WithIngestorRepository(ir database.IngestorRepository) Config {
	return func(us *userService) *userService {
		us.ingestorRepository = ir

		return us
	}
}
