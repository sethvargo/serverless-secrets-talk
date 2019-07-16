#!/usr/bin/env bash
set -eEuo pipefail

# Copyright 2019 Seth Vargo
# Copyright 2019 Google, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

apt-get -yqq update
apt-get -yqq upgrade

apt-get -yqq install curl redis-server vim

sed -i 's/bind .*/# bind 127.0.0.1/' /etc/redis/redis.conf
sed -i 's/# requirepass.*/requirepass super-secret/' /etc/redis/redis.conf

systemctl enable redis-server.service
systemctl restart redis-server.service
