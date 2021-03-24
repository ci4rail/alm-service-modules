package gpsd

// GenericClass is used to determine the class of the incoming messages
type GenericClass struct {
	Class string `json:"class"`
}

// Please have a look at https://gpsd.gitlab.io/gpsd/gpsd_json.html#_core_protocol_responses for reference on the meaning

// Tpv is a Time-Position-Velocity Object
type Tpv struct {
	Class       string  `json:"class"`
	Device      string  `json:"device,omitempty"`
	Mode        float64 `json:"mode"`
	Status      float64 `json:"status,omitempty"`
	Time        string  `json:"time,omitempty"`
	Althae      float64 `json:"altHAE,omitempty"`
	Altmsl      float64 `json:"altMSL,omitempty"`
	Alt         float64 `json:"alt,omitempty"`
	Climb       float64 `json:"climb,omitempty"`
	Datum       string  `json:"datum,omitempty"`
	Depth       float64 `json:"depth,omitempty"`
	Dgpsage     float64 `json:"dgpsAge,omitempty"`
	Dgpssta     float64 `json:"dgpsSta,omitempty"`
	Epc         float64 `json:"epc,omitempty"`
	Epd         float64 `json:"epd,omitempty"`
	Eph         float64 `json:"eph,omitempty"`
	Eps         float64 `json:"eps,omitempty"`
	Ept         float64 `json:"ept,omitempty"`
	Epx         float64 `json:"epx,omitempty"`
	Epy         float64 `json:"epy,omitempty"`
	Epv         float64 `json:"epv,omitempty"`
	Geoidsep    float64 `json:"geoidSep,omitempty"`
	Lat         float64 `json:"lat,omitempty"`
	Leapseconds int     `json:"leapseconds,omitempty"`
	Lon         float64 `json:"lon,omitempty"`
	Track       float64 `json:"track,omitempty"`
	Magtrack    float64 `json:"magtrack,omitempty"`
	Magvar      float64 `json:"magvar,omitempty"`
	Speed       float64 `json:"speed,omitempty"`
	Ecefx       float64 `json:"ecefx,omitempty"`
	Ecefy       float64 `json:"ecefy,omitempty"`
	Ecefz       float64 `json:"ecefz,omitempty"`
	Ecefpacc    float64 `json:"ecefpAcc,omitempty"`
	Ecefvx      float64 `json:"ecefvx,omitempty"`
	Ecefvy      float64 `json:"ecefvy,omitempty"`
	Ecefvz      float64 `json:"ecefvz,omitempty"`
	Ecefvacc    float64 `json:"ecefvAcc,omitempty"`
	Sep         float64 `json:"sep,omitempty"`
	Reld        float64 `json:"relD,omitempty"`
	Rele        float64 `json:"relE,omitempty"`
	Reln        float64 `json:"relN,omitempty"`
	Veld        float64 `json:"velD,omitempty"`
	Vele        float64 `json:"velE,omitempty"`
	Veln        float64 `json:"velN,omitempty"`
	Wanglem     float64 `json:"wanglem,omitempty"`
	Wangler     float64 `json:"wangler,omitempty"`
	Wanglet     float64 `json:"wanglet,omitempty"`
	Wspeedr     float64 `json:"wspeedr,omitempty"`
	Wspeedt     float64 `json:"wspeedt,omitempty"`
}

// Sky reports a sky view of the GPS satellite positions
type Sky struct {
	Class      string      `json:"class"`
	Device     string      `json:"device,omitempty"`
	Time       string      `json:"time,omitempty"`
	Gdop       float64     `json:"gdop,omitempty"`
	Hdop       float64     `json:"hdop,omitempty"`
	Pdop       float64     `json:"pdop,omitempty"`
	Tdop       float64     `json:"tdop,omitempty"`
	Vdop       float64     `json:"vdop,omitempty"`
	Xdop       float64     `json:"xdop,omitempty"`
	Ydop       float64     `json:"ydop,omitempty"`
	Nsat       float64     `json:"nSat,omitempty"`
	Usat       float64     `json:"uSat,omitempty"`
	Satellites []Satellite `json:"satellites,omitempty"`
}

// Satellite is always shipped with a Sky object
type Satellite struct {
	Prn    float64 `json:"PRN"`
	Az     float64 `json:"az,omitempty"`
	El     float64 `json:"el,omitempty"`
	Ss     float64 `json:"ss,omitempty"`
	Used   bool    `json:"used"`
	Gnssid float64 `json:"gnssid,omitempty"`
	Svid   float64 `json:"svid,omitempty"`
	Sigid  float64 `json:"sigid,omitempty"`
	Freqid float64 `json:"freqid,omitempty"`
	Health float64 `json:"health,omitempty"`
}

// Gst is a pseudorange noise report
type Gst struct {
	Class  string  `json:"class"`
	Device string  `json:"device,omitempty"`
	Time   string  `json:"time,omitempty"`
	Rms    float64 `json:"rms,omitempty"`
	Major  float64 `json:"major,omitempty"`
	Minor  float64 `json:"minor,omitempty"`
	Orient float64 `json:"orient,omitempty"`
	Lat    float64 `json:"lat,omitempty"`
	Lon    float64 `json:"lon,omitempty"`
	Alt    float64 `json:"alt,omitempty"`
}

// Att is a vehicle-attitude report
type Att struct {
	Class   string  `json:"class"`
	Device  string  `json:"device,omitempty"`
	Time    string  `json:"time,omitempty"`
	Heading float64 `json:"heading,omitempty"`
	MagSt   string  `json:"mag_st,omitempty"`
	Pitch   float64 `json:"pitch,omitempty"`
	PitchSt string  `json:"pitch_st,omitempty"`
	Yaw     float64 `json:"yaw,omitempty"`
	YawSt   string  `json:"yaw_st,omitempty"`
	Roll    float64 `json:"roll,omitempty"`
	RollSt  string  `json:"roll_st,omitempty"`
	Dip     float64 `json:"dip,omitempty"`
	MagLen  float64 `json:"mag_len,omitempty"`
	MagX    float64 `json:"mag_x,omitempty"`
	MagY    float64 `json:"mag_y,omitempty"`
	MagZ    float64 `json:"mag_z,omitempty"`
	AccLen  float64 `json:"acc_len,omitempty"`
	AccX    float64 `json:"acc_x,omitempty"`
	AccY    float64 `json:"acc_y,omitempty"`
	AccZ    float64 `json:"acc_z,omitempty"`
	GyroX   float64 `json:"gyro_x,omitempty"`
	GyroY   float64 `json:"gyro_y,omitempty"`
	Depth   float64 `json:"depth,omitempty"`
	Temp    float64 `json:"temp,omitempty"`
}

// Toff reports the GPS time as derived from the GPS serial data stream
type Toff struct {
	Class     string  `json:"class"`
	Device    string  `json:"device"`
	RealSec   float64 `json:"real_sec"`
	RealNsec  float64 `json:"real_nsec"`
	ClockSec  float64 `json:"clock_sec"`
	ClockNsec float64 `json:"clock_nsec"`
}

// Pps reports the GPS time as derived from the GPS PPS pulse
type Pps struct {
	Class     string  `json:"class"`
	Device    string  `json:"device"`
	RealSec   float64 `json:"real_sec"`
	RealNsec  float64 `json:"real_nsec"`
	ClockSec  float64 `json:"clock_sec"`
	ClockNsec float64 `json:"clock_nsec"`
	Precision float64 `json:"precision"`
	Qerr      float64 `json:"qErr,omitempty"`
}

// Osc reports the status of a GPS-disciplined oscillator
type Osc struct {
	Class       string  `json:"class"`
	Device      string  `json:"device"`
	Running     bool    `json:"running"`
	Reference   bool    `json:"reference"`
	Disciplined bool    `json:"disciplined"`
	Delta       float64 `json:"delta"`
}

// Version reports protocol specific versioning information
type Version struct {
	Class      string  `json:"class"`
	Release    string  `json:"release"`
	Rev        string  `json:"rev"`
	ProtoMajor float64 `json:"proto_major"`
	ProtoMinor float64 `json:"proto_minor"`
	Remote     string  `json:"remote"`
}

// Devices contains a list of devices
type Devices struct {
	Class   string   `json:"class"`
	Devices []Device `json:"devices"`
	Remote  string   `json:"remote,omitempty"`
}

// Device contains device specific information
type Device struct {
	Class     string `json:"class"`
	Path      string `json:"path,omitempty"`
	Activated string `json:"activated,omitempty"`
	Flags     int    `json:"flags,omitempty"`
	Driver    string `json:"driver,omitempty"`
	Subtype   string `json:"subtype,omitempty"`
	Subtype1  string `json:"subtype1,omitempty"`
	Bps       int    `json:"bps,omitempty"`
	Parity    string `json:"parity,omitempty"`
	Stopbits  string `json:"stopbits"`
	Native    int    `json:"native,omitempty"`
	Cycle     int    `json:"cycle,omitempty"`
	Mincycle  int    `json:"mincycle,omitempty"`
}

// Watch sets the watcher mode
type Watch struct {
	Class   string `json:"class,omitempty"`
	Enable  bool   `json:"enable,omitempty"`
	JSON    bool   `json:"json,omitempty"`
	Nmea    bool   `json:"nmea,omitempty"`
	Raw     int    `json:"raw,omitempty"`
	Scaled  bool   `json:"scaled,omitempty"`
	Split24 bool   `json:"split24,omitempty"`
	Pps     bool   `json:"pps,omitempty"`
	Device  string `json:"device,omitempty"`
	Remote  string `json:"remote,omitempty"`
}

// Error contains error messages coming from gpsd daemon
type Error struct {
	Class   string `json:"class"`
	Message string `json:"message"`
}
