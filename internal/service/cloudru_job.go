package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"memes-generator/internal/config"
)

// CloudRuJobService handles interactions with Cloud.ru Container Apps Jobs
type CloudRuJobService struct {
	baseURL    string
	httpClient *http.Client
	projectID  string
	keyID      string
	keySecret  string
}

// NewCloudRuJobService creates a new Cloud.ru job service
func NewCloudRuJobService() *CloudRuJobService {
	return &CloudRuJobService{
		baseURL:   "https://containers.api.cloud.ru",
		projectID: config.GetProjectID(),
		keyID:     config.GetKeyID(),
		keySecret: config.GetKeySecret(),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Job represents a Cloud.ru job
type Job struct {
	ProjectID     string           `json:"projectId"`
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	Status        string           `json:"status"`
	CreatedAt     string           `json:"createdAt"`
	CreatedBy     string           `json:"createdBy"`
	UpdatedAt     string           `json:"updatedAt"`
	UpdatedBy     string           `json:"updatedBy"`
	Configuration JobConfiguration `json:"configuration"`
	Template      JobTemplate      `json:"template"`
}

// JobConfiguration represents job configuration
type JobConfiguration struct {
	Privileged     bool              `json:"privileged"`
	LoggingService JobLoggingService `json:"loggingService"`
}

// JobLoggingService represents logging service configuration
type JobLoggingService struct {
	IsEnabled bool   `json:"isEnabled"`
	GroupID   string `json:"groupId"`
}

// JobTemplate represents job template
type JobTemplate struct {
	MaxRetries     uint32             `json:"maxRetries"`
	Timeout        uint32             `json:"timeout"`
	Containers     []JobContainer     `json:"containers"`
	InitContainers []JobInitContainer `json:"initContainers,omitempty"`
	Volumes        []JobVolume        `json:"volumes,omitempty"`
}

// JobContainer represents a job container
type JobContainer struct {
	Name          string                `json:"name"`
	Image         string                `json:"image"`
	Env           []JobEnv              `json:"env,omitempty"`
	Command       []string              `json:"command,omitempty"`
	Args          []string              `json:"args,omitempty"`
	Resources     JobContainerResources `json:"resources,omitempty"`
	VolumeMounts  []JobVolumeMount      `json:"volumeMounts,omitempty"`
	LivenessProbe *JobLivenessProbe     `json:"livenessProbe,omitempty"`
}

// JobEnv represents environment variables
type JobEnv struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}

// JobContainerResources represents container resources
type JobContainerResources struct {
	CPU    string                    `json:"cpu"`
	Memory string                    `json:"memory"`
	GPU    *JobContainerResourcesGPU `json:"gpu,omitempty"`
}

// JobContainerResourcesGPU represents GPU resources
type JobContainerResourcesGPU struct {
	Count uint32 `json:"count"`
	SKU   string `json:"sku"`
}

// JobVolumeMount represents volume mount
type JobVolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
	SubPath   string `json:"subPath,omitempty"`
	ReadOnly  bool   `json:"readOnly"`
}

// JobLivenessProbe represents liveness probe
type JobLivenessProbe struct {
	ExecProbe           *JobExecProbe `json:"execProbe,omitempty"`
	InitialDelaySeconds int32         `json:"initialDelaySeconds,omitempty"`
	PeriodSeconds       int32         `json:"periodSeconds,omitempty"`
	TimeoutSeconds      int32         `json:"timeoutSeconds,omitempty"`
	FailureThreshold    int32         `json:"failureThreshold,omitempty"`
}

// JobExecProbe represents exec probe
type JobExecProbe struct {
	Command []string `json:"command"`
}

// JobInitContainer represents init container
type JobInitContainer struct {
	Name         string                `json:"name"`
	Image        string                `json:"image"`
	Env          []JobEnv              `json:"env,omitempty"`
	Command      []string              `json:"command,omitempty"`
	Args         []string              `json:"args,omitempty"`
	Resources    JobContainerResources `json:"resources,omitempty"`
	VolumeMounts []JobVolumeMount      `json:"volumeMounts,omitempty"`
}

// JobVolume represents volume
type JobVolume struct {
	Name             string            `json:"name"`
	Type             string            `json:"type"`
	VolumeAttributes map[string]string `json:"volumeAttributes"`
}

// getAccessToken gets an access token using KEY_ID and KEY_SECRET
func (s *CloudRuJobService) getAccessToken() (string, error) {
	url := "https://iam.api.cloud.ru/api/v1/auth/token"

	payload := strings.NewReader(fmt.Sprintf(`{"keyId": "%s","secret": "%s"}`, s.keyID, s.keySecret))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read containerapps response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Check if body is empty
	if len(body) == 0 {
		return "", fmt.Errorf("authentication API returned empty response body with status %d", resp.StatusCode)
	}

	// Parse response to get token
	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w body length: %d body: %s", err, len(body), string(body))
	}

	return result.AccessToken, nil
}

// GetJob retrieves a job by name
func (s *CloudRuJobService) GetJob(jobName string) (*Job, error) {
	url := fmt.Sprintf("%s/v2/jobs/%s", s.baseURL, jobName)

	// Add project ID as query parameter
	if s.projectID != "" {
		url = fmt.Sprintf("%s?projectId=%s", url, s.projectID)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Get access token
	token, err := s.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Add authorization header
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var job Job
	if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &job, nil
}

// ExecuteJob executes a job with the specified parameters
func (s *CloudRuJobService) ExecuteJob(jobName, memePath string) error {
	// First get the job to retrieve its current configuration
	job, err := s.GetJob(jobName)
	if err != nil {
		return fmt.Errorf("failed to get job: %w", err)
	}

	// Modify the job template to set the command and args for meme generation
	if len(job.Template.Containers) > 0 {
		// Set the command to "generate-meme"
		job.Template.Containers[0].Command = []string{"generate-meme"}

		// Set the args to include the meme path
		job.Template.Containers[0].Args = []string{"--meme-path", memePath}
	}

	// Prepare the patch request
	patchURL := fmt.Sprintf("%s/v2/jobs/%s", s.baseURL, jobName)

	// Add project ID as query parameter
	if s.projectID != "" {
		patchURL = fmt.Sprintf("%s?projectId=%s", patchURL, s.projectID)
	}

	// Create the request body
	requestBody := map[string]interface{}{
		"projectId":      s.projectID,
		"name":           job.Name,
		"description":    job.Description,
		"runImmediately": true,
		"configuration":  job.Configuration,
		"template":       job.Template,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create PATCH request
	req, err := http.NewRequest("PATCH", patchURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create patch request: %w", err)
	}

	// Get access token
	token, err := s.getAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// Add authorization header
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send patch request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("patch request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
