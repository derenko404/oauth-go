package controllers

import (
	"fmt"
	"time"

	"github.com/derenko404/ipapi-go"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"oauth-go/internal/app"
	"oauth-go/internal/middleware"
	"oauth-go/internal/store"
	"oauth-go/pkg/cookieutils"
	"oauth-go/pkg/response"
)

var (
	DeviceIdCookieName = "device_id"
)

type authController struct {
	app *app.App
}

func NewAuthController(app *app.App) *authController {
	return &authController{
		app: app,
	}
}

func getDeviceID(ctx *gin.Context) string {
	deviceID := cookieutils.SetIfNotExists(
		ctx,
		DeviceIdCookieName,
		uuid.New().String(),
		cookieutils.OneHour,
		"/",
		"",
		true,
		true,
	)

	return deviceID
}

func getLocation(ip string) (string, error) {
	location := "Unknown, Unknown"
	resp, err := ipapi.GetIpLocation(ip)

	if err != nil {
		return location, nil
	}

	return fmt.Sprintf("%s, %s", resp.CountryName, resp.City), nil
}

func generateState(secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * 15).Unix(), // 15 min expiry
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func validateState(jwtSecret []byte, tokenString string) bool {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check that the signing method is HMAC (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})

	return err == nil
}

type signInResponse struct {
	URL string `json:"url"`
}

// @Summary     Sign In
// @Description Redirects to selected OAuth provider login URL, not working in swagger
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       provider path string true "Selected provider, available options: google, github"
// @Success     200 {object} response.APISuccessResponse{data=signInResponse}
// @Failure     400 {object} response.APIErrorResponse
// @Router      /auth/sign-in/{provider} [get]
func (controller *authController) SignIn(ctx *gin.Context) {
	provider := ctx.Param("provider")

	state, err := generateState([]byte(controller.app.Config.JwtSecret))

	if err != nil {
		controller.app.Logger.Error("cannot generate state", "error", err)
		response.RespondError(ctx, response.ErrOAuth)
		return
	}

	url, err := controller.app.Services.OAuth.GetSignInUrl(provider, state)

	if err != nil {
		controller.app.Logger.Error("error during sign-in", "error", err)
		response.RespondError(ctx, response.ErrOAuth)
		return
	}

	response.RespondSuccess(ctx, &signInResponse{
		URL: url,
	})
}

type handleCallbackResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	DeviceID     string `json:"device_id"`
}

// @Summary		Endpoint for OAuth providers
// @Description	This endpoint should be called only by OAuth providers
// @Tags			  auth
// @Accept			json
// @Produce		  json
// @Param state path string true "OAuth state string"
// @Param code path string true "OAuth code"
// @Success     200 {object} response.APISuccessResponse{data=handleCallbackResponse}
// @Failure		  400	{object} response.APIErrorResponse
// @Failure		  422	{object} response.APIErrorResponse
// @Failure		  500	{object} response.APIErrorResponse
// @Router			/auth/handle-callback [get]
func (controller *authController) HandleCallback(ctx *gin.Context) {
	provider := ctx.Param("provider")

	var query struct {
		State string `form:"state" binding:"required"`
		Code  string `form:"code" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&query); err != nil {
		controller.app.Logger.Error("error binding query", "error", err)
		response.RespondError(ctx, response.ErrInvalidInput)
		return
	}

	valid := validateState([]byte(controller.app.Config.JwtSecret), query.State)

	if !valid {
		controller.app.Logger.Info("invalid oauth state")
		response.RespondError(ctx, response.ErrInvalidInput)
		return
	}

	profile, err := controller.app.Services.OAuth.GetProfile(ctx.Request.Context(), provider, query.Code)

	if err != nil {
		controller.app.Logger.Error("cannot get profile", "error", err)
		response.RespondError(ctx, response.ErrOAuth)
		return
	}

	var user *store.User

	filters := map[string]any{"email": profile.Email}
	user, err = controller.app.Store.User.GetUserBy(ctx.Request.Context(), filters)

	if err != nil {
		controller.app.Logger.Error("cannot get user", "error", err)

		user, err = controller.app.Store.User.CreateUser(ctx.Request.Context(), &store.UserDto{
			Name:           profile.Name,
			Email:          profile.Email,
			AvatarURL:      profile.AvatarURL,
			Provider:       provider,
			ProviderUserID: profile.ID,
		})

		if err != nil {
			controller.app.Logger.Error("failed to create user", "error", err)
			response.RespondError(ctx, response.ErrInternalServerError)
			return
		}
	}

	clientIP := ctx.ClientIP()
	location, err := getLocation(clientIP)

	if err != nil {
		controller.app.Logger.Info("cannot get location for", "ip", clientIP, "error", err.Error())
	}

	deviceID := getDeviceID(ctx)

	filters = map[string]any{
		"user_id":   user.ID,
		"device_id": deviceID,
	}

	// @TODO
	// add session version
	// increase it each time when user logs in to invalidate old tokens
	// issued for that session
	session, err := controller.app.Store.Session.GetSessionBy(ctx.Request.Context(), filters)

	if err != nil {
		session, err = controller.app.Store.Session.CreateSession(ctx.Request.Context(), &store.UserSessionDto{
			UserID:    user.ID,
			IPAddress: clientIP,
			UserAgent: ctx.GetHeader("User-Agent"),
			Location:  location,
			DeviceID:  deviceID,
		})

		if err != nil {
			controller.app.Logger.Error("failed to create session", "error", err)
			response.RespondError(ctx, response.ErrInternalServerError)
			return
		}
	}

	accessToken, refreshToken := controller.app.Services.Jwt.IssueTokensPair(user.ID, session.ID, user.Email)

	controller.app.Logger.Info("user signed in", "user", user, "session", session)

	response.RespondSuccess(ctx, &handleCallbackResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		DeviceID:     deviceID,
	})
}

type getMeResponse struct {
	User *store.User `json:"user"`
}

// @Summary		Me
// @Description	Returns current user information
// @Tags			  auth
// @Security BearerAuth
// @Accept			json
// @Produce		  json
// @Success     200 {object} response.APISuccessResponse{data=getMeResponse}
// @Failure		  400	{object} response.APIErrorResponse
// @Failure		  422	{object} response.APIErrorResponse
// @Failure		  500	{object} response.APIErrorResponse
// @Router			/auth/me [get]
func (controller *authController) GetMe(ctx *gin.Context) {
	user, err := middleware.MustGetUserFromContext(ctx)

	if err != nil {
		response.RespondError(ctx, response.ErrUnauthorized)
		return
	}

	response.RespondSuccess(ctx, &getMeResponse{User: user})
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type refreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// @Summary		Refresh Token
// @Description	Refreshe jwt token
// @Tags			  auth
// @Accept			json
// @Produce		  json
// @Param refresh_token body refreshTokenRequest true "jwt refresh token"
// @Success     200 {object} response.APISuccessResponse{data=refreshTokenResponse}
// @Failure		  403	{object} response.APIErrorResponse
// @Failure		  422	{object} response.APIErrorResponse
// @Router			/auth/refresh [post]
func (controller *authController) RefreshToken(ctx *gin.Context) {
	var req refreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.RespondError(ctx, response.ErrInvalidInput)
		return
	}

	token, err := controller.app.Services.Jwt.VerifyToken(req.RefreshToken)
	if err != nil {
		response.RespondError(ctx, response.ErrUnauthorized)
		return
	}

	claims, err := controller.app.Services.Jwt.GetClaims(token)
	if err != nil {
		controller.app.Logger.Error("error during request processing", "error", err)
		response.RespondError(ctx, response.ErrUnauthorized)
		return
	}

	filters := map[string]any{
		"id": claims.SessionID,
	}

	session, err := controller.app.Store.Session.GetSessionBy(ctx.Request.Context(), filters)

	if err != nil {
		controller.app.Logger.Error("error during request processing", "error", err)
		response.RespondError(ctx, response.ErrUnauthorized)
		return
	}

	accessToken, refreshToken := controller.app.Services.Jwt.IssueTokensPair(claims.UserID, session.ID, claims.Email)

	response.RespondSuccess(ctx, &refreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

type signOutResponse struct{}

// @Summary		Sign Out
// @Description	Sign out current user
// @Tags			  auth
// @Security BearerAuth
// @Accept			json
// @Produce		  json
// @Success     200 {object} response.APISuccessResponse
// @Failure		  403	{object} response.APIErrorResponse
// @Failure		  500	{object} response.APIErrorResponse
// @Router			/sign-out [get]
func (controller *authController) SignOut(ctx *gin.Context) {
	user, err := middleware.MustGetUserFromContext(ctx)

	if err != nil {
		response.RespondError(ctx, response.ErrUnauthorized)
		return
	}

	filters := map[string]any{
		"user_id":   user.ID,
		"device_id": getDeviceID(ctx),
	}

	err = controller.app.Store.Session.DeleteSessionBy(ctx.Request.Context(), filters)
	if err != nil {
		controller.app.Logger.Error("error deleting session", "error", err)
		response.RespondError(ctx, response.ErrInternalServerError)
		return
	}

	response.RespondSuccess(ctx, &signOutResponse{})
}
