# Local infrastructure

Local environment for GoJudge: a Kubernetes cluster (k3d) with RabbitMQ and
PostgreSQL in the namespace `gojudge`.

| Part | Where | Details |
|------|-------|---------|
| Cluster | `infra/k3d/k3d-config.yaml` | k3d cluster `gojudge`, k3s `v1.35.5` |
| Namespace | `infra/k8s/namespace.yaml` | `gojudge`, shared by all local services |
| RabbitMQ 4.4.0 | `infra/k8s/rabbitmq/` | [RabbitMQ README](/infra/k8s/rabbitmq/README.md) |
| PostgreSQL 17.11 | `infra/k8s/postgres/` | [PostgreSQL README](/infra/k8s/postgres/README.md) |

Run all commands from the repository root.

## You need
- Docker: runs the cluster. [Install Docker](https://docs.docker.com/get-docker/)
- k3d: creates the cluster in Docker. [Install k3d](https://k3d.io/stable/#installation)
- kubectl: talks to the cluster. [Install kubectl](https://kubernetes.io/docs/tasks/tools/)

Check that everything is installed:
```bash
docker version
k3d version
kubectl version --client
```

## Full setup, end to end

These steps take a clean machine to a running environment. Each service README
has more checks and notes.

### 1. Create the cluster
Make sure Docker is running.
```bash
k3d cluster create --config infra/k3d/k3d-config.yaml
kubectl get nodes
```
The node status must show `Ready`.

### 2. Set credentials
Credentials are never committed. Create a local `.env` from the example (an
existing `.env` is kept) and set your own passwords:
```bash
[ -f .env ] || cp .env.example .env
```
Edit `.env`: `RABBITMQ_USER`, `RABBITMQ_PASSWORD`, `POSTGRES_USER`,
`POSTGRES_PASSWORD`, `POSTGRES_DB`. `.env` is ignored by git.

### 3. Create the namespace and Secrets
```bash
set -a; . ./.env; set +a
kubectl apply -f infra/k8s/namespace.yaml
kubectl -n gojudge create secret generic rabbitmq-credentials \
  --from-literal=username="$RABBITMQ_USER" \
  --from-literal=password="$RABBITMQ_PASSWORD" \
  --dry-run=client -o yaml | kubectl apply -f -
kubectl -n gojudge create secret generic postgres-credentials \
  --from-literal=username="$POSTGRES_USER" \
  --from-literal=password="$POSTGRES_PASSWORD" \
  --from-literal=database="$POSTGRES_DB" \
  --dry-run=client -o yaml | kubectl apply -f -
```

### 4. Start RabbitMQ and PostgreSQL
```bash
kubectl apply -f infra/k8s/rabbitmq/rabbitmq.yaml
kubectl apply -f infra/k8s/postgres/postgres.yaml
kubectl -n gojudge rollout status statefulset/rabbitmq --timeout=180s
kubectl -n gojudge rollout status statefulset/postgres --timeout=180s
```

### 5. Check that everything works
```bash
kubectl -n gojudge get pods,pvc
kubectl -n gojudge exec rabbitmq-0 -- rabbitmq-diagnostics -q check_port_connectivity
kubectl -n gojudge exec postgres-0 -- sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "select 1;"'
```
Expected: `rabbitmq-0` and `postgres-0` are `1/1 Running`, both volumes are
`Bound`, the RabbitMQ check prints `Successfully connected to ports ...`, and
PostgreSQL returns `1`.

### 6. Connect
| Service | From inside the cluster | From your machine |
|---------|-------------------------|-------------------|
| RabbitMQ (AMQP) | `rabbitmq.gojudge.svc:5672` | `kubectl -n gojudge port-forward svc/rabbitmq 5672:5672 15672:15672`, then `localhost:5672` |
| RabbitMQ (management UI) | `rabbitmq.gojudge.svc:15672` | same port-forward, then `http://localhost:15672` |
| PostgreSQL | `postgres.gojudge.svc:5432` | `kubectl -n gojudge port-forward svc/postgres 5432:5432`, then `localhost:5432` |

Keep `port-forward` running in a second terminal while you use the service.

## Check that a Job works
The cluster runs Kubernetes Jobs, which GoJudge uses for evaluation:
```bash
kubectl create job hello --image=busybox -- echo "hello job"
kubectl wait --for=condition=complete job/hello --timeout=60s
kubectl logs job/hello
kubectl delete job hello
```
You will see `hello job`.

## Stop and delete
Stop the cluster and keep its data:
```bash
k3d cluster stop gojudge
k3d cluster start gojudge
```
Delete the cluster with all data:
```bash
k3d cluster delete gojudge
```