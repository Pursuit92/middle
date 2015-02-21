package middle

import "strings"

func splitAddr(addr string) (string, string) {
	splitAddr := strings.Split(addr, ":")
	host := strings.Join(splitAddr[:len(splitAddr)-1], ":")
	port := splitAddr[len(splitAddr)-1]
	return host, port
}
