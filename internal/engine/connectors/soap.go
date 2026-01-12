package connectors

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SOAPConnector handles SOAP to REST conversion
// Allows modern REST clients to interact with legacy SOAP services
type SOAPConnector struct {
	SOAPEndpoint string
	SOAPAction   string // Optional SOAP action header
}

// SOAPConfig represents SOAP connector configuration
type SOAPConfig struct {
	Endpoint   string                 `json:"endpoint"`    // SOAP endpoint URL
	Action     string                 `json:"action"`      // SOAP action (optional)
	Method     string                 `json:"method"`      // SOAP method name
	Namespace  string                 `json:"namespace"`   // XML namespace
	Parameters map[string]interface{} `json:"parameters"`  // Method parameters
	Headers    map[string]string      `json:"headers"`     // Custom HTTP headers
}

// SOAPEnvelope represents a standard SOAP 1.1/1.2 envelope
type SOAPEnvelope struct {
	XMLName xml.Name `xml:"soap:Envelope"`
	SOAP    string   `xml:"xmlns:soap,attr"`
	Body    SOAPBody
}

// SOAPBody represents the SOAP body
type SOAPBody struct {
	XMLName xml.Name    `xml:"soap:Body"`
	Content interface{} `xml:",innerxml"`
}

// SOAPFault represents a SOAP fault response
type SOAPFault struct {
	XMLName     xml.Name `xml:"Fault"`
	FaultCode   string   `xml:"faultcode"`
	FaultString string   `xml:"faultstring"`
	Detail      string   `xml:"detail"`
}

// ExecuteWithContext converts REST request to SOAP, calls legacy service, converts response back
func (s *SOAPConnector) ExecuteWithContext(ctx context.Context, config SOAPConfig) Result {
	start := time.Now()

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled before SOAP request: " + ctx.Err().Error())
	default:
	}

	// Build SOAP envelope
	soapRequest, err := buildSOAPRequest(config)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to build SOAP request: %v", err), start)
	}

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "POST", config.Endpoint, bytes.NewBuffer(soapRequest))
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create HTTP request: %v", err), start)
	}

	// Set SOAP headers
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	if config.Action != "" {
		req.Header.Set("SOAPAction", config.Action)
	}

	// Add custom headers
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	// Execute request with timeout
	client := &http.Client{
		Timeout: 30 * time.Second, // SOAP services can be slow
	}
	resp, err := client.Do(req)

	// Check if context was cancelled during request
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during SOAP request: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("SOAP request failed: %v", err), start)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to read SOAP response: %v", err), start)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		// Try to parse SOAP fault
		fault := parseSOAPFault(body)
		if fault != nil {
			return NewFailureResult(fmt.Sprintf("SOAP Fault: %s - %s", fault.FaultCode, fault.FaultString), start)
		}
		return NewFailureResult(fmt.Sprintf("SOAP returned HTTP error: %d", resp.StatusCode), start)
	}

	// Parse SOAP response
	parsedResponse, err := parseSOAPResponse(body)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to parse SOAP response: %v", err), start)
	}

	return NewSuccessResult("SOAP request completed successfully", map[string]interface{}{
		"status_code": resp.StatusCode,
		"response":    parsedResponse,
		"raw_xml":     string(body),
	}, start)
}

// buildSOAPRequest creates a SOAP envelope from the config
func buildSOAPRequest(config SOAPConfig) ([]byte, error) {
	// Build the method call XML
	var methodXML string
	if config.Namespace != "" {
		methodXML = fmt.Sprintf(`<%s xmlns="%s">`, config.Method, config.Namespace)
	} else {
		methodXML = fmt.Sprintf(`<%s>`, config.Method)
	}

	// Add parameters
	for key, value := range config.Parameters {
		methodXML += fmt.Sprintf(`<%s>%v</%s>`, key, value, key)
	}

	methodXML += fmt.Sprintf(`</%s>`, config.Method)

	// Build SOAP envelope
	envelope := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    %s
  </soap:Body>
</soap:Envelope>`, methodXML)

	return []byte(envelope), nil
}

// parseSOAPResponse extracts data from SOAP response
func parseSOAPResponse(body []byte) (map[string]interface{}, error) {
	// Parse the SOAP envelope
	var envelope struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			Content string `xml:",innerxml"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}

	// Return the body content as a simple map
	// In production, you'd want more sophisticated XML to JSON conversion
	return map[string]interface{}{
		"body": envelope.Body.Content,
	}, nil
}

// parseSOAPFault tries to parse a SOAP fault from the response
func parseSOAPFault(body []byte) *SOAPFault {
	var envelope struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			Fault SOAPFault `xml:"Fault"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(body, &envelope); err != nil {
		return nil
	}

	if envelope.Body.Fault.FaultCode != "" {
		return &envelope.Body.Fault
	}

	return nil
}

// DryRunSOAP simulates a SOAP call without actually making the request
func (s *SOAPConnector) DryRunSOAP(config SOAPConfig) Result {
	start := time.Now()

	// Build SOAP request to show what would be sent
	soapRequest, err := buildSOAPRequest(config)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to build SOAP request: %v", err), start)
	}

	return NewSuccessResult("SOAP dry run completed", map[string]interface{}{
		"endpoint":     config.Endpoint,
		"soap_action":  config.Action,
		"soap_request": string(soapRequest),
		"note":         "This is a dry run - no actual SOAP call was made",
	}, start)
}

