// © 2025 Platform Engineering Labs Inc.
//
// SPDX-License-Identifier: FSL-1.1-ALv2

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/masterminds/semver"
	"github.com/tidwall/gjson"

	"github.com/platform-engineering-labs/formae/pkg/model"
	"github.com/platform-engineering-labs/formae/pkg/plugin"
	"github.com/platform-engineering-labs/formae/pkg/plugin/resource"
	"github.com/platform-engineering-labs/formae/plugins/metakube/pkg/client"
	"github.com/platform-engineering-labs/formae/plugins/metakube/pkg/config"
)

type MetaKube struct{}

// Version set at compile time
var Version = "0.0.0"

// Compile time checks to satisfy protocol
var _ plugin.Plugin = MetaKube{}
var _ plugin.ResourcePlugin = MetaKube{}

// Plugin maintains the known symbol reference
var Plugin = MetaKube{}

var supportedResources = []plugin.ResourceDescriptor{
	{
		Type:         "MetaKube::Kubernetes::Cluster",
		Discoverable: true,
	},
	{
		Type:         "MetaKube::Kubernetes::NodeDeployment",
		Discoverable: true,
		ParentResourceTypesWithMappingProperties: map[string][]plugin.ListParameter{
			"MetaKube::Kubernetes::Cluster": {
				{
					ParentProperty: "id",
					ListProperty:   "clusterID",
				},
			},
		},
	},
	{
		Type:         "MetaKube::Project::SSHKey",
		Discoverable: true,
	},
}

func (m MetaKube) Name() string {
	return "metakube"
}

func (m MetaKube) Version() *semver.Version {
	return semver.MustParse(Version)
}

func (m MetaKube) Type() plugin.Type {
	return plugin.Resource
}

func (m MetaKube) Namespace() string {
	return "MetaKube"
}

func (m MetaKube) SupportedResources() []plugin.ResourceDescriptor {
	return supportedResources
}

func (m MetaKube) MaxRequestsPerSecond() int {
	return 5
}

func (m MetaKube) SchemaForResourceType(resourceType string) (model.Schema, error) {
	switch resourceType {
	case "MetaKube::Kubernetes::Cluster":
		return model.Schema{
			Identifier: "id",
			Fields: []string{
				"id", "name", "creationTimestamp", "spec", "status",
			},
		}, nil
	case "MetaKube::Project::SSHKey":
		return model.Schema{
			Identifier: "id",
			Fields: []string{
				"id", "name", "publicKey", "fingerprint",
			},
		}, nil
	case "MetaKube::Kubernetes::NodeDeployment":
		return model.Schema{
			Identifier: "id",
			Fields: []string{
				"id", "name", "spec", "status",
			},
		}, nil
	default:
		return model.Schema{}, fmt.Errorf("unknown resource type: %s", resourceType)
	}
}

func (m MetaKube) Create(ctx context.Context, request *resource.CreateRequest) (*resource.CreateResult, error) {
	cfg := config.FromTarget(request.Target)
	apiClient, err := client.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create MetaKube client: %w", err)
	}

	switch request.Resource.Type {
	case "MetaKube::Kubernetes::Cluster":
		return m.createCluster(ctx, apiClient, request)
	case "MetaKube::Kubernetes::NodeDeployment":
		return m.createNodeDeployment(ctx, apiClient, request)
	case "MetaKube::Project::SSHKey":
		return m.createSSHKey(ctx, apiClient, request)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", request.Resource.Type)
	}
}

func (m MetaKube) Update(ctx context.Context, request *resource.UpdateRequest) (*resource.UpdateResult, error) {
	cfg := config.FromTarget(request.Target)
	apiClient, err := client.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create MetaKube client: %w", err)
	}

	switch request.Resource.Type {
	case "MetaKube::Kubernetes::Cluster":
		return m.updateCluster(ctx, apiClient, request)
	case "MetaKube::Kubernetes::NodeDeployment":
		return m.updateNodeDeployment(ctx, apiClient, request)
	case "MetaKube::Project::SSHKey":
		return m.updateSSHKey(ctx, apiClient, request)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", request.Resource.Type)
	}
}

func (m MetaKube) Status(ctx context.Context, request *resource.StatusRequest) (*resource.StatusResult, error) {
	cfg := config.FromTarget(request.Target)
	apiClient, err := client.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create MetaKube client: %w", err)
	}

	// Get resource status based on type
	switch request.ResourceType {
	case "MetaKube::Kubernetes::Cluster":
		return m.statusCluster(ctx, apiClient, request)
	case "MetaKube::Kubernetes::NodeDeployment":
		return m.statusNodeDeployment(ctx, apiClient, request)
	case "MetaKube::Project::SSHKey":
		return m.statusSSHKey(ctx, apiClient, request)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", request.ResourceType)
	}
}

func (m MetaKube) Delete(ctx context.Context, request *resource.DeleteRequest) (*resource.DeleteResult, error) {
	cfg := config.FromTarget(request.Target)
	apiClient, err := client.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create MetaKube client: %w", err)
	}

	switch request.ResourceType {
	case "MetaKube::Kubernetes::Cluster":
		return m.deleteCluster(ctx, apiClient, request)
	case "MetaKube::Kubernetes::NodeDeployment":
		return m.deleteNodeDeployment(ctx, apiClient, request)
	case "MetaKube::Project::SSHKey":
		return m.deleteSSHKey(ctx, apiClient, request)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", request.ResourceType)
	}
}

func (m MetaKube) Read(ctx context.Context, request *resource.ReadRequest) (*resource.ReadResult, error) {
	cfg := config.FromTarget(request.Target)
	apiClient, err := client.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create MetaKube client: %w", err)
	}

	switch request.ResourceType {
	case "MetaKube::Kubernetes::Cluster":
		return m.readCluster(ctx, apiClient, request)
	case "MetaKube::Kubernetes::NodeDeployment":
		return m.readNodeDeployment(ctx, apiClient, request)
	case "MetaKube::Project::SSHKey":
		return m.readSSHKey(ctx, apiClient, request)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", request.ResourceType)
	}
}

func (m MetaKube) List(ctx context.Context, request *resource.ListRequest) (*resource.ListResult, error) {
	cfg := config.FromTarget(request.Target)
	apiClient, err := client.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create MetaKube client: %w", err)
	}

	switch request.ResourceType {
	case "MetaKube::Kubernetes::Cluster":
		return m.listClusters(ctx, apiClient, request)
	case "MetaKube::Kubernetes::NodeDeployment":
		return m.listNodeDeployments(ctx, apiClient, request)
	case "MetaKube::Project::SSHKey":
		return m.listSSHKeys(ctx, apiClient, request)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", request.ResourceType)
	}
}

func (m MetaKube) TargetBehavior() resource.TargetBehavior {
	return TargetBehavior
}

// GetResourceFilters returns a map of resource types to filter functions
func (m MetaKube) GetResourceFilters() map[string]plugin.ResourceFilter {
	// MetaKube doesn't need any resource filters currently
	return make(map[string]plugin.ResourceFilter)
}

// Helper function to extract field from JSON
func getStringField(data json.RawMessage, path string) string {
	return gjson.GetBytes(data, path).String()
}

// Cluster operations
func (m MetaKube) createCluster(ctx context.Context, apiClient *client.Client, request *resource.CreateRequest) (*resource.CreateResult, error) {
	cfg := config.FromTarget(request.Target)

	projectID := getStringField(request.Resource.Properties, "projectID")
	if projectID == "" {
		return nil, fmt.Errorf("projectID is required")
	}

	// Parse the properties to transform dcName -> spec.cloud.dc
	var props map[string]interface{}
	if err := json.Unmarshal(request.Resource.Properties, &props); err != nil {
		return nil, fmt.Errorf("failed to parse properties: %w", err)
	}

	// Extract dcName and inject it into spec.cloud.dc
	dcName, _ := props["dcName"].(string)
	if dcName == "" {
		return nil, fmt.Errorf("dcName is required")
	}

	// Ensure spec.cloud.dc is set and inject OpenStack credentials from config
	if spec, ok := props["spec"].(map[string]interface{}); ok {
		if cloud, ok := spec["cloud"].(map[string]interface{}); ok {
			cloud["dc"] = dcName

			// Inject OpenStack credentials if available in config
			if cfg.ApplicationCredentialID != "" && cfg.ApplicationCredentialSecret != "" {
				if openstack, ok := cloud["openstack"].(map[string]interface{}); ok {
					openstack["applicationCredentialID"] = cfg.ApplicationCredentialID
					openstack["applicationCredentialSecret"] = cfg.ApplicationCredentialSecret
				}
			}
		} else {
			spec["cloud"] = map[string]interface{}{"dc": dcName}
		}
	}

	// Remove fields that shouldn't be sent to API
	delete(props, "projectID")
	delete(props, "dcName")

	// Wrap cluster spec in "cluster" object as expected by CreateClusterSpec API
	requestBody := map[string]interface{}{
		"cluster": props,
	}

	// Debug: log what we're sending
	debugJSON, _ := json.MarshalIndent(requestBody, "", "  ")
	fmt.Printf("DEBUG: Creating cluster with body:\n%s\n", string(debugJSON))

	// Create cluster via API (Post will marshal requestBody to JSON)
	path := fmt.Sprintf("/api/v2/projects/%s/clusters", projectID)
	respBody, err := apiClient.Post(ctx, path, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster: %w", err)
	}

	clusterID := gjson.GetBytes(respBody, "id").String()
	if clusterID == "" {
		return nil, fmt.Errorf("cluster ID not found in response")
	}

	return &resource.CreateResult{
		ProgressResult: &resource.ProgressResult{
			Operation:          resource.OperationCreate,
			OperationStatus:    resource.OperationStatusInProgress,
			NativeID:           clusterID,
			ResourceType:       "MetaKube::Kubernetes::Cluster",
			ResourceProperties: respBody,
			Metadata:           respBody,
			StartTs:            time.Now(),
			ModifiedTs:         time.Now(),
		},
	}, nil
}

func (m MetaKube) readCluster(ctx context.Context, apiClient *client.Client, request *resource.ReadRequest) (*resource.ReadResult, error) {
	cfg := config.FromTarget(request.Target)
	if cfg.ProjectID == "" {
		return nil, fmt.Errorf("projectID is required in target configuration for reading cluster")
	}
	path := fmt.Sprintf("/api/v2/projects/%s/clusters/%s", cfg.ProjectID, request.NativeID)

	respBody, err := apiClient.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read cluster: %w", err)
	}

	return &resource.ReadResult{
		ResourceType: "MetaKube::Kubernetes::Cluster",
		Properties:   string(respBody),
	}, nil
}

func (m MetaKube) updateCluster(ctx context.Context, apiClient *client.Client, request *resource.UpdateRequest) (*resource.UpdateResult, error) {
	projectID := getStringField(request.Resource.Properties, "projectID")
	path := fmt.Sprintf("/api/v2/projects/%s/clusters/%s", projectID, *request.NativeID)

	respBody, err := apiClient.Patch(ctx, path, request.Resource.Properties)
	if err != nil {
		return nil, fmt.Errorf("failed to update cluster: %w", err)
	}

	return &resource.UpdateResult{
		ProgressResult: &resource.ProgressResult{
			Operation:          resource.OperationUpdate,
			OperationStatus:    resource.OperationStatusInProgress,
			NativeID:           *request.NativeID,
			ResourceType:       "MetaKube::Kubernetes::Cluster",
			ResourceProperties: respBody,
			Metadata:           respBody,
			StartTs:            time.Now(),
			ModifiedTs:         time.Now(),
		},
	}, nil
}

func (m MetaKube) deleteCluster(ctx context.Context, apiClient *client.Client, request *resource.DeleteRequest) (*resource.DeleteResult, error) {
	projectID := getStringField(request.Metadata, "projectID")
	path := fmt.Sprintf("/api/v2/projects/%s/clusters/%s", projectID, *request.NativeID)

	err := apiClient.Delete(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to delete cluster: %w", err)
	}

	return &resource.DeleteResult{
		ProgressResult: &resource.ProgressResult{
			Operation:       resource.OperationDelete,
			OperationStatus: resource.OperationStatusSuccess,
			NativeID:        *request.NativeID,
			ResourceType:    "MetaKube::Kubernetes::Cluster",
			StartTs:         time.Now(),
			ModifiedTs:      time.Now(),
		},
	}, nil
}

func (m MetaKube) statusCluster(ctx context.Context, apiClient *client.Client, request *resource.StatusRequest) (*resource.StatusResult, error) {
	projectID := getStringField(request.Metadata, "projectID")
	nativeID := request.RequestID // RequestID contains the native ID for status checks
	path := fmt.Sprintf("/api/v2/projects/%s/clusters/%s/health", projectID, nativeID)

	respBody, err := apiClient.Get(ctx, path)
	if err != nil {
		return &resource.StatusResult{
			ProgressResult: &resource.ProgressResult{
				Operation:       resource.OperationCheckStatus,
				OperationStatus: resource.OperationStatusInProgress,
				NativeID:        nativeID,
				ResourceType:    "MetaKube::Kubernetes::Cluster",
				StartTs:         time.Now(),
				ModifiedTs:      time.Now(),
			},
		}, nil
	}

	// Check health status
	health := gjson.GetBytes(respBody, "status").String()
	if health == "Running" || health == "Healthy" {
		return &resource.StatusResult{
			ProgressResult: &resource.ProgressResult{
				Operation:       resource.OperationCheckStatus,
				OperationStatus: resource.OperationStatusSuccess,
				NativeID:        nativeID,
				ResourceType:    "MetaKube::Kubernetes::Cluster",
				StartTs:         time.Now(),
				ModifiedTs:      time.Now(),
			},
		}, nil
	}

	return &resource.StatusResult{
		ProgressResult: &resource.ProgressResult{
			Operation:       resource.OperationCheckStatus,
			OperationStatus: resource.OperationStatusInProgress,
			NativeID:        nativeID,
			ResourceType:    "MetaKube::Kubernetes::Cluster",
			StartTs:         time.Now(),
			ModifiedTs:      time.Now(),
		},
	}, nil
}

func (m MetaKube) listClusters(ctx context.Context, apiClient *client.Client, request *resource.ListRequest) (*resource.ListResult, error) {
	cfg := config.FromTarget(request.Target)
	if cfg.ProjectID == "" {
		return nil, fmt.Errorf("projectID is required in target configuration for listing clusters")
	}

	path := fmt.Sprintf("/api/v2/projects/%s/clusters", cfg.ProjectID)
	respBody, err := apiClient.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list clusters: %w", err)
	}

	var resources []resource.Resource
	result := gjson.ParseBytes(respBody)
	result.ForEach(func(key, value gjson.Result) bool {
		resources = append(resources, resource.Resource{
			NativeID:   value.Get("id").String(),
			Properties: value.Raw,
		})
		return true
	})

	return &resource.ListResult{
		ResourceType: "MetaKube::Kubernetes::Cluster",
		Resources:    resources,
	}, nil
}

// NodeDeployment operations (Kubermatic calls them MachineDeployments)
func (m MetaKube) createNodeDeployment(ctx context.Context, apiClient *client.Client, request *resource.CreateRequest) (*resource.CreateResult, error) {
	projectID := getStringField(request.Resource.Properties, "projectID")
	clusterID := getStringField(request.Resource.Properties, "clusterID")

	path := fmt.Sprintf("/api/v2/projects/%s/clusters/%s/machinedeployments", projectID, clusterID)
	respBody, err := apiClient.Post(ctx, path, request.Resource.Properties)
	if err != nil {
		return nil, fmt.Errorf("failed to create node deployment: %w", err)
	}

	nodeDeploymentID := gjson.GetBytes(respBody, "id").String()
	return &resource.CreateResult{
		ProgressResult: &resource.ProgressResult{
			Operation:          resource.OperationCreate,
			OperationStatus:    resource.OperationStatusInProgress,
			NativeID:           nodeDeploymentID,
			ResourceType:       "MetaKube::Kubernetes::NodeDeployment",
			ResourceProperties: respBody,
			Metadata:           respBody,
			StartTs:            time.Now(),
			ModifiedTs:         time.Now(),
		},
	}, nil
}

func (m MetaKube) readNodeDeployment(ctx context.Context, apiClient *client.Client, request *resource.ReadRequest) (*resource.ReadResult, error) {
	cfg := config.FromTarget(request.Target)
	if cfg.ProjectID == "" {
		return nil, fmt.Errorf("projectID is required in target configuration for reading node deployment")
	}
	clusterID := getStringField(request.Metadata, "clusterID")
	if clusterID == "" {
		return nil, fmt.Errorf("clusterID is required in metadata for reading node deployment")
	}
	path := fmt.Sprintf("/api/v2/projects/%s/clusters/%s/machinedeployments/%s", cfg.ProjectID, clusterID, request.NativeID)

	respBody, err := apiClient.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read node deployment: %w", err)
	}

	return &resource.ReadResult{
		ResourceType: "MetaKube::Kubernetes::NodeDeployment",
		Properties:   string(respBody),
	}, nil
}

func (m MetaKube) updateNodeDeployment(ctx context.Context, apiClient *client.Client, request *resource.UpdateRequest) (*resource.UpdateResult, error) {
	projectID := getStringField(request.Resource.Properties, "projectID")
	clusterID := getStringField(request.Resource.Properties, "clusterID")
	path := fmt.Sprintf("/api/v2/projects/%s/clusters/%s/machinedeployments/%s", projectID, clusterID, *request.NativeID)

	respBody, err := apiClient.Patch(ctx, path, request.Resource.Properties)
	if err != nil {
		return nil, fmt.Errorf("failed to update node deployment: %w", err)
	}

	return &resource.UpdateResult{
		ProgressResult: &resource.ProgressResult{
			Operation:          resource.OperationUpdate,
			OperationStatus:    resource.OperationStatusInProgress,
			NativeID:           *request.NativeID,
			ResourceType:       "MetaKube::Kubernetes::NodeDeployment",
			ResourceProperties: respBody,
			Metadata:           respBody,
			StartTs:            time.Now(),
			ModifiedTs:         time.Now(),
		},
	}, nil
}

func (m MetaKube) deleteNodeDeployment(ctx context.Context, apiClient *client.Client, request *resource.DeleteRequest) (*resource.DeleteResult, error) {
	projectID := getStringField(request.Metadata, "projectID")
	clusterID := getStringField(request.Metadata, "clusterID")
	path := fmt.Sprintf("/api/v2/projects/%s/clusters/%s/machinedeployments/%s", projectID, clusterID, *request.NativeID)

	err := apiClient.Delete(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to delete node deployment: %w", err)
	}

	return &resource.DeleteResult{
		ProgressResult: &resource.ProgressResult{
			Operation:       resource.OperationDelete,
			OperationStatus: resource.OperationStatusSuccess,
			NativeID:        *request.NativeID,
			ResourceType:    "MetaKube::Kubernetes::NodeDeployment",
			StartTs:         time.Now(),
			ModifiedTs:      time.Now(),
		},
	}, nil
}

func (m MetaKube) statusNodeDeployment(ctx context.Context, apiClient *client.Client, request *resource.StatusRequest) (*resource.StatusResult, error) {
	readResult, err := m.readNodeDeployment(ctx, apiClient, &resource.ReadRequest{
		ResourceType: request.ResourceType,
		NativeID:     request.RequestID,
		Metadata:     request.Metadata,
	})
	if err != nil {
		return &resource.StatusResult{
			ProgressResult: &resource.ProgressResult{
				Operation:       resource.OperationCheckStatus,
				OperationStatus: resource.OperationStatusInProgress,
				NativeID:        request.RequestID,
				ResourceType:    "MetaKube::Kubernetes::NodeDeployment",
				StartTs:         time.Now(),
				ModifiedTs:      time.Now(),
			},
		}, nil
	}

	// Check if node deployment has ready replicas
	availableReplicas := gjson.Get(readResult.Properties, "status.availableReplicas").Int()
	desiredReplicas := gjson.Get(readResult.Properties, "spec.replicas").Int()

	if availableReplicas >= desiredReplicas && desiredReplicas > 0 {
		return &resource.StatusResult{
			ProgressResult: &resource.ProgressResult{
				Operation:       resource.OperationCheckStatus,
				OperationStatus: resource.OperationStatusSuccess,
				NativeID:        request.RequestID,
				ResourceType:    "MetaKube::Kubernetes::NodeDeployment",
				StartTs:         time.Now(),
				ModifiedTs:      time.Now(),
			},
		}, nil
	}

	return &resource.StatusResult{
		ProgressResult: &resource.ProgressResult{
			Operation:       resource.OperationCheckStatus,
			OperationStatus: resource.OperationStatusInProgress,
			NativeID:        request.RequestID,
			ResourceType:    "MetaKube::Kubernetes::NodeDeployment",
			StartTs:         time.Now(),
			ModifiedTs:      time.Now(),
		},
	}, nil
}

func (m MetaKube) listNodeDeployments(ctx context.Context, apiClient *client.Client, request *resource.ListRequest) (*resource.ListResult, error) {
	cfg := config.FromTarget(request.Target)
	if cfg.ProjectID == "" {
		return nil, fmt.Errorf("projectID is required in target configuration for listing node deployments")
	}

	// clusterID comes from AdditionalProperties (parent resource)
	clusterID := ""
	if request.AdditionalProperties != nil {
		if cid, ok := request.AdditionalProperties["clusterID"]; ok {
			clusterID = cid
		}
	}

	if clusterID == "" {
		return nil, fmt.Errorf("clusterID is required for listing node deployments")
	}

	path := fmt.Sprintf("/api/v2/projects/%s/clusters/%s/machinedeployments", cfg.ProjectID, clusterID)
	respBody, err := apiClient.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list node deployments: %w", err)
	}

	var resources []resource.Resource
	result := gjson.ParseBytes(respBody)
	result.ForEach(func(key, value gjson.Result) bool {
		resources = append(resources, resource.Resource{
			NativeID:   value.Get("id").String(),
			Properties: value.Raw,
		})
		return true
	})

	return &resource.ListResult{
		ResourceType: "MetaKube::Kubernetes::NodeDeployment",
		Resources:    resources,
	}, nil
}

// SSHKey operations (uses v1 API)
func (m MetaKube) createSSHKey(ctx context.Context, apiClient *client.Client, request *resource.CreateRequest) (*resource.CreateResult, error) {
	projectID := getStringField(request.Resource.Properties, "projectID")
	path := fmt.Sprintf("/api/v1/projects/%s/sshkeys", projectID)

	respBody, err := apiClient.Post(ctx, path, request.Resource.Properties)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH key: %w", err)
	}

	sshKeyID := gjson.GetBytes(respBody, "id").String()
	return &resource.CreateResult{
		ProgressResult: &resource.ProgressResult{
			Operation:          resource.OperationCreate,
			OperationStatus:    resource.OperationStatusSuccess,
			NativeID:           sshKeyID,
			ResourceType:       "MetaKube::Project::SSHKey",
			ResourceProperties: respBody,
			Metadata:           respBody,
			StartTs:            time.Now(),
			ModifiedTs:         time.Now(),
		},
	}, nil
}

func (m MetaKube) readSSHKey(ctx context.Context, apiClient *client.Client, request *resource.ReadRequest) (*resource.ReadResult, error) {
	cfg := config.FromTarget(request.Target)
	if cfg.ProjectID == "" {
		return nil, fmt.Errorf("projectID is required in target configuration for reading SSH key")
	}
	path := fmt.Sprintf("/api/v1/projects/%s/sshkeys/%s", cfg.ProjectID, request.NativeID)

	respBody, err := apiClient.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH key: %w", err)
	}

	return &resource.ReadResult{
		ResourceType: "MetaKube::Project::SSHKey",
		Properties:   string(respBody),
	}, nil
}

func (m MetaKube) updateSSHKey(ctx context.Context, apiClient *client.Client, request *resource.UpdateRequest) (*resource.UpdateResult, error) {
	// SSH keys are typically immutable, but we'll support name updates
	projectID := getStringField(request.Resource.Properties, "projectID")
	path := fmt.Sprintf("/api/v1/projects/%s/sshkeys/%s", projectID, *request.NativeID)

	respBody, err := apiClient.Put(ctx, path, request.Resource.Properties)
	if err != nil {
		return nil, fmt.Errorf("failed to update SSH key: %w", err)
	}

	return &resource.UpdateResult{
		ProgressResult: &resource.ProgressResult{
			Operation:          resource.OperationUpdate,
			OperationStatus:    resource.OperationStatusSuccess,
			NativeID:           *request.NativeID,
			ResourceType:       "MetaKube::Project::SSHKey",
			ResourceProperties: respBody,
			Metadata:           respBody,
			StartTs:            time.Now(),
			ModifiedTs:         time.Now(),
		},
	}, nil
}

func (m MetaKube) deleteSSHKey(ctx context.Context, apiClient *client.Client, request *resource.DeleteRequest) (*resource.DeleteResult, error) {
	projectID := getStringField(request.Metadata, "projectID")
	path := fmt.Sprintf("/api/v1/projects/%s/sshkeys/%s", projectID, *request.NativeID)

	err := apiClient.Delete(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to delete SSH key: %w", err)
	}

	return &resource.DeleteResult{
		ProgressResult: &resource.ProgressResult{
			Operation:       resource.OperationDelete,
			OperationStatus: resource.OperationStatusSuccess,
			NativeID:        *request.NativeID,
			ResourceType:    "MetaKube::Project::SSHKey",
			StartTs:         time.Now(),
			ModifiedTs:      time.Now(),
		},
	}, nil
}

func (m MetaKube) statusSSHKey(ctx context.Context, apiClient *client.Client, request *resource.StatusRequest) (*resource.StatusResult, error) {
	// SSH keys are simple resources, if they exist they're ready
	_, err := m.readSSHKey(ctx, apiClient, &resource.ReadRequest{
		ResourceType: request.ResourceType,
		NativeID:     request.RequestID,
		Metadata:     request.Metadata,
	})
	if err != nil {
		return &resource.StatusResult{
			ProgressResult: &resource.ProgressResult{
				Operation:       resource.OperationCheckStatus,
				OperationStatus: resource.OperationStatusFailure,
				NativeID:        request.RequestID,
				ResourceType:    "MetaKube::Project::SSHKey",
				StartTs:         time.Now(),
				ModifiedTs:      time.Now(),
			},
		}, nil
	}

	return &resource.StatusResult{
		ProgressResult: &resource.ProgressResult{
			Operation:       resource.OperationCheckStatus,
			OperationStatus: resource.OperationStatusSuccess,
			NativeID:        request.RequestID,
			ResourceType:    "MetaKube::Project::SSHKey",
			StartTs:         time.Now(),
			ModifiedTs:      time.Now(),
		},
	}, nil
}

func (m MetaKube) listSSHKeys(ctx context.Context, apiClient *client.Client, request *resource.ListRequest) (*resource.ListResult, error) {
	cfg := config.FromTarget(request.Target)
	if cfg.ProjectID == "" {
		return nil, fmt.Errorf("projectID is required in target configuration for listing SSH keys")
	}

	path := fmt.Sprintf("/api/v1/projects/%s/sshkeys", cfg.ProjectID)
	respBody, err := apiClient.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list SSH keys: %w", err)
	}

	var resources []resource.Resource
	result := gjson.ParseBytes(respBody)
	result.ForEach(func(key, value gjson.Result) bool {
		resources = append(resources, resource.Resource{
			NativeID:   value.Get("id").String(),
			Properties: value.Raw,
		})
		return true
	})

	return &resource.ListResult{
		ResourceType: "MetaKube::Project::SSHKey",
		Resources:    resources,
	}, nil
}
