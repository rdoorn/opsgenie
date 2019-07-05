package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
)

func replaceNonAlphanum(s, replacement string) string {
	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9._-]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(s, replacement)
}

func filter2query(filter *og.Filter) (string, error) {
	matchType := ""
	switch filter.ConditionMatchType {
	case og.MatchAll:
		matchType = " AND "
	case og.MatchAllConditions:
		matchType = " AND "
	case og.MatchAnyCondition:
		matchType = " OR "
	}

	str := []string{}
	for _, f := range filter.Conditions {

		var key string

		switch f.Field {
		case og.Message:
			key = "message"
		case og.Alias:
			key = "alias"
		case og.Description:
			key = "description"
		case og.Source:
			key = "source"
		case og.Entity:
			key = "entity"
		case og.Tags:
			key = "tags"
		case og.Actions:
			key = "actions"
		case og.Details:
			key = "details"
		case og.Recipients:
			key = "recipients"
		case og.Teams:
			key = "teams"
		case og.Priority:
			key = "priority"
		case og.ExtraProperties:
			key = f.Key
		default:
			return "", fmt.Errorf("unknown field in search result: %s", f.Field)
		}

		not := ""
		if f.IsNot {
			not = "NOT "
		}

		value := ""
		expectedValue := escapeValue(f.ExpectedValue)
		switch f.Operation {
		case og.Contains:
			value = fmt.Sprintf(": *%s*", expectedValue)
		case og.Matches:
			value = expectedValue
		case og.StartsWith:
			value = fmt.Sprintf(": %s*", expectedValue)
		case og.EndsWith:
			value = fmt.Sprintf(": %s*", expectedValue)
		case og.IsEmpty:
			value = ": \"\""
		case og.GreaterThan:
			value = fmt.Sprintf("> %s*", expectedValue)
		case og.LessThan:
			value = fmt.Sprintf("< %s*", expectedValue)
		}

		str = append(str, fmt.Sprintf("%s%s%s", not, key, value))
	}
	result := strings.Join(str, matchType)

	return result, nil
}

func escapeValue(v string) string {
	v = strings.Replace(v, ":", "\\:", -1)
	return v
}
