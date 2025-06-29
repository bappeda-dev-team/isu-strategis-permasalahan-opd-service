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

type JenisMasalah string

const (
	MASALAH_POKOK JenisMasalah = "MASALAH_POKOK"
	MASALAH       JenisMasalah = "MASALAH"
	AKAR_MASALAH  JenisMasalah = "AKAR_MASALAH"
)

func (j JenisMasalah) IsValid() bool {
	return j == MASALAH_POKOK || j == MASALAH || j == AKAR_MASALAH
}
