package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	shared_pb "github.com/reversersed/go-grpc/tree/main/api_user/pkg/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

type ValidationErrors validator.ValidationErrors
type Validator struct {
	*validator.Validate
}

func New() *Validator {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	v.RegisterValidation("primitiveid", validate_PrimitiveId)
	v.RegisterValidation("lowercase", validate_LowercaseCharacter)
	v.RegisterValidation("uppercase", validate_UppercaseCharacter)
	v.RegisterValidation("digitrequired", validate_AtLeastOneDigit)
	v.RegisterValidation("specialsymbol", validate_SpecialSymbol)
	v.RegisterValidation("onlyenglish", validate_OnlyEnglish)
	v.RegisterValidation("eqfield", validate_FieldsEqual)
	return &Validator{v}
}
func (v *Validator) Register(data any, rules map[string]string) {
	v.Validate.RegisterStructValidationMapRules(rules, data)
}
func (v *Validator) StructValidation(data any) error {
	result := v.Validate.Struct(data)

	if result == nil {
		return nil
	}
	if er, ok := result.(*validator.InvalidValidationError); ok {
		return status.Error(codes.Internal, er.Error())
	}
	details := make([]protoadapt.MessageV1, 0)
	for _, i := range result.(validator.ValidationErrors) {
		tag := i.Tag()
		if len(i.Param()) > 0 {
			tag = fmt.Sprintf("%s:%s", i.Tag(), i.Param())
		}
		details = append(details, &shared_pb.ErrorDetail{
			Field:       i.Field(),
			Struct:      i.StructNamespace(),
			Tag:         tag,
			Description: errorToStringByTag(i),
		})
	}
	stat, err := status.New(codes.InvalidArgument, "validation failed, see the details").WithDetails(details...)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	return stat.Err()
}
func errorToStringByTag(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s: field is required", err.Field())
	case "oneof":
		return fmt.Sprintf("%s: field can only be: %s", err.Field(), err.Param())
	case "min":
		return fmt.Sprintf("%s: must be at least %s characters length", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s: can't be more than %s characters length", err.Field(), err.Param())
	case "lte":
		return fmt.Sprintf("%s: must be less or equal than %s", err.Field(), err.Param())
	case "gte":
		return fmt.Sprintf("%s: must be greater or equal than %s", err.Field(), err.Param())
	case "lt":
		return fmt.Sprintf("%s: must be less than %s", err.Field(), err.Param())
	case "gt":
		return fmt.Sprintf("%s: must be greater than %s", err.Field(), err.Param())
	case "email":
		return fmt.Sprintf("%s: must be a valid email", err.Field())
	case "jwt":
		return fmt.Sprintf("%s: must be a JWT token", err.Field())
	case "lowercase":
		return fmt.Sprintf("%s: must contain at least one lowercase character", err.Field())
	case "uppercase":
		return fmt.Sprintf("%s: must contain at least one uppercase character", err.Field())
	case "digitrequired":
		return fmt.Sprintf("%s: must contain at least one digit", err.Field())
	case "specialsymbol":
		return fmt.Sprintf("%s: must contain at least one special symbol", err.Field())
	case "onlyenglish":
		return fmt.Sprintf("%s: must contain only latin characters", err.Field())
	case "primitiveid":
		return fmt.Sprintf("%s: must be a primitive id type", err.Field())
	case "eqfield":
		return fmt.Sprintf("%s: field must be equal to %s field's value", err.Field(), err.Param())
	default:
		return err.Error()
	}
}
func validate_FieldsEqual(fl validator.FieldLevel) bool {
	return fl.Field().String() == fl.Parent().FieldByName(fl.Param()).String()
}
func validate_PrimitiveId(field validator.FieldLevel) bool {
	var obj primitive.ObjectID

	_, err := primitive.ObjectIDFromHex(field.Field().String())
	return (err == nil) || (field.Field().Kind() == reflect.TypeOf(obj).Kind())
}
func validate_OnlyEnglish(field validator.FieldLevel) bool {
	mathed, err := regexp.MatchString(`^[a-zA-Z]+$`, field.Field().String())
	if err != nil {
		return false
	}
	if !mathed {
		return false
	}
	return true
}
func validate_LowercaseCharacter(field validator.FieldLevel) bool {
	mathed, err := regexp.MatchString("[a-z]+", field.Field().String())
	if err != nil {
		return false
	}
	if !mathed {
		return false
	}
	return true
}
func validate_UppercaseCharacter(field validator.FieldLevel) bool {
	mathed, err := regexp.MatchString("[A-Z]+", field.Field().String())
	if err != nil {
		return false
	}
	if !mathed {
		return false
	}
	return true
}
func validate_AtLeastOneDigit(field validator.FieldLevel) bool {
	mathed, err := regexp.MatchString("[0-9]+", field.Field().String())
	if err != nil {
		return false
	}
	if !mathed {
		return false
	}
	return true
}
func validate_SpecialSymbol(field validator.FieldLevel) bool {
	mathed, err := regexp.MatchString("[!@#\\$%\\^&*()_\\+-.,]+", field.Field().String())
	if err != nil {
		return false
	}
	if !mathed {
		return false
	}
	return true
}
