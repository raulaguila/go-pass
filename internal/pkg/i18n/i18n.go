package i18n

import (
	"errors"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var I18nTranslations map[string]*Translation = map[string]*Translation{}

func NewTranslation(localizer *i18n.Localizer) *Translation {
	translation := &Translation{}
	translation.loadTranslations(localizer)
	return translation
}

type Translation struct {
	ErrGeneric              error
	ErrInvalidId            error
	ErrInvalidDatas         error
	ErrManyRequest          error
	ErrorNonexistentRoute   error
	ErrUndefinedColumn      error
	ErrExpiredToken         error
	ErrDisabledUser         error
	ErrIncorrectPassword    error
	ErrPassUnmatch          error
	ErrUserHasPass          error
	ErrInvalidIpAssociation error

	ErrProfileUsed       error
	ErrProfileNotFound   error
	ErrProfileRegistered error

	ErrUserUsed       error
	ErrUserNotFound   error
	ErrUserRegistered error

	ErrSiteUsed       error
	ErrSiteNotFound   error
	ErrSiteRegistered error

	ErrPhoneUsed       error
	ErrPhoneNotFound   error
	ErrPhoneRegistered error

	ErrOperatorUsed       error
	ErrOperatorNotFound   error
	ErrOperatorRegistered error

	ErrAccountUsed       error
	ErrAccountNotFound   error
	ErrAccountRegistered error
}

func (s *Translation) loadTranslations(localizer *i18n.Localizer) {
	s.ErrGeneric = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrGeneric"}, PluralCount: 1}))
	s.ErrInvalidId = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrInvalidId"}, PluralCount: 1}))
	s.ErrInvalidDatas = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrInvalidDatas"}, PluralCount: 1}))
	s.ErrManyRequest = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrManyRequest"}, PluralCount: 1}))
	s.ErrorNonexistentRoute = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrorNonexistentRoute"}, PluralCount: 1}))
	s.ErrUndefinedColumn = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrUndefinedColumn"}, PluralCount: 1}))
	s.ErrExpiredToken = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrExpiredToken"}, PluralCount: 1}))
	s.ErrDisabledUser = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrDisabledUser"}, PluralCount: 1}))
	s.ErrIncorrectPassword = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrIncorrectPassword"}, PluralCount: 1}))
	s.ErrPassUnmatch = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrPassUnmatch"}, PluralCount: 1}))
	s.ErrUserHasPass = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrUserHasPass"}, PluralCount: 1}))
	s.ErrInvalidIpAssociation = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrInvalidIpAssociation"}, PluralCount: 1}))

	s.ErrProfileUsed = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrProfileUsed"}, PluralCount: 1}))
	s.ErrProfileNotFound = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrProfileNotFound"}, PluralCount: 1}))
	s.ErrProfileRegistered = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrProfileRegistered"}, PluralCount: 1}))

	s.ErrUserUsed = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrUserUsed"}, PluralCount: 1}))
	s.ErrUserNotFound = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrUserNotFound"}, PluralCount: 1}))
	s.ErrUserRegistered = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrUserRegistered"}, PluralCount: 1}))

	s.ErrSiteUsed = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrSiteUsed"}, PluralCount: 1}))
	s.ErrSiteNotFound = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrSiteNotFound"}, PluralCount: 1}))
	s.ErrSiteRegistered = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrSiteRegistered"}, PluralCount: 1}))

	s.ErrPhoneUsed = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrPhoneUsed"}, PluralCount: 1}))
	s.ErrPhoneNotFound = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrPhoneNotFound"}, PluralCount: 1}))
	s.ErrPhoneRegistered = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrPhoneRegistered"}, PluralCount: 1}))

	s.ErrOperatorUsed = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrOperatorUsed"}, PluralCount: 1}))
	s.ErrOperatorNotFound = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrOperatorNotFound"}, PluralCount: 1}))
	s.ErrOperatorRegistered = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrOperatorRegistered"}, PluralCount: 1}))

	s.ErrAccountUsed = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrAccountUsed"}, PluralCount: 1}))
	s.ErrAccountNotFound = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrAccountNotFound"}, PluralCount: 1}))
	s.ErrAccountRegistered = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrAccountRegistered"}, PluralCount: 1}))
}
