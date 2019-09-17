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

module "vault" {
  source         = "terraform-google-modules/vault/google"
  project_id     = "sethvargo-devsecconseattle-19"
  region         = "us-central1"
  kms_keyring    = "vault"
  kms_crypto_key = "vault-init"
  vault_version  = "1.2.2"

  vault_instance_base_image    = "debian-cloud/debian-9"
  storage_bucket_force_destroy = "true"
}

output "vault_addr" {
  value = "${module.vault.vault_addr}"
}

output "service_account_email" {
  value = "${module.vault.service_account_email}"
}

output "ca_cert_pem" {
  value = "${module.vault.ca_cert_pem[0]}"
}
