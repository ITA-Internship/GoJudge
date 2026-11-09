# RabbitMQ

RabbitMQ 4.4.0 (with the management UI) runs in the namespace `gojudge`.
Manifests: `infra/k8s/rabbitmq/rabbitmq.yaml`. The local cluster must be
running (see [`infra/README.md`](/infra/README.md), which also has the full
setup of all services).

### Credentials

Credentials are never committed. If you have no local `.env` file yet, create
it from the example. The command does not overwrite an existing `.env`, so
values you already set (for example for PostgreSQL) are kept:

```bash
[ -f .env ] || cp .env.example .env
```

Edit `.env` and set `RABBITMQ_USER` and `RABBITMQ_PASSWORD`.

The Secret shape is shown in `rabbitmq-secret.example.yaml`. It is an example
only: do not apply it. The real Secret is created from `.env` in the next step.

### Start

```bash
kubectl get namespace gojudge || kubectl apply -f infra/k8s/namespace.yaml
set -a; . ./.env; set +a
kubectl -n gojudge delete secret rabbitmq-credentials --ignore-not-found
kubectl -n gojudge create secret generic rabbitmq-credentials \
  --from-literal=username="$RABBITMQ_USER" \
  --from-literal=password="$RABBITMQ_PASSWORD"
kubectl apply -f infra/k8s/rabbitmq/rabbitmq.yaml
kubectl -n gojudge rollout status statefulset/rabbitmq --timeout=180s
```

Expected: `partitioned roll out complete: 1 new pods have been updated...`.

### Check from inside the cluster

```bash
kubectl -n gojudge exec statefulset/rabbitmq -- rabbitmq-diagnostics -q ping
kubectl -n gojudge exec statefulset/rabbitmq -- rabbitmq-diagnostics -q check_port_connectivity
kubectl -n gojudge run amqp-test --rm -it --restart=Never --image=busybox -- sh -c 'nc -z -w 3 rabbitmq 5672 && echo open'
```

Expected: `Ping succeeded`, no error from the connectivity check, and `open`. Inside the cluster the service is
`rabbitmq.gojudge.svc:5672`.

### Check from your machine

In a second terminal, keep this running:

```bash
kubectl -n gojudge port-forward svc/rabbitmq 5672:5672 15672:15672
```

```bash
set -a; . ./.env; set +a
curl -s -u "$RABBITMQ_USER:$RABBITMQ_PASSWORD" -X PUT -H 'content-type: application/json' -d '{"durable":true}' http://localhost:15672/api/queues/%2F/smoke
curl -s -u "$RABBITMQ_USER:$RABBITMQ_PASSWORD" -X POST -H 'content-type: application/json' -d '{"properties":{},"routing_key":"smoke","payload":"hello","payload_encoding":"string"}' http://localhost:15672/api/exchanges/%2F/amq.default/publish
curl -s -u "$RABBITMQ_USER:$RABBITMQ_PASSWORD" -X POST -H 'content-type: application/json' -d '{"count":1,"ackmode":"ack_requeue_false","encoding":"auto"}' http://localhost:15672/api/queues/%2F/smoke/get
curl -s -u "$RABBITMQ_USER:$RABBITMQ_PASSWORD" -X DELETE http://localhost:15672/api/queues/%2F/smoke
```

Expected: the second command returns `{"routed":true}` and the third returns
the message with `"payload":"hello"`. The last command removes the test queue.

### Notes

- The readiness probe runs `rabbitmq-diagnostics check_port_connectivity`. Unlike
  `ping`, it fails when the RabbitMQ application is stopped or a listener does
  not accept connections, so the pod is not marked ready too early.

- Queues must be durable. RabbitMQ 4.x rejects transient non-exclusive queues.
- Data is stored on a 1Gi volume (`data-rabbitmq-0`), so durable queues and
  definitions survive a pod restart.
- If you have the earlier Deployment running, remove it before you apply the
  StatefulSet. The Service must be recreated too, because `clusterIP` cannot be
  changed: `kubectl -n gojudge delete deployment/rabbitmq service/rabbitmq`.
- If you applied an earlier version that used an `emptyDir`, delete the
  StatefulSet before you apply this one, because its volumes cannot be changed
  in place: `kubectl -n gojudge delete statefulset/rabbitmq`.
- RabbitMQ creates the user from the Secret only on the first start with an
  empty volume. To change the password, or to reset RabbitMQ, delete it with
  its data, then run the commands from "Start" again:

  ```bash
  kubectl delete -f infra/k8s/rabbitmq/rabbitmq.yaml
  kubectl -n gojudge delete pvc data-rabbitmq-0
  ```

- Remove only RabbitMQ: `kubectl delete -f infra/k8s/rabbitmq/rabbitmq.yaml`.