package auth

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/cantylv/service-happy-birthday/internal/entity"
	uAuth "github.com/cantylv/service-happy-birthday/internal/usecase/auth"
	"github.com/cantylv/service-happy-birthday/internal/utils/functions"
	"github.com/cantylv/service-happy-birthday/internal/utils/myconstants"
	"github.com/cantylv/service-happy-birthday/internal/utils/myerrors"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

type DeliveryLayer struct {
	u uAuth.Usecase
}

func NewDeliveryLayer(usecase uAuth.Usecase) DeliveryLayer {
	return DeliveryLayer{
		u: usecase,
	}
}

func (h *DeliveryLayer) SignUp(w http.ResponseWriter, r *http.Request) {
	logger := zap.Must(zap.NewProduction())

	requestId := functions.GetCtxRequestId(r)
	jwtToken, err := functions.GetJWtToken(r)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		logger.Error(fmt.Sprintf("error while getting jwt token: %v", err), zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrInternal.Error(),
			CodeStatus: http.StatusInternalServerError,
		})
		return
	}
	if jwtToken != "" {
		logger.Info("user is already registered", zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrAlreadyRegistered.Error(),
			CodeStatus: http.StatusUnauthorized,
		})
		return
	}

	body, err := io.ReadAll(r.Body)
	defer func(rBody io.ReadCloser) {
		err = rBody.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("error while close request body: %v", err), zap.String(myconstants.RequestId, requestId))
		}
	}(r.Body)
	if err != nil {
		logger.Info(fmt.Sprintf("error while reading request body: %v", err), zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrInvalidRequestData.Error(),
			CodeStatus: http.StatusBadRequest,
		})
		return
	}

	var signUpData entity.SignUpForm
	err = easyjson.Unmarshal(body, &signUpData)
	if err != nil {
		logger.Info(fmt.Sprintf("error while unmarshalling request body: %v", err), zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrInvalidRequestData.Error(),
			CodeStatus: http.StatusBadRequest,
		})
		return
	}

	_, err = signUpData.Validate()
	if err != nil {
		logger.Info("error while struct validate: "+err.Error(), zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrBadCredentials.Error(),
			CodeStatus: http.StatusBadRequest,
		})
		return
	}

	uId, err := h.u.SignUpUser(r.Context(), signUpData)
	if err != nil {
		if errors.Is(err, myerrors.ErrUserAlreadyExist) {
			logger.Info("user with this email already exist, failed to register", zap.String(myconstants.RequestId, requestId))
			functions.ErrorResponse(functions.ErrorResponseProps{
				W:          w,
				Msg:        myerrors.ErrEmailIsReserved.Error(),
				CodeStatus: http.StatusBadRequest,
			})
			return
		}
		logger.Error(fmt.Sprintf("internal error: %v", err), zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrInternal.Error(),
			CodeStatus: http.StatusInternalServerError,
		})
		return
	}

	w, err = functions.SetCookieAndHeaders(functions.SetCookieProps{
		W:   w,
		Uid: uId,
	})
	if err != nil {
		logger.Error(fmt.Sprintf("error while setting cookie: %v", err), zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrInternal.Error(),
			CodeStatus: http.StatusInternalServerError,
		})
		return
	}
	functions.JsonResponse(functions.JsonResponseProps{
		W:          w,
		Payload:    entity.ResponseDetail{Detail: "you've succesful signed up"},
		CodeStatus: http.StatusOK,
	})
	logger.Info("user successful has signed up", zap.String(myconstants.RequestId, requestId))
}

func (h *DeliveryLayer) SignIn(w http.ResponseWriter, r *http.Request) {
	logger := zap.Must(zap.NewProduction())

	requestId := functions.GetCtxRequestId(r)
	jwtToken, err := functions.GetJWtToken(r)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		logger.Error(fmt.Sprintf("error while getting jwt token: %v", err), zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrInternal.Error(),
			CodeStatus: http.StatusInternalServerError,
		})
		return
	}
	if jwtToken != "" {
		logger.Info("user is already authenticated", zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrAlreadyAuthenticated.Error(),
			CodeStatus: http.StatusUnauthorized,
		})
		return
	}

	body, err := io.ReadAll(r.Body)
	defer func(rBody io.ReadCloser) {
		err = rBody.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("error while close request body: %v", err), zap.String(myconstants.RequestId, requestId))
		}
	}(r.Body)
	if err != nil {
		logger.Info(fmt.Sprintf("error while reading request body: %v", err), zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrInvalidRequestData.Error(),
			CodeStatus: http.StatusBadRequest,
		})
		return
	}

	var signInData entity.SignInForm
	err = easyjson.Unmarshal(body, &signInData)
	if err != nil {
		logger.Info(fmt.Sprintf("error while unmarshalling request body: %v", err), zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrInvalidRequestData.Error(),
			CodeStatus: http.StatusBadRequest,
		})
		return
	}

	_, err = signInData.Validate()
	if err != nil {
		logger.Info("error while struct validate: "+err.Error(), zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrBadCredentials.Error(),
			CodeStatus: http.StatusBadRequest,
		})
		return
	}

	uId, err := h.u.SignInUser(r.Context(), signInData)
	if err != nil {
		if errors.Is(err, myerrors.ErrUserNotExist) || errors.Is(err, myerrors.ErrPwdMismatch) {
			logger.Info("user provided wrong credentials", zap.String(myconstants.RequestId, requestId))
			functions.ErrorResponse(functions.ErrorResponseProps{
				W:          w,
				Msg:        myerrors.ErrBadCredentials.Error(),
				CodeStatus: http.StatusBadRequest,
			})
			return
		}
		logger.Error(fmt.Sprintf("internal error: %v", err), zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrInternal.Error(),
			CodeStatus: http.StatusInternalServerError,
		})
		return
	}

	w, err = functions.SetCookieAndHeaders(functions.SetCookieProps{
		W:   w,
		Uid: uId,
	})
	if err != nil {
		logger.Error(fmt.Sprintf("error while setting cookie: %v", err), zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrInternal.Error(),
			CodeStatus: http.StatusInternalServerError,
		})
		return
	}
	functions.JsonResponse(functions.JsonResponseProps{
		W:          w,
		Payload:    entity.ResponseDetail{Detail: "you've succesful signed in"},
		CodeStatus: http.StatusOK,
	})
	logger.Info("user successful has signed in", zap.String(myconstants.RequestId, requestId))
}

func (h *DeliveryLayer) SignOut(w http.ResponseWriter, r *http.Request) {
	logger := zap.Must(zap.NewDevelopment())

	requestId := functions.GetCtxRequestId(r)
	jwtToken, err := functions.GetJWtToken(r)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		logger.Error(fmt.Sprintf("error while getting jwt token: %v", err), zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrInternal.Error(),
			CodeStatus: http.StatusInternalServerError,
		})
		return
	}
	if jwtToken == "" {
		logger.Info("user is not authenticated", zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrAuth.Error(),
			CodeStatus: http.StatusUnauthorized,
		})
		return
	}

	functions.FlashCookie(w, r)
	functions.JsonResponse(functions.JsonResponseProps{
		W:          w,
		Payload:    entity.ResponseDetail{Detail: "you've succesful signed out"},
		CodeStatus: http.StatusOK,
	})
	logger.Info("user have succesful signed out", zap.String(myconstants.RequestId, requestId))
}
