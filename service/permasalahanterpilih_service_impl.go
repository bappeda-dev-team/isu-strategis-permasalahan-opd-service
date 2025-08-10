package service

import (
	"context"
	"database/sql"
	"errors"
	"permasalahanService/helper"
	"permasalahanService/model/domain"
	"permasalahanService/model/web"
	"permasalahanService/repository"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type PermasalahanTerpilihServiceImpl struct {
	permasalahanTerpilihRepository repository.PermasalahanTerpilihRepository
	permasalahanRepository         repository.PermasalahanRepository
	db                             *sql.DB
	validate                       *validator.Validate
}

func NewPermasalahanTerpilihServiceImpl(
	permasalahanTerpilihRepository repository.PermasalahanTerpilihRepository,
	permasalahanRepository repository.PermasalahanRepository,
	db *sql.DB,
	validate *validator.Validate,
) *PermasalahanTerpilihServiceImpl {
	return &PermasalahanTerpilihServiceImpl{
		permasalahanTerpilihRepository: permasalahanTerpilihRepository,
		permasalahanRepository:         permasalahanRepository,
		db:                             db,
		validate:                       validate,
	}
}

func (service *PermasalahanTerpilihServiceImpl) Create(ctx context.Context, request web.PermasalahanTerpilihRequest) (web.ChildResponse, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return web.ChildResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi request
	err = service.validate.Struct(request)
	if err != nil {
		return web.ChildResponse{}, err
	}

	// Ambil data permasalahan dari tb_akar_permasalahan
	permasalahan, err := service.permasalahanRepository.FindById(ctx, tx, strconv.Itoa(request.AkarPermasalahanId))
	if err != nil {
		return web.ChildResponse{}, err
	}
	if permasalahan.Id == 0 {
		return web.ChildResponse{}, errors.New("permasalahan tidak ditemukan")
	}

	// Cek apakah permasalahan ini sudah terpilih
	existingTerpilih, err := service.permasalahanTerpilihRepository.FindByPermasalahanOpdId(ctx, tx, request.AkarPermasalahanId)
	if err != nil {
		return web.ChildResponse{}, err
	}
	if existingTerpilih.Id != 0 {
		return web.ChildResponse{}, errors.New("permasalahan ini sudah terpilih")
	}

	// Buat data permasalahan terpilih baru
	permasalahanTerpilih := domain.PermasalahanTerpilih{
		PermasalahanOpdId: request.AkarPermasalahanId,
		KodeOpd:           permasalahan.KodeOpd,
		Tahun:             permasalahan.Tahun,
	}

	// Simpan ke database
	permasalahanTerpilih, err = service.permasalahanTerpilihRepository.Create(ctx, tx, permasalahanTerpilih)
	if err != nil {
		return web.ChildResponse{}, err
	}

	// Buat response
	response := web.ChildResponse{
		Id:             permasalahanTerpilih.Id,
		IdPermasalahan: permasalahanTerpilih.PermasalahanOpdId,
		NamaPohon:      permasalahan.Permasalahan,
		LevelPohon:     permasalahan.LevelPohon,
		PerangkatDaerah: web.PerangkatDaerah{
			KodeOpd: permasalahan.KodeOpd,
			NamaOpd: permasalahan.NamaOpd,
		},
		IsPermasalahan:       true,
		PermasalahanTerpilih: true,
	}

	return response, nil
}
func (service *PermasalahanTerpilihServiceImpl) FindAll(ctx context.Context, kodeOpd string, tahun string) ([]web.ChildResponse, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	permasalahanTerpilihs, err := service.permasalahanTerpilihRepository.FindAll(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}

	responses := []web.ChildResponse{}
	for _, permasalahanTerpilih := range permasalahanTerpilihs {
		permasalahan, err := service.permasalahanRepository.FindById(ctx, tx, strconv.Itoa(permasalahanTerpilih.PermasalahanOpdId))
		if err != nil {
			return nil, err
		}
		if permasalahan.Id == 0 {
			return nil, errors.New("permasalahan tidak ditemukan")
		}
		responses = append(responses, web.ChildResponse{
			Id:         permasalahanTerpilih.Id,
			NamaPohon:  permasalahan.Permasalahan,
			LevelPohon: permasalahan.LevelPohon,
			PerangkatDaerah: web.PerangkatDaerah{
				KodeOpd: permasalahan.KodeOpd,
				NamaOpd: permasalahan.NamaOpd,
			},
			JenisMasalah: permasalahan.JenisMasalah,
		})
	}
	return responses, nil
}

func (service *PermasalahanTerpilihServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.db.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	return service.permasalahanTerpilihRepository.Delete(ctx, tx, id)
}
