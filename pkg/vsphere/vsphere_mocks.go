package vsphere

import (
	"context"
	"strings"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/vim25/mo"

	"github.com/validator-labs/validator-plugin-vsphere/api/vcenter"
)

// MockVsphereDriver is a mock implementation of the Driver interface
type MockVsphereDriver struct {
	Datacenters        []string
	Clusters           []string
	VMs                []vcenter.VM
	VMFolders          []string
	HostSystems        map[string][]vcenter.HostSystem
	VApps              []mo.VirtualApp
	ResourcePools      []*object.ResourcePool
	HostClusterMapping map[string]string
	ResourceTags       map[string]tags.AttachedTags
}

// ensure that MockVsphereDriver implements the Driver interface
var _ Driver = &MockVsphereDriver{}

// GetVMFolders returns a mocked response
func (d MockVsphereDriver) GetVMFolders(_ context.Context, _ string) ([]string, error) {
	return d.VMFolders, nil
}

// GetK8sDatacenters returns a mocked response
func (d MockVsphereDriver) GetK8sDatacenters(_ context.Context) ([]string, error) {
	return d.Datacenters, nil
}

// GetK8sClusters returns a mocked response
func (d MockVsphereDriver) GetK8sClusters(_ context.Context, _ string) ([]string, error) {
	return d.Clusters, nil
}

// GetHostSystems returns a mocked response
func (d MockVsphereDriver) GetHostSystems(_ context.Context, datacenter, cluster string) ([]vcenter.HostSystem, error) {
	return d.HostSystems[concat(datacenter, cluster)], nil
}

// ValidateCredentials returns a mocked response
func (d MockVsphereDriver) ValidateCredentials() (bool, error) {
	return true, nil
}

// ValidateVersion returns a mocked response
func (d MockVsphereDriver) ValidateVersion(_ string) error {
	return nil
}

// GetHostClusterMapping returns a mocked response
func (d MockVsphereDriver) GetHostClusterMapping(_ context.Context) (map[string]string, error) {
	return d.HostClusterMapping, nil
}

// GetVMs returns a mocked response
func (d MockVsphereDriver) GetVMs(_ context.Context, _ string) ([]vcenter.VM, error) {
	return d.VMs, nil
}

// GetResourcePools returns a mocked response
func (d MockVsphereDriver) GetResourcePools(_ context.Context, _ string, _ string) ([]*object.ResourcePool, error) {
	return d.ResourcePools, nil
}

// GetVApps returns a mocked response
func (d MockVsphereDriver) GetVApps(_ context.Context) ([]mo.VirtualApp, error) {
	return d.VApps, nil
}

// GetResourceTags returns a mocked response
func (d MockVsphereDriver) GetResourceTags(_ context.Context, _ string) (map[string]tags.AttachedTags, error) {
	return d.ResourceTags, nil
}

func concat(ss ...string) string {
	return strings.Join(ss, "_")
}
