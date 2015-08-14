package main

import (
	"regexp"
)

type ProxyParams struct {
	Listen 	string
	Target	string
	Errors	map[string]string
}

func (params *ProxyParams) ValidateForMethod(method string) bool {
	params.Errors = make(map[string]string)

	re := regexp.MustCompile(":[0-9]+$")
	matched := re.Match([]byte(params.Listen))
	if matched == false {
		params.Errors["listen"] = "A 'host:port' `listen` parameter must be provided"
	}

	re = regexp.MustCompile("^.+:[0-9]+$")
	matched = re.Match([]byte(params.Target))
	if method != "DELETE" && matched == false {
		params.Errors["target"] = "A 'host:port' `target` parameter must be provided"
	}

	return len(params.Errors) == 0
}
