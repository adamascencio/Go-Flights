package googleflights

import (
	"encoding/base64"

	"google.golang.org/protobuf/encoding/protowire"
)

// EncodeTFS serializes a query as protobuf followed by unpadded URL-safe Base64.
// Nil fields are omitted; non-nil fields are encoded even when their value is zero.
func EncodeTFS(query TFS) string {
	var w wireEncoder
	w.uint32(1, query.QueryMode)
	w.uint32(2, query.QueryContext)
	for _, leg := range query.Legs {
		w.bytes(3, encodeLeg(leg))
	}
	for _, passenger := range query.Passengers {
		w.varint(8, uint64(int64(passenger)))
	}
	if query.Cabin != nil {
		w.varint(9, uint64(int64(*query.Cabin)))
	}
	w.uint32(12, query.MaxPriceAmount)
	if query.Baggage != nil {
		var baggage wireEncoder
		baggage.uint32(2, query.Baggage.Field2)
		baggage.uint32(3, query.Baggage.Field3)
		w.bytes(13, baggage)
	}
	w.uint32(14, query.DisplayFlag)
	if query.AllResultsFlag != nil {
		var flag wireEncoder
		if query.AllResultsFlag.Sentinel != nil {
			flag.varint(1, *query.AllResultsFlag.Sentinel)
		}
		w.bytes(16, flag)
	}
	if query.HideSeparateSelfTransferTickets != nil {
		var value uint64
		if *query.HideSeparateSelfTransferTickets {
			value = 1
		}
		w.varint(17, value)
	}
	if query.TripType != nil {
		w.varint(19, uint64(int64(*query.TripType)))
	}
	return base64.RawURLEncoding.EncodeToString(w)
}

func encodeLeg(leg FlightLeg) []byte {
	var w wireEncoder
	w.string(2, leg.DepartureDate)
	w.uint32(5, leg.MaxStops)
	for _, airline := range leg.IncludeAirlines {
		w.bytes(6, []byte(airline))
	}
	for _, airline := range leg.ExcludeAirlines {
		w.bytes(7, []byte(airline))
	}
	w.uint32(8, leg.DepartureStartHour)
	w.uint32(9, leg.DepartureEndHourInclusive)
	w.uint32(10, leg.ArrivalStartHour)
	w.uint32(11, leg.ArrivalEndHourInclusive)
	w.uint32(12, leg.MaxDurationMinutes)
	if leg.Origin != nil {
		w.bytes(13, encodePlace(*leg.Origin))
	}
	if leg.Destination != nil {
		w.bytes(14, encodePlace(*leg.Destination))
	}
	for _, airport := range leg.IncludeConnectionAirports {
		w.bytes(15, []byte(airport))
	}
	for _, airport := range leg.ExcludeConnectionAirports {
		w.bytes(16, []byte(airport))
	}
	w.uint32(18, leg.MaxLayoverMinutes)
	if len(leg.EmissionsFlags) > 0 {
		var packed []byte
		for _, flag := range leg.EmissionsFlags {
			packed = protowire.AppendVarint(packed, uint64(flag))
		}
		w.bytes(19, packed)
	}
	return w
}

func encodePlace(place Place) []byte {
	var w wireEncoder
	w.uint32(1, place.EntityType)
	w.string(2, place.EntityID)
	return w
}

type wireEncoder []byte

func (w *wireEncoder) varint(field protowire.Number, value uint64) {
	*w = protowire.AppendTag(*w, field, protowire.VarintType)
	*w = protowire.AppendVarint(*w, value)
}

func (w *wireEncoder) bytes(field protowire.Number, value []byte) {
	*w = protowire.AppendTag(*w, field, protowire.BytesType)
	*w = protowire.AppendBytes(*w, value)
}

func (w *wireEncoder) uint32(field protowire.Number, value *uint32) {
	if value != nil {
		w.varint(field, uint64(*value))
	}
}

func (w *wireEncoder) string(field protowire.Number, value *string) {
	if value != nil {
		w.bytes(field, []byte(*value))
	}
}
