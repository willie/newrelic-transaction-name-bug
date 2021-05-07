# newrelic-transaction-name-bug
sample code to demonstrate https://github.com/newrelic/go-agent/issues/310

To run alice and bob, shell into their respective directories and run: `NEWRELIC_LICENSE="YOUR LICENSE KEY HERE" go run .`

Call bob directly: `curl "http://localhost:8082/header/L-999-70000.vbk"`

To have alice call bob: `curl "http://localhost:8081/header/L-999-70000.vbk"`
