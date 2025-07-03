package web

// @Description Request Isu Strategis Update
type IsuStrategisUpdateRequest struct {
	Id               int                                     `json:"id" validate:"required"`
	KodeOpd          string                                  `json:"kode_opd" validate:"required"`
	NamaOpd          string                                  `json:"nama_opd" validate:"required"`
	KodeBidangUrusan string                                  `json:"kode_bidang_urusan" validate:"required"`
	NamaBidangUrusan string                                  `json:"nama_bidang_urusan" validate:"required"`
	TahunAwal        string                                  `json:"tahun_awal" validate:"required"`
	TahunAkhir       string                                  `json:"tahun_akhir" validate:"required"`
	IsuStrategis     string                                  `json:"isu_strategis" validate:"required"`
	PermasalahanOpd  []PermasalahanIsuStrategisUpdateRequest `json:"permasalahan_opd"`
}

type PermasalahanIsuStrategisUpdateRequest struct {
	PermasalahanOpdId int                       `json:"permasalahan_opd_id"`
	DataDukung        []DataDukungUpdateRequest `json:"data_dukung"`
}

type DataDukungUpdateRequest struct {
	Id               int                       `json:"id"`
	DataDukung       string                    `json:"data_dukung"`
	NarasiDataDukung string                    `json:"narasi_data_dukung"`
	JumlahData       []JumlahDataUpdateRequest `json:"jumlah_data"`
}

type JumlahDataUpdateRequest struct {
	Id         int     `json:"id"`
	Tahun      string  `json:"tahun"`
	JumlahData float64 `json:"jumlah_data"`
	Satuan     string  `json:"satuan"`
}
