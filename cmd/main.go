package main

import (
	"log"

	"github.com/RbPyer/Gofetch/internal/models"
	p "github.com/RbPyer/Gofetch/internal/parsers"
	"sync"
	"github.com/RbPyer/Gofetch/internal/presentator"
)

func main() {
	response := &models.Response{}
	wg := new(sync.WaitGroup)
	wg.Add(5)

	go func() {
		if err := p.GetDiskInfo(response); err != nil {
			log.Fatalf("Some errors with getting disk info: %s", err.Error())
		}
		wg.Done()
	}()

	go func() {
		if err := p.GetGPUInfo(response); err != nil {
			log.Fatalf("Some errors with getting disk info: %s", err.Error())
		}
		wg.Done()
	}()

	go func() {
		if err := p.GetCPUInfo(response); err != nil {
			log.Fatalf("Some errors with getting cpu info: %s", err.Error())
		}
		wg.Done()
	}()

	go func() {
		if err := p.GetRAMInfo(response); err != nil {
			log.Fatalf("Some errors with getting ram info: %s", err.Error())
		}
		wg.Done()
	}()

	go func() {
		if err := p.GetOsVersion(response); err != nil {
			log.Fatalf("Some errors with getting OS version: %s", err.Error())
		}
		wg.Done()
	}()

	if err := p.GetUserInfo(response); err != nil {
		log.Fatalf("Some errors with getting user: %s", err.Error())
	}

	if err := p.GetKernelVersion(response); err != nil {
		log.Fatalf("Some errors with getting kernel version: %s", err.Error())
	}


	if err := p.GetUptime(response); err != nil {
		log.Fatalf("Some errors with getting system uptime: %s", err.Error())
	}

	p.GetShell(response)
	wg.Wait()

	presentator.Present(response)
}
