package queries

import (
	"blockchain-smart-tender-platform/pkg/helpers"
	ch "blockchain-smart-tender-platform/platform/channel"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func Execute(chaincodeID, fcn string, args []string) (str string, err error) {
	response, err := ch.ChannelClient.Execute(
		channel.Request{
			ChaincodeID:     chaincodeID,
			Fcn:             fcn,
			Args:            helpers.ConvertArgs(args),
			InvocationChain: nil,
			IsInit:          false,
		},
	)
	if err != nil {
		return "", err
	}
	str = string(response.Payload)
	return str, nil
}
