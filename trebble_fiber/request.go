package treblle_fiber

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RequestInfo struct {
	Timestamp string          `json:"timestamp"`
	Ip        string          `json:"ip"`
	Url       string          `json:"url"`
	UserAgent string          `json:"user_agent"`
	Method    string          `json:"method"`
	Headers   json.RawMessage `json:"headers"`
	Body      json.RawMessage `json:"body"`
}

var ErrNotJson = errors.New("request body is not JSON")

// Get details about the request
func getRequestInfo(r *fiber.Ctx, startTime time.Time) (RequestInfo, error) {
	defer dontPanic()

	headers := make(map[string]string)
	for k := range r.GetReqHeaders() {
		headers[k] = r.GetReqHeaders()[k]
	}

	ip := extractIP(r.Context().RemoteAddr().String())

	body, err := json.Marshal(headers)
	if err != nil {
		return RequestInfo{}, err
	}

	ri := RequestInfo{
		Timestamp: startTime.Format("2006-01-02 15:04:05"),
		Ip:        ip,
		Url:       r.Context().URI().String(),
		UserAgent: string(r.Request().Header.UserAgent()),
		Method:    string(r.Request().Header.Method()),
		Headers:   body,
	}

	requestBody := make([]byte, len(r.Context().Request.Body()))
	copy(requestBody, r.Context().Request.Body())

	if requestBody != nil && len(requestBody) > 0 {
		buf := new(bytes.Buffer)
		buf.Write(requestBody)
		buf_bytes := buf.Bytes()

		// open 2 NopClosers over the buffer to allow buffer to be read and still passed on
		bodyReaderOriginal := io.NopCloser(bytes.NewBuffer(buf_bytes))

		body, err := io.ReadAll(bodyReaderOriginal)
		if err != nil {
			return ri, err
		}

		// mask all the JSON fields listed in Config.FieldsToMask
		sanitizedBody, err := getMaskedJSON(body)
		if err != nil {
			return ri, err
		}

		ri.Body = sanitizedBody

	}

	headersJson, err := json.Marshal(headers)
	if err != nil {
		return ri, err
	}

	sanitizedHeaders, err := getMaskedJSON(headersJson)

	if err != nil {
		return ri, err
	}
	ri.Headers = sanitizedHeaders
	return ri, nil
}

func recoverBody(r *fiber.Ctx, bodyReaderCopy io.ReadCloser) {
	buf := []byte{}
	bodyReaderCopy.Read(buf)
	r.Context().Request.SetBody(buf)
}

func getMaskedJSON(payloadToMask []byte) (json.RawMessage, error) {
	jsonMap := make(map[string]interface{})
	if err := json.Unmarshal(payloadToMask, &jsonMap); err != nil {
		// probably a JSON array so let's return it.
		return payloadToMask, nil
	}

	sanitizedJson := make(map[string]interface{})
	copyAndMaskJson(jsonMap, sanitizedJson)
	jsonData, err := json.Marshal(sanitizedJson)
	if err != nil {
		return nil, err
	}

	rawMessage := json.RawMessage(jsonData)

	return rawMessage, nil
}

func copyAndMaskJson(src map[string]interface{}, dest map[string]interface{}) {
	for key, value := range src {
		switch src[key].(type) {
		case map[string]interface{}:
			dest[key] = map[string]interface{}{}
			copyAndMaskJson(src[key].(map[string]interface{}), dest[key].(map[string]interface{}))
		default:
			// if JSON key is in the list of keys to mask, replace it with a * string of the same length
			_, exists := Config.FieldsMap[key]
			if exists {
				maskedValue := maskValue(value.(string), key)
				dest[key] = maskedValue
			} else {
				dest[key] = value
			}
		}
	}
}
func maskValue(valueToMask string, key string) string {
	keyLower := strings.ToLower(key)

	if keyLower == "authorization" && regexp.MustCompile(`(?i)^(bearer|basic)\s+`).MatchString(valueToMask) {
		authParts := strings.SplitN(valueToMask, " ", 2)
		authPrefix := authParts[0]
		authToken := authParts[1]
		maskedAuthToken := strings.Repeat("*", len(authToken))
		maskedValue := authPrefix + " " + maskedAuthToken
		return maskedValue
	}

	return strings.Repeat("*", len(valueToMask))
}

func extractIP(remoteAddr string) string {
	// If RemoteAddr contains both IP and port, split and return the IP
	if strings.Contains(remoteAddr, ":") {
		ip, _, err := net.SplitHostPort(remoteAddr)
		if err == nil {
			return ip
		}
	}

	return remoteAddr
}
