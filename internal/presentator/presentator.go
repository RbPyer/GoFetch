package presentator

import (
	"github.com/RbPyer/Gofetch/internal/models"
	"fmt"
)

func Present(r *models.Response) {
	
	fmt.Printf(models.Template, r.Hostname, r.Username, r.OSRelease, r.KernelVersion, r.Uptime, 
		(r.Total+r.Shared-r.Buffers-r.Cached-r.Free-r.SReclaimable)/1024, r.Total/1024,	float64(r.TrueFree)/float64(r.Total)*100,
		r.ModelName, r.Cores, r.Siblings, r.Temperatures, 
		float64(r.Used)/models.GB, float64(r.All)/models.GB, float32(r.Used)/float32(r.All)*100, 
		r.GPUModel,
		r.Shell,
	)
}
