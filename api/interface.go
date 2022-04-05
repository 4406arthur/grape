package api

import (
	"context"

	"github.com/4406arthur/grape/repository"
)

type QueryUsecase interface {
	GetLatestNBlock(context.Context, int) ([]repository.Block, error)
	GetFullBlockByID(context.Context, int) (repository.Block, []repository.Tx, error)
	GetTXByHash(context.Context, string) (repository.Tx, error)
}
