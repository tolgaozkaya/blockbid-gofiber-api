package Kullanicilar

import (
	"blockchain-smart-tender-platform/app/entities"
	"blockchain-smart-tender-platform/internal/constants"
	"blockchain-smart-tender-platform/pkg/queries"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// MessageResponse struct for success and error responses
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse struct for error responses
type ErrorResponse struct {
	Error string `json:"error"`
}

// Define the request structure
type TCKimlikNoDogrulaRequest struct {
	XMLName   xml.Name `xml:"soap:Envelope"`
	XMLNsSoap string   `xml:"xmlns:soap,attr"`
	XMLNsXsi  string   `xml:"xmlns:xsi,attr"`
	XMLNsXsd  string   `xml:"xmlns:xsd,attr"`
	Body      TCKimlikNoDogrulaRequestBody
}

type TCKimlikNoDogrulaRequestBody struct {
	XMLName xml.Name `xml:"soap:Body"`
	Request TCKimlikNoDogrula
}

type TCKimlikNoDogrula struct {
	XMLName    xml.Name `xml:"http://tckimlik.nvi.gov.tr/WS TCKimlikNoDogrula"`
	TCKimlikNo int64    `xml:"TCKimlikNo"`
	Ad         string   `xml:"Ad"`
	Soyad      string   `xml:"Soyad"`
	DogumYili  int      `xml:"DogumYili"`
}

// Define the response structure
type TCKimlikNoDogrulaResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		Response struct {
			Result bool `xml:"TCKimlikNoDogrulaResult"`
		} `xml:"TCKimlikNoDogrulaResponse"`
	}
}

func ValidateTCKimlikNo(tcKimlikNo int64, ad, soyad string, dogumYili int) (bool, error) {
	// Construct the request payload
	requestPayload := TCKimlikNoDogrulaRequest{
		XMLNsSoap: "http://schemas.xmlsoap.org/soap/envelope/",
		XMLNsXsi:  "http://www.w3.org/2001/XMLSchema-instance",
		XMLNsXsd:  "http://www.w3.org/2001/XMLSchema",
		Body: TCKimlikNoDogrulaRequestBody{
			Request: TCKimlikNoDogrula{
				TCKimlikNo: tcKimlikNo,
				Ad:         ad,
				Soyad:      soyad,
				DogumYili:  dogumYili,
			},
		},
	}

	// Marshal the request into XML
	xmlRequest, err := xml.Marshal(requestPayload)
	if err != nil {
		return false, fmt.Errorf("error marshalling request: %v", err)
	}

	// Send the request
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://tckimlik.nvi.gov.tr/Service/KPSPublic.asmx", bytes.NewReader(xmlRequest))
	if err != nil {
		return false, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("SOAPAction", "http://tckimlik.nvi.gov.tr/WS/TCKimlikNoDogrula")

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error reading response: %v", err)
	}

	var response TCKimlikNoDogrulaResponse
	err = xml.Unmarshal(body, &response)
	if err != nil {
		return false, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return true, nil
}

// SignUpUser godoc
// @Summary Kullanıcı Kaydı
// @Description Yeni bir kullanıcı kaydeder.
// @Tags kullanicilar
// @Accept json
// @Produce json
// @Param signUpRequest body entities.CreateUserRequest true "Kayıt Bilgileri"
// @Success 201 {object} MessageResponse "Kullanıcı başarıyla kaydedildi"
// @Failure 400 {object} ErrorResponse "İstek ayrıştırması başarısız"
// @Failure 500 {object} ErrorResponse "Ledger üzerinde kullanıcı kaydedilemedi"
// @Router /signup [post]
func SignUpUser(c *fiber.Ctx) error {
	var signUpRequest entities.CreateUserRequest
	if err := c.BodyParser(&signUpRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request body: " + err.Error()})
	}

	// TC Kimlik No'nun int64'e dönüştürülmesi
	tckInt64, err := strconv.ParseInt(signUpRequest.TCKimlikNo, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "TCKimlikNo format error: " + err.Error()})
	}

	// TC Kimlik No doğrulaması
	isValid, err := ValidateTCKimlikNo(tckInt64, signUpRequest.Isim, signUpRequest.Soyisim, signUpRequest.DogumTarihi.Year())
	if err != nil {
		log.Printf("Validation error for TCKimlikNo %d: %v", tckInt64, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error validating TCKimlikNo: " + err.Error()})
	}
	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid TCKimlikNo" + err.Error()})
	}

	// Şifrenin hash'lenmesi
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signUpRequest.Sifre), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password: " + err.Error()})
	}
	signUpRequest.Sifre = string(hashedPassword)

	// Doğum tarihinin string'e dönüştürülmesi
	dogumTarihiISO := signUpRequest.DogumTarihi.Format(time.RFC3339) // ISO 8601 format

	// Zincir kodu fonksiyonunun argümanlarının hazırlanması
	args := []string{
		signUpRequest.KullaniciAdi,
		signUpRequest.Sifre,
		signUpRequest.Isim,
		signUpRequest.Soyisim,
		signUpRequest.Eposta,
		signUpRequest.Telefon,
		signUpRequest.TCKimlikNo,
		dogumTarihiISO,
	}

	// Yeni kullanıcının kaydedilmesi için zincir kodu fonksiyonunun çağrılması
	_, err = queries.Execute(constants.ChaincodeID, constants.KULLANICI_OLUSTUR_FCN, args)
	if err != nil {
		log.Printf("Failed to execute chaincode function: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user on ledger: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User successfully registered"})
}

type LoginRequest struct {
	KullaniciAdi string `json:"kullaniciAdi" validate:"required"`
	Sifre        string `json:"sifre" validate:"required"`
}

// LoginUser godoc
// @Summary User login
// @Description Logs in a user with a username and password.
// @Tags kullanicilar
// @Accept json
// @Produce json
// @Param loginRequest body LoginRequest true "Login Request"
// @Success 200 {string} MessageResponse "login successful"
// @Failure 400 {object} ErrorResponse "Request parsing failed"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /login [post]
func LoginUser(c *fiber.Ctx) error {
	var loginRequest LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request parsing failed"})
	}

	// Execute chaincode function
	response, err := queries.Query(constants.ChaincodeID, "GetUserByUsername", []string{loginRequest.KullaniciAdi})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if the response contains user data
	var user entities.User
	err = json.Unmarshal(response, &user)
	if err != nil || user.KullaniciAdi == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Sifre), []byte(loginRequest.Sifre))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// User authentication successful, generate JWT token
	token, err := GenerateJWTToken(user.KullaniciAdi, user.KullaniciID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating JWT token"})
	}

	// Return the token
	return c.JSON(fiber.Map{"token": token})
}

var jwtSecretKey = "gazi"

func GenerateJWTToken(username string, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"userID":   userID,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(jwtSecretKey))
}

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	// Extract the token from the Authorization header
	// Assuming the header follows the format "Bearer <token>"
	splits := strings.Split(authHeader, " ")
	if len(splits) != 2 || splits[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Authorization header format"})
	}
	tokenString := splits[1]

	// Parse and validate the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Confirm that the token's signing method matches the expected algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token: " + err.Error()})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Store user information from the token into local variables for later use
		c.Locals("username", claims["username"])
		c.Locals("userID", claims["userID"])
		return c.Next()
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}
}

// UpdateUser updates a user's information in the system.
// @Summary Update user information
// @Description Updates the information for a user by their userID.
// @Tags users
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Param user body entities.UpdateUserRequest true "User Update Request"
// @Success 200 {object} map[string]interface{} "message: User successfully updated"
// @Failure 400 {object} map[string]interface{} "error: Request parsing failed: [error detail]"
// @Failure 500 {object} map[string]interface{} "error: Blockchain query failed: [error detail]"
// @Router /update [put]
func UpdateUser(c *fiber.Ctx) error {
	// Retrieve userID from the local context set by previous middleware
	userID := c.Locals("userID")
	if userID == nil || userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request: userID cannot be empty"})
	}

	// Since c.Locals() returns an interface{}, you need to assert the type to string
	userIDStr, ok := userID.(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid userID format"})
	}

	var user entities.UpdateUserRequest
	if err := c.BodyParser(&user); err != nil {
		// Handle error if parsing the request body fails
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request parsing failed: " + err.Error()})
	}

	params := []string{
		userIDStr,
		user.KullaniciAdi,
		user.Sifre,
		user.Isim,
		user.Soyisim,
		user.Eposta,
		user.Telefon,
		user.DogumTarihi.Format("2006-01-02T00:00:00Z"), // Assuming DogumTarihi is time.Time and formatted as "YYYY-MM-DD"
	}

	// Call the chaincode function with individual parameters
	_, err := queries.Query(constants.ChaincodeID, "UpdateUser", params)
	if err != nil {
		// Handle error if the query to the blockchain fails
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Blockchain query failed: " + err.Error()})
	}

	// If everything is successful, respond with a success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User successfully updated"})
}

// DeleteUser deletes a user from the system.
// @Summary Delete a user
// @Description Deletes a user by their unique userID.
// @Tags users
// @Accept json
// @Produce json
// @Param userID path string true "User ID" "The ID of the user to be deleted"
// @Success 200 {object} map[string]interface{} "message: Kullanıcı başarıyla silindi" "User successfully deleted"
// @Failure 500 {object} map[string]interface{} "error: [error detail]" "Error occurred during user deletion"
// @Router /delete [delete]
func DeleteUser(c *fiber.Ctx) error {
	// Retrieve userID from the local context set by previous middleware
	userID := c.Locals("userID")
	if userID == nil || userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request: userID cannot be empty"})
	}

	// Since c.Locals() returns an interface{}, you need to assert the type
	userIDStr, ok := userID.(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid userID format"})
	}

	_, err := queries.Execute(constants.ChaincodeID, "DeleteUser", []string{userIDStr})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user: " + err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Kullanıcı başarıyla silindi"})
}

// GetUserByID godoc
// @Summary Retrieve user by ID
// @Description Retrieves user information from the local context and queries the database or blockchain to return user details.
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} User "Successfully retrieved user information"
// @Failure 400 {object} map[string]interface{} "error: User ID not found in request context"
// @Failure 500 {object} map[string]interface{} "error: Failed to retrieve user or error parsing user data"
// @Router /get [get]
func GetUserByID(c *fiber.Ctx) error {
	// Assuming the user ID is stored in the local context under the key "userID"
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID not found in request context"})
	}

	// Now, use the userID to perform your operations, such as querying the database or blockchain
	response, err := queries.Query(constants.ChaincodeID, "QueryUser", []string{userID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve user: " + err.Error()})
	}

	var user entities.User
	err = json.Unmarshal(response, &user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing user data"})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func GetUserDashboardData(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID not found in request context"})
	}

	// Query the chaincode function to get the user's dashboard data
	response, err := queries.Query(constants.ChaincodeID, "GetUserDashboardData", []string{userID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve user dashboard data: " + err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Kullanıcı girdilerinin güvenli bir şekilde işlenmesini sağlamak için basit bir karakter filtresi uygulayın.
func ValidateString(input string) bool {
	// Sadece alfanümerik karakterler, tire ve alt çizgiye izin verin.
	// Bu, çoğu Docker imaj adı ve etiket için yeterli olmalıdır.
	matched, err := regexp.MatchString("^[a-zA-Z0-9_-]+$", input)
	if err != nil || !matched {
		log.Printf("Güvensiz girdi tespit edildi: %s", input)
		return false
	}
	return true
}

// CreateDockerImage, Docker imajını oluşturan ve Docker Hub'a yükleyen bir fonksiyondur.
func CreateDockerImage(c *fiber.Ctx) error {
	dockerKullanıcıAdı := "200420"                                                       // Docker Hub kullanıcı adınız
	dockerImajAdı := "blockbid"                                                          // Docker imaj adınız
	dockerEtiket := "v7"                                                                 // İmajınızın etiketi
	dosyaYolu := "/Users/tolgaozkaya/fabric/blockchain-smart-tender-platform/chaincode/" // Dockerfile'ın dizini

	// Kullanıcı girdilerini doğrula
	if !ValidateString(dockerKullanıcıAdı) || !ValidateString(dockerImajAdı) || !ValidateString(dockerEtiket) {
		return fmt.Errorf("Güvensiz girdi tespit edildi")
	}

	// Docker build komutunu güvenli bir şekilde oluştur
	dockerBuildKomut := fmt.Sprintf("docker build -t %s/%s:%s .", dockerKullanıcıAdı, dockerImajAdı, dockerEtiket)
	buildCmd := exec.Command("bash", "-c", dockerBuildKomut)
	buildCmd.Dir = dosyaYolu // Dockerfile'ın bulunduğu dizin

	// Docker build komutunu çalıştır
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("Docker imajı oluşturulamadı: %w", err)
	}

	// Docker push komutunu güvenli bir şekilde oluştur
	dockerPushKomut := fmt.Sprintf("docker push %s/%s:%s", dockerKullanıcıAdı, dockerImajAdı, dockerEtiket)
	pushCmd := exec.Command("bash", "-c", dockerPushKomut)

	// Docker push komutunu çalıştır
	if err := pushCmd.Run(); err != nil {
		return fmt.Errorf("Docker imajı Docker Hub'a yüklenemedi: %w", err)
	}

	log.Printf("Docker imajı başarıyla oluşturuldu ve Docker Hub'a yüklendi: %s/%s:%s\n", dockerKullanıcıAdı, dockerImajAdı, dockerEtiket)
	return nil
}

func QueryUser(c *fiber.Ctx) error {
	userID := c.Params("userID")
	if userID == "" {
		return sendErrorResponse(c, fiber.StatusBadRequest, "User ID is required")
	}

	response, err := queries.Query(constants.ChaincodeID, "QueryUser", []string{userID})
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	var user entities.User
	err = json.Unmarshal(response, &user)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Error parsing user data")
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// sendErrorResponse, hata mesajlarını tutarlı bir formatla göndermek için yardımcı bir fonksiyondur.
func sendErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"error": message})
}
