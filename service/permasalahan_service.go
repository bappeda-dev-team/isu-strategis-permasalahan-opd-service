package service

import (
	"context"
	"permasalahanService/model/web"
)

type PermasalahanService interface {
	Create(ctx context.Context, request web.PermasalahanCreateRequest) (web.ChildResponse, error)
	Update(ctx context.Context, request web.PermasalahanUpdateRequest) (web.ChildResponse, error)
	Delete(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) (web.PermasalahanResponsesbyId, error)
	FindAllPohonKinerja(ctx context.Context, kodeOpd string, tahun string) (*web.PohonKinerjaDataResponse, error)
}
