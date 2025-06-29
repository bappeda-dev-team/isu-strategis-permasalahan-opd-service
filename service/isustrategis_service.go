package service

import (
	"context"
	"permasalahanService/model/web"
)

type IsuStrategisService interface {
	Create(ctx context.Context, request web.IsuStrategisCreateRequest) (web.IsuStrategisResponse, error)
	Update(ctx context.Context, request web.IsuStrategisUpdateRequest) (web.IsuStrategisResponse, error)
	Delete(ctx context.Context, id int) error
	FindById(ctx context.Context, id int) (web.IsuStrategisResponse, error)
	FindAll(ctx context.Context, kodeOpd string, tahunAwal string, tahunAkhir string) ([]web.IsuStrategisResponse, error)
}
