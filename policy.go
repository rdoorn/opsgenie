package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/opsgenie/opsgenie-go-sdk-v2/policy"
)

// policyList shows all the policies
func (h handler) policyList() error {
	//create a policy client
	policyClient, err := policy.NewClient(h.client)

	if err != nil {
		return fmt.Errorf("error occured while creating policy client")
	}

	result, err := policyClient.ListAlertPolicies(nil, &policy.ListAlertPoliciesRequest{TeamId: h.config.TeamID})
	if err != nil {
		return err
	}

	fmt.Printf("Policies:\n")
	for id, r := range result.Policies {
		enabled := ""
		if r.Enabled {
			enabled = "enabled"
		} else {
			enabled = "disabled"
		}
		fmt.Printf("%2d: [%8s] %s \n", id, enabled, r.Name)
	}

	return nil
}

// policyTest can test a policy against a filter
func (h handler) policyTest(policyID int, timeframe string) error {
	//create a policy client
	policyClient, err := policy.NewClient(h.client)

	if err != nil {
		return fmt.Errorf("error occured while creating policy client")
	}

	result, err := policyClient.ListAlertPolicies(nil, &policy.ListAlertPoliciesRequest{TeamId: h.config.TeamID})
	if err != nil {
		return err
	}

	if len(result.Policies) < policyID {
		return fmt.Errorf("policy %d not found\n", policyID)
	}

	if policyID < 0 {
		return fmt.Errorf("policy not specified\n")
	}

	p := result.Policies[policyID]
	details, err := policyClient.GetAlertPolicy(nil, &policy.GetAlertPolicyRequest{Id: p.Id, TeamId: h.config.TeamID})
	if err != nil {
		return err
	}

	log.Printf("type:%+v conditions:%+v\n", details.Filter.ConditionMatchType, details.Filter.Conditions)
	query, err := filter2query(details.Filter)
	if err != nil {
		return err
	}
	log.Printf("string filter: %s", query)

	if timeframe != "" {
		history := subTime(time.Now(), timeframe)
		query = fmt.Sprintf("(%s) AND createdAt > %d", query, history.Unix())
	}

	return h.listAlertsQuery(query)
}

func (h handler) policyDisable(policyID int) error {
	//create a policy client
	policyClient, err := policy.NewClient(h.client)
	if err != nil {
		return fmt.Errorf("error occured while creating policy client")
	}

	// get policy details by the id provided
	result, err := policyClient.ListAlertPolicies(nil, &policy.ListAlertPoliciesRequest{TeamId: h.config.TeamID})
	if err != nil {
		return err
	}

	if len(result.Policies) < policyID {
		return fmt.Errorf("policy %d not found\n", policyID)
	}

	if policyID < 0 {
		return fmt.Errorf("policy not specified\n")
	}

	p := result.Policies[policyID]

	//disable policy
	_, err = policyClient.DisablePolicy(nil, &policy.DisablePolicyRequest{TeamId: h.config.TeamID, Id: p.Id, Type: policy.AlertPolicy})
	if err != nil {
		return err
	}
	fmt.Printf("Policy disabled successfuly\n")
	return nil
}

func (h handler) policyEnable(policyID int, timeStr string) error {
	//create a policy client
	policyClient, err := policy.NewClient(h.client)
	if err != nil {
		return fmt.Errorf("error occured while creating policy client")
	}

	result, err := policyClient.ListAlertPolicies(nil, &policy.ListAlertPoliciesRequest{TeamId: h.config.TeamID})
	if err != nil {
		return err
	}

	if len(result.Policies) < policyID {
		return fmt.Errorf("policy %d not found\n", policyID)
	}

	if policyID < 0 {
		return fmt.Errorf("policy not specified\n")
	}

	p := result.Policies[policyID]

	if timeStr == "" {

		//enable policy
		_, err = policyClient.EnablePolicy(nil, &policy.EnablePolicyRequest{TeamId: h.config.TeamID, Id: p.Id, Type: policy.AlertPolicy})
		if err != nil {
			return err
		}
		fmt.Printf("Policy enabled successfuly\n")
	} else {
		return h.maintenanceCreate(p.Id, p.Name, timeStr)
	}
	return nil
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
	//date := time.Now().Add(-2 * time.Hour).Unix()
	result := strings.Join(str, matchType)

	return result, nil
	//return fmt.Sprintf("%s AND createdAt > %d", result, date)
	//return fmt.Sprintf("%s AND lastOccurredAt > %d", result, date)
}

func escapeValue(v string) string {
	v = strings.Replace(v, ":", "\\:", -1)
	return v
}
