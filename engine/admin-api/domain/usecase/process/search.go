package process

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
)

func (ps *Handler) Search(
	ctx context.Context,
	user *entity.User,
	productID string,
	filter *repository.SearchFilter,
) ([]*entity.RegisteredProcess, error) {
	if filter == nil || *filter == (repository.SearchFilter{}) {
		ps.logger.Info("Retrieving process with no filter", "productID", productID)
	} else {
		ps.logger.Info(
			"Retrieving process with filter",
			"productID", productID, "processType", filter.ProcessType, "processName", filter.ProcessName, "version", filter.Version,
		)
	}

	if filter != nil {
		if err := filter.Validate(); err != nil {
			return nil, err
		}
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
