package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"permasalahanService/model/domain"
	"sort"
	"strconv"
	"time"
)

type IsuStrategisRepositoryImpl struct {
}

func NewIsuStrategisRepositoryImpl() *IsuStrategisRepositoryImpl {
	return &IsuStrategisRepositoryImpl{}
}

func (repository *IsuStrategisRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, isuStrategis domain.IsuStrategis) (domain.IsuStrategis, error) {
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

	// Insert ke junction table dan update status permasalahan
	for _, permasalahan := range isuStrategis.PermasalahanOpd {
		// 1. Insert ke junction table
		scriptJunction := `INSERT INTO tb_permasalahan_isu_strategis 
		                   (id_permasalahan, id_isu_strategis) 
		                   VALUES (?, ?)`
		_, err = tx.ExecContext(ctx, scriptJunction, permasalahan.Id, isuStrategis.Id)
		if err != nil {
			return domain.IsuStrategis{}, err
		}

		// 2. Ambil permasalahan_opd_id dari tb_permasalahan_terpilih
		var permasalahanOpdId int
		scriptGetOpdId := `SELECT permasalahan_opd_id FROM tb_permasalahan_terpilih WHERE id = ?`
		err = tx.QueryRowContext(ctx, scriptGetOpdId, permasalahan.Id).Scan(&permasalahanOpdId)
		if err != nil {
			return domain.IsuStrategis{}, fmt.Errorf("gagal mendapatkan permasalahan_opd_id: %v", err)
		}

		// 3. Update status_permasalahan di tb_permasalahan_opd
		scriptUpdateStatus := `UPDATE tb_permasalahan_opd 
		                       SET status_permasalahan = 'digunakan' 
		                       WHERE id = ?`
		_, err = tx.ExecContext(ctx, scriptUpdateStatus, permasalahanOpdId)
		if err != nil {
			return domain.IsuStrategis{}, fmt.Errorf("gagal update status permasalahan: %v", err)
		}

		// 4. Insert data dukung dengan id_permasalahan DAN id_isu_strategis
		for _, dataDukung := range permasalahan.DataDukung {
			scriptDD := `INSERT INTO tb_data_dukung 
			             (id_permasalahan, id_isustrategis, nama_data_dukung, narasi_data_dukung) 
			             VALUES (?, ?, ?, ?)`

			resultDD, err := tx.ExecContext(ctx, scriptDD,
				permasalahanOpdId,
				isuStrategis.Id,
				dataDukung.DataDukung,
				dataDukung.NarasiDataDukung,
			)
			if err != nil {
				return domain.IsuStrategis{}, err
			}

			dataDukungId, err := resultDD.LastInsertId()
			if err != nil {
				return domain.IsuStrategis{}, err
			}

			// 5. Insert jumlah data
			for _, jumlahData := range dataDukung.JumlahData {
				if jumlahData.Tahun != "" {
					scriptJD := `INSERT INTO tb_jumlah_data 
					             (id_data_dukung, tahun, jumlah, satuan) 
					             VALUES (?, ?, ?, ?)`

					_, err = tx.ExecContext(ctx, scriptJD,
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

	// 2. Update isu strategis basic info
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

	// 3. Dapatkan permasalahan terpilih yang sudah ada di junction table
	scriptGetExisting := `SELECT pt.id, pt.permasalahan_opd_id 
	                      FROM tb_permasalahan_isu_strategis pis
	                      INNER JOIN tb_permasalahan_terpilih pt ON pis.id_permasalahan = pt.id
	                      WHERE pis.id_isu_strategis = ?`

	rows, err := tx.QueryContext(ctx, scriptGetExisting, isuStrategis.Id)
	if err != nil {
		log.Printf("[Repository] Error getting existing permasalahan: %v", err)
		return domain.IsuStrategis{}, err
	}

	type ExistingPermasalahan struct {
		PermasalahanTerpilihId int
		PermasalahanOpdId      int
	}

	existingPermasalahan := make([]ExistingPermasalahan, 0)
	for rows.Next() {
		var ep ExistingPermasalahan
		err := rows.Scan(&ep.PermasalahanTerpilihId, &ep.PermasalahanOpdId)
		if err != nil {
			rows.Close()
			return domain.IsuStrategis{}, err
		}
		existingPermasalahan = append(existingPermasalahan, ep)
	}
	rows.Close()

	// 4. Buat map untuk tracking permasalahan yang akan dipertahankan
	keepPermasalahanIds := make(map[int]bool)
	for _, p := range isuStrategis.PermasalahanOpd {
		keepPermasalahanIds[p.Id] = true
	}

	// 5. Hapus relasi dan reset status untuk permasalahan yang tidak ada di request
	for _, ep := range existingPermasalahan {
		if !keepPermasalahanIds[ep.PermasalahanTerpilihId] {
			log.Printf("[Repository] Removing permasalahan terpilih ID %d from junction table", ep.PermasalahanTerpilihId)

			// Hapus dari junction table
			scriptDeleteJunction := `DELETE FROM tb_permasalahan_isu_strategis 
			                         WHERE id_permasalahan = ? AND id_isu_strategis = ?`
			_, err := tx.ExecContext(ctx, scriptDeleteJunction, ep.PermasalahanTerpilihId, isuStrategis.Id)
			if err != nil {
				log.Printf("[Repository] Error deleting from junction table: %v", err)
				return domain.IsuStrategis{}, err
			}

			// Reset status_permasalahan di tb_permasalahan_opd
			scriptResetStatus := `UPDATE tb_permasalahan_opd 
			                      SET status_permasalahan = '' 
			                      WHERE id = ?`
			_, err = tx.ExecContext(ctx, scriptResetStatus, ep.PermasalahanOpdId)
			if err != nil {
				log.Printf("[Repository] Error resetting status: %v", err)
				return domain.IsuStrategis{}, err
			}

			// Hapus data dukung dan jumlah data
			scriptDeleteDD := `DELETE FROM tb_data_dukung 
			WHERE id_permasalahan = ? 
			AND id_isustrategis = ?`
			_, err = tx.ExecContext(ctx, scriptDeleteDD, ep.PermasalahanOpdId, isuStrategis.Id)
			if err != nil {
				log.Printf("[Repository] Error deleting data dukung: %v", err)
				return domain.IsuStrategis{}, err
			}
		}
	}

	// 6. Proses permasalahan baru atau yang sudah ada
	for _, permasalahan := range isuStrategis.PermasalahanOpd {
		log.Printf("[Repository] Processing permasalahan terpilih ID: %d", permasalahan.Id)

		// Cek apakah sudah ada di junction table
		var existsInJunction int
		scriptCheckJunction := `SELECT COUNT(*) FROM tb_permasalahan_isu_strategis 
		                        WHERE id_permasalahan = ? AND id_isu_strategis = ?`
		err := tx.QueryRowContext(ctx, scriptCheckJunction, permasalahan.Id, isuStrategis.Id).Scan(&existsInJunction)
		if err != nil {
			return domain.IsuStrategis{}, err
		}

		// Jika belum ada, insert ke junction table
		if existsInJunction == 0 {
			// Insert ke junction table
			scriptInsertJunction := `INSERT INTO tb_permasalahan_isu_strategis 
			                         (id_permasalahan, id_isu_strategis) 
			                         VALUES (?, ?)`
			_, err = tx.ExecContext(ctx, scriptInsertJunction, permasalahan.Id, isuStrategis.Id)
			if err != nil {
				log.Printf("[Repository] Error inserting to junction table: %v", err)
				return domain.IsuStrategis{}, err
			}
		}

		// Ambil permasalahan_opd_id dari tb_permasalahan_terpilih
		var permasalahanOpdId int
		scriptGetOpdId := `SELECT permasalahan_opd_id FROM tb_permasalahan_terpilih WHERE id = ?`
		err = tx.QueryRowContext(ctx, scriptGetOpdId, permasalahan.Id).Scan(&permasalahanOpdId)
		if err != nil {
			return domain.IsuStrategis{}, fmt.Errorf("gagal mendapatkan permasalahan_opd_id: %v", err)
		}

		// Update status_permasalahan
		scriptUpdateStatus := `UPDATE tb_permasalahan_opd 
		                       SET status_permasalahan = 'digunakan' 
		                       WHERE id = ?`
		_, err = tx.ExecContext(ctx, scriptUpdateStatus, permasalahanOpdId)
		if err != nil {
			return domain.IsuStrategis{}, fmt.Errorf("gagal update status permasalahan: %v", err)
		}

		// 7. Handle data dukung (sama seperti sebelumnya, tapi gunakan permasalahanOpdId)
		for _, dataDukung := range permasalahan.DataDukung {
			log.Printf("[Repository] Processing data dukung for permasalahan OPD ID: %d", permasalahanOpdId)

			var dataDukungId int64
			if dataDukung.Id != 0 {
				// 🔥 PERBAIKAN: Validasi dengan id_permasalahan DAN id_isu_strategis
				scriptValidateDD := `SELECT id FROM tb_data_dukung 
									 WHERE id = ? 
									 AND id_permasalahan = ? 
									 AND id_isustrategis = ?`
				var ddId int
				err := tx.QueryRowContext(ctx, scriptValidateDD,
					dataDukung.Id,
					permasalahanOpdId,
					isuStrategis.Id).Scan(&ddId)
				if err != nil {
					if err == sql.ErrNoRows {
						log.Printf("[Repository] Data dukung ID %d not found", dataDukung.Id)
						continue
					}
					log.Printf("[Repository] Error validating data dukung: %v", err)
					return domain.IsuStrategis{}, err
				}

				// Update existing data dukung
				scriptDataDukung := `UPDATE tb_data_dukung 
									 SET nama_data_dukung = ?, 
										 narasi_data_dukung = ?,
										 updated_at = CURRENT_TIMESTAMP
									 WHERE id = ? 
									 AND id_permasalahan = ? 
									 AND id_isustrategis = ?`

				log.Printf("[Repository] Updating existing data dukung ID: %d", dataDukung.Id)
				resultDD, err := tx.ExecContext(ctx, scriptDataDukung,
					dataDukung.DataDukung,
					dataDukung.NarasiDataDukung,
					dataDukung.Id,
					permasalahanOpdId,
					isuStrategis.Id, // 🔥 Tambahan
				)
				if err != nil {
					log.Printf("[Repository] Error updating data dukung: %v", err)
					return domain.IsuStrategis{}, err
				}
				dataDukungId = int64(dataDukung.Id)

				affected, _ := resultDD.RowsAffected()
				log.Printf("[Repository] Updated %d rows in data dukung", affected)
			} else {
				// 🔥 PERBAIKAN: Insert dengan id_isu_strategis
				scriptNewDD := `INSERT INTO tb_data_dukung 
								(id_permasalahan, id_isustrategis, nama_data_dukung, narasi_data_dukung, created_at, updated_at) 
								VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

				log.Printf("[Repository] Inserting new data dukung")
				resultNewDD, err := tx.ExecContext(ctx, scriptNewDD,
					permasalahanOpdId,
					isuStrategis.Id, // 🔥 Tambahan
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

			// 8. Handle jumlah data (sama seperti sebelumnya)
			for _, jumlahData := range dataDukung.JumlahData {
				log.Printf("[Repository] Processing jumlah data for data dukung ID: %d", dataDukungId)

				if jumlahData.Id != 0 {
					// Update existing jumlah data
					scriptValidateJD := "SELECT id FROM tb_jumlah_data WHERE id = ? AND id_data_dukung = ?"
					var jdId int
					err := tx.QueryRowContext(ctx, scriptValidateJD, jumlahData.Id, dataDukungId).Scan(&jdId)
					if err != nil {
						if err == sql.ErrNoRows {
							log.Printf("[Repository] Jumlah data ID %d not found for data dukung ID %d", jumlahData.Id, dataDukungId)
							continue
						}
						log.Printf("[Repository] Error validating jumlah data: %v", err)
						return domain.IsuStrategis{}, err
					}

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
	// 1. Ambil semua permasalahan_opd_id yang terkait
	scriptGetOpdIds := `
		SELECT DISTINCT p.id 
		FROM tb_permasalahan_isu_strategis pis
		INNER JOIN tb_permasalahan_terpilih pt ON pis.id_permasalahan = pt.id
		INNER JOIN tb_permasalahan_opd p ON pt.permasalahan_opd_id = p.id
		WHERE pis.id_isu_strategis = ?`

	rows, err := tx.QueryContext(ctx, scriptGetOpdIds, id)
	if err != nil {
		return fmt.Errorf("gagal mendapatkan permasalahan opd: %v", err)
	}

	var opdIds []int
	for rows.Next() {
		var opdId int
		if err := rows.Scan(&opdId); err != nil {
			rows.Close()
			return fmt.Errorf("gagal scan opd id: %v", err)
		}
		opdIds = append(opdIds, opdId)
	}
	rows.Close()

	// 2. Reset status_permasalahan untuk semua permasalahan yang terkait
	for _, opdId := range opdIds {
		scriptResetStatus := "UPDATE tb_permasalahan_opd SET status_permasalahan = '' WHERE id = ?"
		_, err := tx.ExecContext(ctx, scriptResetStatus, opdId)
		if err != nil {
			return fmt.Errorf("gagal reset status permasalahan: %v", err)
		}
	}

	// 3. Hapus dari junction table (CASCADE akan menghapus isu strategis)
	scriptDeleteJunction := "DELETE FROM tb_permasalahan_isu_strategis WHERE id_isu_strategis = ?"
	result, err := tx.ExecContext(ctx, scriptDeleteJunction, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus dari junction table: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Berhasil menghapus %d relasi di junction table", rowsAffected)

	// 4. Hapus isu strategis
	scriptDelete := "DELETE FROM tb_isu_strategis_opd WHERE id = ?"
	result, err = tx.ExecContext(ctx, scriptDelete, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus isu strategis: %v", err)
	}

	rowsAffected, _ = result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("isu strategis dengan ID %d tidak ditemukan", id)
	}

	return nil
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
		pt.id as permasalahan_terpilih_id,
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
		tb_permasalahan_isu_strategis pis ON pis.id_isu_strategis = iso.id
	LEFT JOIN 
		tb_permasalahan_terpilih pt ON pis.id_permasalahan = pt.id
	LEFT JOIN 
		tb_permasalahan_opd p ON pt.permasalahan_opd_id = p.id
	LEFT JOIN 
		tb_data_dukung dd ON dd.id_permasalahan = p.id AND dd.id_isustrategis = iso.id
	LEFT JOIN 
		tb_jumlah_data jd ON jd.id_data_dukung = dd.id
	WHERE 
		iso.id = ?
	ORDER BY 
		iso.id, pt.id, p.id, dd.id, jd.tahun DESC`

	rows, err := tx.QueryContext(ctx, query, isuStrategisId)
	if err != nil {
		return domain.IsuStrategis{}, err
	}
	defer rows.Close()

	var isuStrategis *domain.IsuStrategis
	permasalahanMap := make(map[int]*domain.Permasalahan)
	dataDukungMap := make(map[int]*domain.DataDukung)
	jumlahDataMap := make(map[int]map[string]domain.JumlahData)

	for rows.Next() {
		var (
			id                     int
			kodeOpd                string
			namaOpd                string
			kodeBidangUrusan       string
			namaBidangUrusan       string
			tahunAwal              string
			tahunAkhir             string
			isuStrategisText       string
			permasalahanTerpilihId sql.NullInt64
			permasalahanId         sql.NullInt64
			permasalahan           sql.NullString
			pKodeOpd               sql.NullString
			pTahun                 sql.NullString
			levelPohon             sql.NullInt64
			jenisMasalah           sql.NullString
			dataDukungId           sql.NullInt64
			namaDataDukung         sql.NullString
			narasiDataDukung       sql.NullString
			jumlahDataId           sql.NullInt64
			jumlahDataTahun        sql.NullString
			jumlah                 sql.NullFloat64
			satuan                 sql.NullString
		)

		err := rows.Scan(
			&id, &kodeOpd, &namaOpd, &kodeBidangUrusan, &namaBidangUrusan,
			&tahunAwal, &tahunAkhir, &isuStrategisText, &permasalahanTerpilihId,
			&permasalahanId, &permasalahan,
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

		// Handle Permasalahan (gunakan permasalahan_id dari tb_permasalahan_opd)
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
						IdIsuStrategis:    id,
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
			}
		}
	}

	if isuStrategis == nil {
		return domain.IsuStrategis{}, nil
	}

	// Setelah semua data terkumpul, isi jumlah data
	for i, perm := range isuStrategis.PermasalahanOpd {
		for j, dd := range perm.DataDukung {
			jumlahDataSlice := make([]domain.JumlahData, 0)

			// Ambil semua jumlah data dari map
			if dataMap, exists := jumlahDataMap[dd.Id]; exists {
				for _, jd := range dataMap {
					jumlahDataSlice = append(jumlahDataSlice, jd)
				}

				// Sort berdasarkan tahun DESC
				sort.Slice(jumlahDataSlice, func(x, y int) bool {
					return jumlahDataSlice[x].Tahun > jumlahDataSlice[y].Tahun
				})
			}

			isuStrategis.PermasalahanOpd[i].DataDukung[j].JumlahData = jumlahDataSlice
		}
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
	script := `SELECT id, id_permasalahan, id_isustrategis, nama_data_dukung, narasi_data_dukung 
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
			&dataDukung.IdIsuStrategis,
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
        tb_data_dukung dd ON dd.id_permasalahan = p.id AND dd.id_isustrategis = iso.id
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
			createdAt        time.Time
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
						IdIsuStrategis:    isuStrategisId,
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

func (repository *IsuStrategisRepositoryImpl) FindDataDukungByPermasalahanIdAndIsuStrategisId(ctx context.Context, tx *sql.Tx, permasalahanId int, isuStrategisId int) ([]domain.DataDukung, error) {
	script := `SELECT id, id_permasalahan, id_isustrategis, nama_data_dukung, narasi_data_dukung 
               FROM tb_data_dukung 
               WHERE id_permasalahan = ? AND id_isustrategis = ?`
	rows, err := tx.QueryContext(ctx, script, permasalahanId, isuStrategisId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.DataDukung
	for rows.Next() {
		var dataDukung domain.DataDukung
		err := rows.Scan(&dataDukung.Id, &dataDukung.PermasalahanOpdId, &dataDukung.IdIsuStrategis, &dataDukung.DataDukung, &dataDukung.NarasiDataDukung)
		if err != nil {
			return nil, err
		}
		result = append(result, dataDukung)
	}
	return result, nil
}

func (repository *IsuStrategisRepositoryImpl) DeleteDataDukungByPermasalahanAndIsuStrategis(ctx context.Context, tx *sql.Tx, permasalahanId int, isuStrategisId int) error {
	log.Printf("[Repository] Deleting data dukung for permasalahan ID: %d and isu strategis ID: %d", permasalahanId, isuStrategisId)

	script := "DELETE FROM tb_data_dukung WHERE id_permasalahan = ? AND id_isustrategis = ?"
	result, err := tx.ExecContext(ctx, script, permasalahanId, isuStrategisId)
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

func (repository *IsuStrategisRepositoryImpl) FindallIsuKebelakang(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.IsuStrategis, error) {
	// Query dengan permasalahan_terpilih_id
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
        iso.created_at,
        pt.id as permasalahan_terpilih_id,
        p.id as permasalahan_opd_id,
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
        jd.jumlah as jumlah,
        jd.satuan as satuan
    FROM 
        tb_isu_strategis_opd iso
	LEFT JOIN 
	 	tb_permasalahan_isu_strategis pis ON pis.id_isu_strategis = iso.id
	LEFT JOIN 
		tb_permasalahan_terpilih pt ON pis.id_permasalahan = pt.id
	LEFT JOIN 
		tb_permasalahan_opd p ON pt.permasalahan_opd_id = p.id
    LEFT JOIN 
        tb_data_dukung dd ON dd.id_permasalahan = p.id AND dd.id_isustrategis = iso.id
    LEFT JOIN 
        tb_jumlah_data jd ON jd.id_data_dukung = dd.id
    WHERE 
        iso.kode_opd = ?
    ORDER BY 
        iso.created_at ASC,
        pt.id, p.id, dd.id, jd.tahun DESC`

	rows, err := tx.QueryContext(ctx, query, kodeOpd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Batch processing dengan map
	isuStrategisMap := make(map[int]*domain.IsuStrategis)
	permasalahanMap := make(map[int]*domain.Permasalahan)
	dataDukungMap := make(map[int]*domain.DataDukung)
	jumlahDataMap := make(map[int]map[string]*domain.JumlahData)

	// Fungsi helper untuk generate 6 tahun ke belakang
	generateYearRangeKebelakang := func(targetYear string) []string {
		if targetYear == "" {
			return []string{}
		}
		targetYearInt, _ := strconv.Atoi(targetYear)
		years := make([]string, 0, 6)
		for i := 0; i < 6; i++ {
			years = append(years, strconv.Itoa(targetYearInt-i))
		}
		return years
	}

	// Batch read dari database
	for rows.Next() {
		var (
			isuStrategisId         int
			kodeOpd                string
			namaOpd                string
			kodeBidangUrusan       string
			namaBidangUrusan       string
			tahunAwal              string
			tahunAkhir             string
			isuStrategis           string
			createdAt              time.Time
			permasalahanTerpilihId sql.NullInt64
			permasalahanOpdId      sql.NullInt64
			permasalahan           sql.NullString
			pKodeOpd               sql.NullString
			pTahun                 sql.NullString
			levelPohon             sql.NullInt64
			jenisMasalah           sql.NullString
			dataDukungId           sql.NullInt64
			namaDataDukung         sql.NullString
			narasiDataDukung       sql.NullString
			jumlahDataId           sql.NullInt64
			jumlahDataTahun        sql.NullString
			jumlah                 sql.NullFloat64
			satuan                 sql.NullString
		)

		err := rows.Scan(
			&isuStrategisId, &kodeOpd, &namaOpd, &kodeBidangUrusan, &namaBidangUrusan,
			&tahunAwal, &tahunAkhir, &isuStrategis, &createdAt,
			&permasalahanTerpilihId, &permasalahanOpdId, &permasalahan,
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
				CreatedAt:        createdAt,
				PermasalahanOpd:  make([]domain.Permasalahan, 0),
			}
			isuStrategisMap[isuStrategisId] = isuStr
		}

		// Handle Permasalahan - GUNAKAN permasalahanTerpilihId sebagai key
		if permasalahanTerpilihId.Valid {
			permTerpilihId := int(permasalahanTerpilihId.Int64)

			// Buat unique key: isuStrategisId + permTerpilihId
			// Karena permasalahan yang sama bisa di multiple isu strategis
			uniqueKey := isuStrategisId*10000 + permTerpilihId

			perm, exists := permasalahanMap[uniqueKey]
			if !exists {
				perm = &domain.Permasalahan{
					Id:           permTerpilihId,
					Permasalahan: permasalahan.String,
					KodeOpd:      pKodeOpd.String,
					Tahun:        pTahun.String,
					LevelPohon:   int(levelPohon.Int64),
					JenisMasalah: jenisMasalah.String,
					DataDukung:   make([]domain.DataDukung, 0),
				}
				permasalahanMap[uniqueKey] = perm
				isuStr.PermasalahanOpd = append(isuStr.PermasalahanOpd, *perm)
			}

			// Handle DataDukung
			if dataDukungId.Valid {
				ddId := int(dataDukungId.Int64)

				// Inisialisasi map jumlah data jika belum ada
				if _, exists := jumlahDataMap[ddId]; !exists {
					jumlahDataMap[ddId] = make(map[string]*domain.JumlahData)
				}

				// Simpan data jumlah data ke map
				if jumlahDataId.Valid && jumlahDataTahun.Valid {
					jumlahDataMap[ddId][jumlahDataTahun.String] = &domain.JumlahData{
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
						PermasalahanOpdId: int(permasalahanOpdId.Int64),
						IdIsuStrategis:    isuStrategisId,
						DataDukung:        namaDataDukung.String,
						NarasiDataDukung:  narasiDataDukung.String,
						JumlahData:        make([]domain.JumlahData, 0),
					}
					dataDukungMap[ddId] = dd

					// Tambahkan ke permasalahan yang benar
					for i := range isuStr.PermasalahanOpd {
						if isuStr.PermasalahanOpd[i].Id == permTerpilihId {
							isuStr.PermasalahanOpd[i].DataDukung = append(isuStr.PermasalahanOpd[i].DataDukung, *dd)
							break
						}
					}
				}
			}
		}
	}

	// Post-processing: Generate rentang tahun dan isi data
	for _, isuStr := range isuStrategisMap {
		for i, perm := range isuStr.PermasalahanOpd {
			for j, dd := range perm.DataDukung {
				var jumlahDataSlice []domain.JumlahData

				if tahun != "" {
					// Skenario 2: Filter by tahun dengan 6 tahun ke belakang
					yearRange := generateYearRangeKebelakang(tahun)
					jumlahDataSlice = make([]domain.JumlahData, 0, len(yearRange))

					for _, year := range yearRange {
						if data, exists := jumlahDataMap[dd.Id][year]; exists && data != nil {
							jumlahDataSlice = append(jumlahDataSlice, *data)
						} else {
							jumlahDataSlice = append(jumlahDataSlice, domain.JumlahData{
								Id:           0,
								IdDataDukung: dd.Id,
								Tahun:        year,
								JumlahData:   0,
								Satuan:       "",
							})
						}
					}
				} else {
					// Skenario 1: Tampilkan semua tahun DESC
					if dataMap, exists := jumlahDataMap[dd.Id]; exists {
						jumlahDataSlice = make([]domain.JumlahData, 0, len(dataMap))
						for _, data := range dataMap {
							if data != nil {
								jumlahDataSlice = append(jumlahDataSlice, *data)
							}
						}
						sort.Slice(jumlahDataSlice, func(x, y int) bool {
							return jumlahDataSlice[x].Tahun > jumlahDataSlice[y].Tahun
						})
					}
				}

				isuStr.PermasalahanOpd[i].DataDukung[j].JumlahData = jumlahDataSlice
			}
		}
	}

	// Convert map to slice
	result := make([]domain.IsuStrategis, 0, len(isuStrategisMap))
	for _, isuStr := range isuStrategisMap {
		result = append(result, *isuStr)
	}

	return result, nil
}

func (repository *IsuStrategisRepositoryImpl) DeleteDataDukungById(ctx context.Context, tx *sql.Tx, dataDukungId int) error {
	log.Printf("[Repository] Deleting data dukung with ID: %d", dataDukungId)

	script := "DELETE FROM tb_data_dukung WHERE id = ?"
	result, err := tx.ExecContext(ctx, script, dataDukungId)
	if err != nil {
		log.Printf("[Repository] Error deleting data dukung: %v", err)
		return err
	}

	affected, _ := result.RowsAffected()
	log.Printf("[Repository] Deleted %d data dukung records", affected)
	return nil
}
