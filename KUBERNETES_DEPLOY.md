# Kubernetes Deployment Guide

This guide covers deploying the Coolmate eCommerce backend to a Kubernetes cluster.

## Prerequisites

- Kubernetes cluster (v1.24+) with `kubectl` configured
- Docker image built and pushed to a container registry (e.g., ghcr.io)
- PostgreSQL database (can use the included StatefulSet)
- Redis (can use external or included in docker-compose)
- S3-compatible storage (AWS S3, MinIO, etc.)

## Quick Start

### 1. Create Namespace

```bash
kubectl apply -f k8s/namespace.yaml
```

### 2. Create ConfigMap and Secrets

```bash
# Apply ConfigMap (non-sensitive config)
kubectl apply -f k8s/configmap.yaml

# Update the Secret values before applying in production!
# Edit k8s/configmap.yaml to set actual secret values
kubectl apply -f k8s/configmap.yaml
```

### 3. Deploy PostgreSQL

```bash
kubectl apply -f k8s/postgres-statefulset.yaml
```

Verify PostgreSQL is running:
```bash
kubectl get statefulset -n coolmate
kubectl get pvc -n coolmate
```

### 4. Deploy API

Update the image reference in `k8s/api-deployment.yaml`:

```yaml
spec:
  template:
    spec:
      containers:
        - name: api
          image: ghcr.io/your-org/coolmate-backend:main  # Update this
```

Deploy:
```bash
kubectl apply -f k8s/api-deployment.yaml
```

Verify deployment:
```bash
kubectl get deployment -n coolmate
kubectl get pods -n coolmate
kubectl logs -n coolmate deployment/coolmate-api
```

## Verification

### Check Service Endpoints

```bash
# Get service IPs
kubectl get svc -n coolmate

# Port forward to test locally
kubectl port-forward -n coolmate svc/coolmate-api-service 8080:80
```

### Test Health Endpoint

```bash
curl http://localhost:8080/health
# Expected response:
# {"status":"healthy"}
```

### Monitor Metrics

```bash
# Port forward metrics
kubectl port-forward -n coolmate svc/coolmate-api-service 9090:9090

# Scrape metrics
curl http://localhost:9090/metrics
```

## Scaling

The deployment includes HorizontalPodAutoscaler that automatically scales based on:
- CPU utilization > 70%
- Memory utilization > 80%

Current ranges: 3-10 replicas

View autoscaler status:
```bash
kubectl get hpa -n coolmate -w
```

## Monitoring

### Prometheus Integration

The metrics endpoint at `:9090/metrics` exposes:
- `coolmate_uptime_seconds` - Service uptime
- `process_resident_memory_bytes` - Memory usage
- `go_goroutines` - Active goroutines

Add to Prometheus scrape config:
```yaml
scrape_configs:
  - job_name: 'coolmate-api'
    kubernetes_sd_configs:
      - role: pod
        namespaces:
          names:
            - coolmate
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_label_app]
        action: keep
        regex: coolmate-api
      - source_labels: [__meta_kubernetes_pod_container_port_number]
        action: keep
        regex: "9090"
```

### Logs

View pod logs:
```bash
kubectl logs -n coolmate deployment/coolmate-api -f

# View PostgreSQL logs
kubectl logs -n coolmate statefulset/postgres -f
```

## Database Backups

PostgreSQL includes a daily backup CronJob at 2 AM UTC.

View backup jobs:
```bash
kubectl get cronjob -n coolmate
kubectl get job -n coolmate
```

## Upgrading

### Rolling Update

The deployment uses RollingUpdate strategy:
- maxSurge: 1 (one extra pod during update)
- maxUnavailable: 0 (no downtime)

Update image:
```bash
kubectl set image deployment/coolmate-api \
  -n coolmate \
  api=ghcr.io/your-org/coolmate-backend:v2.0
```

Monitor rollout:
```bash
kubectl rollout status deployment/coolmate-api -n coolmate
```

Rollback if needed:
```bash
kubectl rollout undo deployment/coolmate-api -n coolmate
```

## Security

The deployment includes:
- **Non-root user** - Runs as UID 1000
- **Read-only filesystem** - Except /tmp and /app/cache volumes
- **No privilege escalation** - Capabilities dropped
- **Pod disruption budget** - Ensures minimum 2 replicas during cluster operations
- **Pod anti-affinity** - Spreads pods across different nodes

## Troubleshooting

### Pod not starting

Check pod status and events:
```bash
kubectl describe pod -n coolmate <pod-name>
kubectl logs -n coolmate <pod-name>
```

Common issues:
- **ImagePullBackOff** - Check container registry credentials
- **Pending** - Check resource requests vs available cluster resources
- **CrashLoopBackOff** - Check application logs

### Database connectivity

Test from pod:
```bash
kubectl exec -it -n coolmate deployment/coolmate-api -- sh
# Inside pod:
psql -h postgres-service -U coolmate_prod_user -d coolmate_ecommerce
```

### Performance issues

Check resource usage:
```bash
kubectl top pods -n coolmate
kubectl top nodes
```

## Production Considerations

1. **Resource Limits** - Adjust based on actual usage
   - Current: 256Mi memory request, 512Mi limit
   - Current: 100m CPU request, 500m limit

2. **Secrets Management** - Use external secrets manager:
   - HashiCorp Vault
   - AWS Secrets Manager
   - Google Secret Manager

3. **Persistent Storage** - Ensure adequate storage class
   - Current: 50Gi for PostgreSQL (fast-ssd)
   - Adjust storageClassName if needed

4. **Network Policies** - Add network policies to restrict traffic

5. **Ingress** - Set up ingress controller for external access
   ```yaml
   apiVersion: networking.k8s.io/v1
   kind: Ingress
   metadata:
     name: coolmate-api-ingress
     namespace: coolmate
   spec:
     ingressClassName: nginx
     rules:
       - host: api.coolmate.com
         http:
           paths:
             - path: /
               pathType: Prefix
               backend:
                 service:
                   name: coolmate-api-service
                   port:
                     number: 80
   ```

## Clean Up

Remove all resources:
```bash
kubectl delete namespace coolmate
```

This will delete all resources in the namespace.
