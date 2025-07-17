package core_config

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

// AppConfig holds the application-specific configurations.
type AppConfig struct {
	Route        string   // App route, e.g., "/files", "/documents"
	AppName      string   // App name, e.g., "files", "documents"
	DisplayName  string   // Display name, e.g., "Files", "Documents"
	Domain       string   // Domain without protocol, e.g., "files.hstles.com"
	AuthMethods  []string // Allowed authentication methods: "google", "microsoftonline", "github"
	Icon         string   // Path to app icon
	Illustration string   // Path to illustration image
}

// appConfigs holds all application configurations.
var AppConfigs = []AppConfig{
	{Route: "/", AppName: "services", DisplayName: "Services", Domain: "files.hstles.com", AuthMethods: []string{"google", "microsoftonline", "github", "email"}, Icon: "assets/media/app/hstles.png", Illustration: "assets/media/app/box.png"},
	{Route: "/files", AppName: "files", DisplayName: "Files", Domain: "files.hstles.com", AuthMethods: []string{"google", "microsoftonline", "github"}, Icon: "assets/media/app/hstles.png", Illustration: "assets/media/app/box-files.svg"},
	// {Route: "/app", AppName: "app", DisplayName: "App", Domain: "files.hstles.com", AuthMethods: []string{"google", "microsoftonline"}, Icon: "assets/media/app/hstles.png", Illustration: "/static/assets/media/app/hstles.png"},
	// {Route: "/sign", AppName: "sign", DisplayName: "Sign", Domain: "sign.hstles.com", AuthMethods: []string{"google", "microsoftonline"}, Icon: "assets/media/app/hstles.png", Illustration: "/static/assets/media/app/hstles.png"},
	// {Route: "/sigil", AppName: "sigil", DisplayName: "Sigil", Domain: "sigil.hstles.com", AuthMethods: []string{"google", "microsoftonline"}, Icon: "assets/media/app/hstles.png", Illustration: "/static/assets/media/app/hstles.png"},
	{Route: "/organisation", AppName: "organisation", DisplayName: "Organisation", Domain: "organisation.hstles.com", AuthMethods: []string{"google", "microsoftonline"}, Icon: "assets/media/app/hstles.png", Illustration: "/static/assets/media/app/hstles.png"},
	{Route: "/support", AppName: "support", DisplayName: "Support", Domain: "support.hstles.com", AuthMethods: []string{"google", "microsoftonline"}, Icon: "assets/media/app/hstles.png", Illustration: "/static/assets/media/app/hstles.png"},
	{Route: "/account", AppName: "account", DisplayName: "My Account", Domain: "account.hstles.com", AuthMethods: []string{"google", "microsoftonline", "github", "email"}, Icon: "assets/media/app/hstles.png", Illustration: "/static/assets/media/app/hstles.png"},
}

// GetAppByRoute retrieves the AppConfig based on the route.
func GetAppByRoute(route string) AppConfig {
	for _, config := range AppConfigs {
		if config.Route == route {
			log.Printf("GetAppByRoute: Matched route '%s' to AppConfig: %+v", route, config)
			return config
		}
	}
	log.Printf("GetAppByRoute: Route '%s' not found, defaulting to root AppConfig: %+v", route, AppConfigs[2])
	return AppConfigs[2] // Default to root configuration
}

// GetAppByName retrieves the AppConfig based on the app name.
func GetAppByName(name string) AppConfig {
	for _, config := range AppConfigs {
		if config.AppName == name {
			log.Printf("GetAppByName: Matched name '%s' to AppConfig: %+v", name, config)
			return config
		}
	}
	log.Printf("GetAppByName: Name '%s' not found, defaulting to root AppConfig: %+v", name, AppConfigs[2])
	return AppConfigs[2] // Default to root configuration
}

// GetAppByDomain retrieves the AppConfig based on the domain.
func GetAppByDomain(domain string) AppConfig {
	for _, config := range AppConfigs {
		if config.Domain == domain {
			log.Printf("GetAppByDomain: Matched domain '%s' to AppConfig: %+v", domain, config)
			return config
		}
	}
	log.Printf("GetAppByDomain: Domain '%s' not found, defaulting to root AppConfig: %+v", domain, AppConfigs[0])
	return AppConfigs[0] // Default to root configuration
}

// GetURL constructs the full URL from the domain.
func (ac AppConfig) GetURL() string {
	return "https://" + ac.Domain
}

// ValidateNextParameter validates the 'next' parameter and defaults to "files" if invalid.
func ValidateNextParameter(next string) (string, error) {
	if next == "" {
		// Default to "files" if 'next' is missing
		appConfig := GetAppByRoute("/")
		return appConfig.GetURL(), nil
	}

	if IsValidURL(next) {
		// If 'next' is a valid URL under *.hstles.com
		return next, nil
	}

	appConfig := GetAppByName(next)
	if appConfig.AppName != "" {
		// If 'next' matches a valid app name
		return appConfig.GetURL(), nil
	}

	// Default to "files" if 'next' is invalid
	appConfig = GetAppByRoute("/")
	return appConfig.GetURL(), fmt.Errorf("invalid 'next' parameter: %s, using default route", next)
}

// isValidURL checks if a string is a valid URL and belongs to the "*.hstles.com" domain.
func IsValidURL(str string) bool {
	parsedURL, err := url.ParseRequestURI(str)
	if err != nil {
		return false // Not a valid URL
	}

	// Check if the URL host ends with ".hstles.com"
	if strings.HasSuffix(parsedURL.Host, ".hstles.com") {
		return true
	}

	return false
}
