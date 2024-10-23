package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
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

// IANAProtocolMap maps IANA protocol numbers to their corresponding names
var IANAProtocolMap = map[int]string{
	6:  "tcp",  // TCP
	17: "udp",  // UDP
	1:  "icmp", // ICMP
	// Add more protocols if needed
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
		port := strings.TrimSpace(record[0])
		protocol := strings.TrimSpace(strings.ToLower(record[1]))
		tag := strings.TrimSpace(record[2])

		key := fmt.Sprintf("%s_%s", port, protocol)
		lookupTable[key] = tag
	}

	return lookupTable, nil
}

// ParseFlowLog parses a flow log line into a FlowLog struct
func ParseFlowLog(line string, lookupTable LookupTable) (FlowLog, error) {
	fields := strings.Fields(line)
	if len(fields) < 14 {
		return FlowLog{}, nil
	}

	dstPort := fields[5]
	protocolNum, err := strconv.Atoi(fields[7])
	if err != nil {
		return FlowLog{}, err
	}
	// Map protocol number to protocol name
	protocol := "unknown"
	if protocolStr, exists := IANAProtocolMap[protocolNum]; exists {
		protocol = protocolStr
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
	}, nil
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
func WriteOutput(filePath string,
	tagCount map[string]int,
	portProtocolCount map[string]int) error {
	// Create or open the file
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer file.Close()

	// Create a new CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := WriteTagCounts(tagCount, writer); err != nil {
		return err
	}
	if err := WritePortProtocolCounts(portProtocolCount, writer); err != nil {
		return err
	}
	return nil
}

// WriteTagCounts writes the tag count results
func WriteTagCounts(tagCount map[string]int, writer *csv.Writer) error {
	title := []string{"Tag Counts:"}
	err := writer.Write(title)
	if err != nil {
		fmt.Println("Error writing title:", err)
		return err
	}

	header := []string{"Tag", "Count"}
	err = writer.Write(header)
	if err != nil {
		fmt.Println("Error writing header:", err)
		return err
	}

	for tag, count := range tagCount {
		row := []string{tag, strconv.Itoa(count)}
		err = writer.Write(row)
		if err != nil {
			fmt.Println("Error writing row:", err)
			return err
		}
	}
	return nil
}

// WritePortProtocolCounts prints the port/protocol combination counts
func WritePortProtocolCounts(portProtocolCount map[string]int,
	writer *csv.Writer) error {
	title := []string{"Port/Protocol Combination Counts:"}
	err := writer.Write(title)
	if err != nil {
		fmt.Println("Error writing title:", err)
		return err
	}
	header := []string{"Port", "Protocol", "Count"}
	err = writer.Write(header)
	if err != nil {
		fmt.Println("Error writing header:", err)
		return err
	}

	for key, count := range portProtocolCount {
		parts := strings.Split(key, "_")
		row := []string{parts[0], parts[1], strconv.Itoa(count)}
		err = writer.Write(row)
		if err != nil {
			fmt.Println("Error writing row:", err)
			return err
		}
	}
	return nil
}

func main() {
	// Ensure that exactly 3 arguments are passed (two inputs and one output file)
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run main.go <path_to_lookup_table> <path_to_flow_log_file> <output file>")
		return
	}

	// Get command-line arguments
	path_to_lookup_table := os.Args[1]
	path_to_flow_log_file := os.Args[2]
	outputFile := os.Args[3]

	// Load lookup table
	lookupTable, err := ParseLookupTable(path_to_lookup_table)
	if err != nil {
		fmt.Println("Error reading lookup table:", err)
		return
	}

	// Open the flow log file
	file, err := os.Open(path_to_flow_log_file)
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
		flowLog, err := ParseFlowLog(line, lookupTable)
		if err != nil {
			fmt.Println("Error parsing flow log:", err.Error())
			continue
		}
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
	// Write output to output file
	WriteOutput(outputFile, tagCount, portProtocolCount)
}
