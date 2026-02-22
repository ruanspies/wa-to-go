// Package watogo provides generated Go protobuf stubs (types) for WhatsApp Business API webhooks.
//
// These structures allow you to strongly and safely unmarshal WhatsApp business account incoming messages
// and verification events into Go structs. It includes types covering WebhookPayload, VerificationRequests, Messages,
// Callbacks, and Media.
//
// # Examples of Usage
//
// You can use `google.golang.org/protobuf/encoding/protojson` to dynamically map incoming requests.
// Here is how you can use it in a typical HTTP Handler.
//
//	package main
//
//	import (
//		"io"
//		"log"
//		"net/http"
//		"os"
//
//		"google.golang.org/protobuf/encoding/protojson"
//		"github.com/ruanspies/wa-to-go" // Automatically uses package name 'watogo'
//	)
//
//	// handshakeHandler handles the verification handshake from Meta.
//	func handshakeHandler(w http.ResponseWriter, r *http.Request) {
//		// Extract query parameters into our protobuf message structure.
//		payload := &watogo.VerificationRequest{
//			Hub: &watogo.VerificationRequest_Hub{
//				Mode:        r.URL.Query().Get("hub.mode"),
//				VerifyToken: r.URL.Query().Get("hub.verify_token"),
//				Challenge:   r.URL.Query().Get("hub.challenge"),
//			},
//		}
//		log.Printf("%v", payload)
//
//		verifyToken := os.Getenv("WHATSAPP_VERIFY_TOKEN")
//		if payload.GetHub().GetMode() == "subscribe" && payload.GetHub().GetVerifyToken() == verifyToken {
//			w.WriteHeader(http.StatusOK)
//			w.Write([]byte(payload.GetHub().GetChallenge()))
//			return
//		}
//
//		// Respond with 403 Forbidden if the verification fails.
//		w.WriteHeader(http.StatusForbidden)
//	}
//
//	// messageHandler handles incoming POST requests with webhook payloads.
//	func messageHandler(w http.ResponseWriter, r *http.Request) {
//		// Read the request body
//		body, err := io.ReadAll(r.Body)
//		if err != nil {
//			log.Printf("Error reading request body: %v", err)
//			http.Error(w, "Error reading request body", http.StatusInternalServerError)
//			return
//		}
//		defer r.Body.Close()
//
//		// Unmarshal the JSON into the WebhookPayload protobuf message
//		var payload watogo.WebhookPayload
//		if err := protojson.Unmarshal(body, &payload); err != nil {
//			log.Printf("Error parsing request body: %v\n\n%s", err, string(body))
//			w.WriteHeader(http.StatusOK)
//			w.Write([]byte("EVENT_RECEIVED"))
//			return
//		}
//
//		// Check if this is a WhatsApp Business Account event
//		if payload.GetObject() == "whatsapp_business_account" {
//			for _, entry := range payload.GetEntry() {
//				for _, change := range entry.GetChanges() {
//					value := change.GetValue()
//					// Business Logic routing
//					if value.GetMessages() != nil {
//						// Do something with messages...
//					} else if value.GetStatuses() != nil {
//						// Do something with statuses...
//					}
//				}
//			}
//			w.WriteHeader(http.StatusOK)
//			w.Write([]byte("EVENT_RECEIVED"))
//		} else {
//			w.WriteHeader(http.StatusNotFound)
//		}
//	}
package watogo
