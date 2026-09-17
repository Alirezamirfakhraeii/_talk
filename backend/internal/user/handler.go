package user

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/auth"
	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/platform/httpresponse"
)

const maxAvatarSize = 5 << 20

var allowedAvatarTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

type Handler struct {
	service *Service
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateProfileRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
}

type registerResponseData struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type currentUserResponseData struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Bio        string `json:"bio"`
	AvatarPath string `json:"avatar_path"`
}

type loginResponseData struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type searchUserResponseData struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	AvatarPath string `json:"avatar_path"`
}

type profileResponseData struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Bio        string    `json:"bio"`
	AvatarPath string    `json:"avatar_path"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (handler *Handler) Register(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	var requestBody registerRequest

	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		httpresponse.Error(
			responseWriter,
			http.StatusBadRequest,
			"INVALID_JSON",
			"invalid request body",
			nil,
		)
		return
	}

	registeredUser, err := handler.service.Register(
		request.Context(),
		RegisterInput{
			Name:     requestBody.Name,
			Username: requestBody.Username,
			Email:    requestBody.Email,
			Password: requestBody.Password,
		},
	)

	if err != nil {
		var validationErrors ValidationErrors

		if errors.As(err, &validationErrors) {
			httpresponse.Error(
				responseWriter,
				http.StatusBadRequest,
				"VALIDATION_ERROR",
				"invalid input",
				validationErrors,
			)
			return
		}

		if errors.Is(err, ErrEmailAlreadyExists) {
			httpresponse.Error(
				responseWriter,
				http.StatusConflict,
				"EMAIL_ALREADY_EXISTS",
				"email already in use",
				nil,
			)
			return
		}

		if errors.Is(err, ErrUsernameAlreadyExists) {
			httpresponse.Error(
				responseWriter,
				http.StatusConflict,
				"USERNAME_ALREADY_EXISTS",
				"username already in use",
				nil,
			)
			return
		}

		httpresponse.Error(
			responseWriter,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"internal server error",
			nil,
		)
		return
	}

	httpresponse.Success(
		responseWriter,
		http.StatusCreated,
		"user registered successfully",
		registerResponseData{
			ID:        registeredUser.ID,
			Name:      registeredUser.Name,
			Username:  registeredUser.Username,
			Email:     registeredUser.Email,
			CreatedAt: registeredUser.CreatedAt,
		},
	)
}

func (handler *Handler) Login(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	var requestBody loginRequest

	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		httpresponse.Error(
			responseWriter,
			http.StatusBadRequest,
			"INVALID_JSON",
			"invalid request body",
			nil,
		)
		return
	}

	loginResult, err := handler.service.Login(
		request.Context(),
		LoginInput{
			Email:    requestBody.Email,
			Password: requestBody.Password,
		},
	)

	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			httpresponse.Error(
				responseWriter,
				http.StatusUnauthorized,
				"INVALID_CREDENTIALS",
				"invalid email or password",
				nil,
			)
			return
		}

		httpresponse.Error(
			responseWriter,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"internal server error",
			nil,
		)
		return
	}

	http.SetCookie(
		responseWriter,
		&http.Cookie{
			Name:     "samatalk_session",
			Value:    loginResult.Token,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			Expires:  loginResult.ExpiresAt,
		},
	)

	httpresponse.Success(
		responseWriter,
		http.StatusOK,
		"login successful",
		loginResponseData{
			ID:       loginResult.User.ID,
			Name:     loginResult.User.Name,
			Username: loginResult.User.Username,
			Email:    loginResult.User.Email,
		},
	)
}

func (handler *Handler) Me(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	currentUser, ok := auth.CurrentUserFromContext(
		request.Context(),
	)

	if !ok {
		httpresponse.Error(
			responseWriter,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"authentication required",
			nil,
		)
		return
	}

	profile, err := handler.service.GetProfile(
		request.Context(),
		currentUser.ID,
	)

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			httpresponse.Error(
				responseWriter,
				http.StatusNotFound,
				"USER_NOT_FOUND",
				"user not found",
				nil,
			)
			return
		}

		httpresponse.Error(
			responseWriter,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"internal server error",
			nil,
		)
		return
	}

	httpresponse.Success(
		responseWriter,
		http.StatusOK,
		"current user retrieved successfully",
		currentUserResponseData{
			ID:         profile.ID,
			Name:       profile.Name,
			Username:   profile.Username,
			Email:      profile.Email,
			Bio:        profile.Bio,
			AvatarPath: profile.AvatarPath,
		},
	)
}

func (handler *Handler) Search(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	currentUser, ok := auth.CurrentUserFromContext(request.Context())
	if !ok {
		httpresponse.Error(
			responseWriter,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"authentication required",
			nil,
		)
		return
	}

	searchTerm := request.URL.Query().Get("q")

	users, err := handler.service.SearchUsers(
		request.Context(),
		currentUser.ID,
		searchTerm,
	)

	if err != nil {
		httpresponse.Error(
			responseWriter,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"internal server error",
			nil,
		)
		return
	}

	results := make([]searchUserResponseData, 0, len(users))

	for _, foundUser := range users {
		results = append(
			results,
			searchUserResponseData{
				ID:         foundUser.ID,
				Name:       foundUser.Name,
				Username:   foundUser.Username,
				AvatarPath: foundUser.AvatarPath,
			},
		)
	}

	httpresponse.Success(
		responseWriter,
		http.StatusOK,
		"users retrieved successfully",
		results,
	)
}

func (handler *Handler) UpdateProfile(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	currentUser, ok := auth.CurrentUserFromContext(request.Context())
	if !ok {
		httpresponse.Error(
			responseWriter,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"authentication required",
			nil,
		)
		return
	}

	var requestBody updateProfileRequest

	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		httpresponse.Error(
			responseWriter,
			http.StatusBadRequest,
			"INVALID_JSON",
			"invalid request body",
			nil,
		)
		return
	}

	updatedUser, err := handler.service.UpdateProfile(
		request.Context(),
		currentUser.ID,
		UpdateProfileInput{
			Name:     requestBody.Name,
			Username: requestBody.Username,
			Bio:      requestBody.Bio,
		},
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrProfileNameRequired):
			httpresponse.Error(
				responseWriter,
				http.StatusUnprocessableEntity,
				"NAME_REQUIRED",
				"name is required",
				nil,
			)

		case errors.Is(err, ErrProfileUsernameRequired):
			httpresponse.Error(
				responseWriter,
				http.StatusUnprocessableEntity,
				"USERNAME_REQUIRED",
				"username is required",
				nil,
			)

		case errors.Is(err, ErrProfileBioTooLong):
			httpresponse.Error(
				responseWriter,
				http.StatusUnprocessableEntity,
				"BIO_TOO_LONG",
				"bio must not exceed 160 characters",
				nil,
			)

		case errors.Is(err, ErrUsernameAlreadyExists):
			httpresponse.Error(
				responseWriter,
				http.StatusConflict,
				"USERNAME_ALREADY_EXISTS",
				"username already in use",
				nil,
			)

		case errors.Is(err, ErrUserNotFound):
			httpresponse.Error(
				responseWriter,
				http.StatusNotFound,
				"USER_NOT_FOUND",
				"user not found",
				nil,
			)

		default:
			httpresponse.Error(
				responseWriter,
				http.StatusInternalServerError,
				"INTERNAL_ERROR",
				"internal server error",
				nil,
			)
		}

		return
	}

	httpresponse.Success(
		responseWriter,
		http.StatusOK,
		"profile updated successfully",
		profileResponseData{
			ID:         updatedUser.ID,
			Name:       updatedUser.Name,
			Username:   updatedUser.Username,
			Email:      updatedUser.Email,
			Bio:        updatedUser.Bio,
			AvatarPath: updatedUser.AvatarPath,
			UpdatedAt:  updatedUser.UpdatedAt,
		},
	)
}

func (handler *Handler) UploadAvatar(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	currentUser, ok := auth.CurrentUserFromContext(request.Context())
	if !ok {
		httpresponse.Error(
			responseWriter,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"authentication required",
			nil,
		)
		return
	}

	request.Body = http.MaxBytesReader(
		responseWriter,
		request.Body,
		maxAvatarSize+(1<<20),
	)

	if err := request.ParseMultipartForm(maxAvatarSize); err != nil {
		httpresponse.Error(
			responseWriter,
			http.StatusBadRequest,
			"INVALID_AVATAR",
			"avatar is too large or invalid",
			nil,
		)
		return
	}

	file, _, err := request.FormFile("avatar")
	if err != nil {
		httpresponse.Error(
			responseWriter,
			http.StatusBadRequest,
			"AVATAR_REQUIRED",
			"avatar file is required",
			nil,
		)
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(
		io.LimitReader(file, maxAvatarSize+1),
	)
	if err != nil {
		httpresponse.Error(
			responseWriter,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"could not read avatar",
			nil,
		)
		return
	}

	if len(fileData) > maxAvatarSize {
		httpresponse.Error(
			responseWriter,
			http.StatusRequestEntityTooLarge,
			"AVATAR_TOO_LARGE",
			"avatar must not exceed 5 MB",
			nil,
		)
		return
	}

	contentType := http.DetectContentType(fileData)

	extension, allowed := allowedAvatarTypes[contentType]
	if !allowed {
		httpresponse.Error(
			responseWriter,
			http.StatusUnsupportedMediaType,
			"INVALID_AVATAR_TYPE",
			"avatar must be JPEG, PNG, or WebP",
			nil,
		)
		return
	}

	randomBytes := make([]byte, 8)

	if _, err := rand.Read(randomBytes); err != nil {
		httpresponse.Error(
			responseWriter,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"could not generate avatar filename",
			nil,
		)
		return
	}

	fileName := fmt.Sprintf(
		"%d_%s%s",
		currentUser.ID,
		hex.EncodeToString(randomBytes),
		extension,
	)

	avatarDirectory := filepath.Join(
		"uploads",
		"avatars",
	)

	if err := os.MkdirAll(avatarDirectory, 0755); err != nil {
		httpresponse.Error(
			responseWriter,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"could not create avatar directory",
			nil,
		)
		return
	}

	filePath := filepath.Join(
		avatarDirectory,
		fileName,
	)

	if err := os.WriteFile(
		filePath,
		fileData,
		0644,
	); err != nil {
		httpresponse.Error(
			responseWriter,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"could not save avatar",
			nil,
		)
		return
	}

	avatarPath := "/uploads/avatars/" + fileName

	updatedUser, err := handler.service.UpdateAvatar(
		request.Context(),
		currentUser.ID,
		avatarPath,
	)

	if err != nil {
		_ = os.Remove(filePath)

		if errors.Is(err, ErrUserNotFound) {
			httpresponse.Error(
				responseWriter,
				http.StatusNotFound,
				"USER_NOT_FOUND",
				"user not found",
				nil,
			)
			return
		}

		httpresponse.Error(
			responseWriter,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"internal server error",
			nil,
		)
		return
	}

	httpresponse.Success(
		responseWriter,
		http.StatusOK,
		"avatar updated successfully",
		profileResponseData{
			ID:         updatedUser.ID,
			Name:       updatedUser.Name,
			Username:   updatedUser.Username,
			Email:      updatedUser.Email,
			Bio:        updatedUser.Bio,
			AvatarPath: updatedUser.AvatarPath,
			UpdatedAt:  updatedUser.UpdatedAt,
		},
	)
}
