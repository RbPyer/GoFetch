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
	defer file.Close()

	if err != nil {
		return err
	}
	
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

// func (p *Parser) GetRAMInfo(r *entities.Response) {
// 	freeRam := p.sysinfo.Freeram / 1024 / 1024
// 	allRam := p.sysinfo.Totalram / 1024 / 1024
// 	bufRam := p.sysinfo.Bufferram / 1024 / 1024

// 	fmt.Println(freeRam, allRam, bufRam)


// 	// fmt.Printf("%d/%d", allRam - bufRam - freeRam - shrRam, allRam)

// }


func (p *Parser) GetCPUInfo(r *entities.Response) error {
	file, err := os.Open(entities.CPU_PATH)
	defer file.Close()
	if err != nil {
		return err
	}
	
	fileScanner := bufio.NewScanner(file)
	cpu := entities.CPU{}

	for ;cpu.Cores == 0;fileScanner.Scan() {
		str := fileScanner.Text()
		if strings.Contains(str, "model name") {
			cpu.ModelName = strings.Replace(str, "model name\t: ", "", 1)
		} else if strings.Contains(str, "siblings") {
			cpu.Siblings, err = strconv.Atoi(strings.Replace(str, "siblings\t: ", "", 1))
			if err != nil {
				return err
			}
		} else if strings.Contains(str, "cpu cores") {
			cpu.Cores, err = strconv.Atoi(strings.Replace(str, "cpu cores\t: ", "", 1))
			if err != nil {
				return err
			}
		}
	}

	r.Info = append(r.Info, fmt.Sprintf("CPU: %-65s|", fmt.Sprintf("%s; %d cores / %d threads", cpu.ModelName, cpu.Cores, cpu.Siblings)))

	return nil
}