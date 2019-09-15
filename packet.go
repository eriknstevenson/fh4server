package fh4server

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/golang/glog"
)

// Packet holds the parsed fields and tags from the raw message bytes.
type Packet struct {
	Fields map[string]interface{}
	Tags   map[string]string
}

// ParseBuf attempts to decode and parse the provided encoded packet buffer.
func ParseBuf(buf *bytes.Buffer, whitelist Whitelist) Packet {
	fields := make(map[string]interface{})
	tags := make(map[string]string)
	for _, element := range fh4PacketDefinition {
		value := element.parse(buf)
		if element.finisher != nil {
			value = element.finisher(value)
		}
		if value == nil {
			continue
		}
		if !whitelist(element.label) {
			continue
		}
		switch element.elementType {
		case field:
			fields[element.label] = value
			break
		case tag:
			tags[element.label] = fmt.Sprint(value)
			break
		case timestamp:
			break
		case none:
			break
		default:
			glog.Infof("unexpected elementType encountered: %v", element.elementType)
		}
	}
	return Packet{Fields: fields, Tags: tags}
}

// Parse attempts to decode and parse the provided encoded packet.
// Deprecated. Use ParseBuf instead.
func Parse(packet []byte, whitelist func(string) bool) Packet {
	// TODO: move this validation to the FH4Game packet source.
	const packetSize = 324
	if len(packet) != packetSize {
		glog.Errorf("unexpected packet size (expected %d, got %d)", packetSize, len(packet))
	}
	buf := bytes.NewBuffer(packet)
	return ParseBuf(buf, whitelist)
}

const (
	none      byte = iota
	timestamp byte = iota
	field     byte = iota
	tag       byte = iota
)

// packetElements are defined as functions which consume n bytes from a buffer,
// and return a key/value pair, the type of influx data (field vs tag) it is, or an error.
type packetElement struct {
	label string
	parse func(*bytes.Buffer) interface{}

	// elementType holds whether this packetElement is a field or a tag in influx.
	// Note that by default, packets are specified as none, meaning they will
	// not go to influx.
	// Also, a packetElement can not be both a field and a tag.
	elementType byte

	// unimplemented
	finisher func(interface{}) interface{}
}

func (p *packetElement) withLabel(label string) *packetElement {
	p.label = label
	return p
}

func (p *packetElement) field() *packetElement {
	p.elementType = field
	return p
}

func (p *packetElement) tag() *packetElement {
	p.elementType = tag
	return p
}

func (p *packetElement) timestamp() *packetElement {
	p.elementType = timestamp
	return p
}

func s8() *packetElement {
	return &packetElement{parse: func(buf *bytes.Buffer) interface{} {
		var parsed int8
		binary.Read(bytes.NewBuffer(buf.Next(1)), binary.LittleEndian, &parsed)
		return parsed
	}}
}

func s32() *packetElement {
	return &packetElement{parse: func(buf *bytes.Buffer) interface{} {
		var parsed int32
		binary.Read(bytes.NewBuffer(buf.Next(4)), binary.LittleEndian, &parsed)
		return parsed
	}}
}

func u8() *packetElement {
	return &packetElement{parse: func(buf *bytes.Buffer) interface{} {
		var parsed uint8
		binary.Read(bytes.NewBuffer(buf.Next(1)), binary.LittleEndian, &parsed)
		return parsed
	}}
}

func u16() *packetElement {
	return &packetElement{parse: func(buf *bytes.Buffer) interface{} {
		var parsed uint16
		binary.Read(bytes.NewBuffer(buf.Next(2)), binary.LittleEndian, &parsed)
		return parsed
	}}
}

func u32() *packetElement {
	return &packetElement{parse: func(buf *bytes.Buffer) interface{} {
		var parsed uint32
		binary.Read(bytes.NewBuffer(buf.Next(4)), binary.LittleEndian, &parsed)
		return parsed
	}}
}

func f32() *packetElement {
	return &packetElement{parse: func(buf *bytes.Buffer) interface{} {
		var parsed float32
		binary.Read(bytes.NewBuffer(buf.Next(4)), binary.LittleEndian, &parsed)
		return parsed
		/*bits := binary.LittleEndian.Uint32(buf.Next(4))
		float := math.Float32frombits(bits)
		return float*/
	}}
}

func skipBytes(count int) *packetElement {
	return &packetElement{parse: func(buf *bytes.Buffer) interface{} {
		buf.Next(count)
		return nil
	}}
}

// TODO: put this crap into a configuration file instead.
var (
	// fh4PacketDefinition contains a list of every field in the message we
	// receive from FH4. The order the fields appear in this list is the order
	// in which they exist in the packet.
	//
	// Packet bytes layout
	// [0]-[231]
	// [232]-[243] FH4 new unknown data
	// [244]-[322] FM7 Car Dash data
	// [323] FH4 new unknown data
	//
	// source: https://forums.forzamotorsport.net/turn10_postsm926839_Forza-Motorsport-7--Data-Out--feature-details.aspx
	fh4PacketDefinition = []*packetElement{

		//
		// Start of "sled" format
		//

		// = 1 when race is on. = 0 when in menus/race stopped …
		s32().tag().withLabel("is_race_on"),
		// Can overflow to 0 eventually
		u32().timestamp().withLabel("timestamp_ms"),
		f32().field().withLabel("engine_max_rpm"),
		f32().field().withLabel("engine_idle_rpm"),
		f32().field().withLabel("current_engine_rpm"),

		// In the car's local space; X = right, Y = up, Z = forward
		f32().field().withLabel("acceleration_x"),
		f32().field().withLabel("acceleration_y"),
		f32().field().withLabel("acceleration_z"),

		// In the car's local space; X = right, Y = up, Z = forward
		f32().field().withLabel("velocity_x"),
		f32().field().withLabel("velocity_y"),
		f32().field().withLabel("velocity_z"),

		// In the car's local space; X = pitch, Y = yaw, Z = roll
		f32().field().withLabel("angular_velocity_x"),
		f32().field().withLabel("angular_velocity_y"),
		f32().field().withLabel("angular_velocity_z"),

		f32().field().withLabel("yaw"),
		f32().field().withLabel("pitch"),
		f32().field().withLabel("roll"),

		// Suspension travel normalized: 0.0f = max stretch; 1.0 = max compression
		f32().field().withLabel("normalized_suspension_travel_front_left"),
		f32().field().withLabel("normalized_suspension_travel_front_right"),
		f32().field().withLabel("normalized_suspension_travel_rear_left"),
		f32().field().withLabel("normalized_suspension_travel_rear_right"),

		// Tire normalized slip ratio, = 0 means 100% grip and |ratio| > 1.0 means loss of grip.
		f32().field().withLabel("tire_slip_ratio_front_left"),
		f32().field().withLabel("tire_slip_ratio_front_Right"),
		f32().field().withLabel("tire_slip_ratio_rear_left"),
		f32().field().withLabel("tire_slip_ratio_rear_right"),

		// Wheel rotation speed radians/sec.
		f32().field().withLabel("wheel_rotation_speed_front_left"),
		f32().field().withLabel("wheel_rotation_speed_front_right"),
		f32().field().withLabel("wheel_rotation_speed_rear_left"),
		f32().field().withLabel("wheel_rotation_speed_rear_right"),

		// = 1 when wheel is on rumble strip, = 0 when off.
		s32().field().withLabel("on_rumble_strip_front_left"),
		s32().field().withLabel("on_rumble_strip_front_right"),
		s32().field().withLabel("on_rumble_strip_rear_left"),
		s32().field().withLabel("on_rumble_strip_rear_right"),

		// = from 0 to 1, where 1 is the deepest puddle
		f32().field().withLabel("puddle_depth_front_left"),
		f32().field().withLabel("puddle_depth_front_right"),
		f32().field().withLabel("puddle_depth_rear_left"),
		f32().field().withLabel("puddle_depth_rear_right"),

		// Non-dimensional surface rumble values passed to controller force feedback
		f32().field().withLabel("surface_rumble_front_left"),
		f32().field().withLabel("surface_rumble_front_right"),
		f32().field().withLabel("surface_rumble_rear_left"),
		f32().field().withLabel("surface_rumble_rear_right"),

		// Tire normalized slip angle, = 0 means 100% grip and |angle| > 1.0 means loss of grip.
		f32().field().withLabel("tire_slip_angle_front_left"),
		f32().field().withLabel("tire_slip_angle_front_right"),
		f32().field().withLabel("tire_slip_angle_rear_left"),
		f32().field().withLabel("tire_slip_angle_rear_right"),

		// Tire normalized combined slip, = 0 means 100% grip and |slip| > 1.0 means loss of grip.
		f32().field().withLabel("tire_combined_slip_front_left"),
		f32().field().withLabel("tire_combined_slip_front_right"),
		f32().field().withLabel("tire_combined_slip_rear_left"),
		f32().field().withLabel("tire_combined_slip_rear_right"),

		// Actual suspension travel in meters
		f32().field().withLabel("suspension_travel_meters_front_left"),
		f32().field().withLabel("suspension_travel_meters_front_right"),
		f32().field().withLabel("suspension_travel_meters_rear_left"),
		f32().field().withLabel("suspension_travel_meters_rear_right"),

		s32().tag().withLabel("car_id"),
		s32().tag().withLabel("car_class"),             // 0 (worst)-7 (best) inclusive
		s32().tag().withLabel("car_performance_index"), // 100 (slowest)-999 (fastest) inclusive
		s32().tag().withLabel("drive_train_type"),      // 0 = FWD, 1 = RWD, 2 = AWD
		s32().tag().withLabel("num_engine_cylinders"),

		//
		// End of "sled" format
		//

		// Apparently, FH4 has some unknown bytes here that need to be
		// skipped before we get to the dashboard properties.
		skipBytes(12),

		//
		// Start of "Dash" format
		//

		f32().field().withLabel("position_x"), // in meters
		f32().field().withLabel("position_y"),
		f32().field().withLabel("position_z"),

		f32().field().withLabel("speed"),  // in meters per second
		f32().field().withLabel("power"),  // in watts
		f32().field().withLabel("torque"), // in newton meters

		f32().field().withLabel("tire_temp_front_right"),
		f32().field().withLabel("tire_temp_front_left"),
		f32().field().withLabel("tire_temp_rear_left"),
		f32().field().withLabel("tire_temp_rear_right"),

		f32().field().withLabel("boost"),
		f32().field().withLabel("fuel"),
		f32().field().withLabel("distance_traveled"),
		f32().field().withLabel("best_lap_time"),
		f32().field().withLabel("last_lap_time"),
		f32().field().withLabel("current_lap_time"),
		f32().field().withLabel("current_race_time"),

		u16().tag().withLabel("lap_number"),
		u8().field().withLabel("race_position"),

		u8().field().withLabel("accel"),
		u8().field().withLabel("brake"),
		u8().field().withLabel("clutch"),
		u8().field().withLabel("hand_brake"),
		u8().field().withLabel("gear"),
		s8().field().withLabel("steer"),

		s8().field().withLabel("normalized_driving_line"),
		s8().field().withLabel("normalized_ai_brake_difference"),
	}
)
