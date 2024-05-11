// Ihaleler.go

package Ihaleler

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"blockchain-smart-tender-platform/app/entities"
	"blockchain-smart-tender-platform/internal/constants"
	"blockchain-smart-tender-platform/pkg/queries"

	"github.com/gofiber/fiber/v2"
)

// sendErrorResponse, hata mesajlarını tutarlı bir formatla göndermek için yardımcı bir fonksiyondur.
func sendErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"error": message})
}

func CreateIhale(c *fiber.Ctx) error {
	var request entities.CreateIhaleRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request parsing failed"})
	}

	// Kullanıcı bilgisini doğrula
	request.OlusturanKullaniciID = c.Locals("userID").(string)

	// Convert all fields to string and prepare as array
	ihaleParams := []string{
		request.IsletmeMudurlugu,
		request.AnaKategori,
		request.AltKategori,
		request.TasfiyeNedeni,
		request.BulunduguYer,
		fmt.Sprintf("%f", request.SatisaEsasBedel),
		fmt.Sprintf("%f", request.TeminatBedeli),
		fmt.Sprintf("%f", request.BaslangicBedeli),
		request.GorusBaslangicTarihi.Format(time.RFC3339),
		request.TeklifBaslangicTarihi.Format(time.RFC3339),
		request.IhaleBitisTarihi.Format(time.RFC3339),
		fmt.Sprintf("%f", request.DamgaVergisiOrani),
		request.Aciklama,
		request.OlusturanKullaniciID,
		fmt.Sprintf("%f", request.GuncelFiyat),
	}

	// Hyperledger Fabric chaincode sorgusu
	resp, err := queries.Execute(constants.ChaincodeID, constants.IHALE_OLUSTUR_FCN, ihaleParams)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to query Ihale: %v", err))
	}

	// Başarılı bir yanıt döndür
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Ihale successfully created", "response": resp})
}

func UpdateIhale(c *fiber.Ctx) error {
	var request entities.UpdateIhaleRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request parsing failed"})
	}

	request.IhaleNumarasi = c.Params("ihaleNumarasi")

	// Sadece belirli alanların güncellenmesine izin ver
	chaincodeArgs := []string{
		request.IhaleNumarasi,
		request.Aciklama,
		fmt.Sprintf("%f", request.GuncelFiyat),
		request.BulunduguYer,
		fmt.Sprintf("%f", request.TeminatBedeli),
		fmt.Sprintf("%f", request.SatisaEsasBedel),
		request.IhaleBitisTarihi.Format(time.RFC3339),
		fmt.Sprintf("%f", request.DamgaVergisiOrani),
	}

	// Chaincode fonksiyonunu çağır
	resp, err := queries.Execute(constants.ChaincodeID, constants.IHALE_GUNCELLE_FCN, chaincodeArgs)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to update Ihale: %v", err))
	}

	// Başarılı yanıt dön
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Ihale successfully updated", "response": resp})
}

func QueryIhale(c *fiber.Ctx) error {
	ihaleNumarasi := c.Params("ihaleNumarasi")

	// Blockchain chaincode sorgu fonksiyonunu çalıştır
	resp, err := queries.Query(constants.ChaincodeID, constants.IHALE_SORGULA_FCN, []string{ihaleNumarasi})
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to query Ihale: %v", err))
	}

	// Başarılı yanıt dön
	return c.JSON(fiber.Map{"message": "Ihale successfully retrieved", "response": json.RawMessage(resp)})
}

func QueryKullaniciIhaleleri(c *fiber.Ctx) error {
	kullaniciID := c.Locals("userID").(string)

	// Blockchain chaincode sorgu fonksiyonunu çalıştır
	resp, err := queries.Query(constants.ChaincodeID, constants.KULLANICI_IHALELERI_SORGULA_FCN, []string{kullaniciID})
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to query user's Ihaleler: %v", err))
	}

	// Başarılı yanıt dön
	return c.JSON(fiber.Map{"message": "User's Ihaleler successfully retrieved", "response": json.RawMessage(resp)})
}

func QueryDigerIhaleler(c *fiber.Ctx) error {
	kullaniciID := c.Locals("userID").(string)

	// Blockchain chaincode sorgu fonksiyonunu çalıştır
	resp, err := queries.Query(constants.ChaincodeID, constants.DIGER_IHALELER_SORGULA_FCN, []string{kullaniciID})
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to query other users' Ihaleler: %v", err))
	}

	// Başarılı yanıt dön
	return c.JSON(fiber.Map{"message": "Other users' Ihaleler successfully retrieved", "response": json.RawMessage(resp)})
}

func ListIhaleler(c *fiber.Ctx) error {
	// Execute chaincode function for listing all tenders
	resp, err := queries.Query(constants.ChaincodeID, constants.IHALE_LISTELE_FCN, []string{})
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to list Ihaleler: %v", err))
	}

	// Return successful response
	return c.JSON(fiber.Map{"message": "Ihaleler successfully listed", "response": json.RawMessage(resp)})
}

// StartIhale bir ihaleyi başlatır
func StartIhale(c *fiber.Ctx) error {
	ihaleNumarasi := c.Params("ihaleNumarasi")
	if ihaleNumarasi == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "İhale numarası gereklidir"})
	}

	// Chaincode sorgusu
	_, err := queries.Execute(constants.ChaincodeID, constants.IHALE_BASLAT_FCN, []string{ihaleNumarasi})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("İhale başlatılamadı: %v", err)})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "İhale başarıyla başlatıldı"})
}

// ExtendIhaleHandler bir ihalenin süresini uzatır
func ExtendIhale(c *fiber.Ctx) error {
	ihaleNumarasi := c.Params("ihaleNumarasi")
	var request struct {
		EkSure time.Duration `json:"ekSure"`
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "İstek ayrıştırılamadı"})
	}

	// Chaincode sorgusu
	_, err := queries.Execute(constants.ChaincodeID, constants.IHALE_SURE_EKLE_FCN, []string{ihaleNumarasi, fmt.Sprintf("%d", request.EkSure)})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("İhale süresi uzatılamadı: %v", err)})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "İhale süresi başarıyla uzatıldı"})
}

func CloseIhale(c *fiber.Ctx) error {
	// Retrieve userID from the local context set by AuthMiddleware
	userID := c.Locals("userID")
	if userID == nil || userID == "" {
		return sendErrorResponse(c, fiber.StatusUnauthorized, "User ID not found in request context")
	}

	// Extract ihaleNumarasi from the request parameters
	ihaleNumarasi := c.Params("ihaleNumarasi")
	if ihaleNumarasi == "" {
		return sendErrorResponse(c, fiber.StatusBadRequest, "IhaleNumarasi is required")
	}

	// Convert userID to string
	userIDStr, ok := userID.(string)
	if !ok {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Invalid userID format")
	}

	// Execute the chaincode function with the userID parameter
	response, err := queries.Execute(constants.ChaincodeID, constants.IHALE_SONLANDIR_FCN, []string{ihaleNumarasi, userIDStr})
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Teklifler başarıyla açıldı", "result": response})
}

func CheckAndCloseIhales(c *fiber.Ctx) error {
	// Call the chaincode function to check and close auctions
	response, err := queries.Execute(constants.ChaincodeID, "CloseExpiredIhaleler", []string{})
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to check and close auctions: %v", err))
	}

	// Return a successful response with the result from the chaincode
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Auctions checked and closed successfully",
		"result":  response,
	})
}

func GetHistoryForIhale(c *fiber.Ctx) error {
	IhaleId := c.Params("ihaleNumarasi")
	if IhaleId == "" {
		return sendErrorResponse(c, fiber.StatusBadRequest, "IhaleId is required")
	}

	resp, err := queries.Query(constants.ChaincodeID, "GetHistoryForIhale", []string{IhaleId})
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Başarılı yanıt dön
	return c.JSON(fiber.Map{"message": "Ihale history successfully retrieved", "response": json.RawMessage(resp)})
}

// HyperledgerInfo stores information about the Hyperledger Fabric network
type HyperledgerInfo struct {
	ChannelInformation     string `json:"channel_information"`
	PeerStatus             string `json:"peer_status"`
	OrdererNodes           string `json:"orderer_nodes"`
	CertificateAuthorities string `json:"certificate_authorities"`
}

func FetchNetworkInfo() (*HyperledgerInfo, error) {
	info := &HyperledgerInfo{}

	// Fetch channel information
	channelInfo, err := exec.Command("kubectl", "get", "FabricMainChannel", "demo", "-o", "yaml").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch channel information: %w", err)
	}
	info.ChannelInformation = string(channelInfo)

	// Fetch peer status
	peerStatus, err := exec.Command("kubectl", "get", "fabricpeers", "-n", "default").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch peer status: %w", err)
	}
	info.PeerStatus = string(peerStatus)

	// Fetch orderer nodes
	ordererNodes, err := exec.Command("kubectl", "get", "fabricorderernodes", "-n", "default").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orderer nodes: %w", err)
	}
	info.OrdererNodes = string(ordererNodes)

	// Fetch certificate authorities
	certificateAuthorities, err := exec.Command("kubectl", "get", "fabriccas", "-n", "default").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch certificate authorities: %w", err)
	}
	info.CertificateAuthorities = string(certificateAuthorities)

	return info, nil
}
