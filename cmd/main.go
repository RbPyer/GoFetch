package main

import (
	"log"

	"github.com/RbPyer/Gofetch/internal/models"
	p "github.com/RbPyer/Gofetch/internal/parsers"
	"github.com/RbPyer/Gofetch/internal/presentator"
)

func main() {
	response := models.New()
	if err := p.GetUserInfo(&response); err != nil {
		log.Fatalf("Some errors with getting user: %s", err.Error())
	}
	if err := p.GetOsVersion(&response); err != nil {
		log.Fatalf("Some errors with getting OS version: %s", err.Error())
	}
	if err := p.GetKernelVersion(&response); err != nil {
		log.Fatalf("Some errors with getting kernel version: %s", err.Error())
	}
	if err := p.GetUptime(&response); err != nil {
		log.Fatalf("Some errors with getting system uptime: %s", err.Error())
	}
	if err := p.GetCPUInfo(&response); err != nil {
		log.Fatalf("Some errors with getting cpu info: %s", err.Error())
	}
	if err := p.GetRAMInfo(&response); err != nil {
		log.Fatalf("Some errors with getting ram info: %s", err.Error())
	}
	if err := p.GetDiskInfo(&response); err != nil {
		log.Fatalf("Some errors with getting disk info: %s", err.Error())
	}
	if err := p.GetGPUInfo(&response); err != nil {
		log.Fatalf("Some errors with getting disk info: %s", err.Error())
	}
	p.GetShell(&response)

	presentator.Present(&response)
}
