package process

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
)

func (ps *Handler) Search(
	ctx context.Context,
	user *entity.User,
	productID, processType string,
) ([]*entity.RegisteredProcess, error) {
	log := fmt.Sprintf("Retrieving process for product %q", productID)
	if processType != "" {
		log = fmt.Sprintf("%s with process type filter %q", log, processType)
	}

	ps.logger.Info(log)

	filter := repository.SearchFilter{
		ProcessType: entity.ProcessType(processType),
	}

	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("validating list filter: %w", err)
	}

	productProcesses, err := ps.processRepository.SearchByProduct(ctx, productID, filter)
	if err != nil {
		return nil, err
	}

	kaiProcesses, err := ps.processRepository.GlobalSearch(ctx, filter)
	if err != nil {
		return nil, err
	}

	return append(productProcesses, kaiProcesses...), nil
}
