package config

import (
	"os"
)

// Environment variables
const (
	GenerateMemeEnv     = "GENERATE_MEME"
	ProjectIDEnv        = "PROJECT_ID"
	KeyIDEnv            = "CLOUDRU_KEY_ID"
	KeySecretEnv        = "CLOUDRU_KEY_SECRET"
	DataDirEnv          = "DATA_DIR"
	ContainerJobNameEnv = "CONTAINER_JOB_NAME"
)

// GetGenerateMemeMode returns the meme generation mode based on environment variable
func GetGenerateMemeMode() string {
	return os.Getenv(GenerateMemeEnv)
}

// GetProjectID returns the project ID from environment variable
func GetProjectID() string {
	return os.Getenv(ProjectIDEnv)
}

// GetKeyID returns the Cloud.ru Key ID from environment variable
func GetKeyID() string {
	return os.Getenv(KeyIDEnv)
}

// GetKeySecret returns the Cloud.ru Key Secret from environment variable
func GetKeySecret() string {
	return os.Getenv(KeySecretEnv)
}

// GetDataDir returns the data directory path from environment variable or default
func GetDataDir() string {
	dataDir := os.Getenv(DataDirEnv)
	if dataDir == "" {
		return "./data"
	}
	return dataDir
}

// GetMemesDir returns the memes directory path
func GetMemesDir() string {
	return GetDataDir() + "/memes"
}

// GetTemplatesDir returns the templates directory path
func GetTemplatesDir() string {
	return GetDataDir() + "/templates"
}

// GetContainerJobName returns the container job name from environment variable or default
func GetContainerJobName() string {
	jobName := os.Getenv(ContainerJobNameEnv)
	if jobName == "" {
		return "generate-meme-job"
	}
	return jobName
}

// IsBackgroundMode checks if meme generation should happen in background
func IsBackgroundMode() bool {
	mode := GetGenerateMemeMode()
	return mode == "background"
}

// IsContainerAppJobMode checks if meme generation should use Container App Job
func IsContainerAppJobMode() bool {
	mode := GetGenerateMemeMode()
	return mode == "containerappjob"
}
