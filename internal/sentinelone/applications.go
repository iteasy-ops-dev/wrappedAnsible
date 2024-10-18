package sentinelone

import (
	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/internal/model"
)

func NewSentinelOneApplications(apikey string) *SentinelOneApplications {
	return &SentinelOneApplications{
		apikey: apikey,
	}
}

func (s *SentinelOneApplications) GetCollection() string {
	return config.GlobalConfig.MongoDB.Collections.SentinelOneApplication
}

func (s *SentinelOneApplications) GetAPIURL() string {
	return config.GlobalConfig.SentinelOne.URLs.ApplicationURL
}

func (s *SentinelOneApplications) GetAPIKey() string {
	return s.apikey
}

func (s *SentinelOneApplications) GetFilterKey() string {
	return "id"
}

// Update method to use the BaseUpdater's logic
func (s *SentinelOneApplications) Update() error {
	requester := model.BaseRequester{Requester: s}
	return requester.Update()
}
