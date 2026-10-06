package handler

import (
	"errors"
	"net/http"

	"github.com/MohammadrezaNadirkhanloo/internal/api/middleware"
	"github.com/MohammadrezaNadirkhanloo/internal/api/response"
	"github.com/MohammadrezaNadirkhanloo/internal/config"
	"github.com/MohammadrezaNadirkhanloo/internal/usecase"
	"github.com/MohammadrezaNadirkhanloo/internal/usecase/dto"
	"github.com/MohammadrezaNadirkhanloo/pkg/appctx"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
	"github.com/MohammadrezaNadirkhanloo/pkg/security"
	"github.com/gin-gonic/gin"
)

const RefreshCookieName = "refresh_token"

func refreshCookiePath(cfg *config.Config) string {
	return cfg.Server.APIBasePath + "/auth"
}

type AuthHandler struct {
	users  *usecase.UserUsecase
	tokens *usecase.TokenUsecase
	authz  *usecase.AuthzUsecase
	csrf   *security.CSRFSigner
	cfg    *config.Config
}

func NewAuthHandler(
	users *usecase.UserUsecase,
	tokens *usecase.TokenUsecase,
	authzUC *usecase.AuthzUsecase,
	csrf *security.CSRFSigner,
	cfg *config.Config,
) *AuthHandler {
	return &AuthHandler{users: users, tokens: tokens, authz: authzUC, csrf: csrf, cfg: cfg}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var in dto.RegisterInput
	if !BindJSON(c, &in) {
		return
	}

	out, err := h.users.Register(c.Request.Context(), in)
	if err != nil {
		response.Fail(c, err)
		return
	}

	respondCreated(c, "", out)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var in dto.LoginInput
	if !BindJSON(c, &in) {
		return
	}

	pair, err := h.users.Login(c.Request.Context(), in)
	if err != nil {
		response.Fail(c, err)
		return
	}

	if err := h.setAuthCookies(c, pair); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, pair)
}

// func (h *AuthHandler) SendOtp(c *gin.Context) {
// 	var in dto.SendOtpInput
// 	if !BindJSON(c, &in) {
// 		return
// 	}

// 	if err := h.users.SendOtp(c.Request.Context(), in); err != nil {
// 		response.Fail(c, err)
// 		return
// 	}

// 	response.OK(c, gin.H{"message": "if the number is valid, a code has been sent."})
// }

// func (h *AuthHandler) LoginByOtp(c *gin.Context) {
// 	var in dto.VerifyOtpInput
// 	if !BindJSON(c, &in) {
// 		return
// 	}

// 	pair, err := h.users.LoginByOtp(c.Request.Context(), in)
// 	if err != nil {
// 		response.Fail(c, err)
// 		return
// 	}

// 	if err := h.setAuthCookies(c, pair); err != nil {
// 		response.Fail(c, err)
// 		return
// 	}
// 	response.OK(c, pair)
// }

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := middleware.SingleCookie(c.Request, RefreshCookieName)
	if errors.Is(err, middleware.ErrCookieDuplicate) {
		h.clearAuthCookies(c)
		response.Fail(c, apperror.New(apperror.CodeTokenInvalid,
			"your session is invalid. Please log in again."))
		return
	}
	if err != nil || refreshToken == "" {
		response.Fail(c, apperror.New(apperror.CodeTokenRequired,
			"no session was found. Please log in again."))
		return
	}

	pair, err := h.tokens.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		h.clearAuthCookies(c)
		response.Fail(c, err)
		return
	}

	if err := h.setAuthCookies(c, pair); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, pair)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, _ := middleware.SingleCookie(c.Request, RefreshCookieName)

	if err := h.tokens.Logout(c.Request.Context(), refreshToken); err != nil {
		response.Fail(c, err)
		return
	}

	h.clearAuthCookies(c)
	response.OK(c, gin.H{"message": "you have been logged out successfully."})
}

func (h *AuthHandler) LogoutAll(c *gin.Context) {
	userID, ok := appctx.UserID(c.Request.Context())
	if !ok {
		response.Fail(c, apperror.Unauthorized("you must log in to perform this operation."))
		return
	}

	if err := h.tokens.RevokeAll(c.Request.Context(), userID); err != nil {
		response.Fail(c, err)
		return
	}

	h.clearAuthCookies(c)
	response.OK(c, gin.H{"message": "you have been logged out of all devices."})
}

func (h *AuthHandler) Me(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := appctx.UserID(ctx)
	if !ok {
		response.Fail(c, apperror.Unauthorized("you must log in to perform this operation."))
		return
	}

	user, err := h.users.Find(ctx, userID)
	if err != nil {
		response.Fail(c, err)
		return
	}

	rules, err := h.authz.RulesFor(ctx, user)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.Raw(c, http.StatusOK, dto.NewMeResponse(user, rules))
}

func (h *AuthHandler) setAuthCookies(c *gin.Context, pair dto.TokenPair) error {
	csrfToken, err := h.csrf.Issue(pair.AccessTokenID)
	if err != nil {
		return apperror.Internal(err)
	}

	c.Header("Cache-Control", "no-store")

	accessMaxAge := h.tokens.AccessTTLSeconds()
	h.setCookie(c, middleware.AccessCookieName(h.cfg.Server.InsecureCookies),
		pair.AccessToken, accessMaxAge, true)
	h.setCookie(c, middleware.CSRFCookieName(h.cfg.Server.InsecureCookies),
		csrfToken, accessMaxAge, false)
	h.setRefreshCookie(c, pair.RefreshToken)
	return nil
}

func (h *AuthHandler) clearAuthCookies(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	h.setCookie(c, middleware.AccessCookieName(h.cfg.Server.InsecureCookies), "", -1, true)
	h.setCookie(c, middleware.CSRFCookieName(h.cfg.Server.InsecureCookies), "", -1, false)
	h.clearRefreshCookie(c)
}

func (h *AuthHandler) setCookie(c *gin.Context, name, value string, maxAge int, httpOnly bool) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		Secure:   !h.cfg.Server.InsecureCookies,
		HttpOnly: httpOnly,
		SameSite: http.SameSiteStrictMode,
	})
}

func (h *AuthHandler) setRefreshCookie(c *gin.Context, token string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     RefreshCookieName,
		Value:    token,
		Path:     refreshCookiePath(h.cfg),
		Domain:   h.cookieDomain(),
		MaxAge:   h.tokens.RefreshTTLSeconds(),
		Secure:   !h.cfg.Server.InsecureCookies,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (h *AuthHandler) clearRefreshCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     RefreshCookieName,
		Value:    "",
		Path:     refreshCookiePath(h.cfg),
		Domain:   h.cookieDomain(),
		MaxAge:   -1,
		Secure:   !h.cfg.Server.InsecureCookies,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (h *AuthHandler) cookieDomain() string {
	if h.cfg.Server.Domain == "localhost" || h.cfg.Server.Domain == "" {
		return ""
	}
	return h.cfg.Server.Domain
}
