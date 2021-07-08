package autoconfig

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"sync"
	"time"

	"strings"

	"github.com/spf13/viper"
)

const sep = "|"
const sepEscaped = "||"
const sepReplaced = "{{{{{}}}}}"

var vipUpdate = &updater{}

type updater struct {
	sync.RWMutex
}

func (u *updater) set(key string, value interface{}) {
	u.Lock()
	defer u.Unlock()
	viper.Set(key, value)
}

func (u *updater) getString(key string) string {
	u.RLock()
	defer u.RUnlock()
	return viper.GetString(key)
}

// OrPanic loads the given interface from the environment variables,
// stores them into viper but panics in case of failure: typically when there is
// a parsing error on the variable's value or default value.
func OrPanic(i interface{}) {
	err := AutoConfigure(i)
	if err != nil {
		panic(fmt.Errorf("unable to auto configure: %v", err))
	}
}

// ValueOrPanic sets the value of the given interface based on the given tag and stores it into viper.
// It doesn't support duration as the name of the variable is unknown, so it is not possible to guess the unit.
func ValueOrPanic(v interface{}, valueTag string) {
	t := reflect.TypeOf(v).Elem()
	if t.String() == "time.Duration" {
		panic(fmt.Errorf("duration not supported by ValueOrPanic"))
	}
	value := reflect.ValueOf(v).Elem()
	err := applyValue(value, t, "", valueTag)
	if err != nil {
		panic(fmt.Errorf("unable to auto configure value: %v", err))
	}
}

// DurationOrPanic returns the duration for the given tag and and stores it into viper.
// Unit must be something like 'millis', 'hours, etc.
// See function 'getUnitFromFieldName'.
func DurationOrPanic(valueTag, unit string) time.Duration {
	var d time.Duration
	value := reflect.ValueOf(&d).Elem()
	t := reflect.TypeOf(&d).Elem()
	err := applyValue(value, t, unit, valueTag)
	if err != nil {
		panic(fmt.Errorf("unable to auto configure duration: %v", err))
	}
	return d
}

// AutoConfigure loads the given interface from the environment variables,
// stores them into viper and returns an error in case of failure: typically when there is
// a parsing error on the variable's value or default value.
// Prefer OrPanic as most of the time it is better to do a panic when the application
// fails to get its configuration.
func AutoConfigure(i interface{}) error {
	values := reflect.ValueOf(i).Elem()
	types := reflect.TypeOf(i).Elem()
	for i := 0; i < values.NumField(); i++ {
		fValue := values.Field(i)
		fType := types.Field(i)
		valueTag := fType.Tag.Get("value")
		err := applyValue(fValue, fType.Type, fType.Name, valueTag)
		if err != nil {
			return err
		}
	}
	return nil
}

func applyValue(fValue reflect.Value, fType reflect.Type, fTypeName string, valueTag string) error {
	if fValue.CanSet() && valueTag != "" {
		property, value, err := getValueFromTag(valueTag)
		if err != nil {
			return err
		}
		switch fType.String() {
		case "string":
			fValue.SetString(value)
			vipUpdate.set(property, value)
		case "[]string":
			values := strings.Split(value, " ")
			if value == "" {
				values = []string{}
			}
			fValue.Set(reflect.ValueOf(values))
			vipUpdate.set(property, values)
		case "bool":
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("error while parsing boolean value %v for tag %v: %v", value, valueTag, err)
			}
			fValue.SetBool(boolValue)
			vipUpdate.set(property, boolValue)
		case "int", "int8", "int16", "int32", "int64":
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("error while parsing int value %v for tag %v: %v", value, valueTag, err)
			}
			fValue.SetInt(intValue)
			vipUpdate.set(property, intValue)
		case "time.Duration":
			unit := getUnitFromFieldName(fTypeName)
			durationValue, err := time.ParseDuration(value + unit)
			if err != nil {
				return fmt.Errorf("error while parsing durationValue value %v for tag %v: %v", value, valueTag, err)
			}
			fValue.SetInt(durationValue.Nanoseconds())
			vipUpdate.set(property, durationValue.Nanoseconds())
		default:
			return fmt.Errorf("unsupported type for autoconfiguration: %s", fType.Name())
		}
	}
	return nil
}

func getUnitFromFieldName(fieldName string) string {
	fieldName = strings.ToLower(fieldName)
	if strings.HasSuffix(fieldName, "nanos") {
		return "ns"
	}
	if strings.HasSuffix(fieldName, "micros") {
		return "us"
	}
	if strings.HasSuffix(fieldName, "millis") {
		return "ms"
	}
	if strings.HasSuffix(fieldName, "seconds") {
		return "s"
	}
	if strings.HasSuffix(fieldName, "minutes") {
		return "m"
	}
	if strings.HasSuffix(fieldName, "hours") {
		return "h"
	}
	return "ms" //Milliseconds is default
}

func getValueFromTag(tag string) (string, string, error) {
	property, env, def, err := parseTag(tag)
	if err != nil {
		return "", "", err
	}

	err = validatePropertyFormat(property)
	if err != nil {
		return "", "", err
	}

	//Highest Property source
	value, isSet := os.LookupEnv(env)

	if !isSet {
		value = vipUpdate.getString(property)
	}

	//Default Property
	if value == "" {
		value = def
	}

	return property, value, nil
}

func validatePropertyFormat(property string) error {
	matched, _ := regexp.MatchString(`^[a-zA-Z\d\.\-\s]+$`, property)
	if !matched {
		return errors.New("unsupported property format, only letters, digits and dots are supported: " + property)
	}
	return nil
}

func parseTag(tag string) (propertyName string, envName string, defaultValue string, err error) {
	//platform.field.string|defaultstring
	tag = escape(tag)
	tokens := strings.Split(tag, sep)
	if len(tokens) > 2 {
		return "", "", "", fmt.Errorf("error while parsing tag. invalid format (property|default): %s", tag)
	}
	propertyName = tokens[0]
	envName = toEnvName(propertyName)
	if len(tokens) > 1 {
		defaultValue = unescape(tokens[1])
	}
	return
}

func escape(s string) string {
	return strings.Replace(s, sepEscaped, sepReplaced, -1)
}

func unescape(s string) string {
	return strings.Replace(s, sepReplaced, sep, -1)
}

func toEnvName(propertyName string) string {
	toUpper := strings.ToUpper(propertyName)
	toUnderScore := strings.Replace(toUpper, ".", "_", -1)
	withoutDash := strings.Replace(toUnderScore, "-", "", -1)
	return withoutDash
}

// ClearEnvironment deletes all the environment variables and resets viper.
// It should be used for test purpose only.
func ClearEnvironment() {
	os.Clearenv()
	viper.Reset()
}
