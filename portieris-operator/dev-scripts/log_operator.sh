#!/bin/bash
#
# Copyright 2020 IBM Corporation.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

if ! [ -x "$(command -v kubectl)" ]; then
    echo 'Error: kubectl is not installed.' >&2
    exit 1
fi

if [ -z "$PORTIERIS_NS" ]; then
    echo "PORTIERIS_NS is empty. Please set namespace name for portieris operator."
    exit 1
fi

PORTIERIS_OPERATOR_POD=`kubectl get pod -n ${PORTIERIS_NS} | grep portieris-operator | grep Running | awk '{print $1}'`
if [ -z "$PORTIERIS_OPERATOR_POD" ]; then
    echo "PORTIERIS_OPERATOR_POD is empty. There is no running portieris-operator"
    exit 1
fi

kubectl logs -f -n ${PORTIERIS_NS} ${PORTIERIS_OPERATOR_POD} -c manager
