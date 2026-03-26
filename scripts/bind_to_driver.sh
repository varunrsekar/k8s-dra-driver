#!/bin/bash

# Copyright The Kubernetes Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Usage: ./bind_to_driver.sh <ssss:bb:dd.f> <driver>
# Bind the GPU specified by the PCI_ID=ssss:bb:dd.f to the given driver.

bind_to_driver()
{
   local gpu=$1
   local driver=$2
   local drivers_path="/sys/bus/pci/drivers"
   local driver_override_file="/sys/bus/pci/devices/$gpu/driver_override"
   local bind_file="$drivers_path/$driver/bind"

   if [ ! -e "$driver_override_file" ]; then
      echo "'$driver_override_file' file does not exist" >&2
      return 1
   fi

   echo "$driver" > "$driver_override_file"
   if [ $? -ne 0 ]; then
      echo "failed to write '$driver' to $driver_override_file" >&2
      return 1
   fi

   if [ ! -e "$bind_file" ]; then
      echo "'$bind_file' file does not exist" >&2
      return 1
   fi

   echo "$gpu" > "$bind_file"
   if [ $? -ne 0 ]; then
      echo "failed to write '$gpu' to $bind_file" >&2
      echo "" > "$driver_override_file"
      return 1
   fi
}

bind_to_driver "$1" "$2" || exit 1