/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/log"

	multiclusterv1alpha1 "github.com/samanthajayasinghe/multi-cluster-operator/api/v1alpha1"
)

// LogForwarderReconciler reconciles a LogForwarder object
type LogForwarderReconciler struct {
	client.Client
	Clients map[string]client.Client
	Scheme  *runtime.Scheme
}

// NewLogForwarderReconcilerReconciler ...
func NewLogForwarderReconciler(mgr ctrl.Manager, clusters map[string]cluster.Cluster) (*LogForwarderReconciler, error) {
	r := LogForwarderReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Clients: map[string]client.Client{
			"main": mgr.GetClient(),
		},
	}
	for name, cluster := range clusters {
		r.Clients[name] = cluster.GetClient()
	}

	err := r.SetupWithManager(mgr)
	return &r, err
}

//+kubebuilder:rbac:groups=multicluster.samanthajayasinghe.github.io,resources=logforwarders,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=multicluster.samanthajayasinghe.github.io,resources=logforwarders/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=multicluster.samanthajayasinghe.github.io,resources=logforwarders/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LogForwarder object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *LogForwarderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var forwarder multiclusterv1alpha1.LogForwarder
	var res ctrl.Result

	err := r.Get(ctx, req.NamespacedName, &forwarder)
	if err != nil {
		return res, client.IgnoreNotFound(err)
	}

	job := forwarder.Job()

	for _, c := range r.Clients {
		j := job.DeepCopy()
		err := c.Get(ctx, client.ObjectKeyFromObject(&job), j)
		if client.IgnoreNotFound(err) != nil {
			return res, err
		}

		if err == nil {
			continue
		}

		err = c.Create(ctx, job.DeepCopy())
		if err != nil {
			return res, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LogForwarderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&multiclusterv1alpha1.LogForwarder{}).
		Complete(r)
}
