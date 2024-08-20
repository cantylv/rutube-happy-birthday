package user

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/cantylv/service-happy-birthday/internal/entity"
	"github.com/cantylv/service-happy-birthday/internal/usecase/user"
	"github.com/cantylv/service-happy-birthday/internal/utils/functions"
	"github.com/cantylv/service-happy-birthday/internal/utils/myconstants"
	"github.com/cantylv/service-happy-birthday/internal/utils/myerrors"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

type DeliveryLayer struct {
	uc user.Usecase
}

func NewDeliveryLayer(usecase user.Usecase) DeliveryLayer {
	return DeliveryLayer{
		uc: usecase,
	}
}

func (d *DeliveryLayer) GetUser(w http.ResponseWriter, r *http.Request) {
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

	employeeIdValue := r.Context().Value(myconstants.UserId)
	if employeeIdValue == nil {
		logger.Error("error while receiving employeeId", zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrInternal.Error(),
			CodeStatus: http.StatusInternalServerError,
		})
		return
	}
	u, err := d.uc.GetData(r.Context(), employeeIdValue.(string))
	if err != nil {
		if errors.Is(err, myerrors.ErrUserNotExist) {
			logger.Info("user is not exist", zap.String(myconstants.RequestId, requestId))
			functions.ErrorResponse(functions.ErrorResponseProps{
				W:          w,
				Msg:        myerrors.ErrUserExist.Error(),
				CodeStatus: http.StatusUnauthorized,
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

	functions.JsonResponse(functions.JsonResponseProps{
		W:          w,
		Payload:    u,
		CodeStatus: http.StatusOK,
	})
	logger.Info("user data was successful received", zap.String(myconstants.RequestId, requestId))
}

func (d *DeliveryLayer) UpdateUser(w http.ResponseWriter, r *http.Request) {
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

	var updateData entity.UserUpdate
	err = easyjson.Unmarshal(body, &updateData)
	if err != nil {
		logger.Info(fmt.Sprintf("error while unmarshalling request body: %v", err), zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrInvalidRequestData.Error(),
			CodeStatus: http.StatusBadRequest,
		})
		return
	}

	_, err = updateData.Validate()
	if err != nil {
		logger.Info("error while struct validate: "+err.Error(), zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrBadCredentials.Error(),
			CodeStatus: http.StatusBadRequest,
		})
		return
	}

	userIdValue := r.Context().Value(myconstants.UserId)
	if userIdValue == nil {
		logger.Error("error while receiving employeeId", zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrInternal.Error(),
			CodeStatus: http.StatusInternalServerError,
		})
		return
	}

	err = d.uc.UpdateData(r.Context(), &updateData, userIdValue.(string))
	if err != nil {
		if errors.Is(err, myerrors.ErrEmailIsReserved) {
			logger.Info("user with this email already exist, failed to update", zap.String(myconstants.RequestId, requestId))
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

	functions.JsonResponse(functions.JsonResponseProps{
		W:          w,
		Payload:    entity.ResponseDetail{Detail: "you've succesful updated data"},
		CodeStatus: http.StatusOK,
	})
	logger.Info("user successful has updated his data", zap.String(myconstants.RequestId, requestId))
}
