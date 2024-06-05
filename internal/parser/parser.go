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


type Parser struct {}


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
	s := syscall.Sysinfo_t{}
	syscall.Sysinfo(&s)
	r.Info = append(r.Info, fmt.Sprintf("Uptime: %-62s|", utils.Int64ToTimeStr(s.Uptime)))
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
			mem.TrueFree = mem.Total + mem.Shared - mem.Buffers - mem.Cached - mem.Free - mem.SReclaimable
		}
	}
	r.Info = append(r.Info, fmt.Sprintf("RAM: %-65s|", fmt.Sprintf("%d/%d MiB [%.2f%%]", (mem.Total + mem.Shared - mem.Buffers - mem.Cached - mem.Free - mem.SReclaimable) / 1024, 
			mem.Total / 1024, float64(mem.TrueFree) / float64(mem.Total) * 100)))

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
		switch {
		case strings.HasPrefix(str, "model name"):
			cpu.ModelName = strings.Replace(str, "model name\t: ", "", 1)
		case strings.HasPrefix(str, "siblings"):
			if cpu.Siblings, err = strconv.Atoi(strings.Replace(str, "siblings\t: ", "", 1)); err != nil {return err}
		case strings.HasPrefix(str, "cpu cores"):
			if cpu.Cores, err = strconv.Atoi(strings.Replace(str, "cpu cores\t: ", "", 1)); err != nil {return err}
		}
	}

	r.Info = append(r.Info, fmt.Sprintf("CPU: %-65s|", fmt.Sprintf("%s; %d cores / %d threads", cpu.ModelName, cpu.Cores, cpu.Siblings)))

	return nil
}


func (p *Parser) GetDiskInfo(r *entities.Response) error {
	fs := syscall.Statfs_t{}
	if err := syscall.Statfs("/", &fs); err != nil {
		return err
	}
	diskInfo := entities.DiskInfo{
		All: fs.Blocks * uint64(fs.Bsize),
	}
	diskInfo.Used = diskInfo.All - fs.Bfree * uint64(fs.Bsize)
	r.Info = append(r.Info, fmt.Sprintf("Disk Info: %-59s|", fmt.Sprintf("%.2f/%.2f GiB", float64(diskInfo.Used) / entities.GB, float64(diskInfo.All) / entities.GB)))

	return nil
}