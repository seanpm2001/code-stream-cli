/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

func getExecutions(id string, status string, name string, nested bool) ([]*CodestreamAPIExecutions, error) {
	var arrExecutions []*CodestreamAPIExecutions
	if id != "" {
		x, err := getExecution("/codestream/api/executions/" + id)
		if err != nil {
			return nil, err
		}
		arrExecutions = append(arrExecutions, x)
		return arrExecutions, err
	}
	client := resty.New()
	var qParams = make(map[string]string)
	qParams["$orderby"] = "_requestTimeInMicros desc"

	var filters []string
	if status != "" {
		filters = append(filters, "(status eq '"+strings.ToUpper(status)+"')")
	}
	if name != "" {
		filters = append(filters, "(name eq '"+name+"')")
	}
	if nested {
		filters = append(filters, "(_nested eq '"+strconv.FormatBool(nested)+"')")
	}
	if project != "" {
		filters = append(filters, "(project eq '"+project+"')")
	}
	if len(filters) > 0 {
		qParams["$filter"] = "(" + strings.Join(filters, ") and (") + ")"
		log.Debugln(qParams["$filter"])
	}

	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&documentsList{}).
		SetError(&CodeStreamException{}).
		SetAuthToken(targetConfig.accesstoken).
		Get("https://" + targetConfig.server + "/pipeline/api/executions")
	if queryResponse.IsError() {
		//return nil, queryResponse.Error().(error)
		return nil, errors.New(queryResponse.Error().(*CodeStreamException).Message)
	}

	for _, value := range queryResponse.Result().(*documentsList).Documents {
		c := CodestreamAPIExecutions{}
		mapstructure.Decode(value, &c)
		arrExecutions = append(arrExecutions, &c)
	}
	return arrExecutions, err
}

func getExecution(executionLink string) (*CodestreamAPIExecutions, error) {
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&CodestreamAPIExecutions{}).
		SetAuthToken(targetConfig.accesstoken).
		Get("https://" + targetConfig.server + executionLink)
	if queryResponse.IsError() {
		return nil, queryResponse.Error().(error)
	}
	return queryResponse.Result().(*CodestreamAPIExecutions), err
}

func deleteExecution(id string) (*CodestreamAPIExecutions, error) {
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&CodestreamAPIExecutions{}).
		SetAuthToken(targetConfig.accesstoken).
		Delete("https://" + targetConfig.server + "/pipeline/api/executions/" + id)
	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*CodeStreamException).Message)
	}
	return queryResponse.Result().(*CodestreamAPIExecutions), err
}

func createExecution(id string, inputs string, comment string) (*CodeStreamCreateExecutionResponse, error) {
	// Convert JSON string to byte array
	var inputBytes = []byte(inputs)
	// Unmarshal inputs using a generic interface
	var inputsInterface interface{}
	err := json.Unmarshal(inputBytes, &inputsInterface)
	if err != nil {
		return nil, err
	}
	// Create CodeStreamCreateExecutionRequest struct
	var execution CodeStreamCreateExecutionRequest
	execution.Comments = comment
	execution.Input = inputsInterface
	//Marshal struct to JSON []byte
	executionBytes, err := json.Marshal(execution)
	if err != nil {
		return nil, err
	}
	client := resty.New()
	queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Content-Type", "application/json").
		SetBody(executionBytes).
		SetResult(&CodeStreamCreateExecutionResponse{}).
		SetAuthToken(targetConfig.accesstoken).
		Post("https://" + targetConfig.server + "/pipeline/api/pipelines/" + id + "/executions")
	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*CodeStreamException).Message)
	}
	return queryResponse.Result().(*CodeStreamCreateExecutionResponse), nil
}
