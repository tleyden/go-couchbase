package couchbase

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// Regex to extract just the function name (and not the module path)
var RE_stripFnPreamble = regexp.MustCompile(`^.*\.(.*)$`)

// CleanupHost returns the hostname with the given suffix removed.
func CleanupHost(h, commonSuffix string) string {
	if strings.HasSuffix(h, commonSuffix) {
		return h[:len(h)-len(commonSuffix)]
	}
	return h
}

// FindCommonSuffix returns the longest common suffix from the given
// strings.
func FindCommonSuffix(input []string) string {
	rv := ""
	if len(input) < 2 {
		return ""
	}
	from := input
	for i := len(input[0]); i > 0; i-- {
		common := true
		suffix := input[0][i:]
		for _, s := range from {
			if !strings.HasSuffix(s, suffix) {
				common = false
				break
			}
		}
		if common {
			rv = suffix
		}
	}
	return rv
}

// ParseURL is a wrapper around url.Parse with some sanity-checking
func ParseURL(urlStr string) (result *url.URL, err error) {
	result, err = url.Parse(urlStr)
	if result != nil && result.Scheme == "" {
		result = nil
		err = fmt.Errorf("invalid URL <%s>", urlStr)
	}
	return
}

// Records the name of a function and the time of entry.  Intended to be used as:
//   defer base.TraceExit(base.TraceEnter())
// So that the base.TraceExit() -- which will be called when the calling function exits --
// will get the function name and time it was entered, so it can calc the delta and log it.
func TraceEnter() (functionName string, timeEntered time.Time) {

	functionName = "<unknown>"
	// Skip this function, and fetch the PC and file for its parent
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		functionName = RE_stripFnPreamble.ReplaceAllString(
			runtime.FuncForPC(pc).Name(),
			"$1",
		)
	}

	return functionName, time.Now()

}

// Like TraceEnter, but allows you to pass an extra identifier.
func TraceEnterExtra(extraIdentifier string) (functionName string, timeEntered time.Time) {

	functionName = "<unknown>"
	// Skip this function, and fetch the PC and file for its parent
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		functionName = RE_stripFnPreamble.ReplaceAllString(
			runtime.FuncForPC(pc).Name(),
			"$1",
		)
	}

	if extraIdentifier != "" {
		// attach the extra identifier if there is one
		functionName = fmt.Sprintf("%v-%v", functionName, extraIdentifier)
	}

	return functionName, time.Now()

}

func TraceExit(functionName string, timeEntered time.Time) {

	delta := time.Since(timeEntered)
	if delta.Seconds() >= 1 {
		log.Printf("%v() took %v seconds", functionName, delta.Seconds())
	}

}
