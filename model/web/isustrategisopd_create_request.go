package web

type IsuStrategisCreateRequest struct {
	Id               int                      `json:"id"`
	KodeOpd          string                   `json:"kode_opd" validate:"required"`
	NamaOpd          string                   `json:"nama_opd" validate:"required"`
	KodeBidangUrusan string                   `json:"kode_bidang_urusan" validate:"required"`
	NamaBidangUrusan string                   `json:"nama_bidang_urusan" validate:"required"`
	TahunAwal        string                   `json:"tahun_awal" validate:"required"`
	TahunAkhir       string                   `json:"tahun_akhir" validate:"required"`
	IsuStrategis     string                   `json:"isu_strategis" validate:"required"`
	PermasalahanOpd  []PermasalahanOpdRequest `json:"permasalahan_opd"`
}

type PermasalahanOpdRequest struct {
	IdPermasalahan int                 `json:"id_permasalahan"`
	DataDukung     []DataDukungRequest `json:"data_dukung"`
}

type DataDukungRequest struct {
	Id                int                 `json:"id"`
	PermasalahanOpdId int                 `json:"permasalahan_opd_id"`
	DataDukung        string              `json:"data_dukung" validate:"required"`
	NarasiDataDukung  string              `json:"narasi_data_dukung" validate:"required"`
	JumlahData        []JumlahDataRequest `json:"jumlah_data"`
}

type JumlahDataRequest struct {
	Id           int     `json:"id"`
	IdDataDukung int     `json:"id_data_dukung"`
	Tahun        string  `json:"tahun"`
	JumlahData   float64 `json:"jumlah_data"`
	Satuan       string  `json:"satuan"`
}
