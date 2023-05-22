module github.com/konstellation-io/kai/libs/krt-utils

go 1.20

require (
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/docker/distribution v2.8.2+incompatible
	github.com/go-playground/validator/v10 v10.11.0
	github.com/konstellation-io/kai/libs/simplelogger v0.0.0-20201224090044-7d2e9c2cfd32
	github.com/mattn/go-zglob v0.0.3
	github.com/stretchr/testify v1.8.2
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
)

replace github.com/konstellation-io/kai/libs/simplelogger => ../simplelogger
