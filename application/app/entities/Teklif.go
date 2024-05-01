package entities

import "time"

type Teklif struct {
	IhaleNumarasi string    `json:"ihaleNumarasi"`
	TeklifVerenID string    `json:"teklifVerenID"` // KullaniciID ile değiştirildi
	Fiyat         float64   `json:"fiyat"`
	TeklifZamani  time.Time `json:"teklifZamani"`
}

// Teklif verme işlemi için kullanılan yapı
type TeklifYapRequest struct {
	IhaleNumarasi string  `json:"ihaleNumarasi"`
	TeklifTutari  float64 `json:"teklifTutari"`
	KullaniciID   string  `json:"kullaniciID"`
}
