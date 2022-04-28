# Changelog v1.33

## Features


 - **[candi]** Bump containerd to v1.5.11. [#1386](https://github.com/deckhouse/deckhouse/pull/1386)
 - **[candi]** Improve candi bundle detection to detect centos-based distros [#1173](https://github.com/deckhouse/deckhouse/pull/1173)
 - **[cloud-provider-azure]** Enable accelerated networking for new machine-controller-manager instances. [#1266](https://github.com/deckhouse/deckhouse/pull/1266)
 - **[cloud-provider-yandex]** Changed default platform to standard-v3 for new instances created by machine-controller-manager. [#1361](https://github.com/deckhouse/deckhouse/pull/1361)
 - **[cni-flannel]** Bump flannel to 0.15.1. [#1173](https://github.com/deckhouse/deckhouse/pull/1173)
 - **[prometheus]** Create table with enabled Deckhouse web interfaces on the Grafana home page [#1415](https://github.com/deckhouse/deckhouse/pull/1415)

## Fixes


 - **[candi]** Migrate to cgroupfs on containerd installations. [#1386](https://github.com/deckhouse/deckhouse/pull/1386)
 - **[log-shipper]** Migrate deprecated elasticsearch fields [#1453](https://github.com/deckhouse/deckhouse/pull/1453)
 - **[log-shipper]** Send reloading signal to all vector processes in a container on config change. [#1430](https://github.com/deckhouse/deckhouse/pull/1430)
 - **[prometheus]** Removed the old prometheus_storage_class_change shell hook which has already been replaced by Go hooks. [#1396](https://github.com/deckhouse/deckhouse/pull/1396)
 - **[upmeter]** UI shows only present data [#1405](https://github.com/deckhouse/deckhouse/pull/1405)
 - **[upmeter]** Use finite timeout in agent insecure HTTP client [#1334](https://github.com/deckhouse/deckhouse/pull/1334)
 - **[upmeter]** Fixed slow data loading in [#1257](https://github.com/deckhouse/deckhouse/pull/1257)

## Chore


 - **[dashboard]** Dashboard upgrade from 2.2.0 to 2.5.1 [#1383](https://github.com/deckhouse/deckhouse/pull/1383)

