package autoconfig_test

import (
	"eurocontrol.io/digital-platform-product-deployment/pkg/autoconfig"
	"fmt"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestAutoConfigure_Env(t *testing.T) {
	autoconfig.ClearEnvironment()

	valueString := "value_string"
	valueBool := true
	valueInt := 4242
	valueInt8 := int8(88)
	valueInt16 := int16(1616)
	valueInt32 := int32(3232)
	valueInt64 := int64(6464)
	valueDuration := 2002
	valueSliceStr := "stringA stringB stringC"
	valueSlice := []string{"stringA", "stringB", "stringC"}

	os.Setenv("PLATFORM_FIELD_STRING", valueString)
	os.Setenv("PLATFORM_FIELD_BOOL", strconv.FormatBool(valueBool))
	os.Setenv("PLATFORM_FIELD_INT", strconv.Itoa(valueInt))
	os.Setenv("PLATFORM_FIELD_INT8", strconv.Itoa(int(valueInt8)))
	os.Setenv("PLATFORM_FIELD_INT16", strconv.Itoa(int(valueInt16)))
	os.Setenv("PLATFORM_FIELD_INT32", strconv.Itoa(int(valueInt32)))
	os.Setenv("PLATFORM_FIELD_INT64", strconv.Itoa(int(valueInt64)))
	os.Setenv("PLATFORM_FIELD_DURATION", strconv.Itoa(valueDuration))
	os.Setenv("PLATFORM_FIELD_DURATIONNANO", strconv.Itoa(valueDuration))
	os.Setenv("PLATFORM_FIELD_DURATIONMICRO", strconv.Itoa(valueDuration))
	os.Setenv("PLATFORM_FIELD_DURATIONMILLI", strconv.Itoa(valueDuration))
	os.Setenv("PLATFORM_FIELD_DURATIONSEC", strconv.Itoa(valueDuration))
	os.Setenv("PLATFORM_FIELD_DURATIONMINUTES", strconv.Itoa(valueDuration))
	os.Setenv("PLATFORM_FIELD_DURATIONHOURS", strconv.Itoa(valueDuration))
	os.Setenv("PLATFORM_FIELD_NESTED_STRING", valueString)
	os.Setenv("PLATFORM_FIELD_SLICE", valueSliceStr)

	conf := &configuration{}

	err := autoconfig.AutoConfigure(conf)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	assert.Equal(t, valueString, conf.FieldString)
	assert.Equal(t, valueBool, conf.FieldBool)
	assert.Equal(t, valueInt, conf.FieldInt)
	assert.Equal(t, valueInt8, conf.FieldInt8)
	assert.Equal(t, valueInt16, conf.FieldInt16)
	assert.Equal(t, valueInt32, conf.FieldInt32)
	assert.Equal(t, valueInt64, conf.FieldInt64)
	assert.Equal(t, time.Duration(valueDuration*1000*1000), conf.FieldDuration) //Milliseconds is default
	assert.Equal(t, time.Duration(valueDuration), conf.FieldDurationNanos)
	assert.Equal(t, time.Duration(valueDuration*1000), conf.FieldDurationMicros)
	assert.Equal(t, time.Duration(valueDuration*1000*1000), conf.FieldDurationMillis)
	assert.Equal(t, time.Duration(valueDuration*1000*1000*1000), conf.FieldDurationSeconds)
	assert.Equal(t, time.Duration(valueDuration*1000*1000*1000*60), conf.FieldDurationMinutes)
	assert.Equal(t, time.Duration(valueDuration*1000*1000*1000*3600), conf.FieldDurationHours)
	assert.Equal(t, valueSlice, conf.FieldStringSlice)

	assert.Equal(t, valueString, viper.GetString("platform.field.string"))
	assert.Equal(t, valueBool, viper.GetBool("platform.field.bool"))
	assert.Equal(t, valueInt, viper.GetInt("platform.field.int"))
	assert.Equal(t, valueInt8, int8(viper.GetInt("platform.field.int-8")))
	assert.Equal(t, valueInt16, int16(viper.GetInt("platform.field.int-16")))
	assert.Equal(t, valueInt32, viper.GetInt32("platform.field.int-32"))
	assert.Equal(t, valueInt64, viper.GetInt64("platform.field.int-64"))
	assert.Equal(t, time.Duration(valueDuration*1000*1000), viper.GetDuration("platform.field.duration")) //Milliseconds is default
	assert.Equal(t, time.Duration(valueDuration), viper.GetDuration("platform.field.duration-nano"))
	assert.Equal(t, time.Duration(valueDuration*1000), viper.GetDuration("platform.field.duration-micro"))
	assert.Equal(t, time.Duration(valueDuration*1000*1000), viper.GetDuration("platform.field.duration-milli"))
	assert.Equal(t, time.Duration(valueDuration*1000*1000*1000), viper.GetDuration("platform.field.duration-sec"))
	assert.Equal(t, time.Duration(valueDuration*1000*1000*1000*60), viper.GetDuration("platform.field.duration-minutes"))
	assert.Equal(t, time.Duration(valueDuration*1000*1000*1000*3600), viper.GetDuration("platform.field.duration-hours"))
	assert.Equal(t, valueSlice, viper.GetStringSlice("platform.field.slice"))
}

func TestAutoConfigure_Properties(t *testing.T) {
	autoconfig.ClearEnvironment()

	valueString := "property_string"
	valueBool := false
	valueInt := 2121
	valueInt8 := int8(89)
	valueInt16 := int16(1617)
	valueInt32 := int32(3233)
	valueInt64 := int64(6465)
	valueDuration := 2003
	valueSliceStr := "stringA stringB stringC"
	valueSlice := []string{"stringA", "stringB", "stringC"}

	viper.Set("platform.field.string", valueString)
	viper.Set("platform.field.bool", strconv.FormatBool(valueBool))
	viper.Set("platform.field.int", strconv.Itoa(valueInt))
	viper.Set("platform.field.int-8", strconv.Itoa(int(valueInt8)))
	viper.Set("platform.field.int-16", strconv.Itoa(int(valueInt16)))
	viper.Set("platform.field.int-32", strconv.Itoa(int(valueInt32)))
	viper.Set("platform.field.int-64", strconv.Itoa(int(valueInt64)))
	viper.Set("platform.field.duration", strconv.Itoa(valueDuration))
	viper.Set("platform.field.duration-nano", strconv.Itoa(valueDuration))
	viper.Set("platform.field.duration-micro", strconv.Itoa(valueDuration))
	viper.Set("platform.field.duration-milli", strconv.Itoa(valueDuration))
	viper.Set("platform.field.duration-sec", strconv.Itoa(valueDuration))
	viper.Set("platform.field.duration-minutes", strconv.Itoa(valueDuration))
	viper.Set("platform.field.duration-hours", strconv.Itoa(valueDuration))
	viper.Set("platform.field.slice", valueSliceStr)

	conf := &configuration{}

	err := autoconfig.AutoConfigure(conf)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	assert.Equal(t, valueString, conf.FieldString)
	assert.Equal(t, valueBool, conf.FieldBool)
	assert.Equal(t, valueInt, conf.FieldInt)
	assert.Equal(t, valueInt8, conf.FieldInt8)
	assert.Equal(t, valueInt16, conf.FieldInt16)
	assert.Equal(t, valueInt32, conf.FieldInt32)
	assert.Equal(t, valueInt64, conf.FieldInt64)
	assert.Equal(t, time.Duration(valueDuration*1000*1000), conf.FieldDuration) //Milliseconds is default
	assert.Equal(t, time.Duration(valueDuration), conf.FieldDurationNanos)
	assert.Equal(t, time.Duration(valueDuration*1000), conf.FieldDurationMicros)
	assert.Equal(t, time.Duration(valueDuration*1000*1000), conf.FieldDurationMillis)
	assert.Equal(t, time.Duration(valueDuration*1000*1000*1000), conf.FieldDurationSeconds)
	assert.Equal(t, time.Duration(valueDuration*1000*1000*1000*60), conf.FieldDurationMinutes)
	assert.Equal(t, time.Duration(valueDuration*1000*1000*1000*3600), conf.FieldDurationHours)
	assert.Equal(t, valueSlice, conf.FieldStringSlice)

	assert.Equal(t, valueString, viper.GetString("platform.field.string"))
	assert.Equal(t, valueBool, viper.GetBool("platform.field.bool"))
	assert.Equal(t, valueInt, viper.GetInt("platform.field.int"))
	assert.Equal(t, valueInt8, int8(viper.GetInt("platform.field.int-8")))
	assert.Equal(t, valueInt16, int16(viper.GetInt("platform.field.int-16")))
	assert.Equal(t, valueInt32, viper.GetInt32("platform.field.int-32"))
	assert.Equal(t, valueInt64, viper.GetInt64("platform.field.int-64"))
	assert.Equal(t, time.Duration(valueDuration*1000*1000), viper.GetDuration("platform.field.duration")) //Milliseconds is default
	assert.Equal(t, time.Duration(valueDuration), viper.GetDuration("platform.field.duration-nano"))
	assert.Equal(t, time.Duration(valueDuration*1000), viper.GetDuration("platform.field.duration-micro"))
	assert.Equal(t, time.Duration(valueDuration*1000*1000), viper.GetDuration("platform.field.duration-milli"))
	assert.Equal(t, time.Duration(valueDuration*1000*1000*1000), viper.GetDuration("platform.field.duration-sec"))
	assert.Equal(t, time.Duration(valueDuration*1000*1000*1000*60), viper.GetDuration("platform.field.duration-minutes"))
	assert.Equal(t, time.Duration(valueDuration*1000*1000*1000*3600), viper.GetDuration("platform.field.duration-hours"))
	assert.Equal(t, valueSlice, viper.GetStringSlice("platform.field.slice"))
}

func TestAutoConfigure_Default(t *testing.T) {
	autoconfig.ClearEnvironment()

	conf := &configuration{}

	err := autoconfig.AutoConfigure(conf)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	assert.Equal(t, "defaultstring", conf.FieldString)
	assert.Equal(t, true, conf.FieldBool)
	assert.Equal(t, 42, conf.FieldInt)
	assert.Equal(t, int8(8), conf.FieldInt8)
	assert.Equal(t, int16(16), conf.FieldInt16)
	assert.Equal(t, int32(32), conf.FieldInt32)
	assert.Equal(t, int64(64), conf.FieldInt64)
	assert.Equal(t, time.Duration(1001*1000*1000), conf.FieldDuration) //Milliseconds is default
	assert.Equal(t, time.Duration(1001), conf.FieldDurationNanos)
	assert.Equal(t, time.Duration(1001*1000), conf.FieldDurationMicros)
	assert.Equal(t, time.Duration(1001*1000*1000), conf.FieldDurationMillis)
	assert.Equal(t, time.Duration(1001*1000*1000*1000), conf.FieldDurationSeconds)
	assert.Equal(t, time.Duration(1001*1000*1000*1000*60), conf.FieldDurationMinutes)
	assert.Equal(t, time.Duration(1001*1000*1000*1000*3600), conf.FieldDurationHours)
	assert.Equal(t, []string{"string1", "string2", "string3"}, conf.FieldStringSlice)
}

func TestAutoConfigure_Viper_Is_Set_With_Env(t *testing.T) {
	autoconfig.ClearEnvironment()

	type simpleConf struct {
		FieldString string `value:"platform.field.string"`
	}

	valueString := "value_string"

	os.Setenv("PLATFORM_FIELD_STRING", valueString)

	conf := &simpleConf{}

	err := autoconfig.AutoConfigure(conf)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	assert.Equal(t, valueString, conf.FieldString)
	//assert.Equal(t, valueString, viper.GetString("platform.field.string"))
}

func TestAutoConfigure_No_Default(t *testing.T) {
	autoconfig.ClearEnvironment()

	type simpleConf struct {
		FieldString string `value:"platform.field.string"`
	}

	conf := &simpleConf{}

	err := autoconfig.AutoConfigure(conf)

	if err != nil {
		assert.FailNow(t, err.Error())
	}

	assert.Empty(t, conf.FieldString)
}

func TestAutoConfigure_Priorities_Env_Duration(t *testing.T) {
	autoconfig.ClearEnvironment()
	os.Setenv("PLATFORM_FIELD_DURATIONSEC", "75")

	conf := &configuration{}

	err := autoconfig.AutoConfigure(conf)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	assert.Equal(t, 75*time.Second, conf.FieldDurationSeconds)
}

func TestAutoConfigure_Priorities_Env_Number1(t *testing.T) {
	valueString := "env_string"

	autoconfig.ClearEnvironment()
	os.Setenv("PLATFORM_FIELD_STRING", valueString)
	viper.Set("platform.field.string", "property_string")

	conf := &configuration{}

	err := autoconfig.AutoConfigure(conf)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	assert.Equal(t, valueString, conf.FieldString)
}

func TestAutoConfigure_Err_Parsing_Tag(t *testing.T) {
	autoconfig.ClearEnvironment()
	type wrongConf struct {
		FieldString string `value:"platform.field.string|defaultstring|unexpected"`
	}

	conf := &wrongConf{}

	err := autoconfig.AutoConfigure(conf)

	assert.EqualError(t, err, "error while parsing tag. invalid format (property|default): platform.field.string|defaultstring|unexpected")
}

func TestAutoConfigure_Err_Unsupported_Field_Type(t *testing.T) {
	autoconfig.ClearEnvironment()
	type unsupportedConf struct {
		Field net.Conn `value:"platform.field.conn"`
	}

	conf := &unsupportedConf{}

	err := autoconfig.AutoConfigure(conf)

	assert.EqualError(t, err, "unsupported type for autoconfiguration: Conn")
}

func TestAutoConfigure_Escape_Separator(t *testing.T) {
	autoconfig.ClearEnvironment()
	type aConf struct {
		ACondition string `value:"platform.field.condition|blue||red"`
	}

	conf := &aConf{}

	err := autoconfig.AutoConfigure(conf)

	if err != nil {
		assert.FailNow(t, err.Error())
	}

	assert.Equal(t, "blue|red", conf.ACondition)
}

func TestAutoConfigure_Property_Format(t *testing.T) {
	autoconfig.ClearEnvironment()
	type supportedConf struct {
		Field string `value:"platform.field-NAME.123"`
	}

	conf := &supportedConf{}

	err := autoconfig.AutoConfigure(conf)

	assert.Nil(t, err)

}

func TestAutoConfigure_Err_Invalid_Property_Format(t *testing.T) {
	autoconfig.ClearEnvironment()
	type unsupportedConf struct {
		Field string `value:"platform.$.conn"`
	}

	conf := &unsupportedConf{}

	err := autoconfig.AutoConfigure(conf)

	assert.EqualError(t, err, "unsupported property format, only letters, digits and dots are supported: platform.$.conn")

}

func TestAutoConfigure_NoDefaultValueSlice(t *testing.T) {
	autoconfig.ClearEnvironment()
	type conf struct {
		List []string `value:"i.am.an.empty.slice"`
	}

	c := &conf{}
	err := autoconfig.AutoConfigure(c)

	assert.Nil(t, err)
	assert.Len(t, c.List, 0)
}

func TestAutoConfigure_BlankInTheWay(t *testing.T) {
	autoconfig.ClearEnvironment()
	type conf struct {
		// Bad syntax for struct tag value
		Blank string `value: "there.is.blank.just.before"`
	}

	os.Setenv("THERE_IS_BLANK_JUST_BEFORE", "space")
	c := &conf{}
	err := autoconfig.AutoConfigure(c)

	assert.Nil(t, err)
	assert.Equal(t, c.Blank, "")
}

func TestValueOrPanic(t *testing.T) {
	autoconfig.ClearEnvironment()

	// overwrite value
	s := "hello"
	autoconfig.ValueOrPanic(&s, "testing.power|bye")
	assert.Equal(t, "bye", s)

	// nil slice
	var slice []string
	assert.Nil(t, slice)
	autoconfig.ValueOrPanic(&slice, "nil.slice|a b")
	assert.Len(t, slice, 2)

	// slice and no default
	autoconfig.ValueOrPanic(&slice, "nil.slice")
	assert.Len(t, slice, 0)

	// slice
	os.Setenv("NIL_SLICE", "10 20")
	autoconfig.ValueOrPanic(&slice, "nil.slice")
	assert.Len(t, slice, 2)
	assert.Contains(t, slice, "10")
	assert.Contains(t, slice, "20")
}

func TestValueOrPanic_duration_not_supported(t *testing.T) {
	defer assertPanic(t, "duration not supported by ValueOrPanic")
	var tenMinutes time.Duration
	autoconfig.ValueOrPanic(&tenMinutes, "ten.minutes")
}

func assertPanic(t *testing.T, expected string) {
	r := recover()
	require.NotNil(t, r)
	require.Equal(t, expected, fmt.Sprintf("%v", r))
}
