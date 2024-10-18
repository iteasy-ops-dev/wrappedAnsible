package sentinelone

import "time"

type SentinelOneApplications struct {
	AgentComputerName     string `json:"agentComputerName" bson:"agentComputerName"`
	AgentDomain           string `json:"agentDomain" bson:"agentDomain"`
	AgentID               string `json:"agentId" bson:"agentId"`
	AgentMachineType      string `json:"agentMachineType" bson:"agentMachineType"`
	AgentNetworkStatus    string `json:"agentNetworkStatus" bson:"agentNetworkStatus"`
	AgentOperationalState string `json:"agentOperationalState" bson:"agentOperationalState"`
	AgentOsType           string `json:"agentOsType" bson:"agentOsType"`
	AgentUUID             string `json:"agentUuid" bson:"agentUuid"`
	AgentVersion          string `json:"agentVersion" bson:"agentVersion"`
	ID                    string `json:"id" bson:"id"`
	Name                  string `json:"name" bson:"name"`
	OsType                string `json:"osType" bson:"osType"`
	Publisher             string `json:"publisher" bson:"publisher"`
	RiskLevel             string `json:"riskLevel" bson:"riskLevel"`
	Type                  string `json:"type" bson:"type"`
	Version               string `json:"version" bson:"version"`

	Size int `json:"size" bson:"size"`

	Signed                bool `json:"signed" bson:"signed"`
	AgentInfected         bool `json:"agentInfected" bson:"agentInfected"`
	AgentIsActive         bool `json:"agentIsActive" bson:"agentIsActive"`
	AgentIsDecommissioned bool `json:"agentIsDecommissioned" bson:"agentIsDecommissioned"`

	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
	InstalledAt time.Time `json:"installedAt" bson:"installedAt"`

	apikey string
}
