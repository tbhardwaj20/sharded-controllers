package main

import (
    "fmt"
    "os"

    shardlease "github.com/timebertt/kubernetes-controller-sharding/pkg/shard/lease"
    shardingv1alpha1 "github.com/timebertt/kubernetes-controller-sharding/pkg/apis/sharding/v1alpha1"
    ctrl "sigs.k8s.io/controller-runtime"

    "k8s.io/apimachinery/pkg/labels"
    "sigs.k8s.io/controller-runtime/pkg/cache"
    "sigs.k8s.io/controller-runtime/pkg/client/config"
    "sigs.k8s.io/controller-runtime/pkg/manager"

    "github.com/tusharbhardwaj/go-sharded-controller/controller"
)

func main() {
    ringName := "go-controller-ring"
    shardName := os.Getenv("POD_NAME") // Will match Lease name in deployment

    restConfig := config.GetConfigOrDie()

    lease, err := shardlease.NewResourceLock(restConfig, shardlease.Options{
        ControllerRingName: ringName,
    })
    if err != nil {
        panic(fmt.Errorf("failed to create lease: %w", err))
    }

    mgr, err := manager.New(restConfig, manager.Options{
        LeaderElection:                      true,
        LeaderElectionResourceLockInterface: lease,
        LeaderElectionReleaseOnCancel:       true,
        Cache: cache.Options{
            DefaultLabelSelector: labels.SelectorFromSet(labels.Set{
                shardingv1alpha1.LabelShard(ringName): shardName,
            }),
        },
    })
    if err != nil {
        panic(fmt.Errorf("failed to create manager: %w", err))
    }

    if err := controller.SetupPodController(mgr, ringName, shardName); err != nil {
        panic(fmt.Errorf("failed to set up pod controller: %w", err))
    }

    fmt.Println("Starting the sharded controller...")
    if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
        panic(fmt.Errorf("error running manager: %w", err))
    }
}
