# Relaxed Configuration

Relaxed Configuration allows you to configure properties for your application using configuration server, environment variables and command-line arguments.r

## Getting Started

### Tag your structure

The Go spec in the Struct types definition defines tags as :

_A field declaration may be followed by an optional string literal tag, which becomes an attribute for all the fields in the corresponding field declaration. An
empty tag string is equivalent to an absent tag. The tags are made visible through a reflection interface and take part in type identity for structs but are
otherwise ignored._

We defined the "value" tag for your configuration.

In this tag we define the property and and its default value. The property type comes from the variable type; by reflection. As the Field are populated via
reflection, they must be exposed, and so stared by a capital letter.

```go
type configuration struct {
	FieldString          string        `value:"platform.field.string|defaultstring"`
	FieldBool            bool          `value:"platform.field.bool|true"`
	FieldInt             int           `value:"platform.field.int|42"`
	FieldInt8            int8          `value:"platform.field.int-8|8"`
	FieldInt16           int16         `value:"platform.field.int-16|16"`
	FieldInt32           int32         `value:"platform.field.int-32|32"`
	FieldInt64           int64         `value:"platform.field.int-64|64"`
	FieldDuration        time.Duration `value:"platform.field.duration|1001"`
	FieldDurationNanos   time.Duration `value:"platform.field.duration-nano|1001"`
	FieldDurationMicros  time.Duration `value:"platform.field.duration-micro|1001"`
	FieldDurationMillis  time.Duration `value:"platform.field.duration-milli|1001"`
	FieldDurationSeconds time.Duration `value:"platform.field.duration-sec|1001"`
	FieldDurationMinutes time.Duration `value:"platform.field.duration-minutes|1001"`
	FieldDurationHours   time.Duration `value:"platform.field.duration-hours|1001"`
	FieldStringSlice     []string      `value:"platform.field.slice|string1 string2 string3"`
}
```

#### The supported types are :

- string
- bool
- int
- int8
- int16
- int32
- int64
- time.Duration
- []string

#### Duration Unit

The duration can be specified in different units. To specify the unit to the Autoconfiguration, just suffix your variable with the Unit, as follow:

| Units | suffix
|-------|------
|Nanoseconds    |Nanos
|Microseconds   |Micros
|Milliseconds   |Millis
|Seconds        |Seconds
|Minutes        |Minutes
|Hours          |Hours

If nothing is specified, it will be Milliseconds by default.

```go
config := &configuration{}
err := config.AutoConfigure(config)
```

The config structure pointer will be populated from the Input Sources using the priorities defined below.

## Properties Format

A property must contains only alphanumeric characters, hyphen and dots.

For example : my.property-name.value

NOTE: We recommend that properties are stored in lowercase kabab format. i.e. `my.property-name=foo`. But uppercase are supported.

## Input Sources

### Configuration Server

Properties are bound by exact matching with the properties in the Configuration Server.

### Environment Variables

Environment variables are bound by

- uppercasing
- replacing `.` with `_`
- removing `-`

For example: The property `my.property-name.1=foo` is loaded from env variable `MY_PROPERTYNAME_1`.

### Command Line Arguments (flags)

Properties are bound by exact matching with the command line arguments.

For example : side-car --my.property-name.1=foo
