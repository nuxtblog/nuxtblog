// Package listener registers all built-in event handlers.
//
// To add a new listener group (e.g. webhook, audit log, task system):
//  1. Create a new file in this package (e.g. webhook.go).
//  2. Write a registerXxxListeners() function that calls event.On / event.OnAsync.
//  3. Call it from init() below.
//
// This package is imported via a blank import in internal/logic/logic.go,
// so handlers are registered before the HTTP server starts serving requests.
package listener

func init() {
	registerNotificationListeners()
	registerWebhookListeners()
	registerCommerceListeners()
}
