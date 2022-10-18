package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"

	"cloud.google.com/go/compute/metadata"
	"golang.org/x/oauth2/google"
)

type MetaData struct {
	projectId string
	region    string
	email     string
	selfUrl   string
}

func NewMetadata() (*MetaData, error) {
	meta := getMetadataFromEnvVars()

	if meta.projectId != "" && meta.email != "" && meta.region != "" && meta.selfUrl != "" {
		return &meta, nil
	}

	projectId, err := metadata.ProjectID()

	if err != nil {
		return nil, fmt.Errorf("metadata server is not accessable. if you wanted to set metadata manually set all the variables.  %w", err)
	}

	region, err := metadata.Zone() // Returns region+ sub zone so we need to trim last two chars "us-central1-b"
	region = region[:len(region)-2]

	if err != nil {
		return nil, err
	}

	email, err := metadata.Get("instance/service-accounts/default/email")

	if err != nil {
		return nil, err
	}

	projectNumber, err := metadata.NumericProjectID()

	if err != nil {
		return nil, err
	}

	selfUrl, err := getCloudRunUrl(region, projectNumber, ServiceName)

	if err != nil {
		return nil, err
	}

	log.Infof("meta.region:%s meta.projectId.region:%s meta.email:%s meta.selfUrl:%s ", region, projectId, email, selfUrl)

	return &MetaData{
		projectId: projectId,
		region:    region,
		email:     email,
		selfUrl:   selfUrl,
	}, nil
}

func (m *MetaData) GetProjectId() string {
	return m.projectId
}

func (m *MetaData) GetEmail() string {
	return m.email
}

func (m *MetaData) GetRegion() string {
	return m.region
}

func (m *MetaData) GetSelfUrl() string {
	return m.selfUrl
}

func getMetadataFromEnvVars() MetaData {
	return MetaData{
		projectId: ProjectId,
		region:    Region,
		email:     SaEmail,
		selfUrl:   SelfUrl,
	}
}

func getCloudRunUrl(region string, projectNumber string, serviceName string) (string, error) {

	ctx := context.Background()

	client, err := google.DefaultClient(ctx)

	if err != nil {
		return "", err
	}

	cloudRunApi := fmt.Sprintf("https://%s-run.googleapis.com/apis/serving.knative.dev/v1/namespaces/%s/services/%s", region, projectNumber, serviceName)

	log.Info("cloud run api url: ", cloudRunApi)

	resp, err := client.Get(cloudRunApi)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	log.Info("Response payload: ", string(body[:]))

	cloudRunResp := &cloudRunAPIUrlOnly{}
	json.Unmarshal(body, cloudRunResp)
	url := cloudRunResp.Status.URL

	log.Info("cloud run selfurl is set to: ", url)

	return url, nil
}

type cloudRunAPIUrlOnly struct {
	Status struct {
		URL string `json:"url"`
	} `json:"status"`
}
