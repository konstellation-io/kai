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

{{/* Fullname suffixed with chronograf */}}
{{- define "chronograf.fullname" -}}
{{- printf "%s-chronograf" (include "kai.fullname" .) -}}
{{- end }}

{{/*
Chronograf labels
*/}}
{{- define "chronograf.labels" -}}
{{ include "kai.labels" . }}
{{ include "chronograf.selectorLabels" . }}
{{- end }}

{{/*
Chronograf selector labels
*/}}
{{- define "chronograf.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kai.name" . }}-chronograf
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

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
nats labels
*/}}
{{- define "nats.labels" -}}
{{ include "kai.labels" . }}
{{ include "nats.selectorLabels" . }}
{{- end }}

{{/*
nats selector labels
*/}}
{{- define "nats.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kai.name" . }}-nats
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
nats host
*/}}
{{- define "nats.host" -}}
{{- printf "%s-nats" .Release.Name -}}
{{- end }}

{{/*
nats url
*/}}
{{- define "nats.url" -}}
{{- printf "%s:%d" (include "nats.host" .) (.Values.nats.client.port | int) -}}
{{- end -}}

{{/* Fullname suffixed with mongo-writer */}}
{{- define "mongo-writer.fullname" -}}
{{- printf "%s-mongo-writer" (include "kai.fullname" .) -}}
{{- end }}

{{/*
mongo-writer labels
*/}}
{{- define "mongo-writer.labels" -}}
{{ include "kai.labels" . }}
{{ include "mongo-writer.selectorLabels" . }}
{{- end }}

{{/*
mongo-writer selector labels
*/}}
{{- define "mongo-writer.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kai.name" . }}-mongo-writer
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Mongo Express name
*/}}
{{- define "mongoExpress.fullname" -}}
{{ printf "%s-mongo-express" $.Release.Name }}
{{- end }}

{{/*
Mongo Express labels
*/}}
{{- define "mongoExpress.labels" -}}
{{ include "kai.labels" . }}
{{ include "mongoExpress.selectorLabels" . }}
{{- end }}

{{/*
Mongo Express selector labels
*/}}
{{- define "mongoExpress.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kai.name" . }}-mongo-express
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create InfluxDB URL.
*/}}
{{- define "kai-influxdb.influxURL" -}}
  {{- printf "http://%s-influxdb:8086" .Release.Name -}}
{{- end -}}
{{/*
Create a default fully qualified InfluxDB service name for InfluxDB.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "kai-influxdb.fullname" -}}
{{- if .Values.influxdb.fullnameOverride -}}
{{- .Values.influxdb.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default "influxdb" .Values.influxdb.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{/*
Create a default fully qualified Kapacitor service name for Chronograph.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "kai-kapacitor.fullname" -}}
{{- $name := default "kapacitor" .Values.kapacitor.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
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
