package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	dsArg := flag.String("downstream", "", "services to contact downstream")
	port := flag.Int("port", 8888, "port for the server")
	flag.Parse()

	dsSvcs := strings.Split(*dsArg, ",")
	log.Printf("Downstreams: %v", dsSvcs)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request from %s", r.RemoteAddr)
		requestID := time.Now().UnixNano()
		finalResp := make(map[string]interface{})
		for _, s := range dsSvcs {
			if len(s) == 0 {
				continue
			}
			log.Printf("[%d] Sending request to %s", requestID, s)
			resp, err := http.Get(fmt.Sprintf("http://%s", s))
			if err != nil {
				finalResp[s] = err.Error()
				log.Printf("[%d] Error while getting response. Reason=%v", requestID, err)
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				finalResp[s] = err.Error()
				log.Printf("[%d] Error while reading body. Reason=%v", requestID, err)
				continue
			}
			var dsData map[string]interface{}
			err = json.Unmarshal(body, &dsData)
			if err != nil {
				finalResp[s] = err.Error()
				log.Printf("[%d] Error while unmarshaling body. Reason=%v", requestID, err)
				continue
			}
			finalResp[s] = dsData
		}
		hostname, _ := os.Hostname()
		finalResp["current"] = hostname

		respBytes, err := json.Marshal(finalResp)
		if err != nil {
			finalResp["current"] = err.Error()
			log.Printf("[%d] Error while marshaling. Reason=%v", requestID, err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	})
	log.Printf("starting server at :%d", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), mux); err != nil {
		log.Printf("error while starting HTTP server. Reason = %v", err)
		os.Exit(-1)
	}
}
