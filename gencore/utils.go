package gencore

import "strings"

func ReceiverName(typeName string) string {
	return strings.ToLower(string(typeName[0]))
}
