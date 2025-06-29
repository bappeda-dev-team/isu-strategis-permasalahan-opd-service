package web

// @Description Update Request Permasalahan
type PermasalahanUpdateRequest struct {
	Id           int    `json:"id"`
	Permasalahan string `json:"permasalahan" validate:"required"`
	LevelPohon   int    `json:"level_pohon" validate:"required"`
	KodeOpd      string `json:"kode_opd" validate:"required"`
	NamaOpd      string `json:"nama_opd" validate:"required"`
	Tahun        string `json:"tahun" validate:"required"`
}
