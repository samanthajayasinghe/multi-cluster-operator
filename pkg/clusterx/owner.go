package clusterx

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	OwnerNameKey      = "owner.name"
	OwnerNamespaceKey = "owner.namespace"
	OwnerClusterKey   = "owner.cluster"
)

const ClusterNameKey = "cluster.name"

func SetOwner(owner, sub metav1.Object) {
	labels := sub.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}
	labels[OwnerNamespaceKey] = owner.GetNamespace()
	labels[OwnerNameKey] = owner.GetName()
	labels[OwnerClusterKey] = GetClusterName(owner)
	sub.SetLabels(labels)
}

func GetOwnerNameNs(sub metav1.Object) types.NamespacedName {
	labels := sub.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}
	return types.NamespacedName{
		Namespace: labels[OwnerNamespaceKey],
		Name:      labels[OwnerNameKey],
	}
}

func GetOwnerClusterName(sub metav1.Object) string {
	return sub.GetLabels()[OwnerClusterKey]
}

func GetClusterName(obj metav1.Object) string {
	return obj.GetAnnotations()[ClusterNameKey]
}

func SetClusterName(obj metav1.Object, cluster string) {
	data := obj.GetAnnotations()
	data[ClusterNameKey] = cluster
}
