package parsers

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
	"github.com/RbPyer/Gofetch/internal/models"
	"github.com/RbPyer/Gofetch/internal/utils"
)

func GetUserInfo(r *models.Response) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	userObject, err := user.Current()
	if err != nil {
		return err
	}

	r.Hostname, r.Username = hostname, userObject.Username

	return nil
}

func GetOsVersion(r *models.Response) error {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		str := fileScanner.Text()
		if strings.HasPrefix(str, models.Prefix) {
			str = strings.ReplaceAll(strings.ReplaceAll(str, models.Prefix, ""), "\"", "")
			r.OSRelease = str
			return nil
		}
	}
	return nil
}

func GetKernelVersion(r *models.Response) error {
	u := syscall.Utsname{}
	err := syscall.Uname(&u)
	if err != nil {
		return err
	}
	r.KernelVersion = utils.UtsnameToStr(u.Release)

	return nil
}

func GetUptime(r *models.Response) error {
	s := syscall.Sysinfo_t{}
	err := syscall.Sysinfo(&s)
	if err != nil {
		return err
	}

	r.Uptime = utils.Int64ToTimeStr(s.Uptime)

	return nil
}

func GetRAMInfo(r *models.Response) error {
	file, err := os.Open(models.RamPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	for ; r.SReclaimable == 0; fileScanner.Scan() {
		str := fileScanner.Text()
		switch {
		case strings.HasPrefix(str, "MemTotal"):
			if r.Total, err = strconv.ParseUint(strings.Fields(str)[1], 10, 64); err != nil {
				return err
			}
		case strings.HasPrefix(str, "MemFree"):
			if r.Free, err = strconv.ParseUint(strings.Fields(str)[1], 10, 64); err != nil {
				return err
			}
		case strings.HasPrefix(str, "Buffers"):
			if r.Buffers, err = strconv.ParseUint(strings.Fields(str)[1], 10, 64); err != nil {
				return err
			}
		case strings.HasPrefix(str, "Cached"):
			if r.Cached, err = strconv.ParseUint(strings.Fields(str)[1], 10, 64); err != nil {
				return err
			}
		case strings.HasPrefix(str, "Shmem"):
			if r.Shared, err = strconv.ParseUint(strings.Fields(str)[1], 10, 64); err != nil {
				return err
			}
		case strings.HasPrefix(str, "SReclaimable"):
			if r.SReclaimable, err = strconv.ParseUint(strings.Fields(str)[1], 10, 64); err != nil {
				return err
			}
			r.TrueFree = r.Total + r.Shared - r.Buffers - r.Cached - r.Free - r.SReclaimable
		}
	}
	return nil
}

func GetCPUInfo(r *models.Response) error {
	file, err := os.Open(models.CpuPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	if err = getTemperatureInfo(&r.CPU); err != nil {
		return err
	}

	for ; r.Cores == 0; fileScanner.Scan() {
		str := fileScanner.Text()
		switch {
		case strings.HasPrefix(str, "model name"):
			r.ModelName = strings.Replace(str, "model name\t: ", "", 1)
		case strings.HasPrefix(str, "siblings"):
			if r.Siblings, _ = strconv.Atoi(strings.Replace(str, "siblings\t: ", "", 1)); err != nil {
				return err
			}
		case strings.HasPrefix(str, "cpu cores"):
			if r.Cores, err = strconv.Atoi(strings.Replace(str, "cpu cores\t: ", "", 1)); err != nil {
				return err
			}
		}
	}
	return nil
}

func getTemperatureInfo(cpu *models.CPU) error {
	if _, err := os.Stat("/sys/class/hwmon"); err != nil {
		return err
	}
	path, err := getHardwareMon()
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

func getHardwareMon() (string, error) {
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

func GetDiskInfo(r *models.Response) error {
	fs := syscall.Statfs_t{}
	if err := syscall.Statfs("/", &fs); err != nil {
		return err
	}

	r.All = fs.Blocks * uint64(fs.Bsize)
	r.Used = r.All - fs.Bfree*uint64(fs.Bsize)
	return nil
}

func GetGPUInfo(r *models.Response) error {
	files, err := os.ReadDir(models.GpuPath)
	if err != nil {
		return err
	}
	var (
		data      []byte
		strBuffer string
	)

	for _, file := range files {
		data, err = os.ReadFile(filepath.Join(models.GpuPath, file.Name(), "/class"))
		if err != nil {
			return err
		}
		if strings.HasPrefix(string(data), "0x03") {
			strBuffer, err = getPciId(filepath.Join(models.GpuPath, file.Name(), "/uevent"))
			if err != nil {
				return err
			}
			strBuffer, err = getGPUModel(strings.ToLower(strBuffer))
			if err != nil {
				return err
			}
			r.GPUModel = strBuffer
		}
	}
	return nil
}

func getPciId(path string) (string, error) {
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

func getGPUModel(pciId string) (string, error) {
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
				return rawResult[strings.Index(rawResult, "[")+1 : strings.Index(rawResult, "]")], nil
			}
			return strings.TrimSpace(rawResult), nil
		}
		if strings.HasPrefix(str, pciId) {
			rawResult := strings.Replace(str, pciId, "", 1)
			if strings.Contains(rawResult, "[") && strings.Contains(rawResult, "]") {
				return rawResult[strings.Index(rawResult, "[")+1 : strings.Index(rawResult, "]")], nil
			}
			return strings.TrimSpace(rawResult), nil
		}
	}

	err = errors.New("no GPU model")
	return "", err
}

func GetShell(r *models.Response) {
	data := strings.Split(os.Getenv("SHELL"), "/")
	r.Shell = data[len(data)-1]
}
