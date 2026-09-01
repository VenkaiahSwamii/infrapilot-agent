# InfraPilot Enterprise - Production Deployment Order

## Prerequisites
- Kubernetes cluster running (v1.24+)
- kubectl configured with cluster access
- NGINX Ingress Controller installed
- cert-manager installed (for TLS certificates)
- Default StorageClass configured

## Step 1: Namespace
```bash
kubectl apply -f namespace.yaml
```

## Step 2: Secrets
```bash
kubectl apply -f secrets.yaml
```
**Important:** Replace placeholder values in `secrets.yaml` with secure production values before applying.

## Step 3: ConfigMap
```bash
kubectl apply -f configmap.yaml
```

## Step 4: PostgreSQL StatefulSet
```bash
kubectl apply -f postgres/statefulset.yaml
```

## Step 5: Redis StatefulSet
```bash
kubectl apply -f redis/deployment.yaml
```

## Step 6: Qdrant StatefulSet
```bash
kubectl apply -f qdrant/statefulset.yaml
```

## Step 7: Backend Deployment
```bash
kubectl apply -f backend/deployment.yaml
kubectl apply -f backend/service.yaml
```

## Step 8: Frontend Deployment
```bash
kubectl apply -f frontend/deployment.yaml
kubectl apply -f frontend/service.yaml
```

## Step 9: Ingress
```bash
kubectl apply -f ingress/ingress.yaml
```

## Step 10: Autoscaling
```bash
kubectl apply -f autoscaling/hpa.yaml
```

## Step 11: Network Policies
```bash
kubectl apply -f network/network-policies.yaml
```

## Step 12: Resource Quotas
```bash
kubectl apply -f resource-quotas/quota.yaml
```

## Validation

After deployment, validate all resources:
```bash
# Check all pods are running
kubectl get pods -n infrapilot

# Check services
kubectl get svc -n infrapilot

# Check ingress
kubectl get ingress -n infrapilot

# Check PVCs are bound
kubectl get pvc -n infrapilot

# Check HPA is active
kubectl get hpa -n infrapilot

# Check for any issues
kubectl get events -n infrapilot --sort-by='.lastTimestamp'
```

## Expected State
- All Pods: Running
- PVC: Bound
- Ingress: Ready (with address)
- HPA: Active
- No CrashLoopBackOff
- No Pending Pods

## Troubleshooting

### PVC Pending
Check StorageClass exists:
```bash
kubectl get storageclass
```

### Ingress No Address
Ensure NGINX Ingress Controller is installed:
```bash
kubectl get pods -n ingress-nginx
```

### Pods CrashLoopBackOff
Check logs:
```bash
kubectl logs <pod-name> -n infrapilot --previous
```

## Rollback
To rollback any deployment:
```bash
kubectl rollout undo deployment/<deployment-name> -n infrapilot