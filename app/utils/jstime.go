/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-25 18:43:52
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-04-13 09:46:48
 * @Description:
 */
package utils

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type Time time.Time

const (
	timeFormart = "2006-01-02 15:04:05"
)

/**
 * @description:
 * @params {type}
 * @return:
 */
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).Format(timeFormart) + `"`), nil

	b := make([]byte, 0, len(timeFormart)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormart)
	b = append(b, '"')
	return b, nil
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (t Time) String() string {

	return time.Time(t).Format(timeFormart)
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (t Time) Unix() int64 {
	return time.Time(t).Unix()
}

/**
 * @description: MarshalBSON marshal bson
 * @params {type}
 * @return:
 */
func (t Time) MarshalBSON() ([]byte, error) {
	txt, err := time.Time(t).MarshalText()
	if err != nil {
		return nil, err
	}
	b, err := bson.Marshal(map[string]string{"t": string(txt)})
	return b, err
}

/**
 * @description:MarshalBSONValue marshal bson value
 * @params {type}
 * @return:
 */
func (t *Time) MarshalBSONValue() (bsontype.Type, []byte, error) {
	b, err := bson.Marshal(time.Time(*t))
	return bson.TypeEmbeddedDocument, b, err
}

/**
 * @description: UnmarshalBSON unmarshal bson
 * @params {type}
 * @return:
 */
func (t *Time) UnmarshalBSON(data []byte) error {
	var err error
	var d bson.D
	err = bson.Unmarshal(data, &d)
	if err != nil {
		return err
	}
	if v, ok := d.Map()["t"]; ok {
		tt := time.Time{}
		err = tt.UnmarshalText([]byte(v.(string)))
		*t = Time(tt)
		return err
	}
	return fmt.Errorf("key 't' missing")
}
