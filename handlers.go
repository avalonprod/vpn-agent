package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

type RegisterRequest struct {
	ClientPublicKey string `json:"client_public_key"`
	ClientName      string `json:"client_name,omitempty"`
}

type RegisterResponse struct {
	IP              string `json:"ip"`
	ServerPublicKey string `json:"server_public_key"`
	Endpoint        string `json:"endpoint"`
	AllowedIPs      string `json:"allowed_ips"`
	DNS             string `json:"dns"`
}

func DeletePeerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Auth-Token") != ServerConfig.AuthToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	key := r.URL.Query().Get("key")

	if key == "" {
		http.Error(w, "missing publicKey", http.StatusBadRequest)
		return
	}

	cmd := exec.Command("wg", "set", "wg0", "peer", key, "remove")
	if err := cmd.Run(); err != nil {
		http.Error(w, fmt.Sprintf("failed to remove peer: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func RegisterNewPeerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Auth-Token") != ServerConfig.AuthToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if isPeerExists(req.ClientPublicKey) {
		http.Error(w, "Peer already exists", http.StatusConflict)
		return
	}

	ip := getNextIP()
	if ip == "" {
		http.Error(w, "No available IPs", http.StatusInternalServerError)
		return
	}

	cmd := exec.Command("wg", "set", ServerConfig.WGInterface, "peer", req.ClientPublicKey, "allowed-ips", ip+"/32")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("wg set failed: %v\nOutput: %s\n", err, output)
		http.Error(w, "wg failed", http.StatusInternalServerError)
		return
	}

	usedIPs[ip] = true
	log.Printf("Peer added: %s (%s)", ip, req.ClientName)

	resp := RegisterResponse{
		IP:              ip,
		ServerPublicKey: ServerConfig.ServerPublicKey,
		Endpoint:        ServerConfig.Endpoint,
		AllowedIPs:      defaultAllowedIPs,
		DNS:             ServerConfig.DNSServer,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
