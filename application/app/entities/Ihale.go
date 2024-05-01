package entities

import "time"

type IhaleDurumu string

type Ihale struct {
	IhaleNumarasi         string      `json:"ihaleNumarasi"`
	IsletmeMudurlugu      string      `json:"isletmeMudurlugu"`
	AnaKategori           string      `json:"anaKategori"`
	AltKategori           string      `json:"altKategori"`
	TasfiyeNedeni         string      `json:"tasfiyeNedeni"`
	BulunduguYer          string      `json:"bulunduguYer"`
	SatisaEsasBedel       float64     `json:"satisaEsasBedel"`
	TeminatBedeli         float64     `json:"teminatBedeli"`
	BaslangicBedeli       float64     `json:"baslangicBedeli"`
	GorusBaslangicTarihi  time.Time   `json:"gorusBaslangicTarihi"`
	TeklifBaslangicTarihi time.Time   `json:"teklifBaslangicTarihi"`
	IhaleBitisTarihi      time.Time   `json:"ihaleBitisTarihi"`
	DamgaVergisiOrani     float64     `json:"damgaVergisiOrani"`
	Aciklama              string      `json:"aciklama"` // Ek notlar ve özel şartlar
	OlusturanKullaniciID  string      `json:"olusturanKullaniciID"`
	GuncelFiyat           float64     `json:"guncelFiyat"`
	Durum                 IhaleDurumu `json:"durum"`
	KazananTeklifID       string      `json:"kazananTeklifID"`
}

// Ihale bilgileri için request yapısı
type CreateIhaleRequest struct {
	IsletmeMudurlugu      string    `json:"isletmeMudurlugu"`
	AnaKategori           string    `json:"anaKategori"`
	AltKategori           string    `json:"altKategori"`
	TasfiyeNedeni         string    `json:"tasfiyeNedeni"`
	BulunduguYer          string    `json:"bulunduguYer"`
	SatisaEsasBedel       float64   `json:"satisaEsasBedel"`
	TeminatBedeli         float64   `json:"teminatBedeli"`
	BaslangicBedeli       float64   `json:"baslangicBedeli"`
	GorusBaslangicTarihi  time.Time `json:"gorusBaslangicTarihi"`
	TeklifBaslangicTarihi time.Time `json:"teklifBaslangicTarihi"`
	IhaleBitisTarihi      time.Time `json:"ihaleBitisTarihi"`
	DamgaVergisiOrani     float64   `json:"damgaVergisiOrani"`
	Aciklama              string    `json:"aciklama"`
	OlusturanKullaniciID  string    `json:"olusturanKullaniciID"`
	GuncelFiyat           float64   `json:"guncelFiyat"`
}

type UpdateIhaleRequest struct {
	IhaleNumarasi         string      `json:"ihaleNumarasi"`
	IsletmeMudurlugu      string      `json:"isletmeMudurlugu"`
	AnaKategori           string      `json:"anaKategori"`
	AltKategori           string      `json:"altKategori"`
	TasfiyeNedeni         string      `json:"tasfiyeNedeni"`
	BulunduguYer          string      `json:"bulunduguYer"`
	SatisaEsasBedel       float64     `json:"satisaEsasBedel"`
	TeminatBedeli         float64     `json:"teminatBedeli"`
	BaslangicBedeli       float64     `json:"baslangicBedeli"`
	GorusBaslangicTarihi  time.Time   `json:"gorusBaslangicTarihi"`
	TeklifBaslangicTarihi time.Time   `json:"teklifBaslangicTarihi"`
	IhaleBitisTarihi      time.Time   `json:"ihaleBitisTarihi"`
	DamgaVergisiOrani     float64     `json:"damgaVergisiOrani"`
	Aciklama              string      `json:"aciklama"` // Ek notlar ve özel şartlar
	OlusturanKullaniciID  string      `json:"olusturanKullaniciID"`
	GuncelFiyat           float64     `json:"guncelFiyat"`
	Durum                 IhaleDurumu `json:"durum"`
	KazananTeklifID       string      `json:"kazananTeklifID"`
}
