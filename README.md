# README

## Project Overview

This project is a Go application that parses flow logs and maps them to predefined tags using a lookup table. The lookup table is defined in a CSV format with columns: `dstport`, `protocol`, and `tag`. The application processes each flow log entry by matching its `dstport` and `protocol` fields against the lookup table and assigns a tag accordingly. Flow logs that don't match any entry in the lookup table are tagged as "Untagged".

The project includes a `main.go` file that implements the flow log parsing logic, and a test suite (`flowlog_test.go`) that covers various scenarios to ensure the correctness of the logic.

## Prerequisites

- Go 1.18+ installed on your machine

## Input Files

The input files are already present in this repository:

- A file for the flow logs, e.g., `input_flow_logs.txt`.
- A CSV file for the lookup table, e.g., `input_lookup_table.csv`.

## Run the Application

To run the `main.go` file with a custom flow log file and lookup table file, follow the steps below:

```bash
go run main.go <path-to-lookup_table.csv> <path-to-flowlogs.txt>
```

Replace `<path-to-flowlogs.txt>` and `<path-to-lookup_table.csv>` with the paths to your respective files.

The output will include:

- **Tag counts**: Counts of how many times each tag appeared.
- **Port/Protocol combination counts**: Counts of how many times each port/protocol combination appeared.

## Running the Tests

To ensure that the functionality is working as expected, there is a test suite implemented in `flowlog_test.go`.

### Run all tests

```bash
go test
```

This will automatically detect all the test functions and run them, outputting the result of each test case.

### Flow Log Parsing and Tagging (`TestParseFlowLog`)

This test ensures that the flow logs are correctly tagged based on the `dstport` and `protocol` fields from the flow log and the lookup table. The following scenarios are covered:

1. **Exact match**: Flow logs with `dstport` and `protocol` combinations that match the lookup table.
2. **Untagged**: Flow logs that don't have a matching `dstport` and `protocol` combination.
3. **Invalid Protocol**: Handling flow logs with invalid or unknown protocol numbers.
4. **UDP and TCP Differentiation**: Ensures the correct protocol (UDP or TCP) is applied when matching against the lookup table.

### Tag Counting (`TestCountTags`)

This test verifies that the tags are correctly counted after parsing the flow logs. Each tag that appears in the flow logs should be counted and compared to the expected values.

### Port/Protocol Combination Counting (`TestCountPortProtocol`)

This test ensures that the counting of port/protocol combinations works as expected. It checks that the number of occurrences of each port/protocol combination is tracked correctly.

## Example

### Sample Flow Logs (`flowlogs.txt`)

```text
2 123456789012 eni-0a1b2c3d 10.0.1.201 198.51.100.2 443 49153 6 25 20000 1620140761 1620140821 ACCEPT OK
2 123456789012 eni-1a2b3c4d 192.168.1.100 203.0.113.101 23 49154 6 15 12000 1620140761 1620140821 REJECT OK
```

### Sample Lookup Table (`lookup_table.csv`)

```csv
dstport,protocol,tag
25,tcp,sv_P1
23,tcp,sv_P1
443,tcp,sv_P2
```

### Running the Application

```bash
go run main.go input_lookup_table.csv input_flow_logs.txt output.csv
```

### Expected Output

The output will be written to `output.csv` and contain the expected results.

```
Tag Counts:
sv_P2,1
sv_P1,1

Port/Protocol Combination Counts:
443,tcp,1
23,tcp,1
```
