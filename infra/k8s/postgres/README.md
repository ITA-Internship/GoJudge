## PostgreSQL

PostgreSQL 17 in the namespace `gojudge`.

### Start

1. Start the local cluster (see [infra/README.md](../../README.md)), if it is not running yet:
```bash
k3d cluster create --config infra/k3d/k3d-config.yaml
```

2. In `.env` set `POSTGRES_USER`, `POSTGRES_PASSWORD` and `POSTGRES_DB`
   (copy `.env.example` to `.env` first if you have no `.env`)

3. Run:
```bash
set -a; . ./.env; set +a
kubectl apply -f infra/k8s/namespace.yaml
kubectl -n gojudge create secret generic postgres-credentials \
  --from-literal=username="$POSTGRES_USER" \
  --from-literal=password="$POSTGRES_PASSWORD" \
  --from-literal=database="$POSTGRES_DB"
kubectl apply -f infra/k8s/postgres/postgres.yaml
kubectl -n gojudge rollout status deploy/postgres
```

### Check
```bash
kubectl -n gojudge exec deploy/postgres -- sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "select 1;"'
```
You must see `1`

### Connect

- From the cluster: `postgres.gojudge.svc:5432`
- From your computer: run `kubectl -n gojudge port-forward svc/postgres 5432:5432`
  in a second terminal, then use `localhost:5432`

### Notes
- Remove PostgreSQL: `kubectl delete -f infra/k8s/postgres/postgres.yaml`
