package sentinelone

import (
	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/internal/model"
)

func NewSentinelOneApplicationRisk(apikey string) *SentinelOneApplicationRisk {
	return &SentinelOneApplicationRisk{
		apikey: apikey,
	}
}

func (s *SentinelOneApplicationRisk) GetCollection() string {
	return config.GlobalConfig.MongoDB.Collections.SentinelOneApplicationRisk
}

func (s *SentinelOneApplicationRisk) GetAPIURL() string {
	return config.GlobalConfig.SentinelOne.URLs.ApplicationRiskURL
}

func (s *SentinelOneApplicationRisk) GetAPIKey() string {
	return s.apikey
}

func (s *SentinelOneApplicationRisk) GetFilterKey() string {
	return "applicationId"
}

// Update method to use the BaseUpdater's logic
func (s *SentinelOneApplicationRisk) Update() error {
	requester := model.BaseRequester{Requester: s}
	return requester.Update()
}
