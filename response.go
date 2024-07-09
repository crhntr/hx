package hx

import (
	"encoding/json"
	"net/http"
	"strings"
)

const (
	responseHeaderLocation = "HX-Location"
	responseHeaderPushURL  = "HX-Push-Url"
	responseHeaderRedirect = "HX-Redirect"
	responseHeaderRefresh  = "HX-Refresh"
	responseHeaderReplace  = "HX-Replace"
	responseHeaderReswap   = "HX-Reswap"
	responseHeaderRetarget = "HX-Retarget"
	responseHeaderReselect = "HX-Reselect"

	responseHeaderTrigger            = "HX-Trigger"
	responseHeaderTriggerAfterSettle = "HX-Trigger-After-Settle"
	responseHeaderTriggerAfterSwap   = "HX-Trigger-After-Settle"
)

// SetLocation sets the HX-Location header to the given URL.
// See https://htmx.org/headers/hx-location/
func SetLocation(res http.ResponseWriter, url string) {
	res.Header().Set(responseHeaderLocation, url)
}

// SetPushURL sets the HX-Push-Url header to the given URL.
// See https://htmx.org/headers/hx-push-url/
func SetPushURL(res http.ResponseWriter, url string) {
	res.Header().Set(responseHeaderPushURL, url)
}

func SetPreventHistoryUpdate(res http.ResponseWriter) {
	res.Header().Set(responseHeaderPushURL, "false")
}

// SetRedirect sets the HX-Redirect header to the given URL.
// See https://htmx.org/reference/#response_headers
func SetRedirect(res http.ResponseWriter, url string) {
	res.Header().Set(responseHeaderRedirect, url)
}

// SetRefresh sets the HX-Refresh header to "true".
// It tells the client to refresh the page.
func SetRefresh(res http.ResponseWriter) {
	res.Header().Set(responseHeaderRefresh, "true")
}

// SetReplaceURL sets the HX-Replace header to the given URL.
// See https://htmx.org/headers/hx-replace/
func SetReplaceURL(res http.ResponseWriter, url string) {
	res.Header().Set(responseHeaderReplace, url)
}

// SetSwap sets the HX-Reswap header to the given value.
func SetSwap(res http.ResponseWriter, s Swap) {
	res.Header().Set(responseHeaderReswap, string(s))
}

// SetTarget sets the HX-Retarget header to the given value.
func SetTarget(res http.ResponseWriter, value string) {
	res.Header().Set(responseHeaderRetarget, value)
}

// SetSelect sets the HX-Reselect header to the given value.
func SetSelect(res http.ResponseWriter, value string) {
	res.Header().Set(responseHeaderReselect, value)
}

func triggeredEventNames(res http.ResponseWriter, name string, events []string) {
	if len(events) > 0 {
		res.Header().Set(name, strings.Join(events, ", "))
	}
}

func SetTriggerEvents(res http.ResponseWriter, events ...string) {
	triggeredEventNames(res, responseHeaderTrigger, events)
}

func SetTriggerAfterSettleEvents(res http.ResponseWriter, events ...string) {
	triggeredEventNames(res, responseHeaderTriggerAfterSettle, events)
}

func SetTriggerAfterSwapEvents(res http.ResponseWriter, events ...string) {
	triggeredEventNames(res, responseHeaderTriggerAfterSwap, events)
}

func marshalJSONHeaderData(res http.ResponseWriter, header string, data any) error {
	buffer, err := json.Marshal(data)
	if err != nil {
		return err
	}
	res.Header().Set(header, string(buffer))
	return nil
}

func SetTriggerWithData(res http.ResponseWriter, data EventsWithData) error {
	return marshalJSONHeaderData(res, responseHeaderTrigger, data)
}

func SetTriggerAfterSettleWithData(res http.ResponseWriter, data EventsWithData) error {
	return marshalJSONHeaderData(res, responseHeaderTriggerAfterSettle, data)
}

func SetTriggerAfterSwapWithData(res http.ResponseWriter, data EventsWithData) error {
	return marshalJSONHeaderData(res, responseHeaderTriggerAfterSwap, data)
}

type EventsWithData struct {
	data map[string]json.RawMessage
	err  error
}

func NewEventsWithData(event string, data any) EventsWithData {
	e := EventsWithData{
		data: make(map[string]json.RawMessage),
	}
	_ = e.Set(event, data)
	return e
}

func (e EventsWithData) Set(event string, data any) EventsWithData {
	if e.err != nil {
		return e
	}
	buf, err := json.Marshal(data)
	if err != nil {
		e.err = err
		delete(e.data, event)
		return e
	}
	e.data[event] = buf
	return e
}

func (e EventsWithData) MarshalJSON() ([]byte, error) {
	if e.err != nil {
		return nil, e.err
	}
	buf, err := json.Marshal(e.data)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
