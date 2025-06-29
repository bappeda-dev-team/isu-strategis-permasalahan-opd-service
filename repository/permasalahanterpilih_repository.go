package repository

import (
	"context"
	"database/sql"
	"permasalahanService/model/domain"
)

type PermasalahanTerpilihRepository interface {
	Create(ctx context.Context, tx *sql.Tx, permasalahanTerpilih domain.PermasalahanTerpilih) (domain.PermasalahanTerpilih, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.PermasalahanTerpilih, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.PermasalahanTerpilih, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindByPermasalahanOpdId(ctx context.Context, tx *sql.Tx, permasalahanOpdId int) (domain.PermasalahanTerpilih, error)
}
