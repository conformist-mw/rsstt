package service

// ServiceConfig holds configuration for services
type ServiceConfig struct {
	HTTPTimeout        int
	MaxRetries         int
	MaxMessagesPerUser int
	MaxConcurrentFeeds int
}

var serviceConfig *ServiceConfig

func SetServiceConfig(cfg *ServiceConfig) {
	serviceConfig = cfg
}

func GetServiceConfig() *ServiceConfig {
	return serviceConfig
}
