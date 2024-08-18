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
		if errors.Is(err, myerrors.ErrSubscribeNonExistUser) {
			logger.Info("user could not subscribe because there is no such user", zap.String(myconstants.RequestId, requestId))
			functions.ErrorResponse(functions.ErrorResponseProps{
				W:          w,
				Msg:        myerrors.ErrSubscribeNonExistUser.Error(),
				CodeStatus: http.StatusBadRequest,
			})
			return
		}
		if errors.Is(err, myerrors.ErrSubscribeYourself) {
			logger.Info("user has tried to subscribed himself", zap.String(myconstants.RequestId, requestId))
			functions.ErrorResponse(functions.ErrorResponseProps{
				W:          w,
				Msg:        myerrors.ErrSubscribeYourself.Error(),
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
		if errors.Is(err, myerrors.ErrUnsubscribeNonExistUser) {
			logger.Info("user could not unsubscribe because there is no such user", zap.String(myconstants.RequestId, requestId))
			functions.ErrorResponse(functions.ErrorResponseProps{
				W:          w,
				Msg:        myerrors.ErrUnsubscribeNonExistUser.Error(),
				CodeStatus: http.StatusBadRequest,
			})
			return
		}
		if errors.Is(err, myerrors.ErrUnsubscribeYourself) {
			logger.Info("user has tried to unsubscribed himself", zap.String(myconstants.RequestId, requestId))
			functions.ErrorResponse(functions.ErrorResponseProps{
				W:          w,
				Msg:        myerrors.ErrUnsubscribeYourself.Error(),
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
	intervalString := vars["interval"]
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
		if errors.Is(err, myerrors.ErrSetIntervalNonExistUser) {
			logger.Info("user has tried to set the interval for non-existent user birthday", zap.String(myconstants.RequestId, requestId))
			functions.ErrorResponse(functions.ErrorResponseProps{
				W:          w,
				Msg:        myerrors.ErrSetIntervalNonExistUser.Error(),
				CodeStatus: http.StatusBadRequest,
			})
			return
		}
		if errors.Is(err, myerrors.ErrSetIntervalYourself) {
			logger.Info("user has tried to set an interval for his birthday", zap.String(myconstants.RequestId, requestId))
			functions.ErrorResponse(functions.ErrorResponseProps{
				W:          w,
				Msg:        myerrors.ErrSetIntervalYourself.Error(),
				CodeStatus: http.StatusBadRequest,
			})
			return
		}
		if errors.Is(err, myerrors.ErrSetIntervalNotSubscribe) {
			logger.Info("user tried to set an interval for the birthday of an employee who is not subscribed to", zap.String(myconstants.RequestId, requestId))
			functions.ErrorResponse(functions.ErrorResponseProps{
				W:          w,
				Msg:        myerrors.ErrSetIntervalNotSubscribe.Error(),
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
