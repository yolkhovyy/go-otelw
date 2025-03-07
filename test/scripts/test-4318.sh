


curl -X POST "http://localhost:4318/v1/logs" -H "Content-Type: application/json" -d '{"resourceLogs":[{"resource":{"attributes":[{"key":"service.name","value":"example-service"}]},"instrumentationLibraryLogs":[{"logs":[{"timeUnixNano":"1677613200000000000","severityText":"INFO","body":{"stringValue":"This is an example log entry"},"attributes":[{"key":"example-attribute","value":"example-value"}]}]}]}]}'
