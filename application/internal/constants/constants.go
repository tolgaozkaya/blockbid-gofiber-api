package constants

const (
	// Ihale (Tender) related operations
	IHALE_OLUSTUR_FCN               = "CreateIhale2"
	IHALE_GUNCELLE_FCN              = "UpdateIhale"
	IHALE_SORGULA_FCN               = "QueryTumIhaleler"
	IHALE_LISTELE_FCN               = "ListIhaleler"
	KULLANICI_IHALELERI_SORGULA_FCN = "QueryKullaniciIhaleleri"
	DIGER_IHALELER_SORGULA_FCN      = "QueryDigerIhaleler"

	IHALE_BASLAT_FCN    = "StartIhale"
	IHALE_SURE_EKLE_FCN = "ExtendIhale"
	IHALE_SONLANDIR_FCN = "CloseIhale"

	// Teklif (Bid) related operations
	TEKLIF_YAP_FCN        = "MakeTeklif"
	TEKLIF_SORGULA_FCN    = "QueryTeklif"
	TUMTEKLIF_SORGULA_FCN = "ListTekliflerByIhale"

	// Kullanici (User) related operations
	KULLANICI_OLUSTUR_FCN  = "RegisterUser"
	KULLANICI_LISTELE_FCN  = "GetUserByUsername"
	KULLANICI_GUNCELLE_FCN = "UpdateUser"
	KULLANICI_SIL_FCN      = "DeleteUser"
)

var (
	ChaincodeID string
)
