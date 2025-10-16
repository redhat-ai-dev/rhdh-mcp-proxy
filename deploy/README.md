# RHDH MCP Proxy - OpenShift Deployment

This directory contains OpenShift manifests for deploying the RHDH MCP Proxy application using Kustomize.

## Prerequisites

- OpenShift cluster (4.8+)
- oc CLI configured to access your cluster
- kustomize (v4.0+) or kubectl with kustomize support

## Quick Start

1. **Create Secret Cluster-side**: The secret will be created cluster-side with your actual values:
   ```bash
   # Create the secret with your actual values
   oc create secret generic rhdh-mcp-proxy-secret \
     --from-literal=mcp-token="your-actual-mcp-token" \
     --from-literal=backstage-url="https://your-backstage-instance.com"
   ```

2. **Deploy the Application**:
   ```bash
   oc apply -k .
   ```

3. **Verify Deployment**:
   ```bash
   oc get pods -n rhdh-mcp-proxy
   oc get svc -n rhdh-mcp-proxy
   ```

## Configuration

### Environment Variables

The application requires the following environment variables:

- `BACKSTAGE_URL`: The URL of your Backstage instance
- `MCP_TOKEN`: The MCP authentication token
- `PORT`: The port to run the application on (default: 8080)

### Secrets

Create the secret cluster-side with your actual MCP token and Backstage URL:

```bash
# Create the secret directly in the cluster
oc create secret generic rhdh-mcp-proxy-secret \
  --from-literal=mcp-token="your-actual-mcp-token" \
  --from-literal=backstage-url="https://your-backstage-instance.com"
```

**Note**: The secret will be created cluster-side rather than through GitOps for security reasons.

## Resources Included

- **ServiceAccount**: Basic service account for the application pods
- **Deployment**: Runs 1 replica of the proxy application
- **Service**: ClusterIP service exposing the application on port 8080

**Note**: Namespace, Secret, ConfigMap, and Route will be created cluster-side rather than through GitOps.

## Security Features

- Non-root user (UID 1001)
- Read-only root filesystem
- Dropped capabilities
- Resource limits and requests
- TLS reencrypt for secure communication

## Customization

### Scaling

To change the number of replicas, update the `kustomization.yaml` file:

```yaml
patchesStrategicMerge:
  - |-
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: rhdh-mcp-proxy
    spec:
      replicas: 3  # Change this number
```

### Image Tag

To use a specific image tag, update the `kustomization.yaml` file:

```yaml
images:
  - name: quay.io/redhat-ai-dev/rhdh-mcp-proxy
    newTag: v1.0.0  # Change this tag
```

## Troubleshooting

### Check Pod Status
```bash
oc describe pod -n rhdh-mcp-proxy -l app=rhdh-mcp-proxy
```

### View Logs
```bash
oc logs -n rhdh-mcp-proxy -l app=rhdh-mcp-proxy
```

### Test Connectivity
```bash
oc port-forward -n rhdh-mcp-proxy svc/rhdh-mcp-proxy 8080:8080
curl http://localhost:8080/api/mcp-actions/
```


## Cleanup

To remove the deployment:

```bash
oc delete -k .
```
