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

	// 1. Get all permasalahan terpilih
	permasalahanTerpilihs, err := service.permasalahanTerpilihRepository.FindAll(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}

	if len(permasalahanTerpilihs) == 0 {
		return []web.ChildResponse{}, nil
	}

	// 2. OPTIMASI: Collect all permasalahan IDs untuk batch query
	permasalahanIds := make([]int, 0, len(permasalahanTerpilihs))
	for _, pt := range permasalahanTerpilihs {
		permasalahanIds = append(permasalahanIds, pt.PermasalahanOpdId)
	}

	// 3. OPTIMASI: Batch query - ambil semua permasalahan sekaligus
	permasalahans, err := service.permasalahanRepository.FindByIds(ctx, tx, permasalahanIds)
	if err != nil {
		return nil, err
	}

	// 4. Buat map untuk quick lookup
	permasalahanMap := make(map[int]domain.Permasalahan)
	for _, p := range permasalahans {
		permasalahanMap[p.Id] = p
	}

	// 5. Build responses
	responses := make([]web.ChildResponse, 0, len(permasalahanTerpilihs))
	for _, permasalahanTerpilih := range permasalahanTerpilihs {
		permasalahan, exists := permasalahanMap[permasalahanTerpilih.PermasalahanOpdId]
		if !exists {
			// Skip jika permasalahan tidak ditemukan
			continue
		}

		responses = append(responses, web.ChildResponse{
			Id:             permasalahanTerpilih.Id,
			IdPermasalahan: permasalahanTerpilih.PermasalahanOpdId,
			NamaPohon:      permasalahan.Permasalahan,
			LevelPohon:     permasalahan.LevelPohon,
			PerangkatDaerah: web.PerangkatDaerah{
				KodeOpd: permasalahan.KodeOpd,
				NamaOpd: permasalahan.NamaOpd,
			},
			JenisMasalah: permasalahan.JenisMasalah,
			Status:       permasalahan.StatusPermasalahan,
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
