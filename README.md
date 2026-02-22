# wa-to-go

Making it easier to build WhatsApp integrations in Go. Provides generated Go protobuf stubs (types) for WhatsApp Business API webhooks.

## Installation

```sh
go get github.com/ruanspies/wa-to-go
```

## Usage Example

Below is an example showing how to set up a basic Go HTTP server to handle WhatsApp Business API Webhook verifications (`GET`) and incoming event payloads (`POST`).

### 1. `main.go` - Setting up the server

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func init() {
	if os.Getenv("WHATSAPP_ACCESS_TOKEN") == "" {
		log.Fatal("WHATSAPP_ACCESS_TOKEN is not set")
	}
	if os.Getenv("WHATSAPP_VERIFY_TOKEN") == "" {
		log.Fatal("WHATSAPP_VERIFY_TOKEN is not set")
	} else {
		log.Printf("WHATSAPP_VERIFY_TOKEN is set: %s", os.Getenv("WHATSAPP_VERIFY_TOKEN"))
	}
	if os.Getenv("WHATSAPP_APP_SECRET") == "" {
		log.Fatal("WHATSAPP_APP_SECRET is not set")
	}
	if os.Getenv("WHATSAPP_APP_ID") == "" {
		log.Fatal("WHATSAPP_APP_ID is not set")
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/webhook", handshakeHandler).Methods("GET")
	r.HandleFunc("/webhook", messageHandler).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	h2s := &http2.Server{}
	h1s := &http.Server{
		Addr:    ":" + port,
		Handler: h2c.NewHandler(r, h2s),
	}

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(h1s.ListenAndServe())
}
```

### 2. `handlers.go` - Handling verification and unmarshalling Webhooks

```go
package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"google.golang.org/protobuf/encoding/protojson"
	"github.com/ruanspies/wa-to-go"
)

// handshakeHandler handles the verification handshake from Meta.
func handshakeHandler(w http.ResponseWriter, r *http.Request) {
	// Extract query parameters into our protobuf message structure.
	payload := &watogo.VerificationRequest{
		Hub: &watogo.VerificationRequest_Hub{
			Mode:        r.URL.Query().Get("hub.mode"),
			VerifyToken: r.URL.Query().Get("hub.verify_token"),
			Challenge:   r.URL.Query().Get("hub.challenge"),
		},
	}
	log.Printf("%v", payload)

	// The verification token should match the one configured in the Meta App Dashboard.
	verifyToken := os.Getenv("WHATSAPP_VERIFY_TOKEN")

	if payload.GetHub().GetMode() == "subscribe" && payload.GetHub().GetVerifyToken() == verifyToken {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(payload.GetHub().GetChallenge()))
		return
	}

	// Respond with 403 Forbidden if the verification fails.
	w.WriteHeader(http.StatusForbidden)
}

// messageHandler handles incoming POST requests with webhook payloads.
func messageHandler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Unmarshal the JSON into the WebhookPayload protobuf message
	var payload watogo.WebhookPayload
	if err := protojson.Unmarshal(body, &payload); err != nil {
		log.Printf("Error parsing request body: %v\n\n%s", err, string(body))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("EVENT_RECEIVED"))
		return
	}

	log.Printf("%v", payload.GetEntry())

	// Check if this is a WhatsApp Business Account event
	if payload.GetObject() == "whatsapp_business_account" {
		for _, entry := range payload.GetEntry() {
			for _, change := range entry.GetChanges() {
				value := change.GetValue()
				switch {
				case value.GetMessages() != nil:
					// For handling messages let's respond immediately
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("EVENT_RECEIVED"))

					// Example: send messages to a handler function
					/*
					handleMessages(&handlerMessage{
						messages: value.GetMessages(),
						metadata: value.GetMetadata(),
					})
					*/
				case value.GetStatuses() != nil:
					for i, status := range value.GetStatuses() {
						log.Println(i, status.GetStatus())
					}
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("EVENT_RECEIVED"))
				case value.GetErrors() != nil:
					log.Printf("Errors: %v", value.GetErrors())
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("EVENT_RECEIVED"))
				default:
					log.Printf("Unknown change type: %v", change.GetField())
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("EVENT_RECEIVED"))
				}
			}
		}

	} else {
		// Respond with 404 if the event is not from a WhatsApp Business Account
		w.WriteHeader(http.StatusNotFound)
	}
}
```
