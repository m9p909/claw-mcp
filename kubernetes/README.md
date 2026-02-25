# Kubernetes Deployment Guide

This directory contains Kubernetes manifests for deploying the Claw MCP Server in a Kubernetes cluster.

## Prerequisites

- A Kubernetes cluster (1.20+)
- kubectl configured to access your cluster
- The Claw Docker image built and available in your image registry

## Files

- `statefulset.yaml` - StatefulSet defining a single Claw server pod with persistent storage
- `service.yaml` - Service exposing the Claw server on port 8080

## Deployment

### 1. Deploy the StatefulSet and Service

```bash
kubectl apply -f statefulset.yaml
kubectl apply -f service.yaml
```

### 2. Verify Deployment

```bash
# Check StatefulSet status
kubectl get statefulsets

# Check pods
kubectl get pods -l app=claw

# Check service
kubectl get service claw-service

# View pod details
kubectl describe pod claw-0
```

### 3. Access the Server

#### Within the cluster:

```bash
kubectl exec -it claw-0 -- curl http://localhost:8080/health
```

#### From outside the cluster (port-forward):

```bash
kubectl port-forward svc/claw-service 8080:8080
curl http://localhost:8080/health
```

#### Via Ingress (optional):

Create an Ingress resource to expose the service externally.

## Persistent Storage

The StatefulSet uses a PersistentVolumeClaim (10Gi) mounted at `/root/.mcpclaw` containing:

- `workspace/` - Shared workspace for agents
- `data/` - SQLite database for memory persistence

Storage is bound to the StatefulSet and persists across pod restarts.

## Configuration

To override the server port, modify the `PORT` environment variable in `statefulset.yaml`:

```yaml
env:
- name: PORT
  value: "9000"  # Change from 8080 to 9000
```

## Troubleshooting

### Pod not starting

```bash
kubectl logs claw-0
kubectl describe pod claw-0
```

### PVC not binding

```bash
kubectl get pvc
kubectl describe pvc mcpclaw-claw-0
```

### Connection issues

```bash
# Test from within pod
kubectl exec -it claw-0 -- curl http://localhost:8080/health

# Check service DNS
kubectl exec -it claw-0 -- nslookup claw-service
```

## Cleanup

```bash
kubectl delete statefulset claw
kubectl delete service claw-service
```
