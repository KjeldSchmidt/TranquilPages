## Repeatable Environments

Each environment that is instantiated has a subfolder in /env/<env_name>, where
environment-specific variables are set and the `env` module is instantiated.
This is the only thing done to specifiy the environment.
**_No resources are created_** in these environment specific folders. Every 
resource  is created directly in the `env` module or in a submodule.

## Base Environment

To bootstrap the infrastructure and achieve limited scope on the service
principals used to deploy this project, we create the resource groups for each
environment in a "base" environment, and create a service principal with
contributor permissions on the scope of these resource groups.

Thus, the base environment must be created with a user principal, while the 
other environments can be deployed from CI/CD pipelines with Service Principal
authentication.