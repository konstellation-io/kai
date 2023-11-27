package repository

import "context"

//go:generate mockery --name PredictionRepository --output ../../mocks --filename prediction_repo.go --structname MockPredictionRepo
type PredictionRepository interface {
	CreateUser(ctx context.Context, product, username, password string) error
}
