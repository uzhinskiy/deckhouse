#!/bin/bash

# Copyright 2017-2020 Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

cni_version="0.9.0"
declare -A cni_sha512
cni_sha512[amd64]="13a5274aa0b146e77a8fe3c2b37485af69b655d114bee46ebf0535bfda578865e76d2e4a2a9736f2168cf1c61348580c14656679ed874dcd3132ad83e6d5c2fb"
cni_sha512[arm64]="d2cd05a8386a52edfe18d2686709c16857c3a0dc5fa3926f9d491a948140b940acb4e8e42f61a0a6a19733c23e27940ee50b71d2c5c8b71197500b41bd6bb619"
