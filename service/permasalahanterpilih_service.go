package service

import (
	"context"
	"permasalahanService/model/web"
)

type PermasalahanTerpilihService interface {
	Create(ctx context.Context, request web.PermasalahanTerpilihRequest) (web.ChildResponse, error)
	FindAll(ctx context.Context, kodeOpd string, tahun string) ([]web.ChildResponse, error)
	Delete(ctx context.Context, id int) error
}
