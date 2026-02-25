## ADDED Requirements

### Requirement: Kubernetes StatefulSet deploys Claw server
Claw SHALL be deployed as a Kubernetes StatefulSet with a single replica, providing stable pod identity and persistent storage binding.

#### Scenario: StatefulSet maintains stable identity
- **WHEN** the pod restarts
- **THEN** the same workspace and database volumes are reattached to the pod

### Requirement: Persistent Claw data volume
The deployment SHALL mount a PersistentVolume at `~/.mcpclaw` within the container, providing persistent storage for both the shared workspace (`~/.mcpclaw/workspace`) and SQLite database (`~/.mcpclaw/data`).

#### Scenario: Agents share workspace files
- **WHEN** multiple agents connect to the server simultaneously
- **THEN** all agents can read and write files to the same `~/.mcpclaw/workspace` directory

#### Scenario: Memory data survives pod restart
- **WHEN** an agent writes memory using write_memory tool
- **THEN** the data persists in `~/.mcpclaw/data` after pod termination and restart

### Requirement: Service exposes server endpoint
A Kubernetes Service SHALL expose the Claw server on port 8080, making it accessible to agents within the cluster or via ingress.

#### Scenario: Agent connects via DNS
- **WHEN** an agent connects to the service DNS name (e.g., `claw-service:8080`)
- **THEN** the connection reaches the running server pod

### Requirement: PersistentVolumeClaim defines storage requirements
The deployment configuration SHALL declare a single PersistentVolumeClaim for `~/.mcpclaw` with sufficient capacity for both agent workspace and SQLite database (e.g., 10Gi or configurable).

#### Scenario: PVC provisioning creates storage
- **WHEN** the StatefulSet is created
- **THEN** the Kubernetes storage provisioner allocates and binds physical storage to the PVC

### Requirement: Environment configuration via ConfigMap/Secrets
Configuration parameters (e.g., port, database location) SHALL be defined in Kubernetes ConfigMap or Secrets, allowing deployment-time customization without rebuilding the image.

#### Scenario: Deployment uses environment variable override
- **WHEN** a ConfigMap is created with `MCP_PORT=9000`
- **THEN** the server listens on port 9000 instead of 8080
