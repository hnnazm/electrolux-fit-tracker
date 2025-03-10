package ingestor

import (
	"context"
	"log"
	"os"

	"fit-tracker/database"
	"fit-tracker/ingestor/handler"
)

type (
	Handler interface {
		CheckHealth() bool
		GenerateAccessToken(input handler.GenerateTokenInput) (string, error)
		GetUserInfo(ctx context.Context, input *handler.GetUserInfoInput) (*handler.GetUserInfoResult, error)
		PollTraces(ctx context.Context, input *handler.PollTracesInput)
	}

	FitTrackerService struct {
		handler Handler
		db      *database.Database
		ch      chan *database.SaveTracesInput
	}
)

func New(handler Handler, db *database.Database) *FitTrackerService {
	return &FitTrackerService{
		handler: handler,
		db:      db,
		ch:      make(chan *database.SaveTracesInput),
	}
}

func (s *FitTrackerService) Run(ctx context.Context) error {
	accessToken, err := s.handler.GenerateAccessToken(handler.GenerateTokenInput{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	})
	if err != nil {
		return err
	}

	go s.handler.PollTraces(ctx, &handler.PollTracesInput{
		AccessToken: accessToken,
		DataCh:      s.ch,
	})

	for {
		select {
		case <-ctx.Done():
			return nil
		case input := <-s.ch:
			if err := s.db.IngestorRepository.SaveTraces(ctx, input); err != nil {
				log.Println("error saving traces", err)
			}
		}
	}
}
