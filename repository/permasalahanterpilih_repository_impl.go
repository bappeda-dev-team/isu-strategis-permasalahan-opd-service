package repository

import (
	"context"
	"database/sql"
	"errors"
	"permasalahanService/model/domain"
)

type PermasalahanTerpilihRepositoryImpl struct{}

func NewPermasalahanTerpilihRepositoryImpl() *PermasalahanTerpilihRepositoryImpl {
	return &PermasalahanTerpilihRepositoryImpl{}
}
func (repository *PermasalahanTerpilihRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, permasalahanTerpilih domain.PermasalahanTerpilih) (domain.PermasalahanTerpilih, error) {
	script := "INSERT INTO tb_permasalahan_terpilih (permasalahan_opd_id, kode_opd, tahun) VALUES (?, ?, ?)"
	result, err := tx.ExecContext(ctx, script, permasalahanTerpilih.PermasalahanOpdId, permasalahanTerpilih.KodeOpd, permasalahanTerpilih.Tahun)
	if err != nil {
		return domain.PermasalahanTerpilih{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.PermasalahanTerpilih{}, err
	}
	permasalahanTerpilih.Id = int(id)

	return permasalahanTerpilih, nil
}

func (repository *PermasalahanTerpilihRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.PermasalahanTerpilih, error) {
	script := "SELECT id, permasalahan_opd_id, kode_opd, tahun FROM tb_permasalahan_terpilih WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return domain.PermasalahanTerpilih{}, err
	}
	defer rows.Close()

	permasalahanTerpilih := domain.PermasalahanTerpilih{}
	if rows.Next() {
		err := rows.Scan(
			&permasalahanTerpilih.Id,
			&permasalahanTerpilih.PermasalahanOpdId,
			&permasalahanTerpilih.KodeOpd,
			&permasalahanTerpilih.Tahun,
		)
		if err != nil {
			return domain.PermasalahanTerpilih{}, err
		}
	}
	return permasalahanTerpilih, nil
}

func (repository *PermasalahanTerpilihRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.PermasalahanTerpilih, error) {
	script := "SELECT id, permasalahan_opd_id, kode_opd, tahun FROM tb_permasalahan_terpilih WHERE kode_opd = ? AND tahun = ?"
	rows, err := tx.QueryContext(ctx, script, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permasalahanTerpilihs := []domain.PermasalahanTerpilih{}
	for rows.Next() {
		permasalahanTerpilih := domain.PermasalahanTerpilih{}
		err := rows.Scan(
			&permasalahanTerpilih.Id,
			&permasalahanTerpilih.PermasalahanOpdId,
			&permasalahanTerpilih.KodeOpd,
			&permasalahanTerpilih.Tahun,
		)
		if err != nil {
			return nil, err
		}
		permasalahanTerpilihs = append(permasalahanTerpilihs, permasalahanTerpilih)
	}
	return permasalahanTerpilihs, nil
}

func (repository *PermasalahanTerpilihRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_permasalahan_terpilih WHERE id = ?"
	result, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("permasalahan terpilih tidak ditemukan")
	}
	return nil
}
func (repository *PermasalahanTerpilihRepositoryImpl) FindByPermasalahanOpdId(ctx context.Context, tx *sql.Tx, permasalahanOpdId int) (domain.PermasalahanTerpilih, error) {
	script := "SELECT id, permasalahan_opd_id, kode_opd, tahun FROM tb_permasalahan_terpilih WHERE permasalahan_opd_id = ?"
	rows, err := tx.QueryContext(ctx, script, permasalahanOpdId)
	if err != nil {
		return domain.PermasalahanTerpilih{}, err
	}
	defer rows.Close()

	permasalahanTerpilih := domain.PermasalahanTerpilih{}
	if rows.Next() {
		err := rows.Scan(
			&permasalahanTerpilih.Id,
			&permasalahanTerpilih.PermasalahanOpdId,
			&permasalahanTerpilih.KodeOpd,
			&permasalahanTerpilih.Tahun,
		)
		if err != nil {
			return domain.PermasalahanTerpilih{}, err
		}
	}
	return permasalahanTerpilih, nil
}
