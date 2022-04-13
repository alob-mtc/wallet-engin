package util

import (
	uuid "github.com/satori/go.uuid"
)

type ContextKey string

const (
	RequestIdentifier = "REQ_ID"
)

func UniqueStringIdentifier() string {
	return uuid.NewV4().String()
}
