# Local cluster
Local Kubernetes cluster for GoJudge

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

## Start
Make sure Docker is running

From the repository root:
```bash
k3d cluster create --config infra/k3d/k3d-config.yaml
kubectl get nodes
```
The node status must show `Ready`

## Check that a Job works
```bash
kubectl create job hello --image=busybox -- echo "hello job"
kubectl wait --for=condition=complete job/hello --timeout=60s
kubectl logs job/hello
kubectl delete job hello
```
You will see `hello job`

## Delete cluster
```bash
k3d cluster delete gojudge
```
