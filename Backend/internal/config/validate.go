package config

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)


func (c *Config) Validate() error {
    var problems []string
    add := func(format string, args ...any) {
        problems = append(problems, fmt.Sprintf(format, args...))
    }
    // env
    switch c.Env {
    case EnvDevelopment, EnvStaging, EnvProduction:
    default:
        add("env must be one of development/staging/production, got: %q", c.Env)
    }
    // server
    if c.Server.Port == "" {
        add("server.port must not be empty")
    }
    if c.Server.MaxBodyBytes <= 0 {
        add("server.maxBodyBytes must be positive")
    }
    if !slices.Contains([]string{"debug", "release", "test"}, c.Server.RunMode) {
        add("server.runMode must be debug/release/test, got: %q", c.Server.RunMode)
    }
    if c.Server.APIBasePath == "" {
        add("server.apiBasePath must not be empty (e.g. /api or /api/v1)")
    } else if !strings.HasPrefix(c.Server.APIBasePath, "/") {
        add("server.apiBasePath must start with /, got: %q", c.Server.APIBasePath)
    } else if strings.HasSuffix(c.Server.APIBasePath, "/") {
        add("server.apiBasePath must not end with /, got: %q", c.Server.APIBasePath)
    }
    // auth / authz
    switch c.Auth.Posture {
    case PostureEnterprise, PosturePublic:
    default:
        add("auth.posture must be enterprise or public, got: %q", c.Auth.Posture)
    }
    switch c.Authz.Mode {
    case AuthzRoles, AuthzRules, AuthzBoth:
    default:
        add("authz.mode must be roles/rules/both, got: %q", c.Authz.Mode)
    }
    if !c.Authz.Enforce && c.IsProduction() {
        add("authz.enforce cannot be false in production")
    }
    // database
    if c.Postgres.Password == "" {
        add("PostgreSQL password is not set (env var APP_POSTGRES_PASSWORD)")
    }
    if c.Postgres.DBName == "" {
        add("postgres.dbName must not be empty")
    }
    // jwt
    if len(c.JWT.AccessSecret) < 32 {
        add("jwt.accessSecret must be at least 32 bytes (env APP_JWT_ACCESSSECRET). Generate with: openssl rand -base64 48")
    }
    if len(c.JWT.RefreshSecret) < 32 {
        add("jwt.refreshSecret must be at least 32 bytes (env APP_JWT_REFRESHSECRET)")
    }
    if c.JWT.AccessSecret != "" && c.JWT.AccessSecret == c.JWT.RefreshSecret {
        add("jwt.accessSecret and jwt.refreshSecret must be different")
    }
    if c.JWT.AccessTTL <= 0 {
        add("jwt.accessTTL must be positive")
    }
    if c.JWT.RefreshTTL <= c.JWT.AccessTTL {
        add("jwt.refreshTTL (%s) must be greater than jwt.accessTTL (%s)", c.JWT.RefreshTTL, c.JWT.AccessTTL)
    }
    // cors
    hasWildcard := slices.Contains(c.Cors.AllowOrigins, "*")
    if hasWildcard && c.Cors.AllowCredentials {
        add(`cors.allowOrigins cannot be "*" when cors.allowCredentials is enabled; list allowed origins explicitly`)
    }
    for _, origin := range c.Cors.AllowOrigins {
        if origin == "*" {
            continue
        }
        if !strings.HasPrefix(origin, "http://") && !strings.HasPrefix(origin, "https://") {
            add("cors.allowOrigins value must start with http:// or https://, invalid: %q", origin)
        }
        if strings.HasSuffix(origin, "/") {
            add("cors.allowOrigins value must not end with /: %q", origin)
        }
    }
    // logger
    if !slices.Contains([]string{"zerolog", "zap"}, c.Logger.Backend) {
        add("logger.backend must be zerolog or zap, got: %q", c.Logger.Backend)
    }
    if !slices.Contains([]string{"debug", "info", "warn", "error", "fatal"}, c.Logger.Level) {
        add("logger.level is invalid: %q", c.Logger.Level)
    }
    // otp
    if c.Otp.Digits < 4 || c.Otp.Digits > 10 {
        add("otp.digits must be between 4 and 10")
    }
    if c.Otp.MaxAttempts <= 0 {
        add("otp.maxAttempts must be positive to prevent brute-force")
    }
    // password
    if c.Password.MinLength < 8 {
        add("password.minLength must not be less than 8 (NIST recommendation: min 8, prefer 12)")
    }
    if c.Password.MaxLength < c.Password.MinLength {
        add("password.maxLength must not be less than password.minLength")
    }
    // pagination
    if c.Pagination.MaxPageSize <= 0 {
        add("pagination.maxPageSize must be positive")
    }
    if c.Pagination.DefaultPageSize > c.Pagination.MaxPageSize {
        add("pagination.defaultPageSize must not exceed maxPageSize")
    }
    // production-only checks
    if c.IsProduction() {
        if c.Server.RunMode != "release" {
            add("in production, server.runMode must be release")
        }
        if c.Server.InsecureCookies {
            add("in production, server.insecureCookies must be false")
        }
        if c.Postgres.SSLMode == "disable" {
            add("in production, postgres.sslMode must not be disable (use require or verify-full)")
        }
        if c.Redis.Password == "" {
            add("in production, Redis password must not be empty")
        }
        if hasWildcard {
            add(`in production, "*" is not allowed for cors.allowOrigins`)
        }
        if len(c.Cors.AllowOrigins) == 0 {
            add("in production, at least one origin must be set in cors.allowOrigins")
        }
        if c.Logger.LogRequestBody {
            add("in production, logger.logRequestBody must be false")
        }
        if c.Logger.Level == "debug" {
            add("debug log level is not recommended in production")
        }
        if len(c.Server.TrustedProxies) == 0 {
            add("in production, set server.trustedProxies so ClientIP cannot be spoofed")
        }
    }
    if len(problems) > 0 {
        return fmt.Errorf("invalid configuration (%d issues):\n  - %s",
            len(problems), strings.Join(problems, "\n  - "))
    }
    return nil
}
var ErrMissingSecret = errors.New("config: required secret is not set")
