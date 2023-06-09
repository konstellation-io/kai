# kai

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| https://helm.influxdata.com/ | influxdb | 4.8.1 |
| https://helm.influxdata.com/ | kapacitor | 1.4.6 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| adminApi.affinity | object | `{}` | Assign custom affinity rules to the Admin API pods |
| adminApi.host | string | `"api.kai.local"` | Hostname |
| adminApi.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| adminApi.image.repository | string | `"konstellation/kai-admin-api"` | Image repository |
| adminApi.image.tag | string | `"0.2.0-develop.2"` | Image tag |
| adminApi.ingress.annotations | object | See `adminApi.ingress.annotations` in [values.yaml](./values.yaml) | Ingress annotations |
| adminApi.ingress.className | string | `"kong"` | The name of the ingress class to use |
| adminApi.logLevel | string | `"INFO"` | Default application log level |
| adminApi.nodeSelector | object | `{}` | Define which Nodes the Pods are scheduled on. |
| adminApi.storage.class | string | `"standard"` | Storage class name |
| adminApi.storage.path | string | `"/admin-api-files"` | Persistent volume mount point. This will define Admin API app workdir too. |
| adminApi.storage.size | string | `"1Gi"` | Storage class size |
| adminApi.tls.enabled | bool | `false` | Whether to enable TLS |
| adminApi.tolerations | list | `[]` | Tolerations for use with node taints |
| chronograf.affinity | object | `{}` | Assign custom affinity rules to the Chronograf pods |
| chronograf.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| chronograf.image.repository | string | `"chronograf"` | Image repository |
| chronograf.image.tag | string | `"1.8.4"` | Image tag |
| chronograf.nodeSelector | object | `{}` | Define which Nodes the Pods are scheduled on. |
| chronograf.persistence.accessMode | string | `"ReadWriteOnce"` | Access mode for the volume |
| chronograf.persistence.enabled | bool | `true` | Whether to enable persistence |
| chronograf.persistence.size | string | `"2Gi"` | Storage size |
| chronograf.persistence.storageClass | string | `"standard"` | Storage class name |
| chronograf.tolerations | list | `[]` | Tolerations for use with node taints |
| config.admin.apiHost | string | `"api.kai.local"` | Api Hostname for Admin UI and Admin API |
| config.admin.corsEnabled | bool | `true` | Whether to enable CORS on Admin API |
| config.admin.userEmail | string | `"dev@local.local"` | Email address for sending notifications |
| config.auth.apiTokenSecret | string | `"api_token_secret"` | API token secret |
| config.auth.cookieDomain | string | `"kai.local"` | Admin API secure cookie domain |
| config.auth.jwtSignSecret | string | `"jwt_secret"` | JWT Sign secret |
| config.auth.secureCookie | bool | `false` | Whether to enable secure cookie for Admin API |
| config.auth.verificationCodeDurationInMinutes | int | `1` | Verification login link duration |
| config.baseDomainName | string | `"local"` | Base domain name for Admin API and K8S Manager apps |
| config.mongodb.connectionString.secretKey | string | `""` | The name of the secret key that contains the MongoDB connection string. |
| config.mongodb.connectionString.secretName | string | `""` | The name of the secret that contains a key with the MongoDB connection string. |
| config.smtp.enabled | bool | `false` | Whether to enable SMTP server connection |
| config.smtp.pass | string | `""` | SMTP server password |
| config.smtp.user | string | `""` | SMTP server user |
| developmentMode | bool | `false` | Whether to setup developement mode |
| influxdb.address | string | `"http://kai-influxdb/"` |  |
| influxdb.affinity | object | `{}` | Assign custom affinity rules to the InfluxDB pods |
| influxdb.config.http | object | `{"auth-enabled":false,"enabled":true,"flux-enabled":true}` | [Details](https://docs.influxdata.com/influxdb/v1.8/administration/config/#http) |
| influxdb.image.tag | string | `"1.8.1"` | Image tag |
| influxdb.initScripts.enabled | bool | `true` | Boolean flag to enable and disable initscripts. See https://github.com/influxdata/helm-charts/tree/master/charts/influxdb#configure-the-chart for more info |
| influxdb.initScripts.scripts | object | `{"init.iql":"CREATE DATABASE \"kai\"\n"}` | Init scripts |
| influxdb.nodeSelector | object | `{}` | Define which Nodes the Pods are scheduled on. |
| influxdb.persistence.accessMode | string | `"ReadWriteOnce"` | Access mode for the volume |
| influxdb.persistence.enabled | bool | `true` | Whether to enable persistence. See https://github.com/influxdata/helm-charts/tree/master/charts/influxdb#configure-the-chart for more info |
| influxdb.persistence.size | string | `"10Gi"` | Storage size |
| influxdb.persistence.storageClass | string | `"standard"` | Storage class name |
| influxdb.tolerations | list | `[]` | Tolerations for use with node taints |
| k8sManager.affinity | object | `{}` | Assign custom affinity rules to the K8S Manager pods |
| k8sManager.generatedEntrypoints.ingress.annotations | object | See `entrypoints.ingress.annotations` in [values.yaml](./values.yaml) | The annotations that all the generated ingresses for the entrypoints will have |
| k8sManager.generatedEntrypoints.ingress.className | string | `"kong"` | The ingressClassName to use for the enypoints' generated ingresses |
| k8sManager.generatedEntrypoints.ingress.tls.secretName | string | If not defined, every created ingress will use an autogenerated certificate name based on the deployed runtimeId and .Values.config.baseDomainName. | TLS certificate secret name. If defined, wildcard for the current application domain must be used. |
| k8sManager.generatedEntrypoints.tls | bool | `false` | Whether to enable tls |
| k8sManager.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| k8sManager.image.repository | string | `"konstellation/kai-k8s-manager"` | Image repository |
| k8sManager.image.tag | string | `"0.2.0-develop.2"` | Image tag |
| k8sManager.krtFilesDownloader.image.pullPolicy | string | `"Always"` | Image pull policy |
| k8sManager.krtFilesDownloader.image.repository | string | `"konstellation/krt-files-downloader"` | Image repository |
| k8sManager.krtFilesDownloader.image.tag | string | `"latest"` | Image tag |
| k8sManager.nodeSelector | object | `{}` | Define which Nodes the Pods are scheduled on. |
| k8sManager.serviceAccount.annotations | object | `{}` | The Service Account annotations |
| k8sManager.serviceAccount.create | bool | `true` | Whether to create the Service Account |
| k8sManager.serviceAccount.name | string | `""` | The name of the service account. @default: A pre-generated name based on the chart relase fullname sufixed by `-k8s-manager` |
| k8sManager.tolerations | list | `[]` | Tolerations for use with node taints |
| kapacitor.enabled | bool | `false` | Whether to enable Kapacitor |
| kapacitor.persistence.enabled | bool | `false` | Whether to enable persistence [Details](https://github.com/influxdata/helm-charts/blob/master/charts/kapacitor/values.yaml) |
| keycloak.adminApi.oidcClient.clientId | string | `"admin-cli"` | The name of the OIDC client in Keycloak for the master realm admin |
| keycloak.affinity | object | `{}` | Assign custom affinity rules to the Keycloak pods |
| keycloak.argsOverride | object | `{}` | Args to pass to the Keycloak startup command. This takes precedence over options passed through env variables |
| keycloak.auth.adminPassword | string | `"123456"` | Keycloak admin password |
| keycloak.auth.adminUser | string | `"admin"` | Keycloak admin username |
| keycloak.auth.existingSecret.name | string | `""` | The name of the secret that contains a key with the Keycloak admin password. Existing secret takes precedence over `adminUser` and `adminPassword` |
| keycloak.auth.existingSecret.passwordKey | string | `""` | The name of the secret key that contains the Keycloak admin password. |
| keycloak.auth.existingSecret.userKey | string | `""` | The name of the secret key that contains the Keycloak admin username. |
| keycloak.config.healthEnabled | string | `"true"` | If the server should expose health check endpoints. If set to "false", container liveness and readiness probes should be disabled. |
| keycloak.config.hostnameStrict | string | `"false"` | Disables dynamically resolving the hostname from request headers. Should always be set to true in production, unless proxy verifies the Host header. |
| keycloak.config.httpEnabled | string | `"true"` | Whether to enable http |
| keycloak.config.metricsEnabled | string | `"false"` | Whether to enable metrics |
| keycloak.config.proxy | string | `"edge"` | The proxy address forwarding mode if the server is behind a reverse proxy. Valid values are `none`, `edge`, `reencrypt` and `passthrough` |
| keycloak.db.auth.database | string | `""` | The database name |
| keycloak.db.auth.host | string | `""` | The database hostname |
| keycloak.db.auth.port | string | `""` | The database port |
| keycloak.db.auth.secretDatabaseKey | string | `""` | The name of the secret key that contains the database name. Takes precedence over `database` |
| keycloak.db.auth.secretHostKey | string | `""` | The name of the secret key that contains the database host. |
| keycloak.db.auth.secretName | string | `""` | The name of the secret that contains the database connection config keys. |
| keycloak.db.auth.secretPasswordKey | string | `""` | The name of the secret key that contains the database password. |
| keycloak.db.auth.secretPortKey | string | `""` | The name of the secret key that contains the database port. Takes precedence over `host` |
| keycloak.db.auth.secretUserKey | string | `""` | The name of the secret key that contains the database username. Takes precedence over `port` |
| keycloak.db.type | string | `"postgres"` | Keycloak database type |
| keycloak.extraEnv | object | `{}` | Keycloak extra env vars in the form of a list of key-value pairs |
| keycloak.extraVolumeMounts | list | `[]` | Extra volume mounts |
| keycloak.extraVolumes | list | `[]` | Extra volumes |
| keycloak.host | string | `"auth.kai.local"` |  |
| keycloak.image.pullPolicy | string | `"IfNotPresent"` | The image pull policy |
| keycloak.image.repository | string | `"quay.io/keycloak/keycloak"` | The image repository |
| keycloak.image.tag | string | `"21.1.1"` | The image tag |
| keycloak.imagePullSecrets | list | `[]` | Image pull secrets |
| keycloak.ingress.annotations | object | See `keycloak.ingress.annotations` in [values.yaml](./values.yaml) | Ingress annotations |
| keycloak.ingress.className | string | `"kong"` | The name of the ingress class to use |
| keycloak.kli.oidcClient.clientId | string | `"kai-kli-oidc"` | The name of the OIDC client in Keycloak for KLI |
| keycloak.kong.oidcClient.clientId | string | `"kong-oidc"` | The name of the OIDC client in Keycloak for Kong |
| keycloak.kong.oidcClient.secret | string | `""` | The secret for the OIDC client that will be created on Keycloak first startup |
| keycloak.kong.oidcPluginName | string | `"oidc"` | The name of the OIDC Kong plugin that should be installed on Kong ingress controller |
| keycloak.livinessProbe | object | `{"failureThreshold":3,"httpGet":{"path":"/health/live","port":"http"},"initialDelaySeconds":30,"periodSeconds":10,"timeoutSeconds":5}` | Container liveness probe |
| keycloak.nodeSelector | object | `{}` | Define which Nodes the Pods are scheduled on. |
| keycloak.podAnnotations | object | `{}` | Pod annotations |
| keycloak.podSecurityContext | object | `{}` | Pod security context |
| keycloak.readinessProbe | object | `{"failureThreshold":3,"httpGet":{"path":"/health/ready","port":"http"},"initialDelaySeconds":30,"periodSeconds":10,"timeoutSeconds":5}` | Container readiness probe |
| keycloak.realmName | string | `"konstellation"` | The name of the realm that will be crated on Keycloak first startup |
| keycloak.resources | object | `{}` | Container resources |
| keycloak.securityContext | object | `{}` |  |
| keycloak.service.ports.http | int | `8080` | The http port the service will listen on. Only |
| keycloak.service.ports.https | int | `8443` | The https port the service will listen on |
| keycloak.service.type | string | `"ClusterIP"` | Service type |
| keycloak.serviceAccount.annotations | object | `{}` |  |
| keycloak.serviceAccount.create | bool | `true` |  |
| keycloak.serviceAccount.name | string | `""` |  |
| keycloak.tls.enabled | bool | `false` | Whether to enable TLS |
| keycloak.tolerations | list | `[]` | Assign custom tolerations to the Keycloak pods |
| mongoExpress.affinity | object | `{}` | Assign custom affinity rules to the Mongo Express pods |
| mongoExpress.connectionString.secretKey | string | `""` | The name of the secret key that contains the MongoDB connection string. |
| mongoExpress.connectionString.secretName | string | `""` | The name of the secret that contains a key with the MongoDB connection string. |
| mongoExpress.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| mongoExpress.image.repository | string | `"mongo-express"` | Image repository |
| mongoExpress.image.tag | string | `"0.54.0"` | Image tag |
| mongoExpress.nodeSelector | object | `{}` | Define which Nodes the Pods are scheduled on. |
| mongoExpress.tolerations | list | `[]` | Tolerations for use with node taints |
| mongoWriter.affinity | object | `{}` | Assign custom affinity rules to the Mongo Writter pods |
| mongoWriter.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| mongoWriter.image.repository | string | `"konstellation/kai-mongo-writer"` | Image repository |
| mongoWriter.image.tag | string | `"0.2.0-develop.2"` | Image tag |
| mongoWriter.nodeSelector | object | `{}` | Define which Nodes the Pods are scheduled on. |
| mongoWriter.tolerations | list | `[]` | Tolerations for use with node taints |
| nameOverride | string | `""` | Provide a name in place of kai for `app.kubernetes.io/name` labels |
| nats.affinity | object | `{}` | Assign custom affinity rules to the NATS pods |
| nats.client.port | int | `4222` | Port for client connections |
| nats.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| nats.image.repository | string | `"nats"` | Image repository |
| nats.image.tag | string | `"2.8.4"` | Image tag |
| nats.jetstream.memStorage.enabled | bool | `true` | Whether to enable memory storage for Jetstream |
| nats.jetstream.memStorage.size | string | `"2Gi"` | Memory storage max size for JetStream |
| nats.jetstream.storage.enabled | bool | `true` | Whether to enable a PersistentVolumeClaim for Jetstream |
| nats.jetstream.storage.size | string | `"5Gi"` | Storage size for the Jetstream PersistentVolumeClaim. Notice this is also used for the Jetstream storage limit configuration even if PVC creation is disabled |
| nats.jetstream.storage.storageClassName | string | `"standard"` | Storage class name for the Jetstream PersistentVolumeClaim |
| nats.jetstream.storage.storageDirectory | string | `"/data"` | Directory to use for JetStream storage when using a PersistentVolumeClaim |
| nats.limits.lameDuckDuration | string | `"30s"` | Duration over which to slowly close close client connections after lameDuckGracePeriod has passed |
| nats.limits.lameDuckGracePeriod | string | `"10s"` | Grace period after pod begins shutdown before starting to close client connections |
| nats.limits.maxConnections | string | 64K | Maximum number of active client connections. |
| nats.limits.maxControlLine | string | 4KB | Maximum length of a protocol line (including combined length of subject and queue group). Increasing this value may require cliet changes. Applies to all traffic |
| nats.limits.maxPayload | string | 1MB | Maximum number of bytes in a message payload. Reducing this size may force you to implement chunking in your clients. Applies to client and leafnode payloads. It is not recommended to use values over 8MB but `max_payload` can be set up to 64MB. The max payload must be equal or smaller to the `max_pending` value. |
| nats.limits.maxPending | string | 64MB | Maximum number of bytes buffered for a connection Applies to client connections. Note that applications can also set 'PendingLimits' (number of messages and total size) for their subscriptions. |
| nats.limits.maxPings | string | 2 | After how many unanswered pings the server will allow before closing the connection. |
| nats.limits.maxSubscriptions | string | 0 (unlimited) | Maximum numbers of subscriptions per client and leafnode accounts connection. |
| nats.limits.pingInterval | string | `nil` |  |
| nats.limits.writeDeadline | string | 10s | Maximum number of seconds the server will block when writing. Once this threshold is exceeded the connection will be closed. |
| nats.logging.debug | bool | `false` | Whether to enable logging debug mode |
| nats.logging.logtime | bool | `true` | Timestamp log entries |
| nats.logging.trace | bool | `false` | Whether to enable logging trace mode |
| nats.nodeSelector | object | `{}` | Define which Nodes the Pods are scheduled on. |
| nats.tolerations | list | `[]` | Tolerations for use with node taints |
| natsManager.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| natsManager.image.repository | string | `"konstellation/kai-nats-manager"` | Image repository |
| natsManager.image.tag | string | `"0.2.0-develop.2"` | Image tag |
| rbac.create | bool | `true` | Whether to create the roles for the services that could use custom Service Accounts |
| registry.affinity | object | `{}` | Assign custom affinity rules to the pods |
| registry.config | string | `""` | A string contaning the config for Docker Registry. Ref: https://docs.docker.com/registry/configuration/. |
| registry.configSecret.key | string | `""` | The name of the secret key that contains the registry config file |
| registry.configSecret.name | string | `""` | Takes precedence over 'registry.config'. The name of the secret that contains the registry config file. |
| registry.containerPort | int | `5000` | The container port |
| registry.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| registry.image.repository | string | `"registry"` | Image repository |
| registry.image.tag | string | `"2.8.2"` | Image tag |
| registry.imagePullSecrets | list | `[]` | Image pull secrets |
| registry.nodeSelector | object | `{}` | Define which Nodes the Pods are scheduled on. |
| registry.podAnnotations | object | `{}` | Pod annotations |
| registry.podSecurityContext | object | `{}` | Pod security context |
| registry.resources | object | `{}` | Container resources |
| registry.securityContext | object | `{}` |  |
| registry.service.ports.http | int | `5000` | The http port the service will listen on. Only |
| registry.service.type | string | `"ClusterIP"` | Service type |
| registry.serviceAccount.annotations | object | `{}` | Annotations to add to the service account |
| registry.serviceAccount.create | bool | `true` | Specifies whether a service account should be created |
| registry.serviceAccount.name | string | `""` | The name of the service account to use. If not set and create is true, a name is generated using the fullname template |
| registry.storage.accessMode | string | `"ReadWriteOnce"` | Access mode for the volume |
| registry.storage.enabled | bool | `true` | Whether to enable persistence |
| registry.storage.path | string | `"/var/lib/registry"` | Persistent volume mount point. This will define Registry app workdir too. |
| registry.storage.size | string | `"10Gi"` | Storage size |
| registry.storage.storageClass | string | `"sandard"` | Storage class name |
| registry.tolerations | list | `[]` | Tolerations for use with node taints |
