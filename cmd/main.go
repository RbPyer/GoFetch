package main

import (
	"fmt"
	"log"
	"strings"
	"github.com/RbPyer/Gofetch/internal/parser"
	"github.com/RbPyer/Gofetch/internal/entities"
)

func main() {
	p := parser.NewParser()
	response := entities.NewResponse()
	if err := p.GetUserInfo(response); err != nil {
		log.Fatalf("Some errors with getting user: %s", err.Error())
	}
	if err := p.GetOsVersion(response); err != nil {
		log.Fatalf("Some errors with getting OS version: %s", err.Error())
	}
	if err := p.GetKernelVersion(response); err != nil {
		log.Fatalf("Some errors with getting kernel version: %s", err.Error())
	}
	if err := p.GetUptime(response); err != nil {
		log.Fatalf("Some errors with getting system uptime: %s", err.Error())
	}

	fmt.Println(strings.Join(response.Info, entities.Delimeter))
}