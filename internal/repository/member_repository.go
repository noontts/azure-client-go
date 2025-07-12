package repository

import (
	"azureclient/internal/model"
	"context"

	"gorm.io/gorm"
)

type MemberRepository interface {
	Create(ctx context.Context, member *model.Member) error
	GetByID(ctx context.Context, id uint) (*model.Member, error)
	Update(ctx context.Context, member *model.Member) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]model.Member, error)
}

type memberRepository struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &memberRepository{db: db}
}

func (r *memberRepository) Create(ctx context.Context, member *model.Member) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *memberRepository) GetByID(ctx context.Context, id uint) (*model.Member, error) {
	var member model.Member
	err := r.db.WithContext(ctx).First(&member, id).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *memberRepository) Update(ctx context.Context, member *model.Member) error {
	return r.db.WithContext(ctx).Save(member).Error
}

func (r *memberRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Member{}, id).Error
}

func (r *memberRepository) List(ctx context.Context) ([]model.Member, error) {
	var members []model.Member
	err := r.db.WithContext(ctx).Find(&members).Error
	return members, err
}
