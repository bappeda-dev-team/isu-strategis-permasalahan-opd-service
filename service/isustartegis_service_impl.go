package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"permasalahanService/helper"
	"permasalahanService/model/domain"
	"permasalahanService/model/web"
	"permasalahanService/repository"
	"sort"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

type IsuStrategisServiceImpl struct {
	IsuStrategisRepository         repository.IsuStrategisRepository
	PermasalahanRepository         repository.PermasalahanRepository
	PermasalahanTerpilihRepository repository.PermasalahanTerpilihRepository
	DB                             *sql.DB
	Validate                       *validator.Validate
}

func NewIsuStrategisServiceImpl(
	isuStrategisRepository repository.IsuStrategisRepository,
	permasalahanRepository repository.PermasalahanRepository,
	permasalahanTerpilihRepository repository.PermasalahanTerpilihRepository,
	db *sql.DB,
	validate *validator.Validate,
) *IsuStrategisServiceImpl {
	return &IsuStrategisServiceImpl{
		IsuStrategisRepository:         isuStrategisRepository,
		PermasalahanRepository:         permasalahanRepository,
		PermasalahanTerpilihRepository: permasalahanTerpilihRepository,
		DB:                             db,
		Validate:                       validate,
	}
}

func (service *IsuStrategisServiceImpl) Create(ctx context.Context, request web.IsuStrategisCreateRequest) (web.IsuStrategisResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return web.IsuStrategisResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return web.IsuStrategisResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Convert request ke domain
	permasalahanOpd := make([]domain.Permasalahan, len(request.PermasalahanOpd))
	for i, p := range request.PermasalahanOpd {
		// 1. Validasi di permasalahan_terpilih
		permasalahanTerpilih, err := service.PermasalahanTerpilihRepository.FindByPermasalahanOpdId(ctx, tx, p.IdPermasalahan)
		if err != nil {
			return web.IsuStrategisResponse{}, err
		}
		if permasalahanTerpilih.Id == 0 {
			return web.IsuStrategisResponse{}, fmt.Errorf("permasalahan dengan ID %d belum dipilih sebagai permasalahan terpilih", p.IdPermasalahan)
		}

		// 2. Validasi isu_strategis_id di permasalahan_opd
		permasalahan, err := service.PermasalahanRepository.FindById(ctx, tx, strconv.Itoa(p.IdPermasalahan))
		if err != nil {
			return web.IsuStrategisResponse{}, fmt.Errorf("permasalahan dengan ID %d tidak ditemukan", p.IdPermasalahan)
		}
		if permasalahan.Id == 0 {
			return web.IsuStrategisResponse{}, fmt.Errorf("permasalahan dengan ID %d tidak ditemukan", p.IdPermasalahan)
		}

		// Cek apakah sudah memiliki isu_strategis_id
		if permasalahan.IsuStrategis != 0 {
			return web.IsuStrategisResponse{}, fmt.Errorf("permasalahan dengan ID %d sudah memiliki isu strategis", p.IdPermasalahan)
		}

		dataDukung := make([]domain.DataDukung, len(p.DataDukung))
		for j, dd := range p.DataDukung {
			jumlahData := make([]domain.JumlahData, len(dd.JumlahData))
			for k, jd := range dd.JumlahData {
				jumlahData[k] = domain.JumlahData{
					Tahun:      jd.Tahun,
					JumlahData: jd.JumlahData,
					Satuan:     jd.Satuan,
				}
			}

			dataDukung[j] = domain.DataDukung{
				DataDukung:       dd.DataDukung,
				NarasiDataDukung: dd.NarasiDataDukung,
				JumlahData:       jumlahData,
			}
		}

		permasalahanOpd[i] = domain.Permasalahan{
			Id:           permasalahan.Id,
			Tahun:        permasalahan.Tahun,
			NamaOpd:      permasalahan.NamaOpd,
			LevelPohon:   permasalahan.LevelPohon,
			JenisMasalah: permasalahan.JenisMasalah,
			DataDukung:   dataDukung,
		}
	}

	isuStrategis := domain.IsuStrategis{
		KodeOpd:          request.KodeOpd,
		NamaOpd:          request.NamaOpd,
		KodeBidangUrusan: request.KodeBidangUrusan,
		NamaBidangUrusan: request.NamaBidangUrusan,
		TahunAwal:        request.TahunAwal,
		TahunAkhir:       request.TahunAkhir,
		IsuStrategis:     request.IsuStrategis,
		PermasalahanOpd:  permasalahanOpd,
	}

	isuStrategis, err = service.IsuStrategisRepository.Create(ctx, tx, isuStrategis)
	if err != nil {
		return web.IsuStrategisResponse{}, err
	}

	return helper.ToIsuStrategisResponse(isuStrategis), nil
}

func (service *IsuStrategisServiceImpl) Update(ctx context.Context, request web.IsuStrategisUpdateRequest) (web.IsuStrategisResponse, error) {
	log.Printf("Start updating isu strategis with ID: %d", request.Id)

	// Validasi basic fields
	if request.Id == 0 {
		return web.IsuStrategisResponse{}, fmt.Errorf("id is required")
	}

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return web.IsuStrategisResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// 1. Validasi isu strategis exists
	existingIsuStrategis, err := service.IsuStrategisRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		log.Printf("Error finding isu strategis ID %d: %v", request.Id, err)
		return web.IsuStrategisResponse{}, err
	}
	if existingIsuStrategis.Id == 0 {
		log.Printf("Isu strategis ID %d not found", request.Id)
		return web.IsuStrategisResponse{}, fmt.Errorf("isu strategis dengan ID %d tidak ditemukan", request.Id)
	}

	// 2. Dapatkan semua permasalahan yang terkait dengan isu strategis ini
	existingPermasalahan, err := service.PermasalahanRepository.FindByIsuStrategisId(ctx, tx, request.Id)
	if err != nil {
		log.Printf("Error getting existing permasalahan: %v", err)
		return web.IsuStrategisResponse{}, err
	}

	// Buat map untuk tracking permasalahan yang akan dipertahankan
	keepPermasalahanIds := make(map[int]bool)
	for _, p := range request.PermasalahanOpd {
		keepPermasalahanIds[p.PermasalahanOpdId] = true
	}

	// Reset isu_strategis_id menjadi 0 untuk permasalahan yang tidak ada di request
	// dan hapus data dukung beserta jumlah datanya
	for _, p := range existingPermasalahan {
		if !keepPermasalahanIds[p.Id] {
			// Hapus data dukung dan jumlah data terlebih dahulu
			err := service.deleteDataDukungAndJumlahData(ctx, tx, p.Id)
			if err != nil {
				log.Printf("Error deleting data dukung and jumlah data for permasalahan %d: %v", p.Id, err)
				return web.IsuStrategisResponse{}, err
			}

			// Kemudian reset isu_strategis_id
			err = service.PermasalahanRepository.ResetIsuStrategisId(ctx, tx, p.Id)
			if err != nil {
				log.Printf("Error resetting isu_strategis_id for permasalahan %d: %v", p.Id, err)
				return web.IsuStrategisResponse{}, err
			}
		}
	}

	// 3. Update isu strategis basic info
	isuStrategis := domain.IsuStrategis{
		Id:               request.Id,
		KodeOpd:          request.KodeOpd,
		NamaOpd:          request.NamaOpd,
		KodeBidangUrusan: request.KodeBidangUrusan,
		NamaBidangUrusan: request.NamaBidangUrusan,
		TahunAwal:        request.TahunAwal,
		TahunAkhir:       request.TahunAkhir,
		IsuStrategis:     request.IsuStrategis,
	}

	// 4. Process permasalahan
	permasalahanOpd := make([]domain.Permasalahan, 0)
	for _, p := range request.PermasalahanOpd {
		if p.PermasalahanOpdId == 0 {
			continue
		}

		// Validasi permasalahan exists
		permasalahan, err := service.PermasalahanRepository.FindById(ctx, tx, strconv.Itoa(p.PermasalahanOpdId))
		if err != nil {
			log.Printf("Error finding permasalahan ID %d: %v", p.PermasalahanOpdId, err)
			continue
		}
		if permasalahan.Id == 0 {
			log.Printf("Permasalahan ID %d not found", p.PermasalahanOpdId)
			continue
		}

		// 5. Get existing data dukung untuk permasalahan ini
		existingDataDukung, err := service.IsuStrategisRepository.FindDataDukungByPermasalahanId(ctx, tx, p.PermasalahanOpdId)
		if err != nil {
			log.Printf("Error getting existing data dukung: %v", err)
			return web.IsuStrategisResponse{}, err
		}

		// Buat map untuk tracking data dukung yang akan dipertahankan
		keepDataDukungIds := make(map[int]bool)
		for _, dd := range p.DataDukung {
			keepDataDukungIds[dd.Id] = true
		}

		// Hapus data dukung yang tidak ada di request
		for _, dd := range existingDataDukung {
			if !keepDataDukungIds[dd.Id] {
				err := service.IsuStrategisRepository.DeleteDataDukungByPermasalahanId(ctx, tx, p.PermasalahanOpdId)
				if err != nil {
					log.Printf("Error deleting data dukung %d: %v", dd.Id, err)
					return web.IsuStrategisResponse{}, err
				}
			}
		}

		// 6. Process data dukung
		dataDukung := make([]domain.DataDukung, 0)
		for _, dd := range p.DataDukung {
			// Get existing jumlah data jika data dukung sudah ada
			var existingJumlahData []domain.JumlahData
			if dd.Id != 0 {
				existingJumlahData, err = service.IsuStrategisRepository.FindJumlahDataByDataDukungId(ctx, tx, dd.Id)
				if err != nil {
					log.Printf("Error getting existing jumlah data: %v", err)
					return web.IsuStrategisResponse{}, err
				}

				// Buat map untuk tracking jumlah data yang akan dipertahankan
				keepJumlahDataIds := make(map[int]bool)
				for _, jd := range dd.JumlahData {
					keepJumlahDataIds[jd.Id] = true
				}

				// Hapus jumlah data yang tidak ada di request
				for _, jd := range existingJumlahData {
					if !keepJumlahDataIds[jd.Id] {
						err := service.IsuStrategisRepository.DeleteJumlahDataByDataDukungId(ctx, tx, dd.Id)
						if err != nil {
							log.Printf("Error deleting jumlah data %d: %v", jd.Id, err)
							return web.IsuStrategisResponse{}, err
						}
					}
				}
			}

			// Process jumlah data
			jumlahData := make([]domain.JumlahData, 0)
			for _, jd := range dd.JumlahData {
				if jd.Tahun != "" && jd.Satuan != "" {
					jumlahData = append(jumlahData, domain.JumlahData{
						Id:         jd.Id,
						Tahun:      jd.Tahun,
						JumlahData: jd.JumlahData,
						Satuan:     jd.Satuan,
					})
				}
			}

			if dd.DataDukung != "" {
				dataDukung = append(dataDukung, domain.DataDukung{
					Id:                dd.Id,
					PermasalahanOpdId: p.PermasalahanOpdId,
					DataDukung:        dd.DataDukung,
					NarasiDataDukung:  dd.NarasiDataDukung,
					JumlahData:        jumlahData,
				})
			}
		}

		permasalahanOpd = append(permasalahanOpd, domain.Permasalahan{
			Id:         permasalahan.Id,
			DataDukung: dataDukung,
		})
	}

	isuStrategis.PermasalahanOpd = permasalahanOpd

	// 7. Update ke database
	updatedIsuStrategis, err := service.IsuStrategisRepository.Update(ctx, tx, isuStrategis)
	if err != nil {
		log.Printf("Error updating isu strategis: %v", err)
		return web.IsuStrategisResponse{}, err
	}

	return helper.ToIsuStrategisResponse(updatedIsuStrategis), nil
}

func (service *IsuStrategisServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi apakah isu strategis ada
	isuStrategis, err := service.IsuStrategisRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}
	if isuStrategis.Id == 0 {
		return fmt.Errorf("isu strategis dengan ID %d tidak ditemukan", id)
	}

	// Proses delete
	err = service.IsuStrategisRepository.Delete(ctx, tx, id)
	if err != nil {
		return err
	}

	return nil
}

func (service *IsuStrategisServiceImpl) FindById(ctx context.Context, id int) (web.IsuStrategisResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return web.IsuStrategisResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	isuStrategis, err := service.IsuStrategisRepository.FindById(ctx, tx, id)
	if err != nil {
		return web.IsuStrategisResponse{}, err
	}

	return helper.ToIsuStrategisResponse(isuStrategis), nil
}

func (service *IsuStrategisServiceImpl) FindAll(ctx context.Context, kodeOpd string, tahunAwal string, tahunAkhir string) ([]web.IsuStrategisResponse, error) {
	// Logging awal
	fmt.Printf("[Service] FindAll - Start with params: kodeOpd=%s, tahunAwal=%s, tahunAkhir=%s\n", kodeOpd, tahunAwal, tahunAkhir)

	// Validasi input
	if kodeOpd == "" {
		fmt.Println("[Service] FindAll - kodeOpd is empty")
		return []web.IsuStrategisResponse{}, nil
	}

	// Set default tahun jika kosong
	if tahunAwal == "" && tahunAkhir == "" {
		currentYear := time.Now().Year()
		tahunAwal = fmt.Sprintf("%d", currentYear)
		tahunAkhir = fmt.Sprintf("%d", currentYear)
		fmt.Printf("[Service] FindAll - Using default year: %s\n", tahunAwal)
	}

	// Buat context dengan timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Mulai transaksi dengan retry mechanism
	var isuStrategiss []web.IsuStrategisResponse
	maxRetries := 3
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Printf("[Service] FindAll - Attempt %d of %d\n", attempt, maxRetries)

		tx, err := service.DB.BeginTx(ctxWithTimeout, &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  true,
		})
		if err != nil {
			lastErr = fmt.Errorf("error starting transaction: %v", err)
			if attempt == maxRetries {
				fmt.Printf("[Service] FindAll - Final attempt failed: %v\n", lastErr)
				return nil, lastErr
			}
			time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
			continue
		}

		// Eksekusi query
		result, err := service.IsuStrategisRepository.FindAll(ctxWithTimeout, tx, kodeOpd, tahunAwal, tahunAkhir)
		if err != nil {
			tx.Rollback()
			lastErr = err
			if attempt == maxRetries {
				fmt.Printf("[Service] FindAll - Final attempt failed: %v\n", lastErr)
				return nil, fmt.Errorf("error after %d retries: %v", maxRetries, lastErr)
			}
			time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
			continue
		}

		// Commit transaksi
		if err := tx.Commit(); err != nil {
			lastErr = fmt.Errorf("error committing transaction: %v", err)
			if attempt == maxRetries {
				fmt.Printf("[Service] FindAll - Final attempt failed: %v\n", lastErr)
				return nil, lastErr
			}
			time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
			continue
		}

		// Sort hasil berdasarkan created_at sebelum konversi ke response
		sort.Slice(result, func(i, j int) bool {
			return result[i].CreatedAt.Before(result[j].CreatedAt)
		})

		// Konversi ke response
		isuStrategiss = helper.ToIsuStrategisResponses(result)
		fmt.Printf("[Service] FindAll - Successfully retrieved and sorted %d records\n", len(isuStrategiss))
		break
	}

	// Handle empty result
	if len(isuStrategiss) == 0 {
		fmt.Println("[Service] FindAll - No records found")
		return []web.IsuStrategisResponse{}, nil
	}

	fmt.Println("[Service] FindAll - Completed successfully")
	return isuStrategiss, nil
}

func (service *IsuStrategisServiceImpl) deleteDataDukungAndJumlahData(ctx context.Context, tx *sql.Tx, permasalahanId int) error {
	// 1. Dapatkan semua data dukung untuk permasalahan ini
	dataDukungs, err := service.IsuStrategisRepository.FindDataDukungByPermasalahanId(ctx, tx, permasalahanId)
	if err != nil {
		return err
	}

	// 2. Untuk setiap data dukung, hapus jumlah data terlebih dahulu
	for _, dd := range dataDukungs {
		// Hapus semua jumlah data
		err = service.IsuStrategisRepository.DeleteJumlahDataByDataDukungId(ctx, tx, dd.Id)
		if err != nil {
			return err
		}
	}

	// 3. Setelah semua jumlah data dihapus, hapus data dukung
	err = service.IsuStrategisRepository.DeleteDataDukungByPermasalahanId(ctx, tx, permasalahanId)
	if err != nil {
		return err
	}

	return nil
}
