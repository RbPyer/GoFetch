package parser

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/RbPyer/Gofetch/internal/entities"
	"github.com/RbPyer/Gofetch/internal/utils"
)

type Parser struct{}

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
			str = strings.ReplaceAll(strings.ReplaceAll(str, entities.Prefix, ""), "\"", "")
			r.Info = append(r.Info, fmt.Sprintf("OS: %-66s|", str))
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
	err := syscall.Sysinfo(&s)
	if err != nil {
		return err
	}
	r.Info = append(r.Info, fmt.Sprintf("Uptime: %-62s|", utils.Int64ToTimeStr(s.Uptime)))
	return nil
}

func (p *Parser) GetRAMInfo(r *entities.Response) error {
	file, err := os.Open(entities.RamPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	mem := entities.RAM{}

	for ; mem.SReclaimable == 0; fileScanner.Scan() {
		str := fileScanner.Text()
		switch {
		case strings.HasPrefix(str, "MemTotal"):
			if mem.Total, err = strconv.ParseUint(strings.Fields(str)[1], 10, 64); err != nil {
				return err
			}
		case strings.HasPrefix(str, "MemFree"):
			if mem.Free, err = strconv.ParseUint(strings.Fields(str)[1], 10, 64); err != nil {
				return err
			}
		case strings.HasPrefix(str, "Buffers"):
			if mem.Buffers, err = strconv.ParseUint(strings.Fields(str)[1], 10, 64); err != nil {
				return err
			}
		case strings.HasPrefix(str, "Cached"):
			if mem.Cached, err = strconv.ParseUint(strings.Fields(str)[1], 10, 64); err != nil {
				return err
			}
		case strings.HasPrefix(str, "Shmem"):
			if mem.Shared, err = strconv.ParseUint(strings.Fields(str)[1], 10, 64); err != nil {
				return err
			}
		case strings.HasPrefix(str, "SReclaimable"):
			if mem.SReclaimable, err = strconv.ParseUint(strings.Fields(str)[1], 10, 64); err != nil {
				return err
			}
			mem.TrueFree = mem.Total + mem.Shared - mem.Buffers - mem.Cached - mem.Free - mem.SReclaimable
		}
	}
	r.Info = append(r.Info, fmt.Sprintf("RAM: %-65s|", fmt.Sprintf("%d/%d MiB [%.2f%%]", (mem.Total+mem.Shared-mem.Buffers-mem.Cached-mem.Free-mem.SReclaimable)/1024,
		mem.Total/1024, float64(mem.TrueFree)/float64(mem.Total)*100)))

	return nil
}

func (p *Parser) GetCPUInfo(r *entities.Response) error {
	file, err := os.Open(entities.CpuPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	cpu := entities.CPU{}
	if err = GetTemperatureInfo(&cpu); err != nil {
		return err
	}

	for ; cpu.Cores == 0; fileScanner.Scan() {
		str := fileScanner.Text()
		switch {
		case strings.HasPrefix(str, "model name"):
			cpu.ModelName = strings.Replace(str, "model name\t: ", "", 1)
		case strings.HasPrefix(str, "siblings"):
			if cpu.Siblings, err = strconv.Atoi(strings.Replace(str, "siblings\t: ", "", 1)); err != nil {
				return err
			}
		case strings.HasPrefix(str, "cpu cores"):
			if cpu.Cores, err = strconv.Atoi(strings.Replace(str, "cpu cores\t: ", "", 1)); err != nil {
				return err
			}
		}
	}

	r.Info = append(r.Info, fmt.Sprintf("CPU: %-65s|",
		fmt.Sprintf("%s; %d cores / %d threads\n\nTemperature zones: %v", cpu.ModelName, cpu.Cores, cpu.Siblings,
			cpu.Temperatures)))

	return nil
}

func GetTemperatureInfo(cpu *entities.CPU) error {
	if _, err := os.Stat("/sys/class/hwmon"); err != nil {
		return err
	}
	path, err := GetHardwareMon()
	if err != nil {
		return err
	}
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	var (
		data []byte
		temp uint64
	)

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "temp") && strings.HasSuffix(file.Name(), "input") {
			data, err = os.ReadFile(filepath.Join(path, file.Name()))
			if err != nil {
				return err
			}
			temp, err = strconv.ParseUint(strings.TrimRight(string(data), "\n"), 10, 32)
			if err != nil {
				return err
			}
			cpu.Temperatures = append(cpu.Temperatures, temp/1000)
		}
	}

	return nil
}

func GetHardwareMon() (string, error) {
	files, err := os.ReadDir("/sys/class/hwmon")
	if err != nil {
		return "", err
	}
	var filename string

	for _, file := range files {
		if file.Name()[:5] == "hwmon" {
			if filename, err = utils.CheckCPUMon(filepath.Join("/sys/class/hwmon", file.Name())); err == nil {
				return filename, nil
			}
		}
	}
	err = errors.New("no cpu-temp mon")
	return "", err
}

func (p *Parser) GetDiskInfo(r *entities.Response) error {
	fs := syscall.Statfs_t{}
	if err := syscall.Statfs("/", &fs); err != nil {
		return err
	}
	diskInfo := entities.DiskInfo{
		All: fs.Blocks * uint64(fs.Bsize),
	}
	diskInfo.Used = diskInfo.All - fs.Bfree*uint64(fs.Bsize)
	r.Info = append(r.Info, fmt.Sprintf("Disk Info: %-59s|", fmt.Sprintf("%.2f/%.2f GiB [%.2f%%]", float64(diskInfo.Used)/entities.GB, float64(diskInfo.All)/entities.GB,
		float32(diskInfo.Used)/float32(diskInfo.All)*100)))

	return nil
}

func (p *Parser) GetGPUInfo(r *entities.Response) error {
	path := "/sys/bus/pci/devices"
	if _, err := os.Stat(path); err != nil {
		return err
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	var (
		data      []byte
		strBuffer string
	)

	for _, file := range files {
		data, err = os.ReadFile(filepath.Join(path, file.Name(), "/class"))
		if strings.HasPrefix(string(data), "0x03") {
			strBuffer, err = GetPciId(filepath.Join(path, file.Name(), "/uevent"))
			if err != nil {
				return err
			}
			strBuffer, err = GetGPUModel(strings.ToLower(strBuffer))
			if err != nil {
				return err
			}
			r.Info = append(r.Info, fmt.Sprintf("GPU: %s", strBuffer))

		}
	}
	return nil
}

func GetPciId(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		str := fileScanner.Text()
		if strings.HasPrefix(str, "PCI_ID=") {
			return strings.Replace(str, "PCI_ID=", "", 1), nil
		}
	}

	err = errors.New("no PCI ID")
	return "", err
}

func GetGPUModel(pciId string) (string, error) {
	var numberId = fmt.Sprintf("\t%s", strings.Split(pciId, ":")[1])
	pciId = fmt.Sprintf("\t\t%s", pciId)
	file, err := os.Open("/usr/share/hwdata/pci.ids")
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		str := fileScanner.Text()
		if strings.HasPrefix(str, numberId) {
			rawResult := strings.Replace(str, numberId, "", 1)
			if strings.Contains(rawResult, "[") && strings.Contains(rawResult, "]") {
				indexStart := strings.Index(rawResult, "[")
				indexEnd := strings.Index(rawResult, "]")
				return rawResult[indexStart+1 : indexEnd], nil
			}
			return strings.TrimSpace(rawResult), nil
		} else if strings.HasPrefix(str, pciId) {
			rawResult := strings.Replace(str, pciId, "", 1)
			if strings.Contains(rawResult, "[") && strings.Contains(rawResult, "]") {
				indexStart := strings.Index(rawResult, "[")
				indexEnd := strings.Index(rawResult, "]")
				return rawResult[indexStart+1 : indexEnd], nil
			}
			return strings.TrimSpace(rawResult), nil
		}
	}

	err = errors.New("no GPU model")
	return "", err
}

func (p *Parser) GetShell(r *entities.Response) {
	data := strings.Split(os.Getenv("SHELL"), "/")
	r.Info = append(r.Info, fmt.Sprintf("Shell: %s", data[len(data)-1]))
}
