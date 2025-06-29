package web

import "time"

type IsuStrategisResponse struct {
	Id               int                    `json:"id"`
	KodeOpd          string                 `json:"kode_opd"`
	NamaOpd          string                 `json:"nama_opd"`
	KodeBidangUrusan string                 `json:"kode_bidang_urusan"`
	NamaBidangUrusan string                 `json:"nama_bidang_urusan"`
	TahunAwal        string                 `json:"tahun_awal"`
	TahunAkhir       string                 `json:"tahun_akhir"`
	IsuStrategis     string                 `json:"isu_strategis"`
	CreatedAt        time.Time              `json:"created_at"`
	PermasalahanOpd  []PermasalahanResponse `json:"permasalahan_opd"`
}

type PermasalahanResponse struct {
	Id           int                  `json:"id"`
	Permasalahan string               `json:"masalah"`
	LevelPohon   int                  `json:"level_pohon"`
	JenisMasalah string               `json:"jenis_masalah"`
	DataDukung   []DataDukungResponse `json:"data_dukung"`
}

type DataDukungResponse struct {
	Id                int                  `json:"id"`
	PermasalahanOpdId int                  `json:"permasalahan_opd_id"`
	DataDukung        string               `json:"data_dukung"`
	NarasiDataDukung  string               `json:"narasi_data_dukung"`
	JumlahData        []JumlahDataResponse `json:"jumlah_data"`
}

type JumlahDataResponse struct {
	Id           int     `json:"id"`
	IdDataDukung int     `json:"id_data_dukung"`
	Tahun        string  `json:"tahun"`
	JumlahData   float64 `json:"jumlah_data"`
	Satuan       string  `json:"satuan"`
}
