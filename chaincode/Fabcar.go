package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract represents the main contract
type SmartContract struct {
	contractapi.Contract
}

type IhaleDurumu string

const (
	IhaleDurumuBaslamadi IhaleDurumu = "BASLAMADI"
	IhaleDurumuAcik      IhaleDurumu = "ACIK"
	IhaleDurumuKapali    IhaleDurumu = "KAPALI"
)

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

type Teklif struct {
	IhaleNumarasi string    `json:"ihaleNumarasi"`
	KullaniciID   string    `json:"kullaniciID"` // KullaniciID ile değiştirildi
	Fiyat         float64   `json:"fiyat"`
	TeklifZamani  time.Time `json:"teklifZamani"`
}

type Kullanici struct {
	KullaniciID  string    `json:"kullaniciID"`
	KullaniciAdi string    `json:"kullaniciAdi"`
	Sifre        string    `json:"sifre"`
	Isim         string    `json:"isim"`
	Soyisim      string    `json:"soyisim"`
	Eposta       string    `json:"eposta"`
	Telefon      string    `json:"telefon"`
	TCKimlikNo   string    `json:"tcKimlikNo"`
	DogumTarihi  time.Time `json:"dogumTarihi"`
}

type KullaniciAktivitesi struct {
	KullaniciID  string `json:"kullaniciID"`
	TeklifSayisi int    `json:"teklifSayisi"`
}

type IhaleHistoryItem struct {
	TxID      string    `json:"txID"`
	Value     Ihale     `json:"value"`
	Timestamp time.Time `json:"timestamp"`
	IsDeleted bool      `json:"isDeleted"`
}

type TeklifHistoryItem struct {
	TxID      string    `json:"txID"`
	Value     Teklif    `json:"value"`
	Timestamp time.Time `json:"timestamp"`
	IsDeleted bool      `json:"isDeleted"`
}

type IhaleIstatistikleri struct {
	ToplamIhaleSayisi      int     `json:"toplamIhaleSayisi"`
	OrtalamaSatisEsasBedel float64 `json:"ortalamaSatisEsasBedel"`
}

func (s *SmartContract) CreateIhale(ctx contractapi.TransactionContextInterface, isletmeMudurlugu string, anaKategori string, altKategori string, tasfiyeNedeni string, bulunduguYer string, satisaEsasBedel float64, teminatBedeli float64, baslangicBedeli float64, gorusBaslangicTarihi, teklifBaslangicTarihi, ihaleBitisTarihi time.Time, damgaVergisiOrani float64, aciklama string, olusturanKullaniciID string, guncelFiyat float64) error {
	// Yeni bir UUID oluştur
	ihaleNumarasi := uuid.New().String()

	ihale := Ihale{
		IhaleNumarasi:         ihaleNumarasi,
		IsletmeMudurlugu:      isletmeMudurlugu,
		AnaKategori:           anaKategori,
		AltKategori:           altKategori,
		TasfiyeNedeni:         tasfiyeNedeni,
		BulunduguYer:          bulunduguYer,
		SatisaEsasBedel:       satisaEsasBedel,
		TeminatBedeli:         teminatBedeli,
		BaslangicBedeli:       baslangicBedeli,
		GorusBaslangicTarihi:  gorusBaslangicTarihi,
		TeklifBaslangicTarihi: teklifBaslangicTarihi,
		IhaleBitisTarihi:      ihaleBitisTarihi,
		DamgaVergisiOrani:     damgaVergisiOrani,
		Aciklama:              aciklama,
		OlusturanKullaniciID:  olusturanKullaniciID,
		GuncelFiyat:           baslangicBedeli,
		Durum:                 IhaleDurumuAcik,
		KazananTeklifID:       "",
	}

	ihaleJSON, err := json.Marshal(ihale)
	if err != nil {
		return fmt.Errorf("ihaleyi JSON'a dönüştürme hatası: %s", err)
	}

	// Ihale bilgisini ledger'a kaydet
	return ctx.GetStub().PutState(ihaleNumarasi, ihaleJSON)
}

func (s *SmartContract) UpdateIhale(ctx contractapi.TransactionContextInterface, ihaleNumarasi string, isletmeMudurlugu string, anaKategori string, altKategori string, tasfiyeNedeni string, bulunduguYer string, satisaEsasBedel float64, teminatBedeli float64, baslangicBedeli float64, gorusBaslangicTarihi, teklifBaslangicTarihi, ihaleBitisTarihi time.Time, damgaVergisiOrani float64, aciklama string, olusturanKullaniciID string, guncelFiyat float64) error {
	// İhaleyi ledger'dan sorgula
	ihaleJSON, err := ctx.GetStub().GetState(ihaleNumarasi)
	if err != nil || ihaleJSON == nil {
		return fmt.Errorf("ihale bulunamadı: %s", ihaleNumarasi)
	}

	// Yeni ihale bilgileri ile ihale nesnesi oluştur
	updatedIhale := Ihale{
		IhaleNumarasi:         ihaleNumarasi,
		IsletmeMudurlugu:      isletmeMudurlugu,
		AnaKategori:           anaKategori,
		AltKategori:           altKategori,
		TasfiyeNedeni:         tasfiyeNedeni,
		BulunduguYer:          bulunduguYer,
		SatisaEsasBedel:       satisaEsasBedel,
		TeminatBedeli:         teminatBedeli,
		BaslangicBedeli:       baslangicBedeli,
		GorusBaslangicTarihi:  gorusBaslangicTarihi,
		TeklifBaslangicTarihi: teklifBaslangicTarihi,
		IhaleBitisTarihi:      ihaleBitisTarihi,
		DamgaVergisiOrani:     damgaVergisiOrani,
		Aciklama:              aciklama,
		GuncelFiyat:           guncelFiyat,
	}

	// İhale nesnesini JSON'a dönüştür
	updatedIhaleJSON, err := json.Marshal(updatedIhale)
	if err != nil {
		return fmt.Errorf("ihaleyi JSON'a dönüştürme hatası: %s", err)
	}

	// Güncellenmiş ihale bilgisini ledger'a yaz
	return ctx.GetStub().PutState(ihaleNumarasi, updatedIhaleJSON)
}

func (s *SmartContract) QueryKullaniciIhaleleri(ctx contractapi.TransactionContextInterface, olusturanKullaniciID string) ([]Ihale, error) {
	queryResult, err := ctx.GetStub().GetQueryResult(fmt.Sprintf(`{"selector":{"olusturanKullaniciID":"%s"}}`, olusturanKullaniciID))
	if err != nil {
		return nil, err
	}
	defer queryResult.Close()

	var ihaleler []Ihale
	for queryResult.HasNext() {
		queryResponse, err := queryResult.Next()
		if err != nil {
			return nil, err
		}

		var ihale Ihale
		err = json.Unmarshal(queryResponse.Value, &ihale)
		if err != nil {
			return nil, err
		}

		ihaleler = append(ihaleler, ihale)
	}

	return ihaleler, nil
}

func (s *SmartContract) QueryDigerIhaleler(ctx contractapi.TransactionContextInterface, olusturanKullaniciID string) ([]Ihale, error) {
	queryResult, err := ctx.GetStub().GetQueryResult(fmt.Sprintf(`{"selector":{"olusturanKullaniciID":{"$ne":"%s"}}}`, olusturanKullaniciID))
	if err != nil {
		return nil, err
	}
	defer queryResult.Close()

	var ihaleler []Ihale
	for queryResult.HasNext() {
		queryResponse, err := queryResult.Next()
		if err != nil {
			return nil, err
		}

		var ihale Ihale
		err = json.Unmarshal(queryResponse.Value, &ihale)
		if err != nil {
			return nil, err
		}

		ihaleler = append(ihaleler, ihale)
	}

	return ihaleler, nil
}

func (s *SmartContract) QueryTumIhaleler(ctx contractapi.TransactionContextInterface, ihaleNumarasi string) (*Ihale, error) {
	ihaleJSON, err := ctx.GetStub().GetState(ihaleNumarasi)
	if err != nil {
		return nil, err
	}

	var ihale Ihale
	err = json.Unmarshal(ihaleJSON, &ihale)
	if err != nil {
		return nil, err
	}

	return &ihale, nil
}

func (s *SmartContract) GetIhaleByTeklif(ctx contractapi.TransactionContextInterface, teklifID string) (*Ihale, error) {
	// Get the bid from the ledger
	teklifJSON, err := ctx.GetStub().GetState(teklifID)
	if err != nil {
		return nil, fmt.Errorf("teklif bulunamadı: %s", err)
	}

	var teklif Teklif
	if err := json.Unmarshal(teklifJSON, &teklif); err != nil {
		return nil, fmt.Errorf("teklif verisi çözümlenirken hata oluştu: %s", err)
	}

	// Retrieve the associated auction using the auction ID stored in the bid
	ihaleJSON, err := ctx.GetStub().GetState(teklif.IhaleNumarasi)
	if err != nil {
		return nil, fmt.Errorf("ihale bulunamadı: %s", err)
	}

	var ihale Ihale
	if err := json.Unmarshal(ihaleJSON, &ihale); err != nil {
		return nil, fmt.Errorf("ihale verisi çözümlenirken hata oluştu: %s", err)
	}

	return &ihale, nil
}

func (s *SmartContract) ListIhaleler(ctx contractapi.TransactionContextInterface) ([]Ihale, error) {
	startKey := ""
	endKey := ""
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var ihaleler []Ihale
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var ihale Ihale
		err = json.Unmarshal(queryResponse.Value, &ihale)
		if err != nil {
			return nil, err
		}
		ihaleler = append(ihaleler, ihale)
	}

	return ihaleler, nil
}

func (s *SmartContract) MakeTeklif(ctx contractapi.TransactionContextInterface, ihaleNumarasi, kullaniciID string, fiyat float64) error {
	// İhaleyi sorgula
	ihaleJSON, err := ctx.GetStub().GetState(ihaleNumarasi)
	if err != nil || ihaleJSON == nil {
		return fmt.Errorf("ihale bulunamadı: %s", ihaleNumarasi)
	}

	// İhaleyi deserialize et
	var ihale Ihale
	err = json.Unmarshal(ihaleJSON, &ihale)
	if err != nil {
		return fmt.Errorf("ihale JSON'unun deserialize edilmesi başarısız: %s", err)
	}

	// İhale durumunu kontrol et
	if ihale.Durum != IhaleDurumuAcik {
		return fmt.Errorf("ihale şu anda teklif kabul etmiyor")
	}

	// Teklif zamanını kontrol et
	currentTime := time.Now()
	if currentTime.Before(ihale.TeklifBaslangicTarihi) || currentTime.After(ihale.IhaleBitisTarihi) {
		return fmt.Errorf("teklif için uygun zaman aralığı dışında")
	}

	// Kullanıcının kendi ihalesine teklif verip vermediğini kontrol et
	if ihale.OlusturanKullaniciID == kullaniciID {
		return fmt.Errorf("kullanıcı kendi ihalesine teklif veremez")
	}

	// Mevcut en yüksek tekliften daha yüksek bir teklif verildiğinden emin ol
	if fiyat <= ihale.GuncelFiyat {
		return fmt.Errorf("teklif, mevcut en yüksek tekliften (%f) daha yüksek olmalıdır", ihale.GuncelFiyat)
	}

	// Teklif işlemleri
	teklif := Teklif{
		IhaleNumarasi: ihaleNumarasi,
		KullaniciID:   kullaniciID,
		Fiyat:         fiyat,
		TeklifZamani:  currentTime,
	}

	teklifKey, err := ctx.GetStub().CreateCompositeKey("Teklif", []string{ihaleNumarasi, kullaniciID})
	if err != nil {
		return err
	}

	teklifJSON, err := json.Marshal(teklif)
	if err != nil {
		return err
	}

	// Güncel fiyatı teklif ile güncelle
	ihale.GuncelFiyat = fiyat

	// İhaleyi güncelle
	ihaleJSON, err = json.Marshal(ihale)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(ihaleNumarasi, ihaleJSON)
	if err != nil {
		return err
	}

	// Teklifi ledger'a kaydet
	return ctx.GetStub().PutState(teklifKey, teklifJSON)
}

func (s *SmartContract) QueryTeklif(ctx contractapi.TransactionContextInterface, ihaleNumarasi, teklifVerenID string) (*Teklif, error) {
	teklifKey, err := ctx.GetStub().CreateCompositeKey("Teklif", []string{ihaleNumarasi, teklifVerenID})
	if err != nil {
		return nil, err
	}

	teklifJSON, err := ctx.GetStub().GetState(teklifKey)
	if err != nil || teklifJSON == nil {
		return nil, fmt.Errorf("teklif bulunamadı veya okunamadı")
	}

	var teklif Teklif
	err = json.Unmarshal(teklifJSON, &teklif)
	if err != nil {
		return nil, err
	}

	return &teklif, nil
}

func (s *SmartContract) ListTekliflerByIhale(ctx contractapi.TransactionContextInterface, ihaleNumarasi string) ([]Teklif, error) {
	queryIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Teklif", []string{ihaleNumarasi})
	if err != nil {
		return nil, err
	}
	defer queryIterator.Close()

	var teklifler []Teklif
	for queryIterator.HasNext() {
		response, err := queryIterator.Next()
		if err != nil {
			return nil, err
		}

		var teklif Teklif
		err = json.Unmarshal(response.Value, &teklif)
		if err != nil {
			return nil, err
		}
		teklifler = append(teklifler, teklif)
	}

	return teklifler, nil
}

func (s *SmartContract) ListTekliflerByKullanici(ctx contractapi.TransactionContextInterface, kullaniciID string) ([]Teklif, error) {
	// Kullanıcı ID'sine göre teklifleri bulmak için sorgu dizgisini oluştur
	queryString := fmt.Sprintf(`{"selector":{"kullaniciID":"%s"}}`, kullaniciID)

	// Seçici kullanarak sorgu sonucunu al
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("kullanıcı ID'si ile teklifler alınırken hata oluştu: %s", kullaniciID) // Türkçe hata mesajı
	}
	defer resultsIterator.Close()

	var teklifler []Teklif
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("bir sonraki teklif alınırken hata oluştu: %s", err) // Türkçe hata mesajı
		}

		var teklif Teklif
		// JSON nesnesini Teklif yapısına çözümle
		if err = json.Unmarshal(response.Value, &teklif); err != nil {
			return nil, fmt.Errorf("teklif verisi çözümlenirken hata oluştu: %s", err) // Türkçe hata mesajı
		}

		// Listeye teklifi ekle
		teklifler = append(teklifler, teklif)
	}

	return teklifler, nil
}

func (s *SmartContract) RegisterUser(ctx contractapi.TransactionContextInterface, kullaniciAdi, sifre, isim, soyisim, eposta, telefon, tcKimlikNo string, dogumTarihi time.Time) error {
	// TC Kimlik No ile var olan bir kullanıcı olup olmadığını kontrol et
	tcQueryResult, err := s.GetUserByTCKimlikNo(ctx, tcKimlikNo)
	if err == nil && tcQueryResult != nil {
		return fmt.Errorf("TC Kimlik No zaten kayıtlı: %s", tcKimlikNo)
	}

	// Kullanıcı adı ile var olan bir kullanıcı olup olmadığını kontrol et
	usernameQueryResult, err := s.GetUserByUsername(ctx, kullaniciAdi)
	if err == nil && usernameQueryResult != nil {
		return fmt.Errorf("Kullanıcı adı zaten kayıtlı: %s", kullaniciAdi)
	}

	// Yeni bir UUID oluştur
	kullaniciID := uuid.New().String()

	kullanici := Kullanici{
		KullaniciID:  kullaniciID,
		KullaniciAdi: kullaniciAdi,
		Sifre:        sifre,
		Isim:         isim,
		Soyisim:      soyisim,
		Eposta:       eposta,
		Telefon:      telefon,
		TCKimlikNo:   tcKimlikNo,
		DogumTarihi:  dogumTarihi,
	}

	kullaniciJSON, err := json.Marshal(kullanici)
	if err != nil {
		return err
	}

	// Kullanıcı bilgisini ledger'a kaydet
	return ctx.GetStub().PutState(kullaniciID, kullaniciJSON)
}

func (s *SmartContract) UpdateUser(ctx contractapi.TransactionContextInterface, kullaniciID, yeniKullaniciAdi, yeniEposta, yeniTelefon, yeniSifre string) error {
	// Mevcut kullanıcı verilerini ledger'dan çek
	kullaniciData, err := ctx.GetStub().GetState(kullaniciID)
	if err != nil {
		return fmt.Errorf("kullanıcı bulunamadı: %s", err.Error())
	}
	if kullaniciData == nil {
		return fmt.Errorf("kullanıcı mevcut değil: %s", kullaniciID)
	}

	var kullanici Kullanici
	if err := json.Unmarshal(kullaniciData, &kullanici); err != nil {
		return fmt.Errorf("kullanıcı bilgileri okunamadı: %s", err.Error())
	}

	// Kullanıcı adı, e-posta, telefon numarası ve şifre güncelle
	if yeniKullaniciAdi != "" {
		kullanici.KullaniciAdi = yeniKullaniciAdi
	}
	if yeniEposta != "" {
		kullanici.Eposta = yeniEposta
	}
	if yeniTelefon != "" {
		kullanici.Telefon = yeniTelefon
	}
	if yeniSifre != "" {
		kullanici.Sifre = yeniSifre
	}

	// Güncellenmiş kullanıcı verilerini JSON'a çevir
	guncellenmisKullaniciData, err := json.Marshal(kullanici)
	if err != nil {
		return fmt.Errorf("kullanıcı bilgileri JSON'a dönüştürülemedi: %s", err.Error())
	}

	// Güncellenmiş kullanıcı verilerini ledger'a kaydet
	err = ctx.GetStub().PutState(kullaniciID, guncellenmisKullaniciData)
	if err != nil {
		return fmt.Errorf("kullanıcı bilgileri güncellenemedi: %s", err.Error())
	}

	return nil
}

func (s *SmartContract) QueryUser(ctx contractapi.TransactionContextInterface, kullaniciID string) (*Kullanici, error) {
	kullaniciJSON, err := ctx.GetStub().GetState(kullaniciID)
	if err != nil || kullaniciJSON == nil {
		return nil, fmt.Errorf("kullanıcı bulunamadı")
	}

	var kullanici Kullanici
	err = json.Unmarshal(kullaniciJSON, &kullanici)
	if err != nil {
		return nil, err
	}

	return &kullanici, nil
}

func (s *SmartContract) GetUserByUsername(ctx contractapi.TransactionContextInterface, kullaniciAdi string) (*Kullanici, error) {
	queryResult, err := ctx.GetStub().GetQueryResult(fmt.Sprintf(`{"selector":{"kullaniciAdi":"%s"}}`, kullaniciAdi))
	if err != nil {
		return nil, err
	}
	defer queryResult.Close()

	for queryResult.HasNext() {
		response, iterErr := queryResult.Next()
		if iterErr != nil {
			return nil, iterErr
		}

		var kullanici Kullanici
		jsonErr := json.Unmarshal(response.Value, &kullanici)
		if jsonErr != nil {
			return nil, jsonErr
		}

		// Kullanıcı adı benzersiz olduğu için, ilk eşleşen sonucu döndür.
		return &kullanici, nil
	}

	return nil, fmt.Errorf("kullanıcı adı ile eşleşen kullanıcı bulunamadı: %s", kullaniciAdi)
}

func (s *SmartContract) GetUserByTCKimlikNo(ctx contractapi.TransactionContextInterface, tcKimlikNo string) (*Kullanici, error) {
	queryString := fmt.Sprintf(`{"selector":{"tcKimlikNo":"%s"}}`, tcKimlikNo)
	queryResult, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer queryResult.Close()

	for queryResult.HasNext() {
		response, iterErr := queryResult.Next()
		if iterErr != nil {
			return nil, iterErr
		}

		var kullanici Kullanici
		jsonErr := json.Unmarshal(response.Value, &kullanici)
		if jsonErr != nil {
			return nil, jsonErr
		}

		// TC Kimlik No benzersiz olduğu için, ilk eşleşen sonucu döndür.
		return &kullanici, nil
	}

	return nil, fmt.Errorf("TC Kimlik No ile eşleşen kullanıcı bulunamadı: %s", tcKimlikNo)
}

func (s *SmartContract) ListUsers(ctx contractapi.TransactionContextInterface) ([]Kullanici, error) {
	startKey := ""
	endKey := ""
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var kullaniciListesi []Kullanici
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var kullanici Kullanici
		err = json.Unmarshal(queryResponse.Value, &kullanici)
		if err != nil {
			return nil, err
		}
		kullaniciListesi = append(kullaniciListesi, kullanici)
	}

	return kullaniciListesi, nil
}

func (s *SmartContract) ExtendIhale(ctx contractapi.TransactionContextInterface, ihaleNumarasi string, kullaniciID string, ekSure time.Duration) error {
	ihaleJSON, err := ctx.GetStub().GetState(ihaleNumarasi)
	if err != nil || ihaleJSON == nil {
		return fmt.Errorf("ihale bulunamadı")
	}

	var ihale Ihale
	err = json.Unmarshal(ihaleJSON, &ihale)
	if err != nil {
		return err
	}

	// Kullanıcı kimlik kontrolü
	if ihale.OlusturanKullaniciID != kullaniciID {
		return fmt.Errorf("sadece ihaleyi oluşturan kullanıcı bu işlemi yapabilir")
	}

	// İhalenin bitiş tarihini uzat
	ihale.IhaleBitisTarihi = ihale.IhaleBitisTarihi.Add(ekSure)
	updatedIhaleJSON, err := json.Marshal(ihale)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(ihaleNumarasi, updatedIhaleJSON)
}

func (s *SmartContract) CloseIhale(ctx contractapi.TransactionContextInterface, ihaleNumarasi string, kullaniciID string) error {
	ihaleJSON, err := ctx.GetStub().GetState(ihaleNumarasi)
	if err != nil || ihaleJSON == nil {
		return fmt.Errorf("ihale bulunamadı")
	}

	var ihale Ihale
	err = json.Unmarshal(ihaleJSON, &ihale)
	if err != nil {
		return err
	}

	// Kullanıcı kimlik kontrolü
	if ihale.OlusturanKullaniciID != kullaniciID {
		return fmt.Errorf("sadece ihaleyi oluşturan kullanıcı bu işlemi yapabilir")
	}

	// İhaleye yapılan tüm teklifleri sorgula
	queryIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Teklif", []string{ihaleNumarasi})
	if err != nil {
		return err
	}
	defer queryIterator.Close()

	var enYuksekTeklif Teklif
	var enYuksekFiyat float64 = 0
	for queryIterator.HasNext() {
		response, err := queryIterator.Next()
		if err != nil {
			return err
		}

		var teklif Teklif
		err = json.Unmarshal(response.Value, &teklif)
		if err != nil {
			return err
		}

		// En yüksek teklifi bul
		if teklif.Fiyat > enYuksekFiyat {
			enYuksekFiyat = teklif.Fiyat
			enYuksekTeklif = teklif
		}
	}

	// İhaleyi sonlandır ve kazanan teklifi belirle
	ihale.KazananTeklifID = enYuksekTeklif.KullaniciID
	ihale.Durum = IhaleDurumuKapali
	updatedIhaleJSON, err := json.Marshal(ihale)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(ihaleNumarasi, updatedIhaleJSON)
}

func (s *SmartContract) CheckAndCloseIhales(ctx contractapi.TransactionContextInterface) error {
	// Get the current timestamp from the transaction context
	stub := ctx.GetStub()
	txTimestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("timestamp alınırken hata oluştu: %s", err)
	}
	currentTime := time.Unix(txTimestamp.Seconds, int64(txTimestamp.Nanos))

	queryString := `{"selector":{"durum":"ACIK", "ihaleBitisTarihi":{"$lt": ` + strconv.FormatInt(currentTime.Unix(), 10) + `}}}` // Query only open auctions that are past the end time

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return fmt.Errorf("sorgulama sırasında hata oluştu: %s", err)
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return fmt.Errorf("sonraki öğeyi alırken hata oluştu: %s", err)
		}

		var ihale Ihale
		if err := json.Unmarshal(queryResponse.Value, &ihale); err != nil {
			return fmt.Errorf("ihale verisi çözümlenirken hata oluştu: %s", err)
		}

		// Retrieve all bids for the auction
		teklifIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Teklif", []string{ihale.IhaleNumarasi})
		if err != nil {
			return fmt.Errorf("teklifler alınırken hata oluştu: %s", err)
		}
		defer teklifIterator.Close()

		var highestBid Teklif
		highestBidPrice := float64(0)
		for teklifIterator.HasNext() {
			teklifResponse, err := teklifIterator.Next()
			if err != nil {
				return fmt.Errorf("teklif sırasında hata oluştu: %s", err)
			}

			var teklif Teklif
			if err := json.Unmarshal(teklifResponse.Value, &teklif); err != nil {
				return fmt.Errorf("teklif verisi çözümlenirken hata oluştu: %s", err)
			}

			if teklif.Fiyat > highestBidPrice {
				highestBidPrice = teklif.Fiyat
				highestBid = teklif
			}
		}

		// Update the auction with the winner and close it
		ihale.KazananTeklifID = highestBid.KullaniciID
		ihale.Durum = IhaleDurumuKapali
		ihaleUpdatedJSON, err := json.Marshal(ihale)
		if err != nil {
			return fmt.Errorf("ihale güncellenirken hata oluştu: %s", err)
		}

		// Save updated auction to the ledger
		if err := ctx.GetStub().PutState(ihale.IhaleNumarasi, ihaleUpdatedJSON); err != nil {
			return fmt.Errorf("ihale durumu kaydedilirken hata oluştu: %s", err)
		}
	}

	return nil
}

func (s *SmartContract) CloseExpiredIhaleler(ctx contractapi.TransactionContextInterface) error {
	stub := ctx.GetStub()
	currentTime := time.Now() // Get the current time to compare with auction end times

	// Query to find all open auctions that have passed their end time
	queryString := fmt.Sprintf(`{"selector":{"durum":"ACIK", "ihaleBitisTarihi":{"$lt":"%s"}}}`, currentTime.Format(time.RFC3339))
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return fmt.Errorf("failed to retrieve expired auctions: %s", err)
	}
	defer resultsIterator.Close()

	// Iterate over expired auctions to close them
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return fmt.Errorf("error iterating over expired auctions: %s", err)
		}

		var ihale Ihale
		if err := json.Unmarshal(queryResponse.Value, &ihale); err != nil {
			return fmt.Errorf("error unmarshalling auction data: %s", err)
		}

		// Determine the highest bid for the auction
		var highestBid *Teklif
		var highestBidPrice float64 = 0
		bidsIterator, err := stub.GetStateByPartialCompositeKey("Teklif", []string{ihale.IhaleNumarasi})
		if err != nil {
			return fmt.Errorf("error retrieving bids for auction %s: %s", ihale.IhaleNumarasi, err)
		}
		defer bidsIterator.Close()

		for bidsIterator.HasNext() {
			response, err := bidsIterator.Next()
			if err != nil {
				return fmt.Errorf("error iterating bids: %s", err)
			}

			var bid Teklif
			if err := json.Unmarshal(response.Value, &bid); err != nil {
				continue // Skip bad data
			}

			if bid.Fiyat > highestBidPrice {
				highestBidPrice = bid.Fiyat
				highestBid = &bid
			}
		}

		// Update the auction's status and, if a bid was found, set the winning bid ID
		ihale.Durum = IhaleDurumuKapali
		if highestBid != nil {
			ihale.KazananTeklifID = highestBid.KullaniciID
		}

		// Marshal the updated auction and save it back to the ledger
		ihaleJSON, err := json.Marshal(ihale)
		if err != nil {
			return fmt.Errorf("error marshalling updated auction data: %s", err)
		}
		err = stub.PutState(ihale.IhaleNumarasi, ihaleJSON)
		if err != nil {
			return fmt.Errorf("error updating auction status on the ledger: %s", err)
		}
	}

	return nil
}

func (s *SmartContract) GetIhaleStatistics(ctx contractapi.TransactionContextInterface) (*IhaleIstatistikleri, error) {
	startKey := ""
	endKey := ""
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var toplamIhaleSayisi int
	var toplamSatisEsasBedel float64
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var ihale Ihale
		err = json.Unmarshal(queryResponse.Value, &ihale)
		if err != nil {
			continue // JSON çözümlemesi başarısız olursa, bir sonraki ihaleye geç
		}

		toplamIhaleSayisi++
		toplamSatisEsasBedel += ihale.SatisaEsasBedel
	}

	ortalamaSatisEsasBedel := toplamSatisEsasBedel / float64(toplamIhaleSayisi)
	istatistikler := &IhaleIstatistikleri{
		ToplamIhaleSayisi:      toplamIhaleSayisi,
		OrtalamaSatisEsasBedel: ortalamaSatisEsasBedel,
	}

	return istatistikler, nil
}

func (s *SmartContract) GetUserActivity(ctx contractapi.TransactionContextInterface, kullaniciID string) (*KullaniciAktivitesi, error) {
	// Kullanıcı tarafından verilen tekliflerin sorgulanması
	queryIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Teklif", []string{kullaniciID})
	if err != nil {
		return nil, err
	}
	defer queryIterator.Close()

	var teklifSayisi int
	for queryIterator.HasNext() {
		_, err := queryIterator.Next()
		if err != nil {
			return nil, err
		}
		teklifSayisi++
	}

	aktivite := &KullaniciAktivitesi{
		KullaniciID:  kullaniciID,
		TeklifSayisi: teklifSayisi,
	}

	return aktivite, nil
}

// GetHistoryForIhale bir ihaleyle ilgili işlem geçmişini getirir.
func (s *SmartContract) GetHistoryForIhale(ctx contractapi.TransactionContextInterface, ihaleNumarasi string) ([]IhaleHistoryItem, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(ihaleNumarasi)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var history []IhaleHistoryItem
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var ihale Ihale
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &ihale)
			if err != nil {
				return nil, err
			}
		}

		historyItem := IhaleHistoryItem{
			TxID:      response.TxId,
			Value:     ihale,
			Timestamp: time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)), // Dönüşüm burada yapılıyor
			IsDeleted: response.IsDelete,
		}
		history = append(history, historyItem)
	}

	return history, nil
}

// GetHistoryForTeklif bir teklifle ilgili işlem geçmişini getirir.
func (s *SmartContract) GetHistoryForTeklif(ctx contractapi.TransactionContextInterface, teklifKey string) ([]TeklifHistoryItem, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(teklifKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var history []TeklifHistoryItem
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var teklif Teklif
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &teklif)
			if err != nil {
				return nil, err
			}
		}

		historyItem := TeklifHistoryItem{
			TxID:      response.TxId,
			Value:     teklif,
			Timestamp: time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)), // Dönüşüm burada yapılıyor
			IsDeleted: response.IsDelete,
		}
		history = append(history, historyItem)
	}

	return history, nil
}

func main() {
	config := ServerConfig{
		CCID:    os.Getenv("CHAINCODE_ID"),
		Address: os.Getenv("CHAINCODE_SERVER_ADDRESS"),
	}

	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error create simple chaincode: %s", err.Error())
		return
	}

	server := &shim.ChaincodeServer{
		CCID:    config.CCID,
		Address: config.Address,
		CC:      chaincode,
		TLSProps: shim.TLSProperties{
			Disabled: true,
		},
	}

	if err := server.Start(); err != nil {
		fmt.Printf("Error starting simple chaincode: %s", err.Error())
	}
}

type ServerConfig struct {
	CCID    string
	Address string
}
