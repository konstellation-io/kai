#!/bin/bash
set -xeou pipefail

export HELM_DOCS_VERSION=1.11.3

curl -LO  https://github.com/norwoodj/helm-docs/releases/download/v${HELM_DOCS_VERSION}/helm-docs_${HELM_DOCS_VERSION}_Linux_x86_64.tar.gz
tar -zxvf helm-docs_${HELM_DOCS_VERSION}_Linux_x86_64.tar.gz && chmod +x helm-docs && rm helm-docs_${HELM_DOCS_VERSION}_Linux_x86_64.tar.gz

./helm-docs \
    --chart-search-root=./helm/kai \
    --template-files=CHART.md.gotmpl \
    --output-file=CHART.md

rm helm-docs README.md LICENSE
