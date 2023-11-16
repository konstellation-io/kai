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
    Name log_parser
    Format regex
    Regex ^(?<logtime>[^\s]*)\s+(?<level>[^\s]*)\s+(?<issuer>\[.*\])\s+(?<origin>[^\s]*)\s+(?<message>\S+(\s\S+)*)\s+(?<metadata_json>{.*})$
    Decode_Field json metadata_json

[PARSER]
    Name json_parser
    Format json
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
    Tag         tail.log
    Buffer_Chunk_Size 1k
    Path        /var/log/app/*.log

[FILTER]
    Name parser
    Match tail.log
    Key_Name log
    Parser json_parser
    Reserve_Data True

[OUTPUT]
    Name stdout
    Match *

[OUTPUT]
    Name loki
    Match tail.log
    Host kai-local-loki
    Port 3100
    labels service=kai-product-version
    label_keys $product_id, $version_id, $workflow_id, $process_id, $request_id, $L, $N
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
