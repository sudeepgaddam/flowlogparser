README
Project Overview
This project is a Go application that parses flow logs and maps them to predefined tags using a lookup table. The lookup table is defined in a CSV format with columns: dstport, protocol, and tag. The application processes each flow log entry by matching its dstport and protocol fields against the lookup table and assigns a tag accordingly. Flow logs that don't match any entry in the lookup table are tagged as "Untagged".

The project includes a main.go file that implements the flow log parsing logic, and a test suite (flowlog_test.go) that covers various scenarios to ensure the correctness of the logic.

Prerequisites
Go 1.18+ installed on your machine
Project Structure
main.go: The main application file that handles the parsing of flow logs and the lookup table.
flowlog_test.go: The test suite that contains various unit tests to validate the functionality.
lookup_table.csv: A CSV file that contains the destination port, protocol, and tag mapping.
Running the Application
To run the main.go file with a custom flow log file and lookup table file, follow the steps below:

Assumptions
For protocols, only tcp, udp and icmp protocol numbers are supported as part of this application
If invalid protocol number is present in flow logs, that particular log line is ignored and remaining log lines are processed
Clone the repository:

bash
git clone <repository-url>
cd <repository-directory>
Input files are already present in this repo:

File for the flow logs, e.g., input_flow_logs.txt.
File for the lookup table, e.g., input_lookup_table.csv.
Run the application:

bash
go run main.go <path-to-lookup_table.csv> <path-to-flowlogs.txt>
Replace <path-to-flowlogs.txt> and <path-to-lookup_table.csv> with the paths to your respective files.

The output will include:

Tag counts: Counts of how many times each tag appeared.
Port/Protocol combination counts: Counts of how many times each port/protocol combination appeared.
Running the Tests
To ensure that the functionality is working as expected, there is a test suite implemented in flowlog_test.go.

Run all tests:
go test
This will automatically detect all the test functions and run them, outputting the result of each test case.

1. Flow Log Parsing and Tagging (TestParseFlowLog)
   This test ensures that the flow logs are correctly tagged based on the dstport and protocol fields from the flow log and the lookup table. The following scenarios are covered:

Exact match: Flow logs with dstport and protocol combinations that match the lookup table.
Untagged: Flow logs that don't have a matching dstport and protocol combination.
Invalid Protocol: Handling flow logs with invalid or unknown protocol numbers.
UDP and TCP Differentiation: Ensures the correct protocol (UDP or TCP) is applied when matching against the lookup table. 2. Tag Counting (TestCountTags)
This test verifies that the tags are correctly counted after parsing the flow logs. Each tag that appears in the flow logs should be counted and compared to the expected values.

3. Port/Protocol Combination Counting (TestCountPortProtocol)
   This test ensures that the counting of port/protocol combinations works as expected. It checks that the number of occurrences of each port/protocol combination is tracked correctly.

Example
Sample Flow Logs (flowlogs.txt):

2 123456789012 eni-0a1b2c3d 10.0.1.201 198.51.100.2 443 49153 6 25 20000 1620140761 1620140821 ACCEPT OK
2 123456789012 eni-1a2b3c4d 192.168.1.100 203.0.113.101 23 49154 6 15 12000 1620140761 1620140821 REJECT OK
Sample Lookup Table (lookup_table.csv):

25,tcp,sv_P1
23,tcp,sv_P1
443,tcp,sv_P2
Running the application:

bash
go run main.go input_lookup_table.csv input_flow_logs.txt output.csv
Expected Output:

output.csv will contain expected output
