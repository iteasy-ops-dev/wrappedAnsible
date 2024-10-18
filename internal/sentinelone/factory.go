package sentinelone

import (
	"errors"

	config "iteasy.wrappedAnsible/configs"
)

func GenerateSentinelOne(gtype, apikey string) (iSentinelOne, error) {
	switch gtype {
	case config.GlobalConfig.MongoDB.Collections.SentinelOneAgents:
		return NewSentinelOneAgents(apikey), nil
	case config.GlobalConfig.MongoDB.Collections.SentinelOneApplication:
		return NewSentinelOneApplications(apikey), nil
	case config.GlobalConfig.MongoDB.Collections.SentinelOneApplicationRisk:
		return NewSentinelOneApplicationRisk(apikey), nil
	default:
		return nil, errors.New("구성 할 수 없는 타입")
	}
}
