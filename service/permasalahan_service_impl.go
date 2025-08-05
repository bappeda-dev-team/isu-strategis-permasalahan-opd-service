package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"permasalahanService/helper"
	"permasalahanService/model/domain"
	"permasalahanService/model/web"
	"permasalahanService/repository"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type PermasalahanServiceImpl struct {
	permasalahanRepository repository.PermasalahanRepository
	db                     *sql.DB
	validate               *validator.Validate
}

func NewPermasalahanServiceImpl(permasalahanRepository repository.PermasalahanRepository, db *sql.DB, validate *validator.Validate) *PermasalahanServiceImpl {
	return &PermasalahanServiceImpl{
		permasalahanRepository: permasalahanRepository,
		db:                     db,
		validate:               validate,
	}
}

func (service *PermasalahanServiceImpl) Create(ctx context.Context, request web.PermasalahanCreateRequest) (web.ChildResponse, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return web.ChildResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi jenis masalah
	jenisMasalah := domain.JenisMasalah(request.JenisMasalah)
	if !jenisMasalah.IsValid() {
		return web.ChildResponse{}, errors.New("jenis masalah tidak valid. Pilihan yang tersedia: MASALAH_POKOK, MASALAH, AKAR_MASALAH")
	}

	existingPermasalahan, err := service.permasalahanRepository.FindByPokinId(ctx, tx, request.PokinId)
	if err != nil {
		return web.ChildResponse{}, err
	}
	if existingPermasalahan.Id != 0 {
		return web.ChildResponse{}, errors.New("pokin_id sudah digunakan")
	}

	permasalahan := domain.Permasalahan{
		PokinId:      request.PokinId,
		Permasalahan: request.Permasalahan,
		LevelPohon:   request.LevelPohon,
		KodeOpd:      request.KodeOpd,
		NamaOpd:      request.NamaOpd,
		Tahun:        request.Tahun,
		JenisMasalah: string(jenisMasalah),
	}

	permasalahan, err = service.permasalahanRepository.Create(ctx, tx, permasalahan)
	if err != nil {
		return web.ChildResponse{}, err
	}

	permasalahanResponse := web.ChildResponse{
		Id:             permasalahan.PokinId,
		IdPermasalahan: permasalahan.Id,
		NamaPohon:      permasalahan.Permasalahan,
		LevelPohon:     permasalahan.LevelPohon,
		PerangkatDaerah: web.PerangkatDaerah{
			KodeOpd: permasalahan.KodeOpd,
			NamaOpd: permasalahan.NamaOpd,
		},
		JenisMasalah:   permasalahan.JenisMasalah,
		IsPermasalahan: true,
	}

	return permasalahanResponse, nil
}

func (service *PermasalahanServiceImpl) Update(ctx context.Context, request web.PermasalahanUpdateRequest) (web.ChildResponse, error) {
	tx, err := service.db.Begin()
	if err != nil {
		log.Printf("Error begin transaction: %v", err)
		return web.ChildResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Log request yang diterima
	log.Printf("Updating permasalahan with ID: %d", request.Id)
	log.Printf("Request data: %+v", request)

	// Cari data existing
	existing, err := service.permasalahanRepository.FindById(ctx, tx, strconv.Itoa(request.Id))
	if err != nil {
		log.Printf("Error finding permasalahan: %v", err)
		return web.ChildResponse{}, err
	}

	// Update field yang dikirim dalam request
	existing.Permasalahan = request.Permasalahan
	existing.KodeOpd = request.KodeOpd
	existing.NamaOpd = request.NamaOpd
	existing.Tahun = request.Tahun

	// Log data yang akan diupdate
	log.Printf("Data to update: %+v", existing)

	updated := service.permasalahanRepository.Update(ctx, tx, existing)

	// Log hasil update
	log.Printf("Update result: %+v", updated)

	if updated.Id == 0 {
		log.Printf("Failed to update permasalahan, got empty result")
		return web.ChildResponse{}, errors.New("failed to update permasalahan")
	}

	permasalahanResponse := web.ChildResponse{
		Id:             updated.PokinId,
		IdPermasalahan: updated.Id,
		NamaPohon:      updated.Permasalahan,
		LevelPohon:     updated.LevelPohon,
		PerangkatDaerah: web.PerangkatDaerah{
			KodeOpd: updated.KodeOpd,
			NamaOpd: updated.NamaOpd,
		},
		IsPermasalahan: true,
		JenisMasalah:   updated.JenisMasalah,
	}

	return permasalahanResponse, nil
}

func (service *PermasalahanServiceImpl) Delete(ctx context.Context, id string) error {
	tx, err := service.db.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	err = service.permasalahanRepository.Delete(ctx, tx, id)
	if err != nil {
		return err
	}

	return nil
}

func (service *PermasalahanServiceImpl) FindById(ctx context.Context, id string) (web.PermasalahanResponsesbyId, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return web.PermasalahanResponsesbyId{}, err
	}
	defer helper.CommitOrRollback(tx)

	permasalahan, err := service.permasalahanRepository.FindById(ctx, tx, id)
	if err != nil {
		return web.PermasalahanResponsesbyId{}, err
	}

	permasalahanResponse := web.PermasalahanResponsesbyId{
		Id:         permasalahan.Id,
		NamaPohon:  permasalahan.Permasalahan,
		LevelPohon: permasalahan.LevelPohon,
	}

	return permasalahanResponse, nil
}

func (service *PermasalahanServiceImpl) FindAllPohonKinerja(ctx context.Context, kodeOpd string, tahun string) (*web.PohonKinerjaDataResponse, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data permasalahan dari database
	permasalahans, err := service.permasalahanRepository.FindByKodeOpdAndTahun(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}

	// Ambil data pohon kinerja dari API
	pohonKinerja, err := service.permasalahanRepository.GetPohonKinerjaFromAPI(ctx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}

	// Gabungkan data pohon kinerja dengan permasalahan
	result := service.permasalahanRepository.MergePohonKinerjaWithPermasalahan(ctx, tx, pohonKinerja, permasalahans)

	return result, nil
}
