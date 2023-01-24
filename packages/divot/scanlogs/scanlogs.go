package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/farwydi/divot/pkg/database/minio"
	"github.com/farwydi/divot/pkg/enums"
	"github.com/farwydi/divot/pkg/warcraftlogs"
)

var (
	dbEndpoint, dbKey, dbSecret, dbBucket, dbRegion string

	warcraftLogsClientID, warcraftLogsClientSecret string
)

func init() {
	dbEndpoint = os.Getenv("DB_ENDPOINT")
	if dbEndpoint == "" {
		panic("no db endpoint provided")
	}
	dbKey = os.Getenv("DB_KEY")
	if dbKey == "" {
		panic("no db key provided")
	}
	dbSecret = os.Getenv("DB_SECRET")
	if dbSecret == "" {
		panic("no db secret provided")
	}
	dbBucket = os.Getenv("DB_BUCKET")
	if dbBucket == "" {
		panic("no db bucket provided")
	}
	dbRegion = os.Getenv("DB_REGION")
	if dbRegion == "" {
		panic("no db region provided")
	}

	warcraftLogsClientID = os.Getenv("WOWLOGS_CLIENT_ID")
	if warcraftLogsClientID == "" {
		panic("no wow logs client id provided")
	}
	warcraftLogsClientSecret = os.Getenv("WOWLOGS_CLIENT_SECRET")
	if warcraftLogsClientSecret == "" {
		panic("no wow logs client secret provided")
	}
}

type Request struct {
	ClassesToScan string `json:"classes_to_scan"`
}

type Response struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}

func Main(in Request) (*Response, error) {
	db, err := minio.New(dbEndpoint, dbKey, dbSecret, dbBucket, dbRegion)
	if err != nil {
		return nil, fmt.Errorf("fail init database %v in %s (%s)", err, dbEndpoint, dbRegion)
	}

	wowLogs := warcraftlogs.NewWarcraftLogs(
		warcraftLogsClientID,
		warcraftLogsClientSecret,
		db,
	)

	classesToScan := strings.Split(in.ClassesToScan, ",")

	for _, className := range classesToScan {
		for _, specName := range enums.SpecsMap[strings.ToUpper(className)] {
			for name, id := range enums.DungeonEncounters {
				fmt.Printf("Load %s, for %s-%s\n", name, className, specName)
				err = wowLogs.ScanWorldDataEncounterCharacterRankings(context.Background(), id, className, specName)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return &Response{
		StatusCode: 200,
		Headers:    nil,
		Body:       "",
	}, nil
}
