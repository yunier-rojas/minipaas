package entities

import "strings"

func ParseServiceTarget(target string) (service, port string) {
	service = target
	port = "80"
	if idx := strings.LastIndex(target, ":"); idx != -1 {
		service = target[:idx]
		port = target[idx+1:]
	}
	return
}
