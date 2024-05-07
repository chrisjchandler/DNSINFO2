package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os/exec"
    "strings"
)

// config holds server configuration
type Config struct {
    UseTLS     bool   // Enable TLS
    CertFile   string // Path to certificate file
    KeyFile    string // Path to key file
    ServerPort string // Server port
}

// DNSRecords structure for DNS reply
type DNSRecords struct {
    A     []string `json:"a,omitempty"`
    AAAA  []string `json:"aaaa,omitempty"`
    CNAME []string `json:"cname,omitempty"`
    MX    []string `json:"mx,omitempty"`
    NS    []string `json:"ns,omitempty"`
    TXT   []string `json:"txt,omitempty"`
}

// handleDNSQuery query handler
func handleDNSQuery(w http.ResponseWriter, r *http.Request) {
    domain := r.URL.Query().Get("domain")
    nameserver := r.URL.Query().Get("nameserver")

    // Use a default recursive if none is provided
    if nameserver == "" {
        nameserver = "8.8.8.8" // Example default resolver
    }

    records, err := queryAllRecordTypes(domain, nameserver)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(records)
}

// queryAllRecordTypes performs dig for all DNS types
func queryAllRecordTypes(domain, nameserver string) (DNSRecords, error) {
    records := DNSRecords{}
    recordTypes := map[string]string{
        "A":     "+short",
        "AAAA":  "+short",
        "CNAME": "+short",
        "MX":    "+short",
        "NS":    "+short",
        "TXT":   "+short",
    }

    for recordType, option := range recordTypes {
        cmd := fmt.Sprintf("dig @%s %s %s %s", nameserver, domain, recordType, option)
        output, err := exec.Command("bash", "-c", cmd).CombinedOutput()
        if err != nil {
            return records, err
        }
        parseDigOutput(recordType, strings.TrimSpace(string(output)), &records)
    }
    return records, nil
}

// parseDigOutput parses output from dig command by record type
func parseDigOutput(recordType string, output string, records *DNSRecords) {
    results := strings.Split(output, "\n")
    switch recordType {
    case "A":
        records.A = append(records.A, results...)
    case "AAAA":
        records.AAAA = append(records.AAAA, results...)
    case "CNAME":
        records.CNAME = append(records.CNAME, results...)
    case "MX":
        records.MX = append(records.MX, results...)
    case "NS":
        records.NS = append(records.NS, results...)
    case "TXT":
        for _, txt := range results {
            records.TXT = append(records.TXT, strings.Trim(txt, "\""))
        }
    }
}

func main() {
    config := Config{
        UseTLS:     false,             // enable or disable TLS
        CertFile:   "server.crt",     // Cert 4 TLS
        KeyFile:    "server.key",     // Key for TLS
        ServerPort: "8080",           // listening port
    }

    http.HandleFunc("/dns-query", handleDNSQuery)

    if config.UseTLS {
        log.Printf("Starting HTTPS server on port %s", config.ServerPort)
        log.Fatal(http.ListenAndServeTLS(":"+config.ServerPort, config.CertFile, config.KeyFile, nil))
    } else {
        log.Printf("Starting HTTP server on port %s", config.ServerPort)
        log.Fatal(http.ListenAndServe(":"+config.ServerPort, nil))
    }
}
