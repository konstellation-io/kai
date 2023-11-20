package configuration

import (
	"fmt"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetAppConfig(version domain.Version, processesConfig map[string]string) apiv1.ConfigMap {
	labels := map[string]string{
		"product": version.Product,
		"version": version.Tag,
		"type":    "configuration",
	}

	configMapName := fmt.Sprintf("%s-%s-conf-files", version.Product, version.Tag)

	return apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:   configMapName,
			Labels: labels,
		},
		Data: mergeConfigs(processesConfig, getFluentBitConfig()),
	}
}

func getFluentBitConfig() map[string]string {
	return map[string]string{
		"parsers.conf": `
[PARSER]
    Name multiline_pattern
    Format regex
    Regex ^(?<logtime>\d{4}\-\d{2}\-\d{2}T\d{1,2}\:\d{1,2}\:\d{1,2}(\.\d+Z|\+0000)) (?<level>(ERROR|WARN|INFO|DEBUG)) (?<capture>.*)
`,

		"fluent-bit.conf": `
[SERVICE]
    Flush        1
    Verbose      1

    Daemon       Off
    Log_Level    info

    Plugins_File plugins.conf
    Parsers_File parsers.conf

    HTTP_Server  Off
    HTTP_Listen  0.0.0.0
    HTTP_Port    2020

[INPUT]
    Name        tail
    Tag         mongo_writer_logs.${KAI_PRODUCT_ID}
    Buffer_Chunk_Size 1k
    Path        /var/log/app/*.log
    Multiline On
    Parser_Firstline multiline_pattern

[FILTER]
    Name record_modifier
    Match *
    Record versionTag ${KAI_VERSION_TAG}
    Record processName ${KAI_PROCESS_NAME}
    Record workflowName ${KAI_WORKFLOW_NAME}

[FILTER]
    Name  stdout
    Match *

[OUTPUT]
    Name  nats
    Match *
    Host  ${KAI_MESSAGING_HOST}
    Port  ${KAI_MESSAGING_PORT}

`,
	}
}

func mergeConfigs(configs ...map[string]string) map[string]string {
	fullConfig := map[string]string{}

	for _, config := range configs {
		for key, val := range config {
			fullConfig[key] = val
		}
	}

	return fullConfig
}
