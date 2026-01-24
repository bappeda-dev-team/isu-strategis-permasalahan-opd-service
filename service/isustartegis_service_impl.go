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

	// 🔥 PERBAIKAN: Validasi duplikat berdasarkan permasalahan_terpilih_id, bukan permasalahan_opd_id
	permasalahanTerpilihIds := make(map[int]bool)

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

		// 🔥 VALIDASI DUPLIKAT DI SINI (SETELAH DAPAT permasalahanTerpilih.Id)
		if permasalahanTerpilihIds[permasalahanTerpilih.Id] {
			return web.IsuStrategisResponse{}, fmt.Errorf("permasalahan tidak boleh sama dalam 1(satu) isu strategis")
		}
		permasalahanTerpilihIds[permasalahanTerpilih.Id] = true

		// 2. Validasi permasalahan_opd exists (untuk ambil detail)
		permasalahan, err := service.PermasalahanRepository.FindById(ctx, tx, strconv.Itoa(p.IdPermasalahan))
		if err != nil {
			return web.IsuStrategisResponse{}, fmt.Errorf("permasalahan dengan ID %d tidak ditemukan", p.IdPermasalahan)
		}
		if permasalahan.Id == 0 {
			return web.IsuStrategisResponse{}, fmt.Errorf("permasalahan dengan ID %d tidak ditemukan", p.IdPermasalahan)
		}

		dataDukung := make([]domain.DataDukung, len(p.DataDukung))
		for j, dd := range p.DataDukung {
			// DEDUPLIKASI: Ambil jumlah data dan deduplikasi berdasarkan tahun
			jumlahData := make([]domain.JumlahData, len(dd.JumlahData))
			for k, jd := range dd.JumlahData {
				jumlahData[k] = domain.JumlahData{
					Tahun:      jd.Tahun,
					JumlahData: jd.JumlahData,
					Satuan:     jd.Satuan,
				}
			}

			// DEDUPLIKASI DI SINI
			jumlahData = deduplicateJumlahData(jumlahData)

			dataDukung[j] = domain.DataDukung{
				DataDukung:       dd.DataDukung,
				NarasiDataDukung: dd.NarasiDataDukung,
				JumlahData:       jumlahData,
			}
		}

		permasalahanOpd[i] = domain.Permasalahan{
			Id:           permasalahanTerpilih.Id, // GUNAKAN ID PERMASALAHAN_TERPILIH
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
		TahunAwal:        "",
		TahunAkhir:       "",
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

	// Validasi duplikat permasalahan terpilih dalam request
	permasalahanTerpilihIds := make(map[int]bool)
	for _, p := range request.PermasalahanOpd {
		if p.PermasalahanOpdId == 0 {
			log.Printf("Warning: Found permasalahan with ID 0 in request")
			continue
		}

		if permasalahanTerpilihIds[p.PermasalahanOpdId] {
			log.Printf("ERROR: Duplicate permasalahan ID %d found in request", p.PermasalahanOpdId)
			return web.IsuStrategisResponse{}, fmt.Errorf("permasalahan tidak boleh sama dalam 1(satu) isu strategis (ID duplikat: %d)", p.PermasalahanOpdId)
		}
		permasalahanTerpilihIds[p.PermasalahanOpdId] = true
	}

	log.Printf("Validation passed: %d unique permasalahan in request", len(permasalahanTerpilihIds))

	// 2. Update isu strategis basic info
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

	// 3. Process permasalahan
	permasalahanOpd := make([]domain.Permasalahan, 0)
	for i, p := range request.PermasalahanOpd {
		log.Printf("========================================")
		log.Printf("Processing permasalahan #%d", i+1)
		log.Printf("Request permasalahan_opd_id: %d", p.PermasalahanOpdId)
		log.Printf("========================================")

		if p.PermasalahanOpdId == 0 {
			log.Printf("Skipping permasalahan #%d: ID is 0", i+1)
			continue
		}

		// 🔥 PERBAIKAN: Gunakan FindByPermasalahanOpdId (seperti di CREATE)
		// p.PermasalahanOpdId adalah ID dari tb_permasalahan_opd
		permasalahanTerpilih, err := service.PermasalahanTerpilihRepository.FindByPermasalahanOpdId(ctx, tx, p.PermasalahanOpdId)
		if err != nil {
			log.Printf("ERROR finding permasalahan terpilih with permasalahan_opd_id %d: %v", p.PermasalahanOpdId, err)
			return web.IsuStrategisResponse{}, fmt.Errorf("error mencari permasalahan terpilih dengan permasalahan_opd_id %d: %v", p.PermasalahanOpdId, err)
		}
		if permasalahanTerpilih.Id == 0 {
			log.Printf("ERROR: Permasalahan dengan permasalahan_opd_id %d belum dipilih", p.PermasalahanOpdId)
			return web.IsuStrategisResponse{}, fmt.Errorf("permasalahan dengan ID %d belum dipilih sebagai permasalahan terpilih", p.PermasalahanOpdId)
		}

		log.Printf("✅ Found permasalahan_terpilih ID %d for permasalahan_opd_id %d", permasalahanTerpilih.Id, p.PermasalahanOpdId)

		// Get permasalahan_opd_id dari permasalahan_terpilih (sebenarnya sama dengan p.PermasalahanOpdId)
		permasalahanOpdId := permasalahanTerpilih.PermasalahanOpdId

		// 4. Get existing data dukung untuk permasalahan ini dan isu strategis ini
		existingDataDukung, err := service.IsuStrategisRepository.FindDataDukungByPermasalahanIdAndIsuStrategisId(ctx, tx, permasalahanOpdId, request.Id)
		if err != nil {
			log.Printf("Error getting existing data dukung: %v", err)
			return web.IsuStrategisResponse{}, err
		}

		log.Printf("Found %d existing data dukung for permasalahan_opd_id %d", len(existingDataDukung), permasalahanOpdId)

		// Buat map untuk tracking data dukung yang akan dipertahankan
		keepDataDukungIds := make(map[int]bool)
		for _, dd := range p.DataDukung {
			if dd.Id > 0 {
				keepDataDukungIds[dd.Id] = true
			}
		}

		// Hapus data dukung yang ID-nya tidak ada di request (SATU PER SATU)
		for _, dd := range existingDataDukung {
			if !keepDataDukungIds[dd.Id] {
				log.Printf("Deleting data dukung ID %d (not in request)", dd.Id)

				// Hapus jumlah data dulu (via repository)
				err := service.IsuStrategisRepository.DeleteJumlahDataByDataDukungId(ctx, tx, dd.Id)
				if err != nil {
					log.Printf("Error deleting jumlah data for data dukung %d: %v", dd.Id, err)
					return web.IsuStrategisResponse{}, err
				}

				// Hapus data dukung spesifik by ID
				err = service.IsuStrategisRepository.DeleteDataDukungById(ctx, tx, dd.Id)
				if err != nil {
					log.Printf("Error deleting data dukung %d: %v", dd.Id, err)
					return web.IsuStrategisResponse{}, err
				}
				log.Printf("Successfully deleted data dukung ID %d", dd.Id)
			}
		}

		// 5. Process data dukung
		dataDukung := make([]domain.DataDukung, 0)
		for _, dd := range p.DataDukung {
			// Process jumlah data dari request
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

			// DEDUPLIKASI: Hapus duplicate tahun dalam 1 data dukung
			jumlahData = deduplicateJumlahData(jumlahData)

			if dd.DataDukung != "" {
				dataDukung = append(dataDukung, domain.DataDukung{
					Id:                dd.Id,
					PermasalahanOpdId: permasalahanOpdId,
					IdIsuStrategis:    request.Id,
					DataDukung:        dd.DataDukung,
					NarasiDataDukung:  dd.NarasiDataDukung,
					JumlahData:        jumlahData,
				})
			}
		}

		log.Printf("Appending permasalahan #%d (ID: %d) with %d data dukung", i+1, p.PermasalahanOpdId, len(dataDukung))

		permasalahanOpd = append(permasalahanOpd, domain.Permasalahan{
			Id:         permasalahanTerpilih.Id,
			DataDukung: dataDukung,
		})
	}

	log.Printf("Total permasalahan processed: %d", len(permasalahanOpd))

	isuStrategis.PermasalahanOpd = permasalahanOpd

	// 6. Update ke database
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

func (service *IsuStrategisServiceImpl) deleteDataDukungAndJumlahData(ctx context.Context, tx *sql.Tx, permasalahanTerpilihId int, isuStrategisId int) error {
	// 🔥 PERBAIKAN: Dapatkan permasalahan_opd_id dari permasalahan_terpilih
	permasalahanTerpilih, err := service.PermasalahanTerpilihRepository.FindById(ctx, tx, permasalahanTerpilihId)
	if err != nil {
		return err
	}
	if permasalahanTerpilih.Id == 0 {
		return fmt.Errorf("permasalahan terpilih ID %d tidak ditemukan", permasalahanTerpilihId)
	}

	permasalahanOpdId := permasalahanTerpilih.PermasalahanOpdId

	// 1. Dapatkan semua data dukung untuk permasalahan ini dan isu strategis ini
	dataDukungs, err := service.IsuStrategisRepository.FindDataDukungByPermasalahanIdAndIsuStrategisId(ctx, tx, permasalahanOpdId, isuStrategisId)
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
	err = service.IsuStrategisRepository.DeleteDataDukungByPermasalahanAndIsuStrategis(ctx, tx, permasalahanOpdId, isuStrategisId)
	if err != nil {
		return err
	}

	return nil
}

func (service *IsuStrategisServiceImpl) FindallIsuKebelakang(ctx context.Context, kodeOpd string, tahun string) ([]web.IsuStrategisKebelakangResponse, error) {
	// Logging awal
	fmt.Printf("[Service] FindallIsuKebelakang - Start with params: kodeOpd=%s, tahun=%s\n", kodeOpd, tahun)

	// Validasi input
	if kodeOpd == "" {
		fmt.Println("[Service] FindallIsuKebelakang - kodeOpd is empty")
		return []web.IsuStrategisKebelakangResponse{}, nil
	}

	// Buat context dengan timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Mulai transaksi dengan retry mechanism
	var isuStrategiss []web.IsuStrategisKebelakangResponse
	maxRetries := 3
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Printf("[Service] FindallIsuKebelakang - Attempt %d of %d\n", attempt, maxRetries)

		tx, err := service.DB.BeginTx(ctxWithTimeout, &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  true,
		})
		if err != nil {
			lastErr = fmt.Errorf("error starting transaction: %v", err)
			if attempt == maxRetries {
				fmt.Printf("[Service] FindallIsuKebelakang - Final attempt failed: %v\n", lastErr)
				return nil, lastErr
			}
			time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
			continue
		}

		// Eksekusi query
		result, err := service.IsuStrategisRepository.FindallIsuKebelakang(ctxWithTimeout, tx, kodeOpd, tahun)
		if err != nil {
			tx.Rollback()
			lastErr = err
			if attempt == maxRetries {
				fmt.Printf("[Service] FindallIsuKebelakang - Final attempt failed: %v\n", lastErr)
				return nil, fmt.Errorf("error after %d retries: %v", maxRetries, lastErr)
			}
			time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
			continue
		}

		// Commit transaksi
		if err := tx.Commit(); err != nil {
			lastErr = fmt.Errorf("error committing transaction: %v", err)
			if attempt == maxRetries {
				fmt.Printf("[Service] FindallIsuKebelakang - Final attempt failed: %v\n", lastErr)
				return nil, lastErr
			}
			time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
			continue
		}

		// Sort hasil berdasarkan created_at sebelum konversi ke response
		sort.Slice(result, func(i, j int) bool {
			return result[i].CreatedAt.Before(result[j].CreatedAt)
		})

		// Konversi ke response dengan tahun sekarang
		isuStrategiss = helper.ToIsuStrategisKebelakangResponses(result, tahun)
		fmt.Printf("[Service] FindallIsuKebelakang - Successfully retrieved and sorted %d records\n", len(isuStrategiss))
		break
	}

	// Handle empty result
	if len(isuStrategiss) == 0 {
		fmt.Println("[Service] FindallIsuKebelakang - No records found")
		return []web.IsuStrategisKebelakangResponse{}, nil
	}

	fmt.Println("[Service] FindallIsuKebelakang - Completed successfully")
	return isuStrategiss, nil
}

func deduplicateJumlahData(jumlahData []domain.JumlahData) []domain.JumlahData {
	// Map untuk tracking: tahun -> jumlah data
	dataMap := make(map[string]domain.JumlahData)

	// Loop data, yang terakhir akan overwrite yang sebelumnya
	for _, jd := range jumlahData {
		if jd.Tahun != "" { // Hanya process jika tahun tidak kosong
			dataMap[jd.Tahun] = jd
		}
	}

	// Convert map kembali ke slice
	result := make([]domain.JumlahData, 0, len(dataMap))
	for _, jd := range dataMap {
		result = append(result, jd)
	}

	// Sort berdasarkan tahun DESC untuk konsistensi
	sort.Slice(result, func(i, j int) bool {
		return result[i].Tahun > result[j].Tahun
	})

	return result
}
