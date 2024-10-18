package sentinelone

type SentinelOneApplicationRisk struct {
	ApplicationId       string `json:"applicationId" bson:"applicationId"`
	CveCount            int    `json:"cveCount" bson:"cveCount"`
	DaysDetected        int    `json:"daysDetected" bson:"daysDetected"`
	DetectionDate       string `json:"detectionDate" bson:"detectionDate"`
	EndpointCount       int    `json:"endpointCount" bson:"endpointCount"`
	Estimate            string `json:"estimate" bson:"estimate"`
	HighestNvdBaseScore string `json:"highestNvdBaseScore" bson:"highestNvdBaseScore"`
	HighestSeverity     string `json:"highestSeverity" bson:"highestSeverity"`
	Name                string `json:"name" bson:"name"`
	Vendor              string `json:"vendor" bson:"vendor"`

	apikey string
}
