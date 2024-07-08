package hx

import (
	"fmt"
	"net/http"
	"net/textproto"
	"net/url"
	"slices"
)

var (
	requestHeaderHXBooted              = textproto.CanonicalMIMEHeaderKey("HX-Boosted")
	requestHeaderHXCurrentURL          = textproto.CanonicalMIMEHeaderKey("HX-Current-URL")
	requestHeaderHistoryRestoreRequest = textproto.CanonicalMIMEHeaderKey("HX-History-Restore-Request")
	requestHeaderHXPrompt              = textproto.CanonicalMIMEHeaderKey("HX-Prompt")
	requestHeaderHXRequest             = textproto.CanonicalMIMEHeaderKey("HX-Request")
	requestHeaderHXTarget              = textproto.CanonicalMIMEHeaderKey("HX-Target")
	requestHeaderHXTriggerName         = textproto.CanonicalMIMEHeaderKey("HX-Trigger-Name")
	requestHeaderHXTriggerID           = textproto.CanonicalMIMEHeaderKey("HX-Trigger")
)

func booleanHeaderIsTrue(h http.Header, canonicalName string) bool {
	values, isSet := h[canonicalName]
	return isSet && slices.Contains(values, "true")
}

func exactlyOneHeaderValue(h http.Header, canonicalName string) (string, bool) {
	values, isSet := h[canonicalName]
	if !isSet || len(values) != 1 {
		return "", false
	}
	return values[0], true
}

func IsBoosted(req *http.Request) bool {
	return booleanHeaderIsTrue(req.Header, requestHeaderHXBooted)
}

// CurrentURL returns the value of the HX-Current-URL header, if present.
func CurrentURL(req *http.Request) (*url.URL, error) {
	if value, isSet := exactlyOneHeaderValue(req.Header, requestHeaderHXCurrentURL); !isSet {
		return nil, fmt.Errorf("header %s is not set exactly one time", requestHeaderHXCurrentURL)
	} else {
		return url.Parse(value)
	}
}

// IsHistoryRestoreRequest reports whether the HX-History-Restore-Request request header set to "true".
func IsHistoryRestoreRequest(req *http.Request) bool {
	return booleanHeaderIsTrue(req.Header, requestHeaderHistoryRestoreRequest)
}

// IsRequest reports whether the HX-Request request header set to "true".
func IsRequest(req *http.Request) bool {
	return booleanHeaderIsTrue(req.Header, requestHeaderHXRequest)
}

// Prompt is the value of the HX-Prompt header, if present.
func Prompt(req *http.Request) (string, bool) {
	return exactlyOneHeaderValue(req.Header, requestHeaderHXPrompt)
}

// Target is the value of the HX-Target header, if present.
func Target(req *http.Request) (string, bool) {
	return exactlyOneHeaderValue(req.Header, requestHeaderHXTarget)
}

// TriggerName is the value of the HX-Trigger-Name header, if present.
func TriggerName(req *http.Request) (string, bool) {
	return exactlyOneHeaderValue(req.Header, requestHeaderHXTriggerName)
}

// TriggerID is the value of the HX-Trigger header, if present.
func TriggerID(req *http.Request) (string, bool) {
	return exactlyOneHeaderValue(req.Header, requestHeaderHXTriggerID)
}
