package main

import (
    "os"
	
    "sigs.k8s.io/controller-runtime/pkg/healthz"
    appsv1 "k8s.io/api/apps/v1"
    "k8s.io/apimachinery/pkg/labels"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/builder"
    "sigs.k8s.io/controller-runtime/pkg/cache"
    "sigs.k8s.io/controller-runtime/pkg/client/config"
    "sigs.k8s.io/controller-runtime/pkg/log"
    "sigs.k8s.io/controller-runtime/pkg/log/zap"

    shardingv1alpha1 "github.com/timebertt/kubernetes-controller-sharding/pkg/apis/sharding/v1alpha1"
    shardlease "github.com/timebertt/kubernetes-controller-sharding/pkg/shard/lease"

    "my-sharded-controller/controller"
)

func main() {
    // ✅ Initialize logger
    log.SetLogger(zap.New(zap.UseDevMode(true)))
    logger := ctrl.Log.WithName("main")

    restConfig := config.GetConfigOrDie()

    // ✅ Create shard lease lock
    shardLease, err := shardlease.NewResourceLock(restConfig, shardlease.Options{
        ControllerRingName: "my-controller-ring", // must match your ControllerRing name
    })
    if err != nil {
        logger.Error(err, "failed to create shard lease")
        os.Exit(1)
    }

    // ✅ Create manager with shard lease and filtered cache
    mgr, err := ctrl.NewManager(restConfig, ctrl.Options{
        LeaderElection:                      true,
        LeaderElectionResourceLockInterface: shardLease,
        LeaderElectionReleaseOnCancel:       true,
        Cache: cache.Options{
            DefaultLabelSelector: labels.SelectorFromSet(labels.Set{
                shardingv1alpha1.LabelShard("my-controller-ring"): shardLease.Identity(),
            }),
        },
    })
    if err != nil {
        logger.Error(err, "failed to create manager")
        os.Exit(1)
    }
	
if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
	 logger.Error(err, "unable to set up health check")
 os.Exit(1)
}
if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
logger.Error(err, "unable to set up ready check")
 os.Exit(1)
}


    // ✅ Register the Deployment controller
    if err := builder.ControllerManagedBy(mgr).
        For(&appsv1.Deployment{}).
        Complete(&controller.DeploymentReconciler{
            Client: mgr.GetClient(),
        }); err != nil {
        logger.Error(err, "failed to register controller")
        os.Exit(1)
    }

    // ✅ Start the manager
    if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
        logger.Error(err, "problem running manager")
        os.Exit(1)
    }
}
