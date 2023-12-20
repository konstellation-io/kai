#!/bin/bash
set -xeou pipefail

echo "Creating release..."
# package chart
cr package \
    helm/kai

cr upload \
    --owner "${REPOSITORY_OWNER}" \
    --git-repo "${REPOSITORY_NAME}" \
    --token "${GITHUB_TOKEN}"

# Update index and push to github pages
cr index \
    --owner "${REPOSITORY_OWNER}" \
    --git-repo "${REPOSITORY_NAME}" \
    --token "${GITHUB_TOKEN}" \
    --index-path "index.yaml" \
    --pages-branch "gh-pages" \
    --push
