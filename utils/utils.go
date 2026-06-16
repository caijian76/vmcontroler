package utils

type CreateVMRequest struct {
	CPU         int    `json:"cpu"`
	Memory      int    `json:"memory"`
	Node        string `json:"node"`
	DiskSize    int    `json:"diskSize"`
	NetworkType string `json:"networkType"`
	AutoStart   bool   `json:"autoStart"`
	MountISO    bool   `json:"mountISO"`
	ISOPath     string `json:"isoPath"`
	VMName      string `json:"name"`
}
