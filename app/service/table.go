package service

import (
	"context"
	"errors"
	"pos-api/app/model"
	"pos-api/app/repository"
)

var (
	ErrInvalidStatus = errors.New("invalid table status")
)

type TableService struct {
	Repo repository.TableRepository
}

func NewTableService(repo repository.TableRepository) *TableService {
	return &TableService{Repo: repo}
}

func (s *TableService) Create(ctx context.Context, seat int, status string) (*model.Table, error) {
	if seat <= 0 {
		seat = 2
	}
	if status == "" {
		status = "available"
	}
	if !isValidStatus(status) {
		return nil, ErrInvalidStatus
	}
	tb := &model.Table{Seat: seat, Status: status}
	if err := s.Repo.Create(ctx, tb); err != nil {
		return nil, err
	}
	return tb, nil
}

func (s *TableService) List(ctx context.Context) ([]model.Table, error) {
	return s.Repo.List(ctx)
}

func (s *TableService) UpdateStatus(ctx context.Context, id uint64, status string) error {
	if !isValidStatus(status) {
		return ErrInvalidStatus
	}
	return s.Repo.UpdateStatus(ctx, id, status)
}

func isValidStatus(st string) bool {
	switch st {
	case "available", "occupied", "closed":
		return true
	default:
		return false
	}
}
