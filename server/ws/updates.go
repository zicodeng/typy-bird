package ws

import (
	"fmt"
	"net/http"
)

//UpdateHandler handles requests for the /notifications resource
type UpdateHandler struct {
	notifier *Notifier
}

//NewUpdateHandler constructs a new UpdateHandler
func NewUpdateHandler(notifier *Notifier) *UpdateHandler {
	return &UpdateHandler{notifier}
}

//ServeHTTP handles HTTP requests for the UpdateHandler
func (nh *UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("this is where we will send game postioning")
	nh.notifier.Notify([]byte(msg))
}
