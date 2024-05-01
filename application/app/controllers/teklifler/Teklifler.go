package Teklifler

import (
	"encoding/json"
	"fmt"

	"blockchain-smart-tender-platform/app/entities"
	"blockchain-smart-tender-platform/internal/constants" // constants modülünüzün yolu
	"blockchain-smart-tender-platform/pkg/queries"        // Bu, queries modülünüzün yolu olmalıdır.

	"github.com/gofiber/fiber/v2"
)

// Helper function to send error response.
func sendErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"error": message})
}

// GetUserFromContext retrieves the user from the request context.
func GetUserFromContext(c *fiber.Ctx) (string, error) {
	istekliKullanici, ok := c.Locals("userID").(string)
	if !ok {
		return "", fmt.Errorf("user not found")
	}
	return istekliKullanici, nil
}

// MessageResponse struct for success and error responses
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse struct for error responses
type ErrorResponse struct {
	Error string `json:"error"`
}

// MakeBidHandler, bir ihaleye teklif vermek için bir REST API endpoint'i sağlar
func MakeTeklif(c *fiber.Ctx) error {
	var request entities.TeklifYapRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request parsing failed"})
	}

	IhaleId := c.Params("ihaleId")
	if IhaleId == "" {
		return sendErrorResponse(c, fiber.StatusBadRequest, "IhaleId is required")
	}

	istekliKullanici, err := GetUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "user not found")
	}

	request.KullaniciID = istekliKullanici
	request.IhaleNumarasi = IhaleId

	// Teklif verme parametrelerini ayrı ayrı string olarak diziye dönüştür ve gönder
	resp, err := queries.Execute(constants.ChaincodeID, constants.TEKLIF_YAP_FCN, []string{
		request.IhaleNumarasi,
		request.KullaniciID,
		fmt.Sprintf("%f", request.TeklifTutari),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to make bid: %v", err)})
	}

	// Başarılı bir yanıt döndür
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Bid successfully made", "response": resp})
}

func QueryTeklif(c *fiber.Ctx) error {
	ihaleId := c.Params("ihaleId")
	if ihaleId == "" {
		return sendErrorResponse(c, fiber.StatusBadRequest, "IhaleId is required")
	}

	istekliKullanici, err := GetUserFromContext(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "User not found")
	}

	response, err := queries.Query(constants.ChaincodeID, constants.TEKLIF_SORGULA_FCN, []string{ihaleId, istekliKullanici})
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to query bid: %v", err))
	}

	return c.JSON(fiber.Map{"message": "Bid successfully retrieved", "response": json.RawMessage(response)})
}

func QueryAllTeklif(c *fiber.Ctx) error {
	ihaleId := c.Params("ihaleId")
	if ihaleId == "" {
		return sendErrorResponse(c, fiber.StatusBadRequest, "IhaleId is required")
	}

	response, err := queries.Query(constants.ChaincodeID, constants.TUMTEKLIF_SORGULA_FCN, []string{ihaleId})
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to query all bids for IhaleId %s: %v", ihaleId, err))
	}

	return c.JSON(fiber.Map{"message": "All bids successfully retrieved", "response": json.RawMessage(response)})
}
