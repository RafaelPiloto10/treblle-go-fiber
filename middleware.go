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

func Middleware() func(*fiber.Ctx) {
	return func(r *fiber.Ctx) {
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
				r.Context().Response.Header.Add(k, v)
			}
		}

		if !errors.Is(errReqInfo, ErrNotJson) {
			ti := MetaData{
				ApiKey:    Config.APIKey,
				ProjectID: Config.ProjectID,
				Version:   treblleVersion,
				Sdk:       sdkName,
				Data: DataInfo{
					Server:   Config.serverInfo,
					Language: Config.languageInfo,
					Request:  requestInfo,
					Response: getResponseInfo(rec, startTime),
				},
			}
			// don't block execution while sending data to Treblle
			go sendToTreblle(ti)
		}
	}
}

// If anything happens to go wrong inside one of treblle-go internals, recover from panic and continue
func dontPanic() {
	if err := recover(); err != nil {
		log.Printf("treblle-go panic: %s", err)
	}
}
