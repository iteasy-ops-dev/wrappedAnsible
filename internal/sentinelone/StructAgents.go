package sentinelone

import "time"

type SentinelOneAgents struct {
	AccountID                      string                 `json:"accountId" bson:"accountId"`
	AccountName                    string                 `json:"accountName" bson:"accountName"`
	ActiveDirectory                ActiveDirectory        `json:"activeDirectory" bson:"activeDirectory"`
	ActiveThreats                  int                    `json:"activeThreats" bson:"activeThreats"`
	AgentVersion                   string                 `json:"agentVersion" bson:"agentVersion"`
	AllowRemoteShell               bool                   `json:"allowRemoteShell" bson:"allowRemoteShell"`
	AppsVulnerabilityStatus        string                 `json:"appsVulnerabilityStatus" bson:"appsVulnerabilityStatus"`
	CloudProviders                 map[string]string      `json:"cloudProviders" bson:"cloudProviders"`
	ComputerName                   string                 `json:"computerName" bson:"computerName"`
	ConsoleMigrationStatus         string                 `json:"consoleMigrationStatus" bson:"consoleMigrationStatus"`
	ContainerizedWorkloadCounts    interface{}            `json:"containerizedWorkloadCounts" bson:"containerizedWorkloadCounts"`
	CoreCount                      int                    `json:"coreCount" bson:"coreCount"`
	CPUCount                       int                    `json:"cpuCount" bson:"cpuCount"`
	CPUID                          string                 `json:"cpuId" bson:"cpuId"`
	CreatedAt                      time.Time              `json:"createdAt" bson:"createdAt"`
	DetectionState                 string                 `json:"detectionState" bson:"detectionState"`
	Domain                         string                 `json:"domain" bson:"domain"`
	EncryptedApplications          bool                   `json:"encryptedApplications" bson:"encryptedApplications"`
	ExternalID                     string                 `json:"externalId" bson:"externalId"`
	ExternalIP                     string                 `json:"externalIp" bson:"externalIp"`
	FirewallEnabled                bool                   `json:"firewallEnabled" bson:"firewallEnabled"`
	FirstFullModeTime              time.Time              `json:"firstFullModeTime" bson:"firstFullModeTime"`
	FullDiskScanLastUpdatedAt      time.Time              `json:"fullDiskScanLastUpdatedAt" bson:"fullDiskScanLastUpdatedAt"`
	GroupID                        string                 `json:"groupId" bson:"groupId"`
	GroupIP                        string                 `json:"groupIp" bson:"groupIp"`
	GroupName                      string                 `json:"groupName" bson:"groupName"`
	HasContainerizedWorkload       bool                   `json:"hasContainerizedWorkload" bson:"hasContainerizedWorkload"`
	ID                             string                 `json:"id" bson:"id"`
	InRemoteShellSession           bool                   `json:"inRemoteShellSession" bson:"inRemoteShellSession"`
	Infected                       bool                   `json:"infected" bson:"infected"`
	InstallerType                  string                 `json:"installerType" bson:"installerType"`
	IsActive                       bool                   `json:"isActive" bson:"isActive"`
	IsAdConnector                  bool                   `json:"isAdConnector" bson:"isAdConnector"`
	IsDecommissioned               bool                   `json:"isDecommissioned" bson:"isDecommissioned"`
	IsPendingUninstall             bool                   `json:"isPendingUninstall" bson:"isPendingUninstall"`
	IsUninstalled                  bool                   `json:"isUninstalled" bson:"isUninstalled"`
	IsUpToDate                     bool                   `json:"isUpToDate" bson:"isUpToDate"`
	LastActiveDate                 time.Time              `json:"lastActiveDate" bson:"lastActiveDate"`
	LastIpToMgmt                   string                 `json:"lastIpToMgmt" bson:"lastIpToMgmt"`
	LastLoggedInUserName           string                 `json:"lastLoggedInUserName" bson:"lastLoggedInUserName"`
	LastSuccessfulScanDate         time.Time              `json:"lastSuccessfulScanDate" bson:"lastSuccessfulScanDate"`
	LicenseKey                     string                 `json:"licenseKey" bson:"licenseKey"`
	LocationEnabled                bool                   `json:"locationEnabled" bson:"locationEnabled"`
	LocationType                   string                 `json:"locationType" bson:"locationType"`
	Locations                      []Location             `json:"locations" bson:"locations"`
	MachineSid                     *string                `json:"machineSid" bson:"machineSid,omitempty"`
	MachineType                    string                 `json:"machineType" bson:"machineType"`
	MissingPermissions             []string               `json:"missingPermissions" bson:"missingPermissions"`
	MitigationMode                 string                 `json:"mitigationMode" bson:"mitigationMode"`
	MitigationModeSuspicious       string                 `json:"mitigationModeSuspicious" bson:"mitigationModeSuspicious"`
	ModelName                      string                 `json:"modelName" bson:"modelName"`
	NetworkInterfaces              []NetworkInterface     `json:"networkInterfaces" bson:"networkInterfaces"`
	NetworkQuarantineEnabled       bool                   `json:"networkQuarantineEnabled" bson:"networkQuarantineEnabled"`
	NetworkStatus                  string                 `json:"networkStatus" bson:"networkStatus"`
	OperationalState               string                 `json:"operationalState" bson:"operationalState"`
	OperationalStateExpiration     *time.Time             `json:"operationalStateExpiration" bson:"operationalStateExpiration,omitempty"`
	OsArch                         string                 `json:"osArch" bson:"osArch"`
	OsName                         string                 `json:"osName" bson:"osName"`
	OsRevision                     string                 `json:"osRevision" bson:"osRevision"`
	OsStartTime                    time.Time              `json:"osStartTime" bson:"osStartTime"`
	OsType                         string                 `json:"osType" bson:"osType"`
	OsUsername                     *string                `json:"osUsername" bson:"osUsername,omitempty"`
	ProxyStates                    interface{}            `json:"proxyStates" bson:"proxyStates"`
	RangerStatus                   string                 `json:"rangerStatus" bson:"rangerStatus"`
	RangerVersion                  string                 `json:"rangerVersion" bson:"rangerVersion"`
	RegisteredAt                   time.Time              `json:"registeredAt" bson:"registeredAt"`
	RemoteProfilingState           string                 `json:"remoteProfilingState" bson:"remoteProfilingState"`
	RemoteProfilingStateExpiration *time.Time             `json:"remoteProfilingStateExpiration" bson:"remoteProfilingStateExpiration,omitempty"`
	ScanAbortedAt                  *time.Time             `json:"scanAbortedAt" bson:"scanAbortedAt,omitempty"`
	ScanFinishedAt                 time.Time              `json:"scanFinishedAt" bson:"scanFinishedAt"`
	ScanStartedAt                  time.Time              `json:"scanStartedAt" bson:"scanStartedAt"`
	ScanStatus                     string                 `json:"scanStatus" bson:"scanStatus"`
	SerialNumber                   string                 `json:"serialNumber" bson:"serialNumber"`
	ShowAlertIcon                  bool                   `json:"showAlertIcon" bson:"showAlertIcon"`
	SiteId                         string                 `json:"siteId" bson:"siteId"`
	SiteName                       string                 `json:"siteName" bson:"siteName"`
	StorageName                    *string                `json:"storageName" bson:"storageName,omitempty"`
	StorageType                    *string                `json:"storageType" bson:"storageType,omitempty"`
	Tags                           map[string]interface{} `json:"tags" bson:"tags"`
	ThreatRebootRequired           bool                   `json:"threatRebootRequired" bson:"threatRebootRequired"`
	TotalMemory                    int                    `json:"totalMemory" bson:"totalMemory"`
	UpdatedAt                      time.Time              `json:"updatedAt" bson:"updatedAt"`
	UserActionsNeeded              []string               `json:"userActionsNeeded" bson:"userActionsNeeded"`
	UUID                           string                 `json:"uuid" bson:"uuid"`

	apikey string
}

// ActiveDirectory represents the structure for Active Directory information.
type ActiveDirectory struct {
	ComputerDistinguishedName *string  `json:"computerDistinguishedName" bson:"computerDistinguishedName,omitempty"`
	ComputerMemberOf          []string `json:"computerMemberOf" bson:"computerMemberOf"`
	LastUserDistinguishedName *string  `json:"lastUserDistinguishedName" bson:"lastUserDistinguishedName,omitempty"`
	LastUserMemberOf          []string `json:"lastUserMemberOf" bson:"lastUserMemberOf"`
	UserPrincipalName         *string  `json:"userPrincipalName" bson:"userPrincipalName,omitempty"`
}

// Location represents a location information structure.
type Location struct {
	ID    string `json:"id" bson:"id"`
	Name  string `json:"name" bson:"name"`
	Scope string `json:"scope" bson:"scope"`
}

// NetworkInterface represents a network interface.
type NetworkInterface struct {
	GatewayIP         string   `json:"gatewayIp" bson:"gatewayIp"`
	GatewayMacAddress string   `json:"gatewayMacAddress" bson:"gatewayMacAddress"`
	ID                string   `json:"id" bson:"id"`
	Inet              []string `json:"inet" bson:"inet"`
	Inet6             []string `json:"inet6" bson:"inet6"`
	Name              string   `json:"name" bson:"name"`
	Physical          string   `json:"physical" bson:"physical"`
}
