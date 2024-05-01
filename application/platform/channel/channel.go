package channel

import (
	"log"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	_ "github.com/joho/godotenv/autoload"
)

var ChannelClient *channel.Client

func init() {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	// Configs
	cfg := config.FromFile(os.Getenv("CONFIG_PATH"))
	sdk, err := fabsdk.New(cfg)
	if err != nil {
		log.Fatalf("Failed init fabsdk: %v", err)
	}
	chContext := sdk.ChannelContext(
		os.Getenv("CHANNEL_NAME"),
		fabsdk.WithUser(os.Getenv("CHANNEL_USER")),
		fabsdk.WithOrg(os.Getenv("CHANNEL_ORG")),
	)
	ChannelClient, err = channel.New(chContext)
	if err != nil {
		log.Fatalf("Failed create channel: %v", err)
	}
}
