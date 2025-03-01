// Package port provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package port

import (
	"time"
)

// Error defines model for Error.
type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

// PubSubMessage defines model for PubSubMessage.
type PubSubMessage struct {
	Message struct {
		Data        string    `json:"data"`
		MessageId   string    `json:"message_id"`
		PublishTime time.Time `json:"publish_time"`
	} `json:"message"`
}

// HandleTelegramMessageJSONBody defines parameters for HandleTelegramMessage.
type HandleTelegramMessageJSONBody = map[string]interface{}

// HandleVideoSaveFailedMessageJSONRequestBody defines body for HandleVideoSaveFailedMessage for application/json ContentType.
type HandleVideoSaveFailedMessageJSONRequestBody = PubSubMessage

// HandleVideoSavedMessageJSONRequestBody defines body for HandleVideoSavedMessage for application/json ContentType.
type HandleVideoSavedMessageJSONRequestBody = PubSubMessage

// HandleTelegramMessageJSONRequestBody defines body for HandleTelegramMessage for application/json ContentType.
type HandleTelegramMessageJSONRequestBody = HandleTelegramMessageJSONBody
