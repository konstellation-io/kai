package service

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/service_${GOFILE} -package=mocks

type ProcessRegistry interface {
	DeleteProcess(image, version string) error
}
