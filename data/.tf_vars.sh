#!/bin/bash
set -euo pipefail

######################################################################
#             PROJECT VARIABLES SECTION. CHANGEABLE.
######################################################################

# assuming that runner image is alpine

apk update && apk add jq

export VAULT_SKIP_VERIFY=true
export VAULT_ADDR=https://localhost:8200

# get secrets

export KC_AUTH_ADMIN_CLI_URL=$(vault kv get -format=json secret/${CI_PROJECT_PATH}/${CI_COMMIT_REF_NAME} | jq -r '.data.data.KC_AUTH_ADMIN_CLI_URL')
export KC_AUTH_ADMIN_CLI_CLIENT_ID=$(vault kv get -format=json secret/${CI_PROJECT_PATH}/${CI_COMMIT_REF_NAME} | jq -r '.data.data.KC_AUTH_ADMIN_CLI_CLIENT_ID')
export KC_AUTH_ADMIN_CLI_USERNAME=$(vault kv get -format=json secret/${CI_PROJECT_PATH}/${CI_COMMIT_REF_NAME} | jq -r '.data.data.KC_AUTH_ADMIN_CLI_USERNAME')
export KC_AUTH_ADMIN_CLI_PASSWORD=$(vault kv get -format=json secret/${CI_PROJECT_PATH}/${CI_COMMIT_REF_NAME} | jq -r '.data.data.KC_AUTH_ADMIN_CLI_PASSWORD')

# setup kubeconfig with access to tfstate ns (located in ci/cd vars)

export TFSTATE_NAMESPACE="tfstate"
export KUBECONFIG="${TFSTATE_KUBECONFIG}"
export KUBE_CONFIG_PATH="${TFSTATE_KUBECONFIG}"

# setup tf vars

export TF_VAR_ci_project_root_namespace="${CI_PROJECT_ROOT_NAMESPACE}"
export TF_VAR_ci_project_name="${CI_PROJECT_NAME}"
export TF_VAR_ci_project_path="${CI_PROJECT_PATH}"
export TF_VAR_ci_project_path_slug="${CI_PROJECT_PATH_SLUG}"
export TF_VAR_ci_commit_ref_name="${CI_COMMIT_REF_NAME}"
export TF_VAR_infra_domain_name="${INFRA_DOMAIN_NAME}"

export TF_VAR_kc_auth_admin_cli_client_id="${KC_AUTH_ADMIN_CLI_CLIENT_ID}"
export TF_VAR_kc_auth_admin_cli_username="${KC_AUTH_ADMIN_CLI_USERNAME}"
export TF_VAR_kc_auth_admin_cli_password="${KC_AUTH_ADMIN_CLI_PASSWORD}"
export TF_VAR_kc_auth_admin_cli_url="${KC_AUTH_ADMIN_CLI_URL}"
