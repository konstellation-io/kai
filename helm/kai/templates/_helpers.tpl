{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "kai.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "kai.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name used by the chart label.
*/}}
{{- define "kai.chart" -}}
{{- printf "%s" .Chart.Name | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "kai.labels" -}}
helm.sh/chart: {{ include "kai.chart" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/* Fullname suffixed with admin-api */}}
{{- define "admin-api.fullname" -}}
{{- printf "%s-admin-api" (include "kai.fullname" .) -}}
{{- end }}

{{/*
Admin API labels
*/}}
{{- define "admin-api.labels" -}}
{{ include "kai.labels" . }}
{{ include "admin-api.selectorLabels" . }}
{{- end }}

{{/*
Admin API selector labels
*/}}
{{- define "admin-api.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kai.name" . }}-admin-api
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Admin API serviceaccount name
*/}}
{{- define "admin-api.serviceAccountName" -}}
{{- if .Values.adminApi.serviceAccount.create -}}
    {{ default (include "admin-api.fullname" .) .Values.adminApi.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.adminApi.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/* Fullname suffixed with k8s-manager */}}
{{- define "k8s-manager.fullname" -}}
{{- printf "%s-k8s-manager" (include "kai.fullname" .) -}}
{{- end }}

{{/*
k8s manager labels
*/}}
{{- define "k8s-manager.labels" -}}
{{ include "kai.labels" . }}
{{ include "k8s-manager.selectorLabels" . }}
{{- end }}

{{/*
k8s manager selector labels
*/}}
{{- define "k8s-manager.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kai.name" . }}-k8s-manager
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/* Create the name of k8s-manager service account to use */}}
{{- define "k8s-manager.serviceAccountName" -}}
{{- if .Values.k8sManager.serviceAccount.create -}}
    {{ default (include "k8s-manager.fullname" .) .Values.k8sManager.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.k8sManager.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/* Create the netrc file secret name for the image builder*/}}
{{- define "k8s-manager.netrc.secretName" -}}
{{- printf "%s-imagebuilder-netrc" (include "k8s-manager.fullname" . ) -}}
{{- end -}}

{{/*
prometheus-additional-scrape-configs labels
*/}}
{{- define "prometheus-additional-scrape-configs.labels" -}}
{{ include "kai.labels" . }}
app.kubernetes.io/name: {{ include "kai.name" . }}-scrape-configs
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/* Fullname suffixed with grafana-datasources */}}
{{- define "grafana-datasources.fullname" -}}
{{- printf "%s-grafana-datasources" (include "kai.fullname" .) -}}
{{- end }}

{{/*
grafana-datasources labels
*/}}
{{- define "grafana-datasources.labels" -}}
{{ include "kai.labels" . }}
app.kubernetes.io/name: {{ include "kai.name" . }}-grafana-datasources
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
redis-stack labels
*/}}
{{- define "redis-stack.labels" -}}
{{ include "kai.labels" . }}
app.kubernetes.io/name: {{ include "kai.name" . }}-redis-stack
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
redis master URL
*/}}
{{- define "redis.master.url" -}}
{{- if .Values.redis.enabled -}}
    {{- tpl .Values.config.redis.master.url . }}
{{- else -}}
    {{- .Values.config.redis.master.url -}}
{{- end -}}
{{- end -}}

{{/*
redis replicas URL
*/}}
{{- define "redis.replicas.url" -}}
{{- if .Values.redis.enabled -}}
    {{- tpl .Values.config.redis.replicas.url . }}
{{- else -}}
    {{- .Values.config.redis.replicas.url -}}
{{- end -}}
{{- end -}}

{{/*
redis secret name
*/}}
{{- define "redis.auth.secretName" -}}
{{- if .Values.redis.enabled }}
    {{- include "redis.secretName" .Subcharts.redis }}
{{- else -}}
    {{- default (printf "%s-redis" (include "kai.fullname" .)) .Values.config.redis.auth.existingSecret }}
{{- end -}}
{{- end -}}

{{/*
redis password key
*/}}
{{- define "redis.auth.secretPasswordKey" -}}
{{- if .Values.redis.enabled }}
    {{- include "redis.secretPasswordKey" .Subcharts.redis }}
{{- else -}}
    {{- default "redis-password" .Values.config.redis.auth.existingSecretPasswordKey }}
{{- end -}}
{{- end -}}

{{/*
redis password
*/}}
{{- define "redis.auth.password" -}}
{{- if .Values.redis.auth.enabled }}
    {{- if not (empty .Values.redis.auth.password) }}
        {{- .Values.redis.auth.password -}}
    {{- else -}}
        {{- include "getValueFromSecret" (dict "Namespace" .Release.Namespace "Name" (include "redis.secretName" .Subcharts.redis) "Length" 10 "Key" (include "redis.secretPasswordKey" .Subcharts.redis))  -}}
    {{- end -}}
{{- end -}}
{{- end -}}

{{/* Fullname suffixed with minio-config */}}
{{- define "minio-config.fullname" -}}
{{- printf "%s-minio-config" (include "kai.fullname" .) -}}
{{- end }}

{{/*
minio-config labels
*/}}
{{- define "minio-config.labels" -}}
{{ include "kai.labels" . }}
{{ include "minio-config.selectorLabels" . }}
{{- end }}

{{/*
minio-config selector labels
*/}}
{{- define "minio-config.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kai.name" . }}-minio-config
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
minio-config secret name
*/}}
{{- define "minio-config.tier.aws.secretName" -}}
{{ default (include "minio-config.fullname" . ) .Values.config.minio.tier.aws.auth.secretName }}
{{- end -}}

{{/*
minio-config access key
*/}}
{{- define "minio-config.tier.aws.accessKey" -}}
{{ default "accessKey" .Values.config.minio.tier.aws.auth.secretKeyNames.accessKey }}
{{- end -}}

{{/*
minio-config secret key
*/}}
{{- define "minio-config.tier.aws.secretKey" -}}
{{ default "secretKey" .Values.config.minio.tier.aws.auth.secretKeyNames.secretKey }}
{{- end -}}

{{/*
minio-config aws S3 endpoint URL
*/}}
{{- define "minio-config.tier.s3.endpointURL" -}}
{{- default "https://s3.amazonaws.com" .Values.config.minio.tier.aws.endpointURL -}}
{{- end }}

{{/*
minio-config remote bucket prefix (path in bucket to object transition)
*/}}
{{- define "minio-config.tier.s3.remotePrefix" -}}
{{- default "DATA" .Values.config.minio.tier.remotePrefix -}}
{{- end }}

{{/*
minio-config Tier name
*/}}
{{- define "minio-config.tier.name" -}}
{{- default "KAI-REMOTE-STORAGE" .Values.config.minio.tier.name -}}
{{- end }}

{{/*
minio-config AWS S3 remote bucket region for Tier
*/}}
{{- define "minio-config.tier.s3.region" -}}
{{- default "us-east-1" .Values.config.minio.tier.aws.region -}}
{{- end }}

{{/*
minio-config default buckets region
*/}}
{{- define "minio-config.region" -}}
{{- default "us-east-1" .Values.config.minio.defaultRegion -}}
{{- end }}

{{/*
nats host
*/}}
{{- define "nats.host" -}}
{{ default (include "nats.fullname" .Subcharts.nats) .Values.nats.service.name -}}
{{- end -}}

{{/*
nats url
*/}}
{{- define "nats.url" -}}
{{- printf "%s:%d" (include "nats.host" .) (.Values.nats.config.nats.port | int) -}}
{{- end -}}

{{/* Fullname suffixed with nats-manager */}}
{{- define "nats-manager.fullname" -}}
{{- printf "%s-nats-manager" (include "kai.fullname" .) -}}
{{- end }}

{{/*
nats manager labels
*/}}
{{- define "nats-manager.labels" -}}
{{ include "kai.labels" . }}
{{ include "nats-manager.selectorLabels" . }}
{{- end }}

{{/*
nats manager selector labels
*/}}
{{- define "nats-manager.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kai.name" . }}-nats-manager
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
nats manager serviceaccount name
*/}}
{{- define "nats-manager.serviceAccountName" -}}
{{- if .Values.natsManager.serviceAccount.create -}}
    {{ default (include "nats-manager.fullname" .) .Values.natsManager.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.natsManager.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/* Fullname suffixed with keycloak */}}
{{- define "keycloak.fullname" -}}
{{- printf "%s-keycloak" (include "kai.fullname" .) -}}
{{- end }}

{{/*
kaycloak labels
*/}}
{{- define "keycloak.labels" -}}
{{ include "kai.labels" . }}
{{ include "keycloak.selectorLabels" . }}
{{- end }}

{{/*
kaycloak selector labels
*/}}
{{- define "keycloak.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kai.name" . }}-keycloak
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
keycloak serviceaccount name
*/}}
{{- define "keycloak.serviceAccountName" -}}
{{- if .Values.keycloak.serviceAccount.create -}}
    {{ default (include "keycloak.fullname" .) .Values.keycloak.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.keycloak.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
keycloak secret name
*/}}
{{- define "keycloak.secretName" -}}
{{ default (include "keycloak.fullname" . ) .Values.keycloak.auth.existingSecret.name }}
{{- end -}}

{{/*
keycloak secret user key
*/}}
{{- define "keycloak.secretUserKey" -}}
{{ default "admin-user" .Values.keycloak.auth.existingSecret.userKey }}
{{- end -}}

{{/*
keycloak secret password key
*/}}
{{- define "keycloak.secretPasswordKey" -}}
{{ default "admin-password" .Values.keycloak.auth.existingSecret.passwordKey }}
{{- end -}}

{{/* Fullname suffixed with registry */}}
{{- define "registry.fullname" -}}
{{- printf "%s-registry" (include "kai.fullname" .) -}}
{{- end }}

{{/*
Registry labels
*/}}
{{- define "registry.labels" -}}
{{ include "kai.labels" . }}
{{ include "registry.selectorLabels" . }}
{{- end }}

{{/*
Registru selector labels
*/}}
{{- define "registry.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kai.name" . }}-registry
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Registry serviceaccount name
*/}}
{{- define "registry.serviceAccountName" -}}
{{- if .Values.registry.serviceAccount.create -}}
    {{ default (include "registry.fullname" .) .Values.registry.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.registry.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Registry user
*/}}
{{- define "registry.auth.user" -}}
    {{ default "user" .Values.registry.auth.user }}
{{- end -}}

{{/*
Registry password
*/}}
{{- define "registry.auth.password" -}}
    {{ default "password" .Values.registry.auth.password }}
{{- end -}}

{{/*
Registry auth secret name
*/}}
{{- define "registry.auth.secretName" -}}
{{- printf "%s-auth" (include "registry.fullname" . ) -}}
{{- end -}}

{{/*
Loki Host
*/}}
{{- define "loki.host" -}}
{{- if .Values.loki.enabled -}}
    {{- tpl .Values.config.loki.host . -}}
{{- else -}}
    {{- .Values.config.loki.host -}}
{{- end -}}
{{- end -}}

{{/*
Loki Port
*/}}
{{- define "loki.port" -}}
{{- if .Values.loki.enabled -}}
    {{- tpl .Values.config.loki.port . | quote -}}
{{- else -}}
    {{- .Values.config.loki.port | quote -}}
{{- end -}}
{{- end -}}

{{/*
Loki URL
*/}}
{{- define "loki.url" -}}
{{- printf "http://%s:%s" (include "loki.host" . ) (trimSuffix "\"" (include "loki.port" . ) | trimPrefix "\"") -}}
{{- end -}}

{{/*
Prometheus URL
*/}}
{{- define "prometheus.url" -}}
{{- if .Values.prometheus.enabled -}}
    {{- tpl .Values.config.prometheus.url . | quote -}}
{{- else -}}
    {{- .Values.config.prometheus.url -}}
{{- end -}}
{{- end -}}
