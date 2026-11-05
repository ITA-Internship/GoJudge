## PostgreSQL

PostgreSQL 17.11 in the namespace `gojudge`.

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
  --from-literal=database="$POSTGRES_DB" \
  --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f infra/k8s/postgres/postgres.yaml
kubectl -n gojudge rollout status statefulset/postgres
```

### Check
```bash
kubectl -n gojudge exec postgres-0 -- sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "select 1;"'
```
You must see `1`

### Connect

- From the cluster: `postgres.gojudge.svc:5432`
- From your computer: run `kubectl -n gojudge port-forward svc/postgres 5432:5432`
  in a second terminal, then use `localhost:5432`

### Notes
- Data is kept in volume and survives a pod restart.
- To reset PostgreSQL (for example, to change the password), delete it with its data, then run the commands from step 3 of "Start" again:
```bash
kubectl delete -f infra/k8s/postgres/postgres.yaml
kubectl -n gojudge delete pvc data-postgres-0
```
- `postgres-secret.example.yaml` is only example. Do not apply it
