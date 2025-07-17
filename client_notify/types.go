package client_notify

// EmailResponse is the standard JSON response from all email endpoints.
type EmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// WelcomeEmailRequest for POST /api/email/welcome
type WelcomeEmailRequest struct {
	To       string `json:"to"`
	UserName string `json:"userName"`
}

// SecurityCodeEmailRequest for POST /api/email/security-code
type SecurityCodeEmailRequest struct {
	To   string `json:"to"`
	Code string `json:"code"`
}

// ServiceAlertEmailRequest for POST /api/email/service-alert
type ServiceAlertEmailRequest struct {
	To           string `json:"to"`
	AlertTitle   string `json:"alertTitle"`
	AlertMessage string `json:"alertMessage"`
}

// LoginLinkEmailRequest for POST /api/email/login-link
type LoginLinkEmailRequest struct {
	To        string `json:"to"`
	LoginLink string `json:"loginLink"`
	UserName  string `json:"userName"`
}

// RecoveryCodeEmailRequest for POST /api/email/recovery-code
type RecoveryCodeEmailRequest struct {
	To       string `json:"to"`
	UserName string `json:"userName"`
	Code     string `json:"code"`
}

// GenericEmailRequest for POST /api/email/generic
type GenericEmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}
