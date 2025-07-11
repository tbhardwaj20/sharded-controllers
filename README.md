## Blockers :-
Controller pod is not serving the /healthz and /readyz endpoints on port 8081, which causes both the liveness and readiness probes to fail. Kubernetes interprets this as the pod being unhealthy and repeatedly restarts it, leading to the CrashLoopBackOff state.
