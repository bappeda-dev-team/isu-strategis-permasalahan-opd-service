package domain

type DataDukung struct {
	Id                int
	PermasalahanOpdId int
	IdIsuStrategis    int
	DataDukung        string
	NarasiDataDukung  string
	JumlahData        []JumlahData
}
