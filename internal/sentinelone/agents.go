package sentinelone

import (
	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/internal/model"
)

func NewSentinelOneAgents(apikey string) *SentinelOneAgents {
	return &SentinelOneAgents{
		apikey: apikey,
	}
}

func (s *SentinelOneAgents) GetCollection() string {
	return config.GlobalConfig.MongoDB.Collections.SentinelOneAgents
}

func (s *SentinelOneAgents) GetAPIURL() string {
	return config.GlobalConfig.SentinelOne.URLs.AgentsURL
}

func (s *SentinelOneAgents) GetAPIKey() string {
	return s.apikey
}

func (s *SentinelOneAgents) GetFilterKey() string {
	return "id"
}

// Update method to use the BaseUpdater's logic
func (s *SentinelOneAgents) Update() error {
	requester := model.BaseRequester{Requester: s}
	return requester.Update()
}
