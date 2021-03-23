package gpsd

// Class used for enum mapping
type Class int

const (
	// Tpv enum value
	Tpv Class = 0
	// Sky enum value
	Sky Class = 1
	// Gst enum value
	Gst Class = 2
	// Att enum value
	Att Class = 3
	// Toff enum value
	Toff Class = 4
	// Pps enum value
	Pps Class = 5
	// Osc enum value
	Osc Class = 6
	// Version enum value
	Version Class = 7
	// Devices enum value
	Devices Class = 8
	// Error enum value
	Error Class = 9
)

// String converts enum value to string
func (c Class) String() string {
	classes := [...]string{
		"TPV",
		"SKY",
		"GST",
		"ATT",
		"TOFF",
		"PPS",
		"OSC",
		"VERSION",
		"DEVICES",
		"ERROR",
	}
	if c < Tpv || c > Error {
		return "Unknown"
	}
	return classes[c]
}
