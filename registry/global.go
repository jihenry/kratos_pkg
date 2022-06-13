package registry

import (
	"sync"

	"github.com/go-kratos/kratos/v2/registry"
)

var (
	pvDiscovery registry.Discovery
	pv          sync.Mutex
)

func SetDiscovery(r registry.Discovery) {
	pv.Lock()
	pvDiscovery = r
	pv.Unlock()
}

func GetDiscovery() registry.Discovery {
	return pvDiscovery
}
