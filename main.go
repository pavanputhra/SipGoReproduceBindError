package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/emiago/sipgo"
	"github.com/emiago/sipgo/sip"
)

const (
	DESTINATION_IP   = "2.2.2.2"
	DESTINATION_PORT = "5061"
	PUBLIC_IP        = "212.128.253.126"
)

func main() {
	ua, _ := sipgo.NewUA()           // Build user agent
	srv, _ := sipgo.NewServer(ua)    // Creating server handle for ua
	client, _ := sipgo.NewClient(ua) // Creating client handle for ua
	srv.OnInvite(func(req *sip.Request, tx sip.ServerTransaction) {

		fmt.Println("Invite received")
		resp := sip.NewResponseFromRequest(req, 100, "Trying", nil)
		tx.Respond(resp)

		viaHeader := &sip.ViaHeader{
			ProtocolName:    "SIP",
			ProtocolVersion: "2.0",
			Transport:       "UDP",
			Host:            PUBLIC_IP,
			Port:            sip.DefaultPort("udp"),
			Params:          sip.NewParams(),
		}
		viaHeader.Params.Add("branch", sip.GenerateBranchN(16))

		clonedReq := req.Clone()
		clonedReq.SetBody(req.Body())
		clonedReq.PrependHeader(viaHeader)

		clonedReq.ViaNAT = true
		clonedReq.SetDestination(fmt.Sprintf("%s:%s", DESTINATION_IP, DESTINATION_PORT))

		fmt.Println("Forwarding request to via public IP")
		ctx := context.Background()
		clTx, err := client.TransactionRequest(ctx, clonedReq)
		if err != nil {
			log.Printf("Failed to create client transaction: %v", err)
			tx.Respond(sip.NewResponseFromRequest(req, sip.StatusInternalServerError, "Server Error", nil))
			return
		}

		<-clTx.Done()
	})

	// For registrars
	// srv.OnRegister(registerHandler)
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	fmt.Println("Starting SIP proxy server on port 5060")
	go func() {
		err := srv.ListenAndServe(ctx, "udp", "0.0.0.0:5060")
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	<-ctx.Done()
}
