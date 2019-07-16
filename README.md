# Serverless Secrets Talk

This repository has the code samples and setup scripts for my live-demo of how
to manage secrets in serverless environments.


## Setup

```text
$ ./bin/setup
```

## Demo

1. Show code...

1. How will we get Redis its password?

### IAM

1. Do you actually need a credential? Can you use IAM, even across clouds
   instead?

### Plaintext envvars

1. Deploy:

    ```text
    $ ./bin/deploy-app
    ```

1. Show basic app functionality.

1. Reset count.

1. Copy link, change to a negative number.

1. Copy link, change to a non-int, get debug page.

    > Reason #1 why you shouldn't store plaintext secrets in environment
    > variables - they are easily leaked by popular webframeworks and debugging
    > tools.

1. Ways to mitigate?

    1. CI/CD tests.

    1. Automated deployments.

    1. Don't store plaintext secrets in environment variables.

1. Add `ENV=prod`, re-deploy app.

1. Show that we get a generic demo page, highlight that our secrets are still in
   plaintext.

1. What about software supply chain? Turns out we accidentially included a
   malicious third-party dependency.

   Malicious dependency dumps our environment to an HTTP endpoint.

1. Show Stackdriver log of malicious-server, higlighting host and password.

1. Telnet into redis host and password:

    ```text
    telnet <IP> 6379
    > AUTH super-secret...

    > PING

    > SET visits 1000000
    ```

    > Reason #2 why you shouldn't sore plaintext secrets in environment
    > variables - they are easily retrieved by a malicious depending in your
    > software supply chain.

1. Ways to mitigate?

    1. Vulnerability scanning (reactive).

    1. Egress firewall rules (don't allow outbound traffic).

    1. Don't store plaintext secrets in environment variables.


### Encrypted envvars

Let's encrypt those envvars!

1. Encrypt our plaintext value with a KMS provider:

    ```text
    $ ./bin/encrypt-string super-secret...
    ```

1. Edit the deployment script to send in the encrypted value:

    ```text
    REDIS_PASS="..."
    ```

1. Just for fun, let's run in development mode again - remove `ENV=production`.

1. Update the code to decrypt first (helper is in `kms.go`):

    ```go
    func main() {
      redisPass, err := kmsDecrypt(redisPass)
      if err != nil {
        panic(err)
      }

      # ...
    }
    ```

1. Deploy app.

    ```text
    $ ./bin/deploy-app
    ```

1. Cause a page crash, show that environment variable is encrypted.

1. Visit malicious server, see that payload is encrypted.


### Central storage

Encrypted envvars lack auditing and central management. Using a central storage
system like a cloud provider's secret manager or storage system centralizes
access, permissions, and logging.

We'll use [berglas](https://github.com/GoogleCloudPlatform/berglas) for these
examples.

1. Create a secret:

    ```text
    $ berglas create sethvargo-oscon19-secrets/redis-pass super-secret... \
        --key projects/sethvargo-oscon19/locations/global/keyRings/serverless/cryptoKeys/secrets
    ```

1. Grant our serverless app the ability to access the value:

    ```text
    $ berglas grant "sethvargo-oscon19-secrets/redis-pass" \
        --member "serviceAccount:myapp-sa@sethvargo-oscon19.iam.gserviceaccount.com"
    ```

1. Update our serverless app to pull from Berglas:

    ```go
    func main() {
      redisPass, err := berglasAccess("redis-pass")
      if err != nil {
        panic(err)
      }
    }
    ```

1. Drop `REDIS_PASS` environment variable.

1. Deploy app.

    ```text
    $ ./bin/deploy-app
    ```

1. Cause a page crash, no environment variable is present.

1. Visit malicious server, no environment variable is present.


### HashiCorp Vault

Vault is already running, we just need to configure it.

1. Enable Vault's KV secret engine:

    ```text
    $ ./bin/vault secrets enable -version=2 kv
    ```

1. Create our secret in Vault:

    ```text
    $ ./bin/vault kv put kv/myapp/redis-pass value=super-secret...
    ```

1. Create a policy granting access to our secret:

    ```text
    $ ./bin/vault policy write myapp-kv-read ./scripts/vault-policy.hcl
    ```

1. Allow serverless app to auth to Vault:

    ```text
    $ ./bin/vault write auth/gcp-serverless/role/myapp \
        type=iam \
        project_id=sethvargo-oscon19 \
        policies=myapp-kv-read \
        bound_service_accounts=myapp-sa@sethvargo-oscon19.iam.gserviceaccount.com \
        max_jwt_exp=60m
    ```

1. Demo `vault.go` code.

1. Update serverless app to pull from Vault:

    ```go
    func main() {
      redisPass, err := vaultAccess("kv/data/myapp/redis-pass")
      if err != nil {
        panic(err)
      }
    }
    ```

1. Add `VAULT_ADDR` to function deployment:

    ```text
    (cd terraform && terraform output vault_addr)
    ```

    ```text
    VAULT_ADDR="..."

    --set-env-vars="VAULT_ADDR="${VAULT_ADDR}..."
    ```
