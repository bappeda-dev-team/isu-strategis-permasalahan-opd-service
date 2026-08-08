package internal

type FindByIdsRequest struct {
	Ids []int `json:"ids" validate:"required,min=1"`
}

type PpdItem struct {
	Id               int    `json:"id"`
	KodeBidangUrusan string `json:"kode_bidang_urusan"`
	KodeOpd          string `json:"kode_opd"`
	Potensi          string `json:"potensi"`
	Tahunu           int    `json:"tahun"`
}
type IsuItem struct {
	Id               int    `json:"id"`
	KodeBidangUrusan string `json:"kode_bidang_urusan"`
	KodeOpd          string `json:"kode_opd"`
	Isu              string `json:"isu"`
	Tahunu           int    `json:"tahun"`
}

type PotensiPerangkatDaerahResponse struct {
	Code   int       `json:"code"`
	Status string    `json:"status"`
	Data   []PpdItem `json:"data"`
}

type IsuResponse struct {
	Code   int       `json:"code"`
	Status string    `json:"status"`
	Data   []IsuItem `json:"data"`
}