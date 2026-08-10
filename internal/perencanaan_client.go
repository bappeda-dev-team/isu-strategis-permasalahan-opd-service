package internal

import "context"

type PerencanaanClient interface {
	GetPotensiPerangkatDaerah(ctx context.Context, ids []int) ([]PpdItem, error)
	GetIsuKlhs(ctx context.Context, ids []int) ([]IsuItem, error)
	GetIsuGlobal(ctx context.Context, ids []int) ([]IsuItem, error)
	GetIsuNasional(ctx context.Context, ids []int) ([]IsuItem, error)
	GetIsuRegional(ctx context.Context, ids []int) ([]IsuItem, error)
}
