package services

import (
	"context"
	"errors"
	"fmt"
	"heimdall/internal/config"
	"heimdall/internal/core/domain/entities"
	"heimdall/internal/core/domain/models"
	"heimdall/internal/core/exceptions"
	ports2 "heimdall/internal/core/ports"
	"heimdall/pkg/helpers"
	"strings"
)

var (
	ErrUserAlreadyExists = errors.New("a user already exists with provided email")
	ErrDisposableEmail   = errors.New("provided email is disposable")
)

const (
	SignUpbYEmailApprovalMethod = "confirmation_required"
	SignUpByEmailType           = "email"
	RandomHashLength            = 64
)

type accountService struct {
	accountDBRepo      ports2.AccountRepository
	config             config.Config
	reputationProvider ports2.EchoIPProvider
	accountCacheRepo   ports2.AccountRepository
}

func NewAccountService(accountDBRepo ports2.AccountRepository, config config.Config, reputation ports2.EchoIPProvider, accountCacheRepo ports2.AccountRepository) ports2.AccountService {
	return &accountService{
		accountDBRepo:      accountDBRepo,
		config:             config,
		reputationProvider: reputation,
		accountCacheRepo:   accountCacheRepo,
	}
}

func (a *accountService) GetAccountByUUID(ctx context.Context, uuid string) (*entities.Account, error) {
	return a.accountDBRepo.GetAccountByUUID(ctx, uuid)
}

func (a *accountService) signUpByEmailApproval(ctx context.Context, email string, options map[string]interface{}) (*models.SignUpEmailApprovalResult, error) {
	correctEmail := helpers.CorrectEmail(email, false)
	err := a.checkEmailDisposable(ctx, email)
	if err != nil {
		return nil, err
	}

	userData, err := a.accountDBRepo.GetAccountByEmail(ctx, correctEmail)
	if err != nil && !errors.Is(err, exceptions.ErrNoRecords) {
		return nil, err
	}

	// to-do: add translation
	if userData != nil {
		return nil, ErrUserAlreadyExists
	}
	hash := helpers.RandomString(RandomHashLength, true)

	err = a.accountCacheRepo.CacheEmailVerificationToken(ctx, hash, email)
	if err != nil {
		return nil, err
	}
	baseURL := strings.Trim(a.config.Business.BusinessDomain, "/")
	endpoint := "/signup/verify-email?code=%s&account=%s"
	endpoint = fmt.Sprintf(endpoint, hash, email)

	link := baseURL + endpoint
	for key, value := range options {
		link += "&" + key + "=" + fmt.Sprintf("%v", value)
	}

	// to-do: send email notification
	return &models.SignUpEmailApprovalResult{
		Email:  email,
		Method: SignUpbYEmailApprovalMethod,
		Type:   SignUpByEmailType,
	}, nil
}

func (a *accountService) checkEmailDisposable(ctx context.Context, email string) error {
	info, err := a.reputationProvider.GetEmailDisposableInfo(ctx, email)
	if err != nil {
		return err
	}
	if info.Disposable {
		return ErrDisposableEmail
	}
	return nil
}
