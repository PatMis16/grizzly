package grafana

import (
	"fmt"

	"github.com/grafana/grizzly/pkg/grizzly"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

/*
 * @TODO
 * 1. The API does not have a GET method, so we have to fake it here
 * 2. The API expects an ID and a tenantId in an update, but these are
 *    generated by the server so cannot be represented in Jsonnet.
 *    Therefore, we have to pre-retrieve the check to get those values
 *    so we can inject them before posting JSON.
 * 3. This means pre-retrieving the check *twice*, once to establish
 *    whether this resource has changed or not (within Grizzly ifself)
 *    and again within this provider to retrieve IDs. Not ideal.
 */

// SyntheticMonitoringHandler is a Grizzly Handler for Grafana Synthetic Monitoring
type SyntheticMonitoringHandler struct {
	Provider Provider
}

// NewSyntheticMonitoringHandler returns a Grizzly Handler for Grafana Synthetic Monitoring
func NewSyntheticMonitoringHandler(provider Provider) *SyntheticMonitoringHandler {
	return &SyntheticMonitoringHandler{
		Provider: provider,
	}
}

// Kind returns the name for this handler
func (h *SyntheticMonitoringHandler) Kind() string {
	return "SyntheticMonitoringCheck"
}

// APIVersion returns the group and version for the provider of which this handler is a part
func (h *SyntheticMonitoringHandler) APIVersion() string {
	return h.Provider.APIVersion()
}

// GetExtension returns the file name extension for a check
func (h *SyntheticMonitoringHandler) GetExtension() string {
	return "json"
}

// Parse parses a manifest object into a struct for this resource type
func (h *SyntheticMonitoringHandler) Parse(m manifest.Manifest) (grizzly.ResourceList, error) {
	resource := grizzly.Resource(m)
	resource.SetSpecString("job", resource.GetMetadata("name"))
	return resource.AsResourceList(), nil
}

// Unprepare removes unnecessary elements from a remote resource ready for presentation/comparison
func (h *SyntheticMonitoringHandler) Unprepare(resource grizzly.Resource) *grizzly.Resource {
	resource.DeleteSpecKey("tenantId")
	resource.DeleteSpecKey("id")
	resource.DeleteSpecKey("modified")
	resource.DeleteSpecKey("created")
	return &resource
}

// Prepare gets a resource ready for dispatch to the remote endpoint
func (h *SyntheticMonitoringHandler) Prepare(existing, resource grizzly.Resource) *grizzly.Resource {
	resource.SetSpecString("tenantId", existing.GetSpecString("tenantId"))
	resource.SetSpecString("id", existing.GetSpecString("id"))
	return &resource
}

// GetByUID retrieves JSON for a resource from an endpoint, by UID
func (h *SyntheticMonitoringHandler) GetByUID(UID string) (*grizzly.Resource, error) {
	return getRemoteCheck(UID)
}

// GetRemote retrieves a datasource as a Resource
func (h *SyntheticMonitoringHandler) GetRemote(resource grizzly.Resource) (*grizzly.Resource, error) {
	uid := fmt.Sprintf("%s.%s", resource.GetMetadata("type"), resource.Name())
	return getRemoteCheck(uid)
}

// Add adds a new check to the SyntheticMonitoring endpoint
func (h *SyntheticMonitoringHandler) Add(resource grizzly.Resource) error {
	url := getSyntheticMonitoringURL("api/v1/check/add")
	return postCheck(url, resource)
}

// Update pushes an updated check to the SyntheticMonitoring endpoing
func (h *SyntheticMonitoringHandler) Update(existing, resource grizzly.Resource) error {
	url := getSyntheticMonitoringURL("api/v1/check/update")
	return postCheck(url, resource)
}
