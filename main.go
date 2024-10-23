package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// LookupTable represents the port and protocol to tag mapping
type LookupTable map[string]string

// FlowLog represents a single flow log entry
type FlowLog struct {
	DstPort  string
	Protocol string
	Tag      string
}

// ParseLookupTable reads the lookup CSV file and returns a LookupTable
func ParseLookupTable(filePath string) (LookupTable, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lookupTable := make(LookupTable)
	reader := csv.NewReader(file)
	_, err = reader.Read() // Skip header row
	if err != nil {
		return nil, err
	}

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		key := fmt.Sprintf("%s_%s", record[0], strings.ToLower(record[1]))
		lookupTable[key] = record[2]
	}

	return lookupTable, nil
}

// ParseFlowLog parses a flow log line into a FlowLog struct
func ParseFlowLog(line string, lookupTable LookupTable) FlowLog {
	fields := strings.Fields(line)
	if len(fields) < 14 {
		return FlowLog{}
	}

	dstPort := fields[5]
	protocolNum := fields[7]

	// Map protocol number to protocol name
	var protocol string
	switch protocolNum {
	case "6":
		protocol = "tcp"
	case "17":
		protocol = "udp"
	case "1":
		protocol = "icmp"
	default:
		protocol = "unknown"
	}

	// Find tag in lookup table
	key := fmt.Sprintf("%s_%s", dstPort, protocol)
	tag, exists := lookupTable[key]
	if !exists {
		tag = "Untagged"
	}

	return FlowLog{
		DstPort:  dstPort,
		Protocol: protocol,
		Tag:      tag,
	}
}

// CountTags counts the occurrences of each tag
func CountTags(flowLogs []FlowLog) map[string]int {
	tagCount := make(map[string]int)
	for _, log := range flowLogs {
		tagCount[log.Tag]++
	}
	return tagCount
}

// CountPortProtocol counts occurrences of each port/protocol combination
func CountPortProtocol(flowLogs []FlowLog) map[string]int {
	portProtocolCount := make(map[string]int)
	for _, log := range flowLogs {
		key := fmt.Sprintf("%s_%s", log.DstPort, log.Protocol)
		portProtocolCount[key]++
	}
	return portProtocolCount
}

// PrintTagCounts prints the tag count results
func PrintTagCounts(tagCount map[string]int) {
	fmt.Println("Tag Counts:")
	fmt.Println("Tag,Count")
	for tag, count := range tagCount {
		fmt.Printf("%s,%d\n", tag, count)
	}
}

// PrintPortProtocolCounts prints the port/protocol combination counts
func PrintPortProtocolCounts(portProtocolCount map[string]int) {
	fmt.Println("\nPort/Protocol Combination Counts:")
	fmt.Println("Port,Protocol,Count")
	for key, count := range portProtocolCount {
		parts := strings.Split(key, "_")
		fmt.Printf("%s,%s,%d\n", parts[0], parts[1], count)
	}
}

func main() {
	// Load lookup table
	lookupTable, err := ParseLookupTable("input_lookup_table.csv")
	if err != nil {
		fmt.Println("Error reading lookup table:", err)
		return
	}

	// Open the flow log file
	file, err := os.Open("input_flow_logs.txt")
	if err != nil {
		fmt.Println("Error reading flow log:", err)
		return
	}
	defer file.Close()

	// Parse flow logs
	var flowLogs []FlowLog
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		flowLog := ParseFlowLog(line, lookupTable)
		if flowLog.DstPort != "" {
			flowLogs = append(flowLogs, flowLog)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning flow log:", err)
		return
	}

	// Count tags and port/protocol combinations
	tagCount := CountTags(flowLogs)
	portProtocolCount := CountPortProtocol(flowLogs)

	// Print results
	PrintTagCounts(tagCount)
	PrintPortProtocolCounts(portProtocolCount)
}
