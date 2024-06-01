package parser

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"strings"
	"syscall"
	"github.com/RbPyer/Gofetch/internal/entities"
	"github.com/RbPyer/Gofetch/internal/utils"
)


type Parser struct {
	
}

func NewParser() *Parser {
	return &Parser{}
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
			r.Info = append(r.Info, fmt.Sprintf("OS: %-56s|", newStr))
			return nil
		} 
	}

	r.Info = append(r.Info, fmt.Sprintf("OS: %-56s|", "OS: no information about your os"))
	return nil
}


func (p *Parser) GetKernelVersion(r *entities.Response) error {
	u := syscall.Utsname{}
	err := syscall.Uname(&u)
	if err != nil {
		return err
	}
	r.Info = append(r.Info, fmt.Sprintf("Kernel: %-117s|", utils.UtsnameToStr(u.Release)))
	return nil
}


func (p *Parser) GetUptime(r *entities.Response) error {
	si := syscall.Sysinfo_t{}
	syscall.Sysinfo(&si)
	r.Info = append(r.Info, fmt.Sprintf("Uptime: %-52s|", utils.Int64ToTimeStr(si.Uptime)))
	return nil
}