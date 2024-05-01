package helpers

import "encoding/json"

func ConvertArgs(args []string) (arr [][]byte) {
	for _, element := range args {
		arr = append(arr, []byte(element))
	}
	return arr
}

func ConvertTransientMap(transientMap map[string]interface{}) (map[string][]byte, error) {
	arr := make(map[string][]byte)
	for key, element := range transientMap {
		str, err := json.Marshal(element)
		if err != nil {
			return nil, err
		}
		arr[key] = []byte(str)
	}
	return arr, nil
}
