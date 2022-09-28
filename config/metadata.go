package config

import (
	"fmt"

	"cloud.google.com/go/compute/metadata"
)

type MetaData struct {
	projectId string
	region    string
	email     string
}

func NewMetadata() (*MetaData, error) {
	meta := getMetadataFromEnvVars()

	if meta.projectId != "" && meta.email != "" && meta.region != "" {
		return &meta, nil
	}

	projectId, err := metadata.ProjectID()

	if err != nil {
		return nil, fmt.Errorf("metadata server is not accessable. if you wanted to set metadata manually set all the variables.  %w", err)
	}

	region, err := metadata.Zone()

	if err != nil {
		return nil, err
	}

	email, err := metadata.Get("instance/service-accounts/default/email")

	if err != nil {
		return nil, err
	}

	return &MetaData{
		projectId: projectId,
		region:    region,
		email:     email,
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

func getMetadataFromEnvVars() MetaData {
	return MetaData{
		projectId: ProjectId,
		region:    Region,
		email:     SaEmail,
	}
}
