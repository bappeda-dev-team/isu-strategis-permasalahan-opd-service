package repository

import (
	"context"
	"database/sql"
	"permasalahanService/model/domain"
)

type IsuStrategisRepository interface {
	Create(ctx context.Context, tx *sql.Tx, isuStrategis domain.IsuStrategis) (domain.IsuStrategis, error)
	Update(ctx context.Context, tx *sql.Tx, isuStrategis domain.IsuStrategis) (domain.IsuStrategis, error)
	FindById(ctx context.Context, tx *sql.Tx, isuStrategisId int) (domain.IsuStrategis, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string, tahunAwal string, tahunAkhir string) ([]domain.IsuStrategis, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindDataDukungById(ctx context.Context, tx *sql.Tx, dataDukungId int) (domain.DataDukung, error)
	FindJumlahDataById(ctx context.Context, tx *sql.Tx, jumlahDataId int) (domain.JumlahData, error)
	FindDataDukungByPermasalahanId(ctx context.Context, tx *sql.Tx, permasalahanId int) ([]domain.DataDukung, error)
	FindJumlahDataByDataDukungId(ctx context.Context, tx *sql.Tx, dataDukungId int) ([]domain.JumlahData, error)
	DeleteJumlahDataByDataDukungId(ctx context.Context, tx *sql.Tx, dataDukungId int) error
	DeleteDataDukungByPermasalahanId(ctx context.Context, tx *sql.Tx, permasalahanId int) error
}
