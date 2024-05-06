package main

import (
    "encoding/json"
    "net/http"
    "os/exec"
    "strings"
    "fmt"
)

// DNSRecords specifies the record sets
type DNSRecords struct {
    A     []string `json:"a,omitempty"`
    AAAA  []string `json:"aaaa,omitempty"`
    CNAME []string `json:"cname,omitempty"`
    MX    []string `json:"mx,omitempty"`
    NS    []string `json:"ns,omitempty"`
    TXT   []string `json:"txt,omitempty"`
}

// set default recursive
const defaultResolver = "8.8.8.8" // Google is the default

// handleDNSQuery 
func handleDNSQuery(w http.ResponseWriter, r *http.Request) {
    domain := r.URL.Query().Get("domain")
    nameserver := r.URL.Query().Get("nameserver")

    // Determine nameserver 
    if nameserver == "" {
        nameserver = defaultResolver // hardcoded default if no resolver specified
    }

    records, err := queryAllRecordTypes(domain, nameserver)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(records)
}

// queryAllRecordTypes digs for all DNS record types & outputs
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

// parseDigOutput  parses by record type 
func parseDigOutput(recordType, output string, records *DNSRecords) {
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
    http.HandleFunc("/dns-query", handleDNSQuery)
    http.ListenAndServe(":8080", nil)
}
