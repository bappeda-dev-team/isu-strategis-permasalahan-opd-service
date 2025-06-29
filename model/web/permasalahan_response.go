package web

type PohonKinerjaResponse struct {
	Code   int                      `json:"code"`
	Status string                   `json:"status"`
	Data   PohonKinerjaDataResponse `json:"data"`
}

// @Description Response Permasalahan OPD
type PohonKinerjaDataResponse struct {
	KodeOpd string          `json:"kode_opd"`
	NamaOpd string          `json:"nama_opd"`
	Tahun   string          `json:"tahun"`
	Childs  []ChildResponse `json:"childs"`
}

// type TujuanOpdResponse struct {
// 	Id        int                 `json:"id"`
// 	KodeOpd   string              `json:"kode_opd"`
// 	Tujuan    string              `json:"tujuan"`
// 	Indikator []IndikatorResponse `json:"indikator"`
// }

type IndikatorResponse struct {
	Indikator string           `json:"indikator"`
	Targets   []TargetResponse `json:"targets"`
}

type TargetResponse struct {
	Tahun  string `json:"tahun"`
	Target string `json:"target"`
	Satuan string `json:"satuan"`
}

type ChildResponse struct {
	Id                   int             `json:"id"`
	IdPermasalahan       int             `json:"id_permasalahan"`
	Parent               *int            `json:"parent"`
	NamaPohon            string          `json:"nama_pohon"`
	LevelPohon           int             `json:"level_pohon"`
	PerangkatDaerah      PerangkatDaerah `json:"perangkat_daerah"`
	IsPermasalahan       bool            `json:"is_permasalahan,omitempty"`
	PermasalahanTerpilih bool            `json:"permasalahan_terpilih,omitempty"`
	JenisMasalah         string          `json:"jenis_masalah"`
	Status               string          `json:"status,omitempty"`
	Childs               []ChildResponse `json:"childs,omitempty"`
}

type PerangkatDaerah struct {
	KodeOpd string `json:"kode_opd"`
	NamaOpd string `json:"nama_opd"`
}

type Pelaksana struct {
	IdPelaksana string `json:"id_pelaksana"`
	PegawaiId   string `json:"pegawai_id"`
	Nip         string `json:"nip"`
	NamaPegawai string `json:"nama_pegawai"`
}

type IndikatorPokin struct {
	IdIndikator   string        `json:"id_indikator"`
	NamaIndikator string        `json:"nama_indikator"`
	Targets       []TargetPokin `json:"targets"`
}

type TargetPokin struct {
	IdTarget    string `json:"id_target"`
	IndikatorId string `json:"indikator_id"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}
