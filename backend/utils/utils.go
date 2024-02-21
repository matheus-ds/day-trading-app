package utils

// Get config path for local or docker
func GetConfigPath(env string) string {
	if env == "docker" {
		return "./config/config.docker"
	}
	return "./config/config.local"
}
