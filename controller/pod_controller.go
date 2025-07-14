package controller

import (
    "context"
    corev1 "k8s.io/api/core/v1"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/builder"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/log"

    shardcontroller "github.com/timebertt/kubernetes-controller-sharding/pkg/shard/controller"
)

// PodReconciler is our controller's logic handler
type PodReconciler struct {
    client.Client
}

// Reconcile is triggered for pods matching the assigned shard
func (r *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    logger := log.FromContext(ctx)

    var pod corev1.Pod
    if err := r.Get(ctx, req.NamespacedName, &pod); err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    logger.Info("Reconciling Pod", "name", pod.Name, "namespace", pod.Namespace)

    // You can place your business logic here
    // Example: Just log the pod labels
    logger.Info("Pod Labels", "labels", pod.Labels)

    return ctrl.Result{}, nil
}

// SetupPodController sets up the controller with drain and shard support
func SetupPodController(mgr ctrl.Manager, ringName, shardName string) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&corev1.Pod{}, builder.WithPredicates(
            shardcontroller.Predicate(ringName, shardName),
        )).
        Complete(
            shardcontroller.NewShardedReconciler(mgr).
                For(&corev1.Pod{}).
                InControllerRing(ringName).
                WithShardName(shardName).
                MustBuild(&PodReconciler{Client: mgr.GetClient()}),
        )
}

