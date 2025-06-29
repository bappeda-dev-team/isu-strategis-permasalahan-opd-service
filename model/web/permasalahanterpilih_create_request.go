package web

type PermasalahanTerpilihRequest struct {
	AkarPermasalahanId int    `json:"masalah_id"`
	KodeOpd            string `json:"kode_opd"`
	Tahun              string `json:"tahun"`
}
