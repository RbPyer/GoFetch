package parser

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"strings"
	"syscall"
	"strconv"

	"github.com/RbPyer/Gofetch/internal/entities"
	"github.com/RbPyer/Gofetch/internal/utils"
)


type Parser struct {
	sysinfo syscall.Sysinfo_t
}

func NewParser(sysinfo *syscall.Sysinfo_t) *Parser {
	syscall.Sysinfo(sysinfo)
	return &Parser{
		sysinfo: *sysinfo,
	}
}


func (p *Parser) GetUserInfo(r *entities.Response) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	userObject, err := user.Current()
	if err != nil {
		return err
	}

	r.Info = append(r.Info, fmt.Sprintf("%s@%s", userObject.Username, hostname))
	return nil
}


func (p *Parser) GetOsVersion(r *entities.Response) error {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return err
	}
	defer file.Close()
	
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		str := fileScanner.Text()
		if strings.HasPrefix(str, entities.Prefix) {
			newStr := strings.ReplaceAll(strings.ReplaceAll(str, entities.Prefix, ""), "\"", "")
			r.Info = append(r.Info, fmt.Sprintf("OS: %-66s|", newStr))
			return nil
		} 
	}

	r.Info = append(r.Info, fmt.Sprintf("OS: %-66s|", "OS: no information about your os"))
	return nil
}


func (p *Parser) GetKernelVersion(r *entities.Response) error {
	u := syscall.Utsname{}
	err := syscall.Uname(&u)
	if err != nil {
		return err
	}
	r.Info = append(r.Info, fmt.Sprintf("Kernel: %-127s|", utils.UtsnameToStr(u.Release)))
	return nil
}


func (p *Parser) GetUptime(r *entities.Response) error {
	r.Info = append(r.Info, fmt.Sprintf("Uptime: %-62s|", utils.Int64ToTimeStr(p.sysinfo.Uptime)))
	return nil
}

func (p *Parser) GetRAMInfo(r *entities.Response) error {
	file, err := os.Open(entities.RAM_PATH)
	if err != nil {
		return err
	}
	defer file.Close()
	
	fileScanner := bufio.NewScanner(file)
	mem := entities.RAM{}
	
	for ;mem.SReclaimable == 0;fileScanner.Scan() {
		str := fileScanner.Text()
		switch {
		case strings.HasPrefix(str, "MemTotal"):
			if mem.Total, err = strconv.ParseUint(strings.Fields(str)[1], 10, 64); err != nil {return err}
		case strings.HasPrefix(str, "MemFree"):
			if mem.Free, err = strconv.ParseUint(strings.Fields(str)[1], 10, 64); err != nil {return err}
		case strings.HasPrefix(str, "Buffers"):
			if mem.Buffers, err = strconv.ParseUint(strings.Fields(str)[1], 10, 64); err != nil {return err}
		case strings.HasPrefix(str, "Cached"):
			if mem.Cached, err = strconv.ParseUint(strings.Fields(str)[1], 10, 64); err != nil {return err}
		case strings.HasPrefix(str, "Shmem"):
			if mem.Shared, err = strconv.ParseUint(strings.Fields(str)[1], 10, 64); err != nil {return err}
		case strings.HasPrefix(str, "SReclaimable"):
			if mem.SReclaimable, err = strconv.ParseUint(strings.Fields(str)[1], 10, 64); err != nil {return err}
		}
	}
	r.Info = append(r.Info, fmt.Sprintf("RAM: %-65s|", fmt.Sprintf("%d/%d MiB", (mem.Total + mem.Shared - mem.Buffers - mem.Cached - mem.Free - mem.SReclaimable) / 1024, 
			mem.Total / 1024)))

	return nil
}


func (p *Parser) GetCPUInfo(r *entities.Response) error {
	file, err := os.Open(entities.CPU_PATH)
	if err != nil {
		return err
	}
	defer file.Close()
	
	fileScanner := bufio.NewScanner(file)
	cpu := entities.CPU{}

	for ;cpu.Cores == 0;fileScanner.Scan() {
		str := fileScanner.Text()
		if strings.HasPrefix(str, "model name") {
			cpu.ModelName = strings.Replace(str, "model name\t: ", "", 1)
		} else if strings.HasPrefix(str, "siblings") {
			cpu.Siblings, err = strconv.Atoi(strings.Replace(str, "siblings\t: ", "", 1))
			if err != nil {
				return err
			}
		} else if strings.HasPrefix(str, "cpu cores") {
			cpu.Cores, err = strconv.Atoi(strings.Replace(str, "cpu cores\t: ", "", 1))
			if err != nil {
				return err
			}
		}
	}

	r.Info = append(r.Info, fmt.Sprintf("CPU: %-65s|", fmt.Sprintf("%s; %d cores / %d threads", cpu.ModelName, cpu.Cores, cpu.Siblings)))

	return nil
}