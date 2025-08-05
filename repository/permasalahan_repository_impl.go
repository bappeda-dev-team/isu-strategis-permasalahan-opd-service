package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"permasalahanService/model/domain"
	"permasalahanService/model/web"
)

type PermasalahanRepositoryImpl struct {
}

func NewPermasalahanRepositoryImpl() *PermasalahanRepositoryImpl {
	return &PermasalahanRepositoryImpl{}
}

func (repository *PermasalahanRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, permasalahan domain.Permasalahan) (domain.Permasalahan, error) {
	script := "INSERT INTO tb_permasalahan_opd (pokin_id, permasalahan, level_pohon, kode_opd, tahun, jenis_masalah, nama_opd) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, script, permasalahan.PokinId, permasalahan.Permasalahan, permasalahan.LevelPohon, permasalahan.KodeOpd, permasalahan.Tahun, permasalahan.JenisMasalah, permasalahan.NamaOpd)
	if err != nil {
		return domain.Permasalahan{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.Permasalahan{}, err
	}
	permasalahan.Id = int(id)

	return permasalahan, nil
}

func (repository *PermasalahanRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, permasalahan domain.Permasalahan) domain.Permasalahan {
	// Cek dulu data yang ada
	existing, err := repository.FindById(ctx, tx, fmt.Sprintf("%d", permasalahan.Id))
	if err != nil {
		log.Printf("Error finding existing data: %v", err)
		return domain.Permasalahan{}
	}

	// Jika data sama persis, langsung return data yang ada
	if existing.Permasalahan == permasalahan.Permasalahan &&
		existing.LevelPohon == permasalahan.LevelPohon &&
		existing.KodeOpd == permasalahan.KodeOpd &&
		existing.NamaOpd == permasalahan.NamaOpd &&
		existing.Tahun == permasalahan.Tahun {
		log.Printf("No changes detected, returning existing data")
		return existing
	}

	script := "UPDATE tb_permasalahan_opd SET permasalahan = ?, kode_opd = ?, tahun = ?, nama_opd = ? WHERE id = ?"

	result, err := tx.ExecContext(ctx, script,
		permasalahan.Permasalahan,
		permasalahan.KodeOpd,
		permasalahan.Tahun,
		permasalahan.NamaOpd,
		permasalahan.Id)

	if err != nil {
		log.Printf("Error executing update: %v", err)
		return domain.Permasalahan{}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return domain.Permasalahan{}
	}

	// Jika tidak ada yang berubah (rowsAffected = 0), itu bukan error
	// Kita tetap return data yang ada
	if rowsAffected == 0 {
		log.Printf("No rows affected, data might be identical")
		return existing
	}

	// Ambil data terbaru jika ada perubahan
	updated, err := repository.FindById(ctx, tx, fmt.Sprintf("%d", permasalahan.Id))
	if err != nil {
		log.Printf("Error getting updated data: %v", err)
		return domain.Permasalahan{}
	}

	return updated
}

func (repository *PermasalahanRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	script := "DELETE FROM tb_permasalahan_opd WHERE id = ?"
	result, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return errors.New("id not found")
	}
	return nil
}
func (repository *PermasalahanRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domain.Permasalahan, error) {
	log.Printf("Finding permasalahan with ID: %s", id)

	script := "SELECT id, pokin_id, permasalahan, level_pohon, kode_opd, nama_opd, tahun, jenis_masalah FROM tb_permasalahan_opd WHERE id = ?"

	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		return domain.Permasalahan{}, err
	}
	defer rows.Close()

	permasalahan := domain.Permasalahan{}
	if rows.Next() {
		err := rows.Scan(
			&permasalahan.Id,
			&permasalahan.PokinId,
			&permasalahan.Permasalahan,
			&permasalahan.LevelPohon,
			&permasalahan.KodeOpd,
			&permasalahan.NamaOpd,
			&permasalahan.Tahun,
			&permasalahan.JenisMasalah,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return domain.Permasalahan{}, err
		}
		log.Printf("Found permasalahan: %+v", permasalahan)
		return permasalahan, nil
	}

	log.Printf("No permasalahan found with ID: %s", id)
	return domain.Permasalahan{}, errors.New("permasalahan not found")
}

func (repository *PermasalahanRepositoryImpl) FindByKodeOpdAndTahun(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.Permasalahan, error) {
	script := "SELECT id, pokin_id, permasalahan, level_pohon, kode_opd, nama_opd, tahun, jenis_masalah FROM tb_permasalahan_opd WHERE kode_opd = ? AND tahun = ?"
	rows, err := tx.QueryContext(ctx, script, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permasalahans []domain.Permasalahan
	for rows.Next() {
		permasalahan := domain.Permasalahan{}
		err := rows.Scan(
			&permasalahan.Id,
			&permasalahan.PokinId,
			&permasalahan.Permasalahan,
			&permasalahan.LevelPohon,
			&permasalahan.KodeOpd,
			&permasalahan.NamaOpd,
			&permasalahan.Tahun,
			&permasalahan.JenisMasalah,
		)
		if err != nil {
			return nil, err
		}
		permasalahans = append(permasalahans, permasalahan)
	}
	return permasalahans, nil
}

func (repository *PermasalahanRepositoryImpl) GetPohonKinerjaFromAPI(ctx context.Context, kodeOpd string, tahun string) (*web.PohonKinerjaDataResponse, error) {
	apiPokinOpd := os.Getenv("API_POKIN_OPD")
	url := fmt.Sprintf("%s/api/pokin_opd/findall/%s/%s", apiPokinOpd, kodeOpd, tahun)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Temporary struct untuk decode response API
	var apiResponse struct {
		Code   int                          `json:"code"`
		Status string                       `json:"status"`
		Data   web.PohonKinerjaDataResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	// Hanya mengembalikan data tanpa wrapper
	return &apiResponse.Data, nil
}

func (repository *PermasalahanRepositoryImpl) MergePohonKinerjaWithPermasalahan(ctx context.Context, tx *sql.Tx, pohonKinerja *web.PohonKinerjaDataResponse, permasalahans []domain.Permasalahan) *web.PohonKinerjaDataResponse {
	// Membuat map untuk mempermudah pencarian permasalahan berdasarkan pokin_id
	permasalahanMap := make(map[int]domain.Permasalahan)
	for _, p := range permasalahans {
		permasalahanMap[p.PokinId] = p
	}

	// Fungsi rekursif untuk memproses child nodes
	var processNode func(node *web.ChildResponse) *web.ChildResponse
	processNode = func(node *web.ChildResponse) *web.ChildResponse {
		// Filter berdasarkan level pohon dan status
		if node.LevelPohon < 4 || node.LevelPohon > 6 {
			return nil
		}

		// Filter berdasarkan status
		if node.Status == "menunggu_disetujui" {
			return nil
		}

		// Cek apakah node memiliki ID yang cocok dengan permasalahan
		if permasalahan, exists := permasalahanMap[node.Id]; exists {
			// Update node dengan data permasalahan
			node.NamaPohon = permasalahan.Permasalahan
			node.IsPermasalahan = true
			node.IdPermasalahan = permasalahan.Id

			// Cek apakah permasalahan ini terpilih
			isTerpilih, err := repository.IsPermasalahanTerpilih(ctx, tx, permasalahan.Id)
			if err == nil && isTerpilih {
				node.PermasalahanTerpilih = true
			} else {
				node.PermasalahanTerpilih = false
			}
		} else {
			node.IsPermasalahan = false
			node.PermasalahanTerpilih = false
			node.IdPermasalahan = 0
		}

		// Proses child nodes jika ada
		var filteredChilds []web.ChildResponse
		for i := range node.Childs {
			if filteredChild := processNode(&node.Childs[i]); filteredChild != nil {
				filteredChilds = append(filteredChilds, *filteredChild)
			}
		}
		node.Childs = filteredChilds

		return node
	}

	// Proses semua child nodes di root
	var filteredRootChilds []web.ChildResponse
	for i := range pohonKinerja.Childs {
		if filteredChild := processNode(&pohonKinerja.Childs[i]); filteredChild != nil {
			filteredRootChilds = append(filteredRootChilds, *filteredChild)
		}
	}
	pohonKinerja.Childs = filteredRootChilds

	return pohonKinerja
}

func (repository *PermasalahanRepositoryImpl) FindByPokinId(ctx context.Context, tx *sql.Tx, pokinId int) (domain.Permasalahan, error) {
	script := "SELECT id, pokin_id, permasalahan, level_pohon, kode_opd, tahun FROM tb_permasalahan_opd WHERE pokin_id = ?"
	rows, err := tx.QueryContext(ctx, script, pokinId)
	if err != nil {
		return domain.Permasalahan{}, err
	}
	defer rows.Close()

	permasalahan := domain.Permasalahan{}
	if rows.Next() {
		err := rows.Scan(
			&permasalahan.Id,
			&permasalahan.PokinId,
			&permasalahan.Permasalahan,
			&permasalahan.LevelPohon,
			&permasalahan.KodeOpd,
			&permasalahan.Tahun,
		)
		if err != nil {
			return domain.Permasalahan{}, err
		}
		return permasalahan, nil
	}
	return permasalahan, nil
}

func (repository *PermasalahanRepositoryImpl) IsPermasalahanTerpilih(ctx context.Context, tx *sql.Tx, idPermasalahan int) (bool, error) {
	script := `
        SELECT COUNT(*) 
        FROM tb_permasalahan_terpilih pt 
        WHERE pt.permasalahan_opd_id = ?
    `
	var count int
	err := tx.QueryRowContext(ctx, script, idPermasalahan).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repository *PermasalahanRepositoryImpl) FindByIsuStrategisId(ctx context.Context, tx *sql.Tx, isuStrategisId int) ([]domain.Permasalahan, error) {
	script := `SELECT id, isu_strategis_id FROM tb_permasalahan_opd WHERE isu_strategis_id = ?`
	rows, err := tx.QueryContext(ctx, script, isuStrategisId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Permasalahan
	for rows.Next() {
		var permasalahan domain.Permasalahan
		err := rows.Scan(&permasalahan.Id, &permasalahan.IsuStrategis)
		if err != nil {
			return nil, err
		}
		result = append(result, permasalahan)
	}
	return result, nil
}

func (repository *PermasalahanRepositoryImpl) ResetIsuStrategisId(ctx context.Context, tx *sql.Tx, id int) error {
	script := `UPDATE tb_permasalahan_opd SET isu_strategis_id = 0 WHERE id = ?`
	_, err := tx.ExecContext(ctx, script, id)
	return err
}
