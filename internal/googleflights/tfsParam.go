package googleflights

// TFS models the payload of the tfs URL parameter. Field-number comments refer
// to the reverse-engineered spec:
// https://gist.github.com/MomoDeve/a18053dea84dd28e320b8b2c489540eb
//
// Optional scalar fields use pointers: nil means absent, while a pointer to zero
// means explicitly present. EncodeTFS serializes these types as protobuf.
// The spec is not an official Google API contract.
type TFS struct {
	QueryMode                       *uint32         // 1; observed: 28, meaning unknown.
	QueryContext                    *uint32         // 2; observed: 2, meaning unknown.
	Legs                            []FlightLeg     // 3
	Passengers                      []PassengerType // 8; one entry per traveler, unpacked.
	Cabin                           *Cabin          // 9
	MaxPriceAmount                  *uint32         // 12; currency is outside tfs.
	Baggage                         *BaggageFilter  // 13
	DisplayFlag                     *uint32         // 14; observed: 1, meaning unknown.
	AllResultsFlag                  *AllResultsFlag // 16
	HideSeparateSelfTransferTickets *bool           // 17
	TripType                        *TripType       // 19
}

// Ptr makes explicitly present optional values convenient to construct.
func Ptr[T any](value T) *T { return new(value) }

// NewTFS returns the observed defaults for a one-adult, economy, one-way search.
// Set Legs before building a URL, and change TripType for other journey types.
func NewTFS() TFS {
	return TFS{
		QueryMode:      Ptr(uint32(28)),
		QueryContext:   Ptr(uint32(2)),
		Passengers:     []PassengerType{PassengerAdult},
		Cabin:          Ptr(CabinEconomy),
		DisplayFlag:    Ptr(uint32(1)),
		AllResultsFlag: &AllResultsFlag{Sentinel: Ptr(^uint64(0))},
		TripType:       Ptr(TripOneWay),
	}
}

// FlightLeg describes a dated route and its filters.
type FlightLeg struct {
	DepartureDate             *string  // 2; YYYY-MM-DD.
	MaxStops                  *uint32  // 5; nil: unrestricted, 0: nonstop.
	IncludeAirlines           []string // 6
	ExcludeAirlines           []string // 7
	DepartureStartHour        *uint32  // 8
	DepartureEndHourInclusive *uint32  // 9
	ArrivalStartHour          *uint32  // 10
	ArrivalEndHourInclusive   *uint32  // 11
	MaxDurationMinutes        *uint32  // 12
	Origin                    *Place   // 13
	Destination               *Place   // 14
	IncludeConnectionAirports []string // 15
	ExcludeConnectionAirports []string // 16
	MaxLayoverMinutes         *uint32  // 18
	EmissionsFlags            []uint32 // 19; packed.
}

// Place identifies an airport or Google entity. EntityType 1 denotes an airport;
// the distinction between entity types 2 and 3 remains uncertain.
type Place struct {
	EntityType *uint32 // 1
	EntityID   *string // 2; IATA code or Google entity ID.
}

// BaggageFilter retains research names because their meanings are uncertain.
type BaggageFilter struct {
	Field2 *uint32 // 2
	Field3 *uint32 // 3
}

// AllResultsFlag retains the observed sentinel without assuming its meaning.
type AllResultsFlag struct {
	Sentinel *uint64 // 1
}

// PassengerType identifies one traveler's category.
type PassengerType int32

const (
	PassengerUnknown      PassengerType = 0
	PassengerAdult        PassengerType = 1
	PassengerChild        PassengerType = 2
	PassengerInfantOnLap  PassengerType = 3
	PassengerInfantInSeat PassengerType = 4
)

// Cabin identifies the requested service class.
type Cabin int32

const (
	CabinUnknown        Cabin = 0
	CabinEconomy        Cabin = 1
	CabinPremiumEconomy Cabin = 2
	CabinBusiness       Cabin = 3
	CabinFirst          Cabin = 4
)

// TripType identifies the search's journey structure.
type TripType int32

const (
	TripUnknown   TripType = 0
	TripRoundTrip TripType = 1
	TripOneWay    TripType = 2
	TripMultiCity TripType = 3
)
