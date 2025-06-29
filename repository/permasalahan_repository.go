package repository

import (
	"context"
	"database/sql"
	"permasalahanService/model/domain"
	"permasalahanService/model/web"
)

type PermasalahanRepository interface {
	Create(ctx context.Context, tx *sql.Tx, permasalahan domain.Permasalahan) (domain.Permasalahan, error)
	Update(ctx context.Context, tx *sql.Tx, permasalahan domain.Permasalahan) domain.Permasalahan
	Delete(ctx context.Context, tx *sql.Tx, id string) error
	FindById(ctx context.Context, tx *sql.Tx, id string) (domain.Permasalahan, error)
	FindByKodeOpdAndTahun(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.Permasalahan, error)
	GetPohonKinerjaFromAPI(ctx context.Context, kodeOpd string, tahun string) (*web.PohonKinerjaDataResponse, error)
	MergePohonKinerjaWithPermasalahan(ctx context.Context, tx *sql.Tx, pohonKinerja *web.PohonKinerjaDataResponse, permasalahans []domain.Permasalahan) *web.PohonKinerjaDataResponse
	FindByPokinId(ctx context.Context, tx *sql.Tx, pokinId int) (domain.Permasalahan, error)
	IsPermasalahanTerpilih(ctx context.Context, tx *sql.Tx, idPermasalahan int) (bool, error)
	FindByIsuStrategisId(ctx context.Context, tx *sql.Tx, isuStrategisId int) ([]domain.Permasalahan, error)
	ResetIsuStrategisId(ctx context.Context, tx *sql.Tx, id int) error
}
