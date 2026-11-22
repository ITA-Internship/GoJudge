# Cilium: SC-04 egress test

Does Cilium block egress from the start of a Pod

## Run

```bash
./spk05/cilium/setup.sh
./spk05/cilium/run-sc04.sh
```

Control without the policy:

```bash
kubectl --context k3d-spk05-cilium -n spk01 delete networkpolicy spk01-deny-egress
./spk05/cilium/run-sc04.sh 1 SC-04-control
kubectl --context k3d-spk05-cilium apply -f spk01/manifests/sc-04-deny-egress.yaml
```

## Result

- With the policy: 5/5 blocked
- Control without the policy: network reached

Logs are in `results/`
