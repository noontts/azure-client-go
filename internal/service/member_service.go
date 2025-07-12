package service

import (
	"azureclient/internal/errs"
	"azureclient/internal/model"
	"azureclient/internal/otel"
	"azureclient/internal/repository"
	"context"

	"go.opentelemetry.io/otel/attribute"
)

type MemberService interface {
	CreateMember(ctx context.Context, member *model.Member) error
	GetMemberByID(ctx context.Context, id uint) (*model.Member, error)
	UpdateMember(ctx context.Context, member *model.Member) error
	DeleteMember(ctx context.Context, id uint) error
	ListMembers(ctx context.Context) ([]model.Member, error)
}

type memberService struct {
	repos repository.Repositories
	// Add other dependencies like Client, AzureClient here as needed
}

func NewMemberService(repos repository.Repositories) MemberService {
	return &memberService{repos: repos}
}

func (s *memberService) CreateMember(ctx context.Context, member *model.Member) error {
	ctx, span := otel.Tracer.Start(ctx, "CreateMember")
	defer span.End()
	span.SetAttributes(attribute.String("member.email", member.Email))
	return s.repos.Member.Create(ctx, member)
}

func (s *memberService) GetMemberByID(ctx context.Context, id uint) (*model.Member, error) {
	ctx, span := otel.Tracer.Start(ctx, "GetMemberByID")
	defer span.End()
	span.SetAttributes(attribute.Int("member.id", int(id)))
	if id == 9999 {
		span.RecordError(errs.UnableToProceed)
		return nil, errs.UnableToProceed
	}
	return s.repos.Member.GetByID(ctx, id)
}

func (s *memberService) UpdateMember(ctx context.Context, member *model.Member) error {
	ctx, span := otel.Tracer.Start(ctx, "UpdateMember")
	defer span.End()
	span.SetAttributes(attribute.Int("member.id", int(member.ID)))
	return s.repos.Member.Update(ctx, member)
}

func (s *memberService) DeleteMember(ctx context.Context, id uint) error {
	ctx, span := otel.Tracer.Start(ctx, "DeleteMember")
	defer span.End()
	span.SetAttributes(attribute.Int("member.id", int(id)))
	return s.repos.Member.Delete(ctx, id)
}

func (s *memberService) ListMembers(ctx context.Context) ([]model.Member, error) {
	ctx, span := otel.Tracer.Start(ctx, "ListMembers")
	defer span.End()
	return s.repos.Member.List(ctx)
}
