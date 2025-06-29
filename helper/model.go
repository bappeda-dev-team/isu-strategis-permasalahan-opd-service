package helper

import (
	"permasalahanService/model/domain"
	"permasalahanService/model/web"
)

func ToIsuStrategisResponse(isuStrategis domain.IsuStrategis) web.IsuStrategisResponse {
	return web.IsuStrategisResponse{
		Id:               isuStrategis.Id,
		KodeOpd:          isuStrategis.KodeOpd,
		NamaOpd:          isuStrategis.NamaOpd,
		KodeBidangUrusan: isuStrategis.KodeBidangUrusan,
		NamaBidangUrusan: isuStrategis.NamaBidangUrusan,
		TahunAwal:        isuStrategis.TahunAwal,
		TahunAkhir:       isuStrategis.TahunAkhir,
		IsuStrategis:     isuStrategis.IsuStrategis,
		CreatedAt:        isuStrategis.CreatedAt,
		PermasalahanOpd:  ToPermasalahanResponses(isuStrategis.PermasalahanOpd),
	}
}

func ToIsuStrategisResponses(isuStrategiss []domain.IsuStrategis) []web.IsuStrategisResponse {
	var responses []web.IsuStrategisResponse
	for _, isuStrategis := range isuStrategiss {
		responses = append(responses, ToIsuStrategisResponse(isuStrategis))
	}
	return responses
}

func ToPermasalahanResponse(permasalahan domain.Permasalahan) web.PermasalahanResponse {
	// Pastikan DataDukung tidak nil dan dikonversi dengan benar
	var dataDukung []web.DataDukungResponse
	if permasalahan.DataDukung != nil {
		dataDukung = ToDataDukungResponses(permasalahan.DataDukung)
	} else {
		dataDukung = make([]web.DataDukungResponse, 0)
	}

	return web.PermasalahanResponse{
		Id:           permasalahan.Id,
		Permasalahan: permasalahan.Permasalahan,
		LevelPohon:   permasalahan.LevelPohon,
		JenisMasalah: permasalahan.JenisMasalah,
		DataDukung:   dataDukung,
	}
}

func ToPermasalahanResponses(permasalahans []domain.Permasalahan) []web.PermasalahanResponse {
	var responses []web.PermasalahanResponse
	for _, permasalahan := range permasalahans {
		responses = append(responses, ToPermasalahanResponse(permasalahan))
	}
	return responses
}

func ToDataDukungResponse(dataDukung domain.DataDukung) web.DataDukungResponse {
	// Pastikan JumlahData tidak nil
	var jumlahData []web.JumlahDataResponse
	if dataDukung.JumlahData != nil {
		jumlahData = ToJumlahDataResponses(dataDukung.JumlahData)
	} else {
		jumlahData = make([]web.JumlahDataResponse, 0)
	}

	return web.DataDukungResponse{
		Id:                dataDukung.Id,
		PermasalahanOpdId: dataDukung.PermasalahanOpdId,
		DataDukung:        dataDukung.DataDukung,
		NarasiDataDukung:  dataDukung.NarasiDataDukung,
		JumlahData:        jumlahData,
	}
}

func ToDataDukungResponses(dataDukungs []domain.DataDukung) []web.DataDukungResponse {
	if dataDukungs == nil {
		return make([]web.DataDukungResponse, 0)
	}

	responses := make([]web.DataDukungResponse, 0)
	for _, dataDukung := range dataDukungs {
		responses = append(responses, ToDataDukungResponse(dataDukung))
	}
	return responses
}

func ToJumlahDataResponse(jumlahData domain.JumlahData) web.JumlahDataResponse {
	return web.JumlahDataResponse{
		Id:           jumlahData.Id,
		IdDataDukung: jumlahData.IdDataDukung,
		Tahun:        jumlahData.Tahun,
		JumlahData:   jumlahData.JumlahData,
		Satuan:       jumlahData.Satuan,
	}
}

func ToJumlahDataResponses(jumlahDatas []domain.JumlahData) []web.JumlahDataResponse {
	if jumlahDatas == nil {
		return make([]web.JumlahDataResponse, 0)
	}

	responses := make([]web.JumlahDataResponse, 0)
	for _, jumlahData := range jumlahDatas {
		responses = append(responses, ToJumlahDataResponse(jumlahData))
	}
	return responses
}
