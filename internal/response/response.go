package response

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Message struct {
	HttpStatus      int    `json:"-"`
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
}

type respHandler struct {
	w http.ResponseWriter
}

func Write(h http.ResponseWriter) *respHandler {
	return &respHandler{w: h}
}

func (h *respHandler) JSON(msg any) {
	h.w.Header().Add("Content-Type", "application/json")
	h.w.WriteHeader(http.StatusOK)
	json.NewEncoder(h.w).Encode(msg)
}

func (h *respHandler) AbortWithJSON(msg *Message) {
	h.w.Header().Add("Content-Type", "application/json")
	h.w.WriteHeader(msg.HttpStatus)
	json.NewEncoder(h.w).Encode(msg)
}

func BadResponse(httpStatus int, serviceCode string, responseCode string, responseMessage string) *Message {
	respCode := strings.Join([]string{
		strconv.Itoa(httpStatus),
		serviceCode,
		responseCode}, "")

	return &Message{
		HttpStatus:      httpStatus,
		ResponseCode:    respCode,
		ResponseMessage: responseMessage,
	}

}
