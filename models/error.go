package models

import "errors"

var (
	ErrorNotFound       = errors.New("not found")
	ErrorDuplicateEntry = errors.New("duplicate entry")
	ErrorWrongParams    = errors.New("wrong parameters")
	ErrorUnsupported    = errors.New("unsupported")
	ErrorNotAllowed     = errors.New("action not allowed")
)

type VerifyCodeError error

var (
	VerifyErrorNil       VerifyCodeError = errors.New("verify_error_nil")
	VerifyServiceError   VerifyCodeError = errors.New("verify_server_error")
	EmailFormatError     VerifyCodeError = errors.New("email_format_error")
	EmailAlreadyUsed     VerifyCodeError = errors.New("email_in_use")
	EmailCompanyNotFound VerifyCodeError = errors.New("email_company_not_found")
	SendEmailRateExceed  VerifyCodeError = errors.New("send_email_rate_exceed")
	VerifyCodeRateExceed VerifyCodeError = errors.New("verify_code_rate_exceed")
	WrongVerifyCode      VerifyCodeError = errors.New("wrong_verify_code")
	VerifyCodeExpired    VerifyCodeError = errors.New("verify_code_expired")
)
