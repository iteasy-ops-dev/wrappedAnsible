package ansible

import (
	"encoding/json"
	"log"
)

type AnsibleOutput struct {
	Plays []Plays `json:"plays"`
}

type Plays struct {
	Play  Play   `json:"play"`
	Tasks []Task `json:"tasks"`
}

type Play struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Task struct {
	Hosts      map[string]HostInfo `json:"hosts"`
	ChildTasks ChildTaskInfo       `json:"task"`
}

type HostInfo struct {
	Action  string                   `json:"action"`
	Facts   GatherFacts              `json:"ansible_facts"`
	Massage string                   `json:"msg"`
	Changed bool                     `json:"changed"`
	Failed  bool                     `json:"failed"`
	Results []map[string]interface{} `json:"results"`
}

type ChildTaskInfo struct {
	Name string `json:"name"`
}

type GatherFacts struct {
	AllIPv4Address           []string `json:"ansible_all_ipv4_addresses"`
	AllIPv6Address           []string `json:"ansible_all_ipv6_addresses"`
	Architecture             string   `json:"ansible_architecture"`
	BiosVersion              string   `json:"ansible_bios_version"`
	Distribution             string   `json:"ansible_distribution"`
	DistributionMajorVersion string   `json:"ansible_distribution_major_version"`
	OSFamily                 string   `json:"ansible_os_family"`

	ServiceManager string `json:"ansible_service_mgr"`
	PackageManager string `json:"ansible_pkg_mgr"`
	System         string `json:"ansible_system"`

	Processor               []string `json:"ansible_processor"`
	ProcessorCores          uint64   `json:"ansible_processor_cores"`
	ProcessorThreadsPerCore uint64   `json:"ansible_processor_threads_per_core"`

	MemFreeMB  uint64          `json:"ansible_memfree_mb"`
	MemTotalMB uint64          `json:"ansible_memtotal_mb"`
	MemoryMB   AnsibleMemoryMB `json:"ansible_memory_mb"`

	Mount []Mount `json:"ansible_mounts"`

	Devices map[string]DeviceInfo `json:"ansible_devices"`

	DateTime DateTime `json:"ansible_date_time"`
}

type DateTime struct {
	TZ string `jsont:"tz"`
}

type AnsibleMemoryMB struct {
	Nocache Nocache `json:"nocache"`
	Real    Real    `json:"real"`
	Swap    Swap    `json:"swap"`
}

type Nocache struct {
	Free uint64 `json:"free"`
	Used uint64 `json:"used"`
}

type Real struct {
	Free  uint64 `json:"free"`
	Total uint64 `json:"total"`
	Used  uint64 `json:"used"`
}

type Swap struct {
	Cached uint64 `json:"cached"`
	Free   uint64 `json:"free"`
	Total  uint64 `json:"total"`
	Used   uint64 `json:"used"`
}

type Mount struct {
	Device          string `json:"device"`
	Fstype          string `json:"fstype"`
	Mount           string `json:"mount"`
	Options         string `json:"options"`
	UUID            string `json:"uuid"`
	BlockAvailable  uint64 `json:"block_available"`
	BlockSize       uint64 `json:"block_size"`
	BlockTotal      uint64 `json:"block_total"`
	BlockUsed       uint64 `json:"block_used"`
	Inode_available uint64 `json:"inode_available"`
	Inode_total     uint64 `json:"inode_total"`
	Inode_used      uint64 `json:"inode_used"`
	SizeAvailable   uint64 `json:"size_available"`
	SizeTotal       uint64 `json:"size_total"`
}

type DeviceInfo struct {
	Host   string `json:"host"`
	Model  string `json:"model"`
	Size   string `json:"size"`
	Vendor string `json:"vendor"`
}

type Response struct {
	ID      string
	Name    string
	Status  bool
	Returns []Return
}

type Return struct {
	Name    string
	Host    string
	Action  string
	Msg     string
	Results []map[string]interface{}
	Failed  bool
}

func Output(outputData []byte) *AnsibleOutput {
	var output AnsibleOutput

	err := json.Unmarshal(outputData, &output)
	if err != nil {
		log.Println("Error parsing JSON:", err)
	}

	return &output
}

// func (o *AnsibleOutput) Info(cliStatus bool) Response {
// 	r := Response{
// 		ID:     o.Plays[0].Play.ID,
// 		Name:   o.Plays[0].Play.Name,
// 		Status: cliStatus,
// 	}

// 	for _, task := range o.Plays[0].Tasks {
// 		for host, hostInfo := range task.Hosts {
// 			if hostInfo.Action != "gather_facts" {
// 				continue
// 			}
// 			r.Returns = append(r.Returns, Return{
// 				Name:   task.ChildTasks.Name,
// 				Host:   host,
// 				Action: hostInfo.Action,
// 				Msg:    hostInfo.Facts,
// 				// Failed:  hostInfo.Failed,
// 				// Results: hostInfo.Results,
// 			})
// 			// fmt.Printf("Host: %s\n", host)
// 			// fmt.Printf("Action: %s\n", taskInfo.Action)
// 			// fmt.Printf("Failed: %+v\n", taskInfo.Failed)
// 			// fmt.Printf("Facts: %+v\n", hostInfo.Facts)
// 		}
// 	}

// 	return r
// }

func (o *AnsibleOutput) Debug(cliStatus bool) Response {
	r := Response{
		ID:     o.Plays[0].Play.ID,
		Name:   o.Plays[0].Play.Name,
		Status: cliStatus,
	}
	for _, task := range o.Plays[0].Tasks {
		for host, hostInfo := range task.Hosts {
			// if hostInfo.Action == "gather_facts" {
			// 	continue
			// }
			if hostInfo.Massage == "" {
				continue
			}
			r.Returns = append(r.Returns, Return{
				Name:   task.ChildTasks.Name,
				Host:   host,
				Action: hostInfo.Action,
				Msg:    hostInfo.Massage,
				// Failed:  hostInfo.Failed,
				Results: hostInfo.Results,
			})
		}
	}
	return r
}
