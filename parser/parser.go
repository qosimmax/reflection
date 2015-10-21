package parser

import (
	"errors"
	"flag"
	"reflect"
	"strconv"
	"time"
)

const (
	//Tag names
	TagRequired = "required"
	TagName     = "name"
	TagDefault  = "default"
	TagDesc     = "description"
)

var (
	ErrNotPtr       = errors.New("interface is not pointer or interface")
	ErrRequired     = errors.New("field required")
	ErrDescRequired = errors.New("field description tag required")

	debug = true
)

//Strcut for parsing flags
type Args struct {
	ConfFile     string        `required:"true" name:"config" default:"/etc/daemon.conf" description:"Конфигурационный файл"`
	Daemon       bool          `required:"false" name:"daemon" default:"true" description:"Запуск приложения в режиме daemon"`
	Pool         uint64        `required:"false" name:"pool" default:"5" description:"кол-во пул потоков"`
	TimeOut      float64       `required:"false" name:"timeout" default:"2.5" description:"time out"`
	DurationTime time.Duration `required:"true" name:"duration" default:"3ms" description:"duration time"`
}

type Argument struct {
	Name     string
	Desc     string
	Default  string
	Required string
	Val      interface{}
	Elem     reflect.Value
	DefVal   interface{}
}

func GetArguments(i interface{}) error {

	var args []Argument

	val := reflect.ValueOf(i)

	//check for interface is pointer and can be addressed
	if !val.CanAddr() && val.Kind() != reflect.Ptr {
		return ErrNotPtr
	}

	val = val.Elem()
	ref_type := val.Type()

	for i := 0; i < val.NumField(); i++ {
		//Get Config tag
		arg := Argument{}

		arg.Name = ref_type.Field(i).Tag.Get(TagName)

		if arg.Name == "" {
			continue
		}

		arg.Default = ref_type.Field(i).Tag.Get(TagDefault)
		arg.Desc = ref_type.Field(i).Tag.Get(TagDesc)
		arg.Required = ref_type.Field(i).Tag.Get(TagRequired)

		if arg.Desc == "" {
			return ErrDescRequired
		}

		field := val.Field(i)
		/*
			if debug {
				fmt.Println("field: ", field.Type())
				fmt.Printf("arg: %#v \n", arg)
			}*/

		arg.Elem = field

		switch field.Interface().(type) {
		case bool:
			//converting default value to bool
			val_bool, err := strconv.ParseBool(arg.Default)
			if err != nil {
				return err
			}

			//getting flag value
			val := flag.Bool(arg.Name, val_bool, arg.Desc)
			arg.Val = val
			arg.DefVal = val_bool

		case string:
			val := flag.String(arg.Name, "", arg.Desc)
			arg.Val = val
			arg.DefVal = arg.Default

		case int, int64:
			//converting default value to int64
			val_int, err := strconv.ParseInt(arg.Default, 10, 64)
			if err != nil {
				return err
			}
			//getting flag value
			val := flag.Int64(arg.Name, val_int, arg.Desc)
			arg.Val = val
			arg.DefVal = val_int

		case uint, uint64:
			//converting default value to uint64
			val_uint, err := strconv.ParseUint(arg.Default, 10, 64)
			if err != nil {
				return err
			}
			//getting flag value
			val := flag.Uint64(arg.Name, 0, arg.Desc)
			arg.Val = val
			arg.DefVal = val_uint

		case float64:
			//converting default value to float64
			val_float, err := strconv.ParseFloat(arg.Default, 64)
			if err != nil {
				return err
			}
			//getting flag value
			val := flag.Float64(arg.Name, 0, arg.Desc)
			arg.Val = val
			arg.DefVal = val_float

		case time.Duration:
			//getting flag value
			val := flag.String(arg.Name, "", arg.Desc)
			arg.Val = val
			arg.DefVal = arg.Default

		default:
		}

		args = append(args, arg)
	}

	return ParseArgs(args)
}

func ParseArgs(args []Argument) error {

	flag.Parse()

	for _, arg := range args {
		switch arg.Elem.Interface().(type) {
		case bool:
			val := arg.Val.(*bool)
			arg.Elem.SetBool(*val)
		case string:
			val := arg.Val.(*string)

			//if string field is required return error
			if arg.Required == "true" && *val == "" {
				return ErrRequired
			}

			//if field is not required set val
			arg.Elem.SetString(*val)

			//if field is not required set default value
			if arg.Required == "false" && *val == "" {
				arg.Elem.SetString(arg.DefVal.(string))
			}

		case int, int64:
			val := arg.Val.(*int64)

			//if string field is required return error
			if arg.Required == "true" && *val == 0 {
				return ErrRequired
			}

			//if field is not required set val
			arg.Elem.SetInt(*val)

			//if field is not required set default value
			if arg.Required == "false" && *val == 0 {
				arg.Elem.SetInt(arg.DefVal.(int64))
			}

		case uint, uint64:
			val := arg.Val.(*uint64)

			//if string field is required return error
			if arg.Required == "true" && *val == 0 {
				return ErrRequired
			}

			//if field is not required set val
			arg.Elem.SetUint(*val)

			//if field is not required set default value
			if arg.Required == "false" && *val == 0 {
				arg.Elem.SetUint(arg.DefVal.(uint64))
			}

		case float64:
			val := arg.Val.(*float64)

			//if string field is required return error
			if arg.Required == "true" && *val == 0.0 {
				return ErrRequired
			}

			//if field is not required set val
			arg.Elem.SetFloat(*val)

			//if field is not required set default value
			if arg.Required == "false" && *val == 0.0 {
				arg.Elem.SetFloat(arg.DefVal.(float64))
			}

		case time.Duration:

			val := arg.Val.(*string)

			//if string field is required return error
			if arg.Required == "true" && *val == "" {
				return ErrRequired
			}

			var str_val string
			if arg.Required == "false" && *val == "" {
				str_val = arg.Default
			} else {
				str_val = *val
			}

			//converting default value to Duration
			drt, err := time.ParseDuration(str_val)
			if err != nil {
				return err
			}

			duration_val := reflect.ValueOf(drt)
			arg.Elem.Set(duration_val)
		}

	}

	return nil
}
