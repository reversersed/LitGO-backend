package copier

import (
	"fmt"

	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type option func(*copier.Option)

func CopyOption(opt ...option) copier.Option {
	var option copier.Option

	for _, v := range opt {
		v(&option)
	}
	return option
}

var (
	WithPrimitiveToStringConverter = func(c *copier.Option) {
		c.Converters = append(c.Converters, copier.TypeConverter{SrcType: primitive.ObjectID{}, DstType: string(""), Fn: func(src interface{}) (dst interface{}, err error) {
			s, ok := src.(primitive.ObjectID)
			if !ok {
				return nil, fmt.Errorf("unable to convert %v to primitive object id", src)
			}
			return s.Hex(), nil
		}})
	}
	WithIgnoreEmptyFields = func(c *copier.Option) {
		c.IgnoreEmpty = true
	}
)
