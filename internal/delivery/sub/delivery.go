package sub

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cantylv/service-happy-birthday/internal/entity"
	"github.com/cantylv/service-happy-birthday/internal/usecase/sub"
	"github.com/cantylv/service-happy-birthday/internal/utils/functions"
	"github.com/cantylv/service-happy-birthday/internal/utils/myconstants"
	"github.com/cantylv/service-happy-birthday/internal/utils/myerrors"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type DeliveryLayer struct {
	uc sub.Usecase
}

func NewDeliveryLayer(u sub.Usecase) DeliveryLayer {
	return DeliveryLayer{
		uc: u,
	}
}

func (d *DeliveryLayer) Sub(w http.ResponseWriter, r *http.Request) {
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

	vars := mux.Vars(r)
	err = d.uc.Subscribe(r.Context(), entity.SubProps{
		IdFollower: r.Context().Value(myconstants.UserId).(string),
		IdEmployee: vars["employee_id"],
	})
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
		Payload:    entity.ResponseDetail{Detail: "you've succesful subed on employee"},
		CodeStatus: http.StatusOK,
	})
	logger.Info("user successful has subed", zap.String(myconstants.RequestId, requestId))
}

func (d *DeliveryLayer) Unsub(w http.ResponseWriter, r *http.Request) {
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

	vars := mux.Vars(r)
	err = d.uc.Unsubscribe(r.Context(), entity.SubProps{
		IdFollower: r.Context().Value(myconstants.UserId).(string),
		IdEmployee: vars["employee_id"],
	})
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
		Payload:    entity.ResponseDetail{Detail: "you've succesful unsubed on employee"},
		CodeStatus: http.StatusOK,
	})
	logger.Info("user successful has unsubed", zap.String(myconstants.RequestId, requestId))
}

func (d *DeliveryLayer) ChangeSubInterval(w http.ResponseWriter, r *http.Request) {
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

	vars := mux.Vars(r)
	intervalString := vars["employee_id"]
	interval, err := strconv.Atoi(intervalString)
	if err != nil {
		logger.Info("user provided wrong value of path parameter 'employee_id'", zap.String(myconstants.RequestId, requestId))
		functions.ErrorResponse(functions.ErrorResponseProps{
			W:          w,
			Msg:        myerrors.ErrBadVarPathEmployeeId.Error(),
			CodeStatus: http.StatusBadRequest,
		})
		return
	}
	err = d.uc.ChangeInterval(r.Context(), entity.SetUpIntervalProps{
		Ids: entity.SubProps{
			IdFollower: r.Context().Value(myconstants.UserId).(string),
			IdEmployee: vars["employee_id"],
		},
		NewInterval: uint16(interval),
	})
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
		if errors.Is(err, myerrors.ErrNoSubscription) {
			logger.Info("user is not subscribed on employee", zap.String(myconstants.RequestId, requestId))
			functions.ErrorResponse(functions.ErrorResponseProps{
				W:          w,
				Msg:        myerrors.ErrNoSubscriptionEmployee.Error(),
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
		Payload:    entity.ResponseDetail{Detail: "you've succesful set up new interval"},
		CodeStatus: http.StatusOK,
	})
	logger.Info("user successful has set up new interval", zap.String(myconstants.RequestId, requestId))
}
