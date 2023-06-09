package model

import (
	"database/sql/driver"
	"errors"
	"strconv"
)

// BitBool is an implementation of a bool for the MySQL type BIT(1).
// This type allows you to avoid wasting an entire byte for MySQL's boolean type TINYINT.
type BitBool bool

// Value implements the driver.Valuer interface,
// and turns the BitBool into a bitfield (BIT(1)) for MySQL storage.
func (b BitBool) Value() (driver.Value, error) {
	if b {
		return []byte{1}, nil
	} else {
		return []byte{0}, nil
	}
}

// Scan implements the sql.Scanner interface,
// and turns the bitfield incoming from MySQL into a BitBool
func (b *BitBool) Scan(src interface{}) error {
	v, ok := src.([]byte)
	if !ok {
		return errors.New("bad []byte type assertion")
	}
	*b = v[0] == 1
	return nil
}

// MarshalJSON
//
//	@Description:
//	@receiver b
//	@return []byte
//	@return error
func (b BitBool) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatBool(bool(b))), nil
}

// UnmarshalJSON
// @Description:
// @receiver b
// @param data
// @return error
func (b *BitBool) UnmarshalJSON(data []byte) error {
	bv, err := strconv.ParseBool(string(data))
	if err != nil {
		return err
	}
	*b = BitBool(bv)
	return nil
}
