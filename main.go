package azservicebus

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/grafana/sobek"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/azservicebus", new(RootModule))
}

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create k6/x/nats module instances for each VU.
type RootModule struct{}

// ModuleInstance represents an instance of the module for every VU.
type ServiceBus struct {
	timeout time.Duration
	cli     *azservicebus.Client
	vu      modules.VU
	exports map[string]interface{}
}

// Configuration represents the configuration for the module.
type Configuration struct {
	ConnectionString   string `js:"connectionString"`
	Timeout            int64  `js:"timeout"`
	InsecureSkipVerify bool   `js:"insecureSkipVerify"`
}

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Instance = &ServiceBus{}
	_ modules.Module   = &RootModule{}
)

// NewModuleInstance implements the modules.Module interface and returns
// a new instance for each VU.
func (r *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	sb := &ServiceBus{
		vu:      vu,
		exports: make(map[string]interface{}),
	}

	sb.exports["ServiceBus"] = sb.client

	return sb
}

// Exports implements the modules.Instance interface and returns the exports
// of the JS module.
func (sb *ServiceBus) Exports() modules.Exports {
	return modules.Exports{
		Named: sb.exports,
	}
}

func (sb *ServiceBus) client(c sobek.ConstructorCall) *sobek.Object {
	rt := sb.vu.Runtime()

	var cfg Configuration
	err := rt.ExportTo(c.Argument(0), &cfg)
	if err != nil {
		common.Throw(rt, fmt.Errorf("AzServicebus constructor expect Configuration as it's argument: %w", err))
	}

	var op azservicebus.ClientOptions
	if cfg.InsecureSkipVerify {
		op.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	client, err := azservicebus.NewClientFromConnectionString(cfg.ConnectionString, &op)
	if err != nil {
		common.Throw(rt, err)
	}

	return rt.ToValue(&ServiceBus{
		timeout: time.Duration(cfg.Timeout) * time.Millisecond,
		vu:      sb.vu,
		cli:     client,
	}).ToObject(rt)
}

func (sb *ServiceBus) Close() {
	if sb.cli != nil {
		sb.cli.Close(context.Background())
	}
}
