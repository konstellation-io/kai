package service

import "context"

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/service_${GOFILE} -package=mocks

type ProcessRegistry interface {
	DeleteProcess(ctx context.Context, image, version string) error
}
