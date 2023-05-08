package tip

import (
	"context"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	repository "github.com/romankravchuk/muerta/internal/repositories/tip"
)

type TipServicer interface {
	FindTipByID(ctx context.Context, id int) (dto.FindTipDTO, error)
	FindTips(ctx context.Context, filter *dto.TipFilterDTO) ([]dto.FindTipDTO, error)
	CreateTip(ctx context.Context, payload *dto.CreateTipDTO) error
	UpdateTip(ctx context.Context, id int, payload *dto.UpdateTipDTO) error
	DeleteTip(ctx context.Context, id int) error
	RestoreTip(ctx context.Context, id int) error
}

type tipService struct {
	repo repository.TipRepositorer
}

// CreateTip implements TipServicer
func (svc *tipService) CreateTip(ctx context.Context, payload *dto.CreateTipDTO) error {
	model := translate.CreateTipDTOToModel(payload)
	if err := svc.repo.Create(ctx, model); err != nil {
		return err
	}
	return nil
}

// DeleteTip implements TipServicer
func (svc *tipService) DeleteTip(ctx context.Context, id int) error {
	if err := svc.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

// FindTipByID implements TipServicer
func (svc *tipService) FindTipByID(ctx context.Context, id int) (dto.FindTipDTO, error) {
	model, err := svc.repo.FindByID(ctx, id)
	result := translate.TipModelToFindTipDTO(&model)
	if err != nil {
		return dto.FindTipDTO{}, err
	}
	return result, nil
}

// FindTips implements TipServicer
func (svc *tipService) FindTips(ctx context.Context, filter *dto.TipFilterDTO) ([]dto.FindTipDTO, error) {
	models, err := svc.repo.FindMany(ctx, filter.Limit, filter.Offset, filter.Description)
	dtos := translate.TipModelsToFindTipDTOs(models)
	if err != nil {
		return nil, err
	}
	return dtos, nil
}

// RestoreTip implements TipServicer
func (svc *tipService) RestoreTip(ctx context.Context, id int) error {
	if err := svc.repo.Restore(ctx, id); err != nil {
		return err
	}
	return nil
}

// UpdateTip implements TipServicer
func (svc *tipService) UpdateTip(ctx context.Context, id int, payload *dto.UpdateTipDTO) error {
	model, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if payload.Description != "" {
		model.Description = payload.Description
	}
	if err := svc.repo.Update(ctx, model); err != nil {
		return err
	}
	return nil
}

func New(repo repository.TipRepositorer) TipServicer {
	return &tipService{
		repo: repo,
	}
}