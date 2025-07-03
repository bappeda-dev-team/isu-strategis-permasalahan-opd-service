package web

// @Description Request Permasalahan Create
type PermasalahanCreateRequest struct {
	PokinId      int    `json:"pokin_id" validate:"required"`
	Permasalahan string `json:"permasalahan" validate:"required"`
	LevelPohon   int    `json:"level_pohon" validate:"required"`
	// enum:MASALAH_POKOK,MASALAH,AKAR_MASALAH
	JenisMasalah string `json:"jenis_masalah" validate:"required"`
	KodeOpd      string `json:"kode_opd" validate:"required"`
	NamaOpd      string `json:"nama_opd" validate:"required"`
	Tahun        string `json:"tahun" validate:"required"`
}
