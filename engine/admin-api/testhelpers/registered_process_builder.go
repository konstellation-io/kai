package testhelpers

import (
	"fmt"
	"time"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/spf13/viper"
)

type RegisteredProcessBuilder struct {
	product           string
	registeredProcess *entity.RegisteredProcess
}

func NewRegisteredProcessBuilder(product string) *RegisteredProcessBuilder {
	return &RegisteredProcessBuilder{
		product: product,
		registeredProcess: &entity.RegisteredProcess{
			Name:       "process-name",
			Version:    "v1.0.0",
			Type:       entity.ProcessTypeTask.String(),
			UploadDate: time.Now().Truncate(time.Microsecond).UTC(),
			Owner:      "user@email.com",
			Status:     entity.ProcessStatusStarting.String(),
		},
	}
}

func (p *RegisteredProcessBuilder) WithName(name string) *RegisteredProcessBuilder {
	p.registeredProcess.Name = name
	return p
}

func (p *RegisteredProcessBuilder) WithType(processType string) *RegisteredProcessBuilder {
	p.registeredProcess.Type = processType
	return p
}

func (p *RegisteredProcessBuilder) WithOwner(owner string) *RegisteredProcessBuilder {
	p.registeredProcess.Owner = owner
	return p
}

func (p *RegisteredProcessBuilder) WithVersion(version string) *RegisteredProcessBuilder {
	p.registeredProcess.Version = version
	return p
}

func (p *RegisteredProcessBuilder) WithStatus(status string) *RegisteredProcessBuilder {
	p.registeredProcess.Status = status
	return p
}

func (p *RegisteredProcessBuilder) Build() *entity.RegisteredProcess {
	p.registeredProcess.ID = fmt.Sprintf("%s_%s:%s", p.product, p.registeredProcess.Name, p.registeredProcess.Version)
	p.registeredProcess.Image = fmt.Sprintf("%s/%s", viper.GetString(config.RegistryHostKey), p.registeredProcess.ID)

	return p.registeredProcess
}
