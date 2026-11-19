# SPK-05: Isolation mechanism comparison

Spike code for [issue #32](https://github.com/ITA-Internship/GoJudge/issues/32). Not merged to `main`
The report is `docs/spikes/SPK-05.md` on the `docs/spk-05-report` branch

Reuses the SPK-01 runner, image and manifests from `spk01/`.

| Folder | Candidate |
|--------|-----------|
| `cilium/` | Cilium |
| `gvisor/` | gVisor (`runsc`) |

Experiment logs saved in `cilium/results/` and `gvisor/results/`, like `spk01/results/`
