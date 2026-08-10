package web

// @Description Request Isu Strategis Create
type IsuStrategisCreateRequest struct {
	Id               	   int                      `json:"id"`
	KodeOpd          	   string                   `json:"kode_opd" validate:"required"`
	NamaOpd          	   string                   `json:"nama_opd" validate:"required"`
	KodeBidangUrusan 	   string                   `json:"kode_bidang_urusan" validate:"required"`
	NamaBidangUrusan 	   string                   `json:"nama_bidang_urusan" validate:"required"`
	TahunAwal        	   string                   `json:"tahun_awal"`
	TahunAkhir       	   string                   `json:"tahun_akhir"`
	IdPpd       	   	   *int                      `json:"id_ppd"`
	IdIsuKlhs       	   *int                      `json:"id_isu_klhs"`
	IdIsuGlobal       	   *int                      `json:"id_isu_global"`
	IdIsuNasional          *int                      `json:"id_isu_nasional"`
	IdIsuRegional          *int                      `json:"id_isu_regional"`
	PotensiPerangkatDaerah string                   `json:"potensi_perangkat_daerah"`
	IsuKlhs       		   string                   `json:"isu_klhs"`
	IsuGlobal       	   string                   `json:"isu_global"`
	IsuNasional       	   string                   `json:"isu_nasional"`
	IsuRegional       	   string                   `json:"isu_regional"`
	IsuStrategis     	   string                   `json:"isu_strategis" validate:"required"`
	PermasalahanOpd  	   []PermasalahanOpdRequest `json:"permasalahan_opd"`
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
