package clustersync

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	hiveinternal "github.com/openshift/hive/apis/hiveinternal/v1alpha1"
)

var log = logf.Log.WithName("controller_clustersync")

// Add creates a new ClusterSync Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileClusterSync{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("clustersync-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ClusterSync
	err = c.Watch(&source.Kind{Type: &hiveinternal.ClusterSync{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileClusterSync implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileClusterSync{}

// ReconcileClusterSync reconciles a ClusterSync object
type ReconcileClusterSync struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile analyzes ClusterSync objects in all namespaces and reports on failures.
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileClusterSync) Reconcile(context context.Context, request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ClusterSync")

	// Fetch the ClusterSync instance
	instance := &hiveinternal.ClusterSync{}
	err := r.client.Get(context, request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, fmt.Sprintf("failed to retrieve ClusterSync %s/%s", request.Namespace, request.Name))
		return reconcile.Result{}, err
	}

	if err := countErrors(instance, reqLogger); err != nil {
		reqLogger.Error(err, "error counting errors")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
