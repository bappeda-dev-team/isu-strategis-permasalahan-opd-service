package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"permasalahanService/model/domain"
	"strconv"
	"time"
)

type IsuStrategisRepositoryImpl struct {
}

func NewIsuStrategisRepositoryImpl() *IsuStrategisRepositoryImpl {
	return &IsuStrategisRepositoryImpl{}
}

func (repository *IsuStrategisRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, isuStrategis domain.IsuStrategis) (domain.IsuStrategis, error) {
	// Insert isu strategis
	script := `INSERT INTO tb_isu_strategis_opd 
               (kode_opd, nama_opd, kode_bidang_urusan, nama_bidang_urusan, tahun_awal, tahun_akhir, isu_strategis) 
               VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := tx.ExecContext(ctx, script,
		isuStrategis.KodeOpd,
		isuStrategis.NamaOpd,
		isuStrategis.KodeBidangUrusan,
		isuStrategis.NamaBidangUrusan,
		isuStrategis.TahunAwal,
		isuStrategis.TahunAkhir,
		isuStrategis.IsuStrategis,
	)
	if err != nil {
		return domain.IsuStrategis{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.IsuStrategis{}, err
	}
	isuStrategis.Id = int(id)

	// Update permasalahan
	for _, permasalahan := range isuStrategis.PermasalahanOpd {
		// Cek ulang isu_strategis_id sebelum update
		var currentIsuStrategisId int
		script = `SELECT isu_strategis_id FROM tb_permasalahan_opd WHERE id = ?`
		err = tx.QueryRowContext(ctx, script, permasalahan.Id).Scan(&currentIsuStrategisId)
		if err != nil {
			return domain.IsuStrategis{}, err
		}

		// Jika isu_strategis_id tidak 0, batalkan proses
		if currentIsuStrategisId != 0 {
			return domain.IsuStrategis{}, fmt.Errorf("permasalahan dengan ID %d sudah memiliki isu strategis", permasalahan.Id)
		}

		// Jika aman, lakukan update
		script = `UPDATE tb_permasalahan_opd SET isu_strategis_id = ? WHERE id = ? AND isu_strategis_id = 0`
		result, err = tx.ExecContext(ctx, script, isuStrategis.Id, permasalahan.Id)
		if err != nil {
			return domain.IsuStrategis{}, err
		}

		// Pastikan update berhasil
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return domain.IsuStrategis{}, err
		}
		if rowsAffected == 0 {
			return domain.IsuStrategis{}, fmt.Errorf("permasalahan dengan ID %d sudah memiliki isu strategis atau tidak ditemukan", permasalahan.Id)
		}

		// Insert data dukung
		for _, dataDukung := range permasalahan.DataDukung {
			script = `INSERT INTO tb_data_dukung 
                      (id_permasalahan, nama_data_dukung, narasi_data_dukung) 
                      VALUES (?, ?, ?)`

			result, err = tx.ExecContext(ctx, script,
				permasalahan.Id,
				dataDukung.DataDukung,
				dataDukung.NarasiDataDukung,
			)
			if err != nil {
				return domain.IsuStrategis{}, err
			}

			dataDukungId, err := result.LastInsertId()
			if err != nil {
				return domain.IsuStrategis{}, err
			}

			// Insert jumlah data
			for _, jumlahData := range dataDukung.JumlahData {
				if jumlahData.Tahun != "" { // hanya insert jika ada tahun
					script = `INSERT INTO tb_jumlah_data 
                              (id_data_dukung, tahun, jumlah, satuan) 
                              VALUES (?, ?, ?, ?)`

					_, err = tx.ExecContext(ctx, script,
						dataDukungId,
						jumlahData.Tahun,
						jumlahData.JumlahData,
						jumlahData.Satuan,
					)
					if err != nil {
						return domain.IsuStrategis{}, err
					}
				}
			}
		}
	}

	return isuStrategis, nil
}

func (repository *IsuStrategisRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, isuStrategis domain.IsuStrategis) (domain.IsuStrategis, error) {
	log.Printf("[Repository] Start updating isu strategis ID: %d", isuStrategis.Id)

	// 1. Validasi isu strategis exists
	existingIsu, err := repository.FindById(ctx, tx, isuStrategis.Id)
	if err != nil {
		log.Printf("[Repository] Error finding isu strategis: %v", err)
		return domain.IsuStrategis{}, err
	}
	if existingIsu.Id == 0 {
		log.Printf("[Repository] Isu strategis not found")
		return domain.IsuStrategis{}, fmt.Errorf("isu strategis dengan ID %d tidak ditemukan", isuStrategis.Id)
	}

	// 1. Update isu strategis
	scriptIsu := `UPDATE tb_isu_strategis_opd 
                  SET kode_opd = ?, 
                      nama_opd = ?, 
                      kode_bidang_urusan = ?, 
                      nama_bidang_urusan = ?, 
                      tahun_awal = ?, 
                      tahun_akhir = ?, 
                      isu_strategis = ?,
                      updated_at = CURRENT_TIMESTAMP
                  WHERE id = ?`

	log.Printf("[Repository] Executing update isu strategis query")
	resultIsu, err := tx.ExecContext(ctx, scriptIsu,
		isuStrategis.KodeOpd,
		isuStrategis.NamaOpd,
		isuStrategis.KodeBidangUrusan,
		isuStrategis.NamaBidangUrusan,
		isuStrategis.TahunAwal,
		isuStrategis.TahunAkhir,
		isuStrategis.IsuStrategis,
		isuStrategis.Id,
	)
	if err != nil {
		log.Printf("[Repository] Error updating isu strategis: %v", err)
		return domain.IsuStrategis{}, err
	}

	affected, err := resultIsu.RowsAffected()
	if err != nil {
		log.Printf("[Repository] Error getting affected rows: %v", err)
		return domain.IsuStrategis{}, err
	}
	if affected == 0 {
		log.Printf("[Repository] No rows affected in isu strategis update")
		return domain.IsuStrategis{}, fmt.Errorf("isu strategis dengan ID %d tidak dapat diupdate", isuStrategis.Id)
	}
	log.Printf("[Repository] Updated %d rows in isu strategis", affected)

	// 2. Update permasalahan
	for _, permasalahan := range isuStrategis.PermasalahanOpd {
		log.Printf("[Repository] Processing permasalahan ID: %d", permasalahan.Id)

		// Validasi permasalahan exists
		scriptValidate := "SELECT id FROM tb_permasalahan_opd WHERE id = ?"
		var permId int
		err := tx.QueryRowContext(ctx, scriptValidate, permasalahan.Id).Scan(&permId)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("[Repository] Permasalahan ID %d not found", permasalahan.Id)
				continue // Skip if not found
			}
			log.Printf("[Repository] Error validating permasalahan: %v", err)
			return domain.IsuStrategis{}, err
		}

		// Update isu_strategis_id di permasalahan
		scriptPermasalahan := `UPDATE tb_permasalahan_opd 
                              SET isu_strategis_id = ? 
                              WHERE id = ?`

		resultPerm, err := tx.ExecContext(ctx, scriptPermasalahan, isuStrategis.Id, permasalahan.Id)
		if err != nil {
			log.Printf("[Repository] Error updating permasalahan: %v", err)
			return domain.IsuStrategis{}, err
		}

		affectedPerm, _ := resultPerm.RowsAffected()
		log.Printf("[Repository] Updated %d rows in permasalahan", affectedPerm)

		// 3. Handle data dukung
		for _, dataDukung := range permasalahan.DataDukung {
			log.Printf("[Repository] Processing data dukung for permasalahan ID: %d", permasalahan.Id)

			var dataDukungId int64
			if dataDukung.Id != 0 {
				// Validasi data dukung exists
				scriptValidateDD := "SELECT id FROM tb_data_dukung WHERE id = ? AND id_permasalahan = ?"
				var ddId int
				err := tx.QueryRowContext(ctx, scriptValidateDD, dataDukung.Id, permasalahan.Id).Scan(&ddId)
				if err != nil {
					if err == sql.ErrNoRows {
						log.Printf("[Repository] Data dukung ID %d not found for permasalahan ID %d", dataDukung.Id, permasalahan.Id)
						continue // Skip if not found
					}
					log.Printf("[Repository] Error validating data dukung: %v", err)
					return domain.IsuStrategis{}, err
				}

				// Update existing data dukung
				scriptDataDukung := `UPDATE tb_data_dukung 
                                   SET nama_data_dukung = ?, 
                                       narasi_data_dukung = ?,
                                       updated_at = CURRENT_TIMESTAMP
                                   WHERE id = ? AND id_permasalahan = ?`

				log.Printf("[Repository] Updating existing data dukung ID: %d", dataDukung.Id)
				resultDD, err := tx.ExecContext(ctx, scriptDataDukung,
					dataDukung.DataDukung,
					dataDukung.NarasiDataDukung,
					dataDukung.Id,
					permasalahan.Id,
				)
				if err != nil {
					log.Printf("[Repository] Error updating data dukung: %v", err)
					return domain.IsuStrategis{}, err
				}
				dataDukungId = int64(dataDukung.Id)

				affected, _ := resultDD.RowsAffected()
				log.Printf("[Repository] Updated %d rows in data dukung", affected)
			} else {
				// Insert new data dukung
				scriptNewDD := `INSERT INTO tb_data_dukung 
                               (id_permasalahan, nama_data_dukung, narasi_data_dukung, created_at, updated_at) 
                               VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

				log.Printf("[Repository] Inserting new data dukung")
				resultNewDD, err := tx.ExecContext(ctx, scriptNewDD,
					permasalahan.Id,
					dataDukung.DataDukung,
					dataDukung.NarasiDataDukung,
				)
				if err != nil {
					log.Printf("[Repository] Error inserting new data dukung: %v", err)
					return domain.IsuStrategis{}, err
				}

				dataDukungId, err = resultNewDD.LastInsertId()
				if err != nil {
					log.Printf("[Repository] Error getting last insert ID for data dukung: %v", err)
					return domain.IsuStrategis{}, err
				}
				log.Printf("[Repository] Inserted new data dukung with ID: %d", dataDukungId)
			}

			// 4. Handle jumlah data
			for _, jumlahData := range dataDukung.JumlahData {
				log.Printf("[Repository] Processing jumlah data for data dukung ID: %d", dataDukungId)

				if jumlahData.Id != 0 {
					// Validasi jumlah data exists
					scriptValidateJD := "SELECT id FROM tb_jumlah_data WHERE id = ? AND id_data_dukung = ?"
					var jdId int
					err := tx.QueryRowContext(ctx, scriptValidateJD, jumlahData.Id, dataDukungId).Scan(&jdId)
					if err != nil {
						if err == sql.ErrNoRows {
							log.Printf("[Repository] Jumlah data ID %d not found for data dukung ID %d", jumlahData.Id, dataDukungId)
							continue // Skip if not found
						}
						log.Printf("[Repository] Error validating jumlah data: %v", err)
						return domain.IsuStrategis{}, err
					}

					// Update existing jumlah data
					scriptJD := `UPDATE tb_jumlah_data 
                                SET tahun = ?, 
                                    jumlah = ?, 
                                    satuan = ?,
                                    updated_at = CURRENT_TIMESTAMP
                                WHERE id = ? AND id_data_dukung = ?`

					log.Printf("[Repository] Updating existing jumlah data ID: %d", jumlahData.Id)
					resultJD, err := tx.ExecContext(ctx, scriptJD,
						jumlahData.Tahun,
						jumlahData.JumlahData,
						jumlahData.Satuan,
						jumlahData.Id,
						dataDukungId,
					)
					if err != nil {
						log.Printf("[Repository] Error updating jumlah data: %v", err)
						return domain.IsuStrategis{}, err
					}

					affected, _ := resultJD.RowsAffected()
					log.Printf("[Repository] Updated %d rows in jumlah data", affected)
				} else {
					// Insert new jumlah data
					scriptNewJD := `INSERT INTO tb_jumlah_data 
                                  (id_data_dukung, tahun, jumlah, satuan, created_at, updated_at) 
                                  VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

					log.Printf("[Repository] Inserting new jumlah data")
					resultNewJD, err := tx.ExecContext(ctx, scriptNewJD,
						dataDukungId,
						jumlahData.Tahun,
						jumlahData.JumlahData,
						jumlahData.Satuan,
					)
					if err != nil {
						log.Printf("[Repository] Error inserting new jumlah data: %v", err)
						return domain.IsuStrategis{}, err
					}

					newJDId, err := resultNewJD.LastInsertId()
					if err != nil {
						log.Printf("[Repository] Error getting last insert ID for jumlah data: %v", err)
						return domain.IsuStrategis{}, err
					}
					log.Printf("[Repository] Inserted new jumlah data with ID: %d", newJDId)
				}
			}
		}
	}

	// Fetch updated data
	updatedIsuStrategis, err := repository.FindById(ctx, tx, isuStrategis.Id)
	if err != nil {
		log.Printf("[Repository] Error fetching updated isu strategis: %v", err)
		return domain.IsuStrategis{}, err
	}

	log.Printf("[Repository] Successfully completed update for isu strategis ID: %d", isuStrategis.Id)
	return updatedIsuStrategis, nil
}

func (repository *IsuStrategisRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_isu_strategis_opd WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	return err
}

func (repository *IsuStrategisRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, isuStrategisId int) (domain.IsuStrategis, error) {
	query := `
        SELECT 
            iso.id,
            iso.kode_opd,
            iso.nama_opd,
            iso.kode_bidang_urusan,
            iso.nama_bidang_urusan,
            iso.tahun_awal,
            iso.tahun_akhir,
            iso.isu_strategis,
            p.id as permasalahan_id,
            p.permasalahan,
            p.kode_opd as p_kode_opd,
            p.tahun as p_tahun,
            p.level_pohon,
            p.jenis_masalah,
            dd.id as data_dukung_id,
            dd.nama_data_dukung,
            dd.narasi_data_dukung,
            jd.id as jumlah_data_id,
            jd.tahun as jumlah_data_tahun,
            COALESCE(jd.jumlah, 0.0) as jumlah,
            COALESCE(jd.satuan, '') as satuan
        FROM 
            tb_isu_strategis_opd iso
        LEFT JOIN 
            tb_permasalahan_opd p ON p.isu_strategis_id = iso.id
        LEFT JOIN 
            tb_data_dukung dd ON dd.id_permasalahan = p.id
        LEFT JOIN 
            tb_jumlah_data jd ON jd.id_data_dukung = dd.id
        WHERE 
            iso.id = ?
        ORDER BY 
            iso.id, p.id, dd.id, jd.tahun DESC`

	rows, err := tx.QueryContext(ctx, query, isuStrategisId)
	if err != nil {
		return domain.IsuStrategis{}, err
	}
	defer rows.Close()

	var isuStrategis *domain.IsuStrategis
	permasalahanMap := make(map[int]*domain.Permasalahan)
	dataDukungMap := make(map[int]*domain.DataDukung)

	for rows.Next() {
		var (
			id               int
			kodeOpd          string
			namaOpd          string
			kodeBidangUrusan string
			namaBidangUrusan string
			tahunAwal        string
			tahunAkhir       string
			isuStrategisText string
			permasalahanId   sql.NullInt64
			permasalahan     sql.NullString
			pKodeOpd         sql.NullString
			pTahun           sql.NullString
			levelPohon       sql.NullInt64
			jenisMasalah     sql.NullString
			dataDukungId     sql.NullInt64
			namaDataDukung   sql.NullString
			narasiDataDukung sql.NullString
			jumlahDataId     sql.NullInt64
			jumlahDataTahun  sql.NullString
			jumlah           sql.NullFloat64
			satuan           sql.NullString
		)

		err := rows.Scan(
			&id, &kodeOpd, &namaOpd, &kodeBidangUrusan, &namaBidangUrusan,
			&tahunAwal, &tahunAkhir, &isuStrategisText, &permasalahanId, &permasalahan,
			&pKodeOpd, &pTahun, &levelPohon, &jenisMasalah, &dataDukungId,
			&namaDataDukung, &narasiDataDukung, &jumlahDataId, &jumlahDataTahun,
			&jumlah, &satuan,
		)
		if err != nil {
			return domain.IsuStrategis{}, err
		}

		// Inisialisasi IsuStrategis jika belum ada
		if isuStrategis == nil {
			isuStrategis = &domain.IsuStrategis{
				Id:               id,
				KodeOpd:          kodeOpd,
				NamaOpd:          namaOpd,
				KodeBidangUrusan: kodeBidangUrusan,
				NamaBidangUrusan: namaBidangUrusan,
				TahunAwal:        tahunAwal,
				TahunAkhir:       tahunAkhir,
				IsuStrategis:     isuStrategisText,
				PermasalahanOpd:  make([]domain.Permasalahan, 0),
			}
		}

		// Handle Permasalahan
		if permasalahanId.Valid {
			permId := int(permasalahanId.Int64)
			perm, exists := permasalahanMap[permId]
			if !exists {
				perm = &domain.Permasalahan{
					Id:           permId,
					Permasalahan: permasalahan.String,
					KodeOpd:      pKodeOpd.String,
					Tahun:        pTahun.String,
					LevelPohon:   int(levelPohon.Int64),
					JenisMasalah: jenisMasalah.String,
					DataDukung:   make([]domain.DataDukung, 0),
				}
				permasalahanMap[permId] = perm
				isuStrategis.PermasalahanOpd = append(isuStrategis.PermasalahanOpd, *perm)
			}

			// Handle DataDukung
			if dataDukungId.Valid {
				ddId := int(dataDukungId.Int64)
				dd, exists := dataDukungMap[ddId]
				if !exists {
					dd = &domain.DataDukung{
						Id:                ddId,
						PermasalahanOpdId: permId,
						DataDukung:        namaDataDukung.String,
						NarasiDataDukung:  narasiDataDukung.String,
						JumlahData:        make([]domain.JumlahData, 0),
					}
					dataDukungMap[ddId] = dd

					// Tambahkan ke permasalahan yang benar
					for i := range isuStrategis.PermasalahanOpd {
						if isuStrategis.PermasalahanOpd[i].Id == permId {
							isuStrategis.PermasalahanOpd[i].DataDukung = append(isuStrategis.PermasalahanOpd[i].DataDukung, *dd)
							break
						}
					}
				}

				// Handle JumlahData - hanya untuk tahun dalam rentang
				if jumlahDataId.Valid && jumlahDataTahun.Valid {
					tahunJumlahData := jumlahDataTahun.String
					// Cek apakah tahun dalam rentang
					if tahunJumlahData >= tahunAwal && tahunJumlahData <= tahunAkhir {
						jd := domain.JumlahData{
							Id:           int(jumlahDataId.Int64),
							IdDataDukung: ddId,
							Tahun:        tahunJumlahData,
							JumlahData:   jumlah.Float64,
							Satuan:       satuan.String,
						}

						// Tambahkan ke data dukung yang benar
						for i := range isuStrategis.PermasalahanOpd {
							for j := range isuStrategis.PermasalahanOpd[i].DataDukung {
								if isuStrategis.PermasalahanOpd[i].DataDukung[j].Id == ddId {
									isuStrategis.PermasalahanOpd[i].DataDukung[j].JumlahData = append(
										isuStrategis.PermasalahanOpd[i].DataDukung[j].JumlahData, jd)
									break
								}
							}
						}
					}
				}
			}
		}
	}

	if isuStrategis == nil {
		return domain.IsuStrategis{}, nil
	}

	return *isuStrategis, nil
}

func (repository *IsuStrategisRepositoryImpl) getJumlahDataByDataDukungId(ctx context.Context, tx *sql.Tx, dataDukungId int) ([]domain.JumlahData, error) {
	script := `SELECT id, id_data_dukung, tahun, jumlah, satuan 
               FROM tb_jumlah_data 
               WHERE id_data_dukung = ?
               ORDER BY tahun DESC`

	rows, err := tx.QueryContext(ctx, script, dataDukungId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jumlahDatas []domain.JumlahData
	for rows.Next() {
		var jumlahData domain.JumlahData
		err := rows.Scan(
			&jumlahData.Id,
			&jumlahData.IdDataDukung,
			&jumlahData.Tahun,
			&jumlahData.JumlahData,
			&jumlahData.Satuan,
		)
		if err != nil {
			return nil, err
		}

		jumlahDatas = append(jumlahDatas, jumlahData)
	}

	if len(jumlahDatas) == 0 {
		return []domain.JumlahData{}, nil
	}

	return jumlahDatas, nil
}

// Fungsi untuk mencari data dukung berdasarkan ID
func (repository *IsuStrategisRepositoryImpl) FindDataDukungById(ctx context.Context, tx *sql.Tx, dataDukungId int) (domain.DataDukung, error) {
	script := `SELECT id, id_permasalahan, nama_data_dukung, narasi_data_dukung 
               FROM tb_data_dukung 
               WHERE id = ?`

	rows, err := tx.QueryContext(ctx, script, dataDukungId)
	if err != nil {
		return domain.DataDukung{}, err
	}
	defer rows.Close()

	var dataDukung domain.DataDukung
	if rows.Next() {
		err := rows.Scan(
			&dataDukung.Id,
			&dataDukung.PermasalahanOpdId,
			&dataDukung.DataDukung,
			&dataDukung.NarasiDataDukung,
		)
		if err != nil {
			return domain.DataDukung{}, err
		}

		// Get Jumlah Data
		jumlahDatas, err := repository.getJumlahDataByDataDukungId(ctx, tx, dataDukung.Id)
		if err != nil {
			return domain.DataDukung{}, err
		}
		dataDukung.JumlahData = jumlahDatas
	}

	return dataDukung, nil
}

// Fungsi untuk mencari jumlah data berdasarkan ID
func (repository *IsuStrategisRepositoryImpl) FindJumlahDataById(ctx context.Context, tx *sql.Tx, jumlahDataId int) (domain.JumlahData, error) {
	script := `SELECT id, id_data_dukung, tahun, jumlah, satuan 
               FROM tb_jumlah_data 
               WHERE id = ?`

	rows, err := tx.QueryContext(ctx, script, jumlahDataId)
	if err != nil {
		return domain.JumlahData{}, err
	}
	defer rows.Close()

	var jumlahData domain.JumlahData
	if rows.Next() {
		err := rows.Scan(
			&jumlahData.Id,
			&jumlahData.IdDataDukung,
			&jumlahData.Tahun,
			&jumlahData.JumlahData,
			&jumlahData.Satuan,
		)
		if err != nil {
			return domain.JumlahData{}, err
		}
	}

	return jumlahData, nil
}

func (repository *IsuStrategisRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string, tahunAwal string, tahunAkhir string) ([]domain.IsuStrategis, error) {
	query := `
    SELECT 
        iso.id,
        iso.kode_opd,
        iso.nama_opd,
        iso.kode_bidang_urusan,
        iso.nama_bidang_urusan,
        iso.tahun_awal,
        iso.tahun_akhir,
        iso.isu_strategis,
        iso.created_at,      -- Tambahkan created_at
        p.id as permasalahan_id,
        p.permasalahan,
        p.kode_opd as p_kode_opd,
        p.tahun as p_tahun,
        p.level_pohon,
        p.jenis_masalah,
        dd.id as data_dukung_id,
        dd.nama_data_dukung,
        dd.narasi_data_dukung,
        jd.id as jumlah_data_id,
        jd.tahun as jumlah_data_tahun,
        COALESCE(jd.jumlah, 0.0) as jumlah,
        COALESCE(jd.satuan, '') as satuan
    FROM 
        tb_isu_strategis_opd iso
    LEFT JOIN 
        tb_permasalahan_opd p ON p.isu_strategis_id = iso.id
    LEFT JOIN 
        tb_data_dukung dd ON dd.id_permasalahan = p.id
    LEFT JOIN 
        tb_jumlah_data jd ON jd.id_data_dukung = dd.id
    WHERE 
        iso.kode_opd = ?
        AND (? = '' OR iso.tahun_awal >= ?)
        AND (? = '' OR iso.tahun_akhir <= ?)
    ORDER BY 
        iso.created_at ASC,   -- Tambahkan ordering
        p.id, dd.id, jd.tahun DESC`

	rows, err := tx.QueryContext(ctx, query, kodeOpd, tahunAwal, tahunAwal, tahunAkhir, tahunAkhir)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	isuStrategisMap := make(map[int]*domain.IsuStrategis)
	permasalahanMap := make(map[int]*domain.Permasalahan)
	dataDukungMap := make(map[int]*domain.DataDukung)
	jumlahDataMap := make(map[int]map[string]domain.JumlahData)

	// Fungsi helper untuk generate tahun-tahun dalam rentang
	generateYearRange := func(start, end string) []string {
		startYear, _ := strconv.Atoi(start)
		endYear, _ := strconv.Atoi(end)
		years := make([]string, 0)
		for year := endYear; year >= startYear; year-- {
			years = append(years, strconv.Itoa(year))
		}
		return years
	}

	for rows.Next() {
		var (
			isuStrategisId   int
			kodeOpd          string
			namaOpd          string
			kodeBidangUrusan string
			namaBidangUrusan string
			tahunAwal        string
			tahunAkhir       string
			isuStrategis     string
			createdAt        time.Time // Tambahkan created_at
			permasalahanId   sql.NullInt64
			permasalahan     sql.NullString
			pKodeOpd         sql.NullString
			pTahun           sql.NullString
			levelPohon       sql.NullInt64
			jenisMasalah     sql.NullString
			dataDukungId     sql.NullInt64
			namaDataDukung   sql.NullString
			narasiDataDukung sql.NullString
			jumlahDataId     sql.NullInt64
			jumlahDataTahun  sql.NullString
			jumlah           sql.NullFloat64
			satuan           sql.NullString
		)

		err := rows.Scan(
			&isuStrategisId, &kodeOpd, &namaOpd, &kodeBidangUrusan, &namaBidangUrusan,
			&tahunAwal, &tahunAkhir, &isuStrategis, &createdAt, &permasalahanId, &permasalahan,
			&pKodeOpd, &pTahun, &levelPohon, &jenisMasalah, &dataDukungId,
			&namaDataDukung, &narasiDataDukung, &jumlahDataId, &jumlahDataTahun,
			&jumlah, &satuan,
		)
		if err != nil {
			return nil, err
		}

		// Get or create IsuStrategis
		isuStr, exists := isuStrategisMap[isuStrategisId]
		if !exists {
			isuStr = &domain.IsuStrategis{
				Id:               isuStrategisId,
				KodeOpd:          kodeOpd,
				NamaOpd:          namaOpd,
				KodeBidangUrusan: kodeBidangUrusan,
				NamaBidangUrusan: namaBidangUrusan,
				TahunAwal:        tahunAwal,
				TahunAkhir:       tahunAkhir,
				IsuStrategis:     isuStrategis,
				CreatedAt:        createdAt, // Tambahkan created_at
				PermasalahanOpd:  make([]domain.Permasalahan, 0),
			}
			isuStrategisMap[isuStrategisId] = isuStr
		}

		// Handle Permasalahan
		if permasalahanId.Valid {
			permId := int(permasalahanId.Int64)
			perm, exists := permasalahanMap[permId]
			if !exists {
				perm = &domain.Permasalahan{
					Id:           permId,
					Permasalahan: permasalahan.String,
					KodeOpd:      pKodeOpd.String,
					Tahun:        pTahun.String,
					LevelPohon:   int(levelPohon.Int64),
					JenisMasalah: jenisMasalah.String,
					DataDukung:   make([]domain.DataDukung, 0),
				}
				permasalahanMap[permId] = perm
				isuStr.PermasalahanOpd = append(isuStr.PermasalahanOpd, *perm)
			}

			// Handle DataDukung
			if dataDukungId.Valid {
				ddId := int(dataDukungId.Int64)

				// Inisialisasi map jumlah data jika belum ada
				if _, exists := jumlahDataMap[ddId]; !exists {
					jumlahDataMap[ddId] = make(map[string]domain.JumlahData)
				}

				// Simpan data jumlah data ke map
				if jumlahDataId.Valid && jumlahDataTahun.Valid {
					jumlahDataMap[ddId][jumlahDataTahun.String] = domain.JumlahData{
						Id:           int(jumlahDataId.Int64),
						IdDataDukung: ddId,
						Tahun:        jumlahDataTahun.String,
						JumlahData:   jumlah.Float64,
						Satuan:       satuan.String,
					}
				}

				dd, exists := dataDukungMap[ddId]
				if !exists {
					dd = &domain.DataDukung{
						Id:                ddId,
						PermasalahanOpdId: permId,
						DataDukung:        namaDataDukung.String,
						NarasiDataDukung:  narasiDataDukung.String,
						JumlahData:        make([]domain.JumlahData, 0),
					}
					dataDukungMap[ddId] = dd

					// Tambahkan ke permasalahan yang benar
					for i := range isuStr.PermasalahanOpd {
						if isuStr.PermasalahanOpd[i].Id == permId {
							isuStr.PermasalahanOpd[i].DataDukung = append(isuStr.PermasalahanOpd[i].DataDukung, *dd)
							break
						}
					}
				}
			}
		}
	}

	// Setelah semua data terkumpul, generate rentang tahun dan isi data
	for _, isuStr := range isuStrategisMap {
		for i, perm := range isuStr.PermasalahanOpd {
			for j, dd := range perm.DataDukung {
				yearRange := generateYearRange(isuStr.TahunAwal, isuStr.TahunAkhir)
				jumlahDataSlice := make([]domain.JumlahData, 0)

				for _, year := range yearRange {
					if data, exists := jumlahDataMap[dd.Id][year]; exists {
						jumlahDataSlice = append(jumlahDataSlice, data)
					} else {
						jumlahDataSlice = append(jumlahDataSlice, domain.JumlahData{
							IdDataDukung: dd.Id,
							Tahun:        year,
							JumlahData:   0,
							Satuan:       "",
						})
					}
				}

				isuStr.PermasalahanOpd[i].DataDukung[j].JumlahData = jumlahDataSlice
			}
		}
	}

	// Convert map to slice
	var result []domain.IsuStrategis
	for _, isuStr := range isuStrategisMap {
		result = append(result, *isuStr)
	}

	return result, nil
}

func (repository *IsuStrategisRepositoryImpl) FindDataDukungByPermasalahanId(ctx context.Context, tx *sql.Tx, permasalahanId int) ([]domain.DataDukung, error) {
	script := `SELECT id, id_permasalahan, nama_data_dukung, narasi_data_dukung 
               FROM tb_data_dukung 
               WHERE id_permasalahan = ?`
	rows, err := tx.QueryContext(ctx, script, permasalahanId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.DataDukung
	for rows.Next() {
		var dataDukung domain.DataDukung
		err := rows.Scan(&dataDukung.Id, &dataDukung.PermasalahanOpdId, &dataDukung.DataDukung, &dataDukung.NarasiDataDukung)
		if err != nil {
			return nil, err
		}
		result = append(result, dataDukung)
	}
	return result, nil
}

func (repository *IsuStrategisRepositoryImpl) DeleteDataDukung(ctx context.Context, tx *sql.Tx, idPermasalahan int) error {
	log.Printf("[Repository] Deleting all data dukung for permasalahan ID: %d", idPermasalahan)

	script := "DELETE FROM tb_data_dukung WHERE id_permasalahan = ?"
	result, err := tx.ExecContext(ctx, script, idPermasalahan)
	if err != nil {
		log.Printf("[Repository] Error deleting data dukung: %v", err)
		return err
	}

	affected, _ := result.RowsAffected()
	log.Printf("[Repository] Deleted %d data dukung records", affected)
	return nil
}

func (repository *IsuStrategisRepositoryImpl) DeleteJumlahData(ctx context.Context, tx *sql.Tx, id int) error {
	log.Printf("[Repository] Deleting all jumlah data for data dukung ID: %d", id)

	script := "DELETE FROM tb_jumlah_data WHERE id_data_dukung = ?"
	result, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		log.Printf("[Repository] Error deleting jumlah data: %v", err)
		return err
	}

	affected, _ := result.RowsAffected()
	log.Printf("[Repository] Deleted %d jumlah data records", affected)
	return nil
}

func (repository *IsuStrategisRepositoryImpl) FindJumlahDataByDataDukungId(ctx context.Context, tx *sql.Tx, dataDukungId int) ([]domain.JumlahData, error) {
	log.Printf("[Repository] Finding jumlah data for data dukung ID: %d", dataDukungId)

	script := `SELECT 
                id,
                id_data_dukung,
                tahun,
                COALESCE(jumlah, 0) as jumlah,
                COALESCE(satuan, '') as satuan
               FROM tb_jumlah_data 
               WHERE id_data_dukung = ?`

	rows, err := tx.QueryContext(ctx, script, dataDukungId)
	if err != nil {
		log.Printf("[Repository] Error querying jumlah data: %v", err)
		return nil, err
	}
	defer rows.Close()

	var result []domain.JumlahData
	for rows.Next() {
		var jumlahData domain.JumlahData
		err := rows.Scan(
			&jumlahData.Id,
			&jumlahData.IdDataDukung,
			&jumlahData.Tahun,
			&jumlahData.JumlahData,
			&jumlahData.Satuan,
		)
		if err != nil {
			log.Printf("[Repository] Error scanning jumlah data: %v", err)
			return nil, err
		}
		result = append(result, jumlahData)
	}

	log.Printf("[Repository] Found %d jumlah data records", len(result))
	return result, nil
}

func (repository *IsuStrategisRepositoryImpl) DeleteJumlahDataByDataDukungId(ctx context.Context, tx *sql.Tx, dataDukungId int) error {
	log.Printf("[Repository] Deleting all jumlah data for data dukung ID: %d", dataDukungId)

	script := "DELETE FROM tb_jumlah_data WHERE id_data_dukung = ?"
	result, err := tx.ExecContext(ctx, script, dataDukungId)
	if err != nil {
		log.Printf("[Repository] Error deleting jumlah data: %v", err)
		return err
	}

	affected, _ := result.RowsAffected()
	log.Printf("[Repository] Deleted %d jumlah data records", affected)
	return nil
}

func (repository *IsuStrategisRepositoryImpl) DeleteDataDukungByPermasalahanId(ctx context.Context, tx *sql.Tx, permasalahanId int) error {
	log.Printf("[Repository] Deleting all data dukung for permasalahan ID: %d", permasalahanId)

	script := "DELETE FROM tb_data_dukung WHERE id_permasalahan = ?"
	result, err := tx.ExecContext(ctx, script, permasalahanId)
	if err != nil {
		log.Printf("[Repository] Error deleting data dukung: %v", err)
		return err
	}

	affected, _ := result.RowsAffected()
	log.Printf("[Repository] Deleted %d data dukung records", affected)
	return nil
}
