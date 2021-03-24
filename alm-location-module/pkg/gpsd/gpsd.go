package gpsd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

// The API is based on this: https://gpsd.gitlab.io/gpsd/client-howto.html#_interfacing_from_the_client_side

// FilterHandler is a generic type for handling all the different messages
type FilterHandler func(interface{})

// Connection is the struct that contains all necessary stuff for a session
type Connection struct {
	socket  net.Conn
	filters map[string]FilterHandler
	reader  *bufio.Reader
	close   chan bool
}

// NewClient creates a new gpsd client
func NewClient(addr string) (*Connection, error) {
	var err error
	socket, err := net.Dial("tcp", addr)
	if err != nil {
		return &Connection{}, err
	}
	// Create bufio reader and wait until a new frame begins
	reader := bufio.NewReader(socket)
	_, err = reader.ReadString('\n')
	if err != nil {
		return &Connection{}, err
	}
	filters := make(map[string]FilterHandler)

	return &Connection{
		socket:  socket,
		reader:  reader,
		filters: filters,
		close:   make(chan bool),
	}, nil
}

// Close closes the client connection
func (c *Connection) Close() {
	c.close <- true
}

// Register connects a `class` with a handler function.
// See https://gpsd.gitlab.io/gpsd/gpsd_json.html for valid classes

// RegisterTpv registers a handler fot TPV messages
func (c *Connection) RegisterTpv(h FilterHandler) {
	c.filters["TPV"] = h
}

// RegisterSky registers a handler fot SKY messages
func (c *Connection) RegisterSky(h FilterHandler) {
	c.filters["SKY"] = h
}

// RegisterGst registers a handler fot GST messages
func (c *Connection) RegisterGst(h FilterHandler) {
	c.filters["GST"] = h
}

// RegisterAtt registers a handler fot ATT messages
func (c *Connection) RegisterAtt(h FilterHandler) {
	c.filters["ATT"] = h
}

// RegisterToff registers a handler fot TOFF messages
func (c *Connection) RegisterToff(h FilterHandler) {
	c.filters["TOFF"] = h
}

// RegisterPps registers a handler fot PPS messages
func (c *Connection) RegisterPps(h FilterHandler) {
	c.filters["PPS"] = h
}

// RegisterOsc registers a handler fot OSC messages
func (c *Connection) RegisterOsc(h FilterHandler) {
	c.filters["OSC"] = h
}

// RegisterVersion registers a handler fot VERSION messages
func (c *Connection) RegisterVersion(h FilterHandler) {
	c.filters["VERSION"] = h
}

// RegisterDevices registers a handler fot DEVICES messages
func (c *Connection) RegisterDevices(h FilterHandler) {
	c.filters["DEVICES"] = h
}

// RegisterError registers a handler fot ERROR messages
func (c *Connection) RegisterError(h FilterHandler) {
	c.filters["ERROR"] = h
}

// func (c *Connection) Register(class Class, h FilterHandler) {
// 	c.filters[class.String()] = h
// }

// Watch enables watching to gpsd server
func (c *Connection) Watch() (chan bool, error) {
	// Send the WATCH command to gpsd
	j, err := json.Marshal(Watch{
		Enable: true,
		JSON:   true,
	})
	if err != nil {
		return nil, err
	}
	err = c.command("WATCH=" + string(j))
	if err != nil {
		return nil, err
	}
	done := make(chan bool)
	go c.watchLoop()
	return done, nil
}

func (c *Connection) watchLoop() {
	for {
		select {
		case <-c.close:
			return
		default:
			if line, err := c.reader.ReadString('\n'); err == nil {
				var class GenericClass
				if err = json.Unmarshal([]byte(line), &class); err == nil {
					// Check if class is registered, otherwise skip message
					if _, ok := c.filters[class.Class]; ok {
						r, err := unmarshallClass(class.Class, line)
						if err != nil {
							fmt.Printf("Error: cannot unmarshal message of class %s: %s\n", class.Class, line)
						}
						handler := c.filters[class.Class]
						if handler != nil {
							handler(r)
						} else {
							fmt.Printf("Error: no filter handler set for class %s\n", class.Class)
						}
					} else {
						continue
					}
				} else {
					fmt.Printf("Error: cannot detect class of line: %s\n", line)
				}
			} else {
				fmt.Println("Error: cannot read from gpsd")
				return
			}
		}
	}
}

func (c *Connection) command(cmd string) error {
	_, err := fmt.Fprintf(c.socket, "?%s;", cmd)
	return err
}

func unmarshallClass(class string, line string) (interface{}, error) {
	bytes := []byte(line)
	switch class {
	case "TPV":
		var data *Tpv
		err := json.Unmarshal(bytes, &data)
		return data, err
	case "SKY":
		var data *Sky
		err := json.Unmarshal(bytes, &data)
		return data, err
	case "GST":
		var data *Gst
		err := json.Unmarshal(bytes, &data)
		return data, err
	case "ATT":
		var data *Att
		err := json.Unmarshal(bytes, &data)
		return data, err
	case "TOFF":
		var data *Toff
		err := json.Unmarshal(bytes, &data)
		return data, err
	case "PPS":
		var data *Pps
		err := json.Unmarshal(bytes, &data)
		return data, err
	case "OSC":
		var data *Osc
		err := json.Unmarshal(bytes, &data)
		return data, err
	case "VERSION":
		var data *Version
		err := json.Unmarshal(bytes, &data)
		return data, err
	case "DEVICES":
		var data *Devices
		err := json.Unmarshal(bytes, &data)
		return data, err
	case "ERROR":
		var data *Error
		err := json.Unmarshal(bytes, &data)
		return data, err
	}
	return nil, fmt.Errorf("Error: unknown class")
}
