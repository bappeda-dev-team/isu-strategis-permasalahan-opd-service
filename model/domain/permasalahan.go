package domain

type Permasalahan struct {
	Id           int
	PokinId      int
	Permasalahan string
	KodeOpd      string
	NamaOpd      string
	Tahun        string
	LevelPohon   int
	JenisMasalah string
	IsuStrategis int
	DataDukung   []DataDukung
}

type PermasalahanTerpilih struct {
	Id                int
	PermasalahanOpdId int
	KodeOpd           string
	Tahun             string
}
