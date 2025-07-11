package controller

import (
    "context"
    "fmt"

    // appsv1 "k8s.io/api/apps/v1"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type DeploymentReconciler struct {
    client.Client
}

func (r *DeploymentReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
    fmt.Printf("Reconciling %s\n", req.NamespacedName)
    return reconcile.Result{}, nil
}
