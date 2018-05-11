package main

import (
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// getStringValue get string value associated to a label
func getStringValue(labels map[string]string, labelName string, defaultValue string) string {
	if value, ok := labels[labelName]; ok && len(value) > 0 {
		return value
	}
	return defaultValue
}

// getBoolValue get bool value associated to a label
func getBoolValue(labels map[string]string, labelName string, defaultValue bool) bool {
	rawValue, ok := labels[labelName]
	if ok {
		v, err := strconv.ParseBool(rawValue)
		if err == nil {
			return v
		}
	}
	return defaultValue
}

// getIntValue get int value associated to a label
func getIntValue(labels map[string]string, labelName string, defaultValue int) int {
	if rawValue, ok := labels[labelName]; ok {
		value, err := strconv.Atoi(rawValue)
		if err == nil {
			return value
		}
		log.Errorf("Unable to parse %q: %q, falling back to %v. %v", labelName, rawValue, defaultValue, err)
	}
	return defaultValue
}

// getInt64Value get int64 value associated to a label
func getInt64Value(labels map[string]string, labelName string, defaultValue int64) int64 {
	if rawValue, ok := labels[labelName]; ok {
		value, err := strconv.ParseInt(rawValue, 10, 64)
		if err == nil {
			return value
		}
		log.Errorf("Unable to parse %q: %q, falling back to %v. %v", labelName, rawValue, defaultValue, err)
	}
	return defaultValue
}

// getSliceStringValue get a slice of string associated to a label
func getSliceStringValue(labels map[string]string, labelName string) []string {
	var value []string

	if values, ok := labels[labelName]; ok {
		value = splitAndTrimString(values, ",")

		if len(value) == 0 {
			log.Debugf("Could not load %q.", labelName)
		}
	}
	return value
}

// has Check if a value is associated to a label
func has(labels map[string]string, labelName string) bool {
	value, ok := labels[labelName]
	return ok && len(value) > 0
}

// hasPrefix Check if a value is associated to a less one label with a prefix
func hasPrefix(labels map[string]string, prefix string) bool {
	for name, value := range labels {
		if strings.HasPrefix(name, prefix) && len(value) > 0 {
			return true
		}
	}
	return false
}

// splitAndTrimString splits separatedString at the separator character and trims each
// piece, filtering out empty pieces. Returns the list of pieces or nil if the input
// did not contain a non-empty piece.
func splitAndTrimString(base string, sep string) []string {
	var trimmedStrings []string

	for _, s := range strings.Split(base, sep) {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			trimmedStrings = append(trimmedStrings, s)
		}
	}

	return trimmedStrings
}
