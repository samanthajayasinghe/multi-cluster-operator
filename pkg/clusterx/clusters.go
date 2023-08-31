package clusterx

import (
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MultiClusterClient struct {
	client  client.Client
	clients map[string]client.Client
}

const (
	ClusterSeparator   = "/"
	DefaultClusterName = "default"
)

// GetClient by cluster name
func (mc MultiClusterClient) GetClient(name string) client.Client {
	if c, ok := mc.clients[name]; ok {
		return c
	}

	return mc.client
}

func (mc MultiClusterClient) GetClientByNs(ns string) client.Client {
	name, _ := GetClusterNameNs(ns)
	return mc.GetClient(name)
}

func (mc MultiClusterClient) GetOwnerClient(obj client.Object) client.Client {
	return mc.GetClient(GetOwnerClusterName(obj))
}

func (mc MultiClusterClient) GetClientByObj(obj client.Object) client.Client {
	return mc.GetClient(GetClusterName(obj))
}

func GetClusterNameNs(ns string) (clusterName, namespace string) {
	ss := strings.Split(ns, ClusterSeparator)
	if len(ss) < 2 {
		return DefaultClusterName, ns
	}

	return ss[0], ss[1]
}
