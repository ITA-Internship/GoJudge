#!/usr/bin/env bash
# Removes the SPK-06 cluster. Results in spk06/results stay
set -euo pipefail

k3d cluster delete spk06
