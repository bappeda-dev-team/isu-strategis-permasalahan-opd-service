package internal

import "context"

type PerencanaanClient interface {
	GetPotensiPerangkatDaerah(ctx context.Context, ids int) ([]PpdItem, error)
}
