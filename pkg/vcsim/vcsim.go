// Package vcsim is used to mock interactions with a vCenter server
package vcsim

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/vmware/govmomi/simulator"
	_ "github.com/vmware/govmomi/vapi/simulator" // Importing the simulator package to enable simulation of vCenter server

	"github.com/validator-labs/validator-plugin-vsphere/api/vcenter"
)

var sig chan os.Signal

func init() {
	// simulator.Trace = true
}

// VCSimulator is used to mock interactions with a vCenter server
type VCSimulator struct {
	Account vcenter.Account
	Options govcOptions
	log     logr.Logger
}

type govcOptions struct {
	InsecureConnection string
	RefreshRestClient  bool

	Cluster                     string
	Datacenter                  string
	Datastore                   string
	DistributedVirtualPortgroup string
	DistributedVirtualSwitch    string
	Folder                      string
	Host                        string
	Network                     network
	ResourcePool                string
	Portgroup                   string
	VM                          string
}

type network struct {
	Name    string
	Network string
}

// NewVCSim creates a new VCSimulator
func NewVCSim(username string, port int, log logr.Logger) *VCSimulator {
	return &VCSimulator{
		Account: NewTestVsphereAccount(username, port),
		log:     log,
	}
}

// NewTestVsphereAccount creates a new vsphere account for testing
func NewTestVsphereAccount(username string, port int) vcenter.Account {
	// Starting & stopping vcsim between test cases appears to work, but govmomi calls
	// throw an auth error on the 2nd iteration unless a unique username is used
	// each time the simulator is instantiated.
	return vcenter.Account{
		Insecure: true,
		Password: "welcome123",
		Username: username,
		Host:     fmt.Sprintf("127.0.0.1:%d", port),
	}
}

// Start starts the mock vcsim server
func (v *VCSimulator) Start() {
	model := simulator.VPX()

	model.Autostart = false
	model.Cluster = 2
	model.ClusterHost = 1
	model.Datacenter = 1
	model.Datastore = 2
	model.Host = 1
	model.Machine = 1
	model.Pool = 2
	model.Portgroup = 1
	model.ServiceContent.About.ApiVersion = "6.7.3"

	v.Options = govcOptions{
		InsecureConnection:          strconv.FormatBool(true),
		Datacenter:                  "DC0",
		DistributedVirtualPortgroup: "DC0_DVPG0",
		DistributedVirtualSwitch:    "DVS0",
		Cluster:                     "DC0_C0",
		ResourcePool:                "DC0_C0_RP0",
		VM:                          "DC0_C0_RP0_VM0",
		Datastore:                   "LocalDS_0",
		Folder:                      "SC_Tyler",
		Network: network{
			Name:    "VM Network",
			Network: "VM Network",
		},
	}

	cleanUp, err := v.createVCenterSimulator(model)
	if err != nil {
		log.Fatalf("failed to create vCenter simulator: %s", err)
	}

	sig = make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer cleanUp()
		<-sig
	}()

	log.Println("started vcsim server")
}

// Shutdown shuts down the vcim mock server
func (v *VCSimulator) Shutdown() {
	log.Println("shutting down vcsim server")
	sig <- syscall.SIGTERM
}

// createvCenterSimulator creates a vCenter simulator
func (v *VCSimulator) createVCenterSimulator(model *simulator.Model) (func(), error) {
	if model == nil {
		model = simulator.VPX()
	}

	err := model.Create()
	if err != nil {
		return func() {}, errors.Wrap(err, "Error while creating model")
	}

	host := v.Account.Host
	if _, err = url.Parse(fmt.Sprintf("https://%s/sdk", host)); err != nil {
		return nil, errors.Errorf("invalid vCenter server")
	}

	model.Service.RegisterEndpoints = true
	model.Service.Listen = &url.URL{
		User: v.Account.Userinfo(),
		Host: host,
	}
	model.Service.TLS = new(tls.Config)

	s := model.Service.NewServer()

	return func() {
		s.Close()
	}, nil
}
