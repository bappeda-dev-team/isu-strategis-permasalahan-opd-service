package domain

import "time"

type IsuStrategis struct {
	Id               int
	KodeOpd          string
	NamaOpd          string
	KodeBidangUrusan string
	NamaBidangUrusan string
	TahunAwal        string
	TahunAkhir       string
	IsuStrategis     string
	CreatedAt        time.Time
	PermasalahanOpd  []Permasalahan
}
