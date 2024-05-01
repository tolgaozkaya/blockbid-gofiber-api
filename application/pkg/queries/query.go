package queries

import (
	"blockchain-smart-tender-platform/pkg/helpers"
	ch "blockchain-smart-tender-platform/platform/channel"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func Query(chaincodeID, fcn string, args []string) (resp []byte, err error) {
	response, err := ch.ChannelClient.Query(
		channel.Request{
			ChaincodeID:     chaincodeID,
			Fcn:             fcn,
			Args:            helpers.ConvertArgs(args),
			InvocationChain: nil,
			IsInit:          false,
		},
	)
	if err != nil {
		return nil, err
	}
	return response.Payload, nil
}
