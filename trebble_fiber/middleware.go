package treblle_fiber

import (
	"errors"
	"log"
	"net/http/httptest"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	treblleVersion = 0.6
	sdkName        = "go"
)

func Middleware() func(*fiber.Ctx) error {
	return func(r *fiber.Ctx) error {
		startTime := time.Now()

		requestInfo, errReqInfo := getRequestInfo(r, startTime)

		// intercept the response so it can be copied
		rec := httptest.NewRecorder()

		// do the actual request as intended
		r.Next()
		// after this finishes, we have the response recorded

		// copy the original headers
		for k, vs := range rec.Header() {
			for _, v := range vs {
				rec.Header().Add(k, v)
			}
		}

		rec.Body.Write(r.Context().Response.Body())

		if !errors.Is(errReqInfo, ErrNotJson) {
			responseInfo := getFiberResponseInfo(rec, startTime)

			ti := MetaData{
				ApiKey:    Config.APIKey,
				ProjectID: Config.ProjectID,
				Version:   treblleVersion,
				Sdk:       sdkName,
				Data: DataInfo{
					Server:   Config.serverInfo,
					Language: Config.languageInfo,
					Request:  requestInfo,
					Response: responseInfo,
				},
			}
			// don't block execution while sending data to Treblle
			go sendToTreblle(ti)
		}

		return nil
	}
}

// If anything happens to go wrong inside one of treblle-go internals, recover from panic and continue
func dontPanic() {
	if err := recover(); err != nil {
		log.Printf("treblle-go panic: %s", err)
	}
}
