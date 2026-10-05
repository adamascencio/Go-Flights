package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/adamascencio/Go-Flights/internal/googleflights"
)

func main() {
	origin := flag.String("from", "LAX", "origin airport IATA code")
	destination := flag.String("to", "MEX", "destination airport IATA code")
	date := flag.String("date", "", "departure date (YYYY-MM-DD)")
	nonstop := flag.Bool("nonstop", false, "restrict results to nonstop flights")
	flag.Parse()
	if _, err := time.Parse(time.DateOnly, *date); err != nil {
		fmt.Fprintln(os.Stderr, "provide -date in YYYY-MM-DD format")
		os.Exit(1)
	}

	query := googleflights.NewTFS()
	query.Legs = []googleflights.FlightLeg{{
		DepartureDate: date,
		Origin:        &googleflights.Place{EntityType: googleflights.Ptr(uint32(1)), EntityID: origin},
		Destination:   &googleflights.Place{EntityType: googleflights.Ptr(uint32(1)), EntityID: destination},
	}}
	if *nonstop {
		query.Legs[0].MaxStops = googleflights.Ptr(uint32(0))
	}
	u, err := googleflights.BuildURL(query, googleflights.URLOptions{
		Language: "en", Region: "us", Currency: "USD",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(u)
}
