## Github Actions secrets

In order for the pipeline to function, the following secrets/variables need to
be set for the project:

Secret name: `AZURE_CREDENTIALS`
Secret value: `{ "clientSecret":  "******", "subscriptionId":  "******", "tenantId":  "******", "clientId":  "******" }`

This service principal is created by manually `terraform apply`ing the 
`./infra/base` terraform root module with a user account that is an owner of
the relevant subscription. The four required values can then be taken from the
terraform output values

---

Secret Name: `DOCKER_API_TOKEN`/`DOCKER_USERNAME`

Username and API token to push the build image to dockerhub. Token must be
created with write permissions.

## Keyvault Secrets

These keyvault secrets need to be set manually, but placeholders for them are created via IaC.

`user-login-oauth-client-id` and `user-login-oauth-client-secret` are Google 
OAuth credentials, to be acquired in the 
[Cloud Console](https://console.cloud.google.com/apis/credentials). 