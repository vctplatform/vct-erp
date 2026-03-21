package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

const moneyScale = 4

// Money stores a fixed-point value scaled to 4 decimal places to mirror NUMERIC(20,4).
type Money struct {
	scaled big.Int
}

// ParseMoney converts a decimal string into a fixed-point Money value.
func ParseMoney(input string) (Money, error) {
	value := strings.TrimSpace(input)
	if value == "" {
		return Money{}, errors.New("money value is required")
	}

	sign := 1
	switch value[0] {
	case '-':
		sign = -1
		value = value[1:]
	case '+':
		value = value[1:]
	}

	parts := strings.Split(value, ".")
	if len(parts) > 2 {
		return Money{}, fmt.Errorf("invalid money value %q", input)
	}

	wholePart := parts[0]
	if wholePart == "" {
		wholePart = "0"
	}

	if !isDigitsOnly(wholePart) {
		return Money{}, fmt.Errorf("invalid money value %q", input)
	}

	fractionalPart := ""
	if len(parts) == 2 {
		fractionalPart = parts[1]
		if fractionalPart == "" {
			fractionalPart = "0"
		}
		if len(fractionalPart) > moneyScale || !isDigitsOnly(fractionalPart) {
			return Money{}, fmt.Errorf("invalid money value %q", input)
		}
	}

	fractionalPart += strings.Repeat("0", moneyScale-len(fractionalPart))
	combined := strings.TrimLeft(wholePart+fractionalPart, "0")
	if combined == "" {
		combined = "0"
	}

	var scaled big.Int
	if _, ok := scaled.SetString(combined, 10); !ok {
		return Money{}, fmt.Errorf("invalid money value %q", input)
	}
	if sign < 0 {
		scaled.Neg(&scaled)
	}

	return Money{scaled: scaled}, nil
}

// MustParseMoney parses a money literal and panics on error.
func MustParseMoney(input string) Money {
	money, err := ParseMoney(input)
	if err != nil {
		panic(err)
	}
	return money
}

// Add returns the arithmetic sum.
func (m Money) Add(other Money) Money {
	var sum big.Int
	sum.Add(&m.scaled, &other.scaled)
	return Money{scaled: sum}
}

// Sub returns the arithmetic difference.
func (m Money) Sub(other Money) Money {
	var difference big.Int
	difference.Sub(&m.scaled, &other.scaled)
	return Money{scaled: difference}
}

// Equal compares two money values.
func (m Money) Equal(other Money) bool {
	return m.scaled.Cmp(&other.scaled) == 0
}

// Cmp compares two money values.
func (m Money) Cmp(other Money) int {
	return m.scaled.Cmp(&other.scaled)
}

// Sign reports whether the value is negative, zero, or positive.
func (m Money) Sign() int {
	return m.scaled.Sign()
}

// IsZero returns true when the amount is zero.
func (m Money) IsZero() bool {
	return m.Sign() == 0
}

// IsPositive returns true when the amount is greater than zero.
func (m Money) IsPositive() bool {
	return m.Sign() > 0
}

// String formats the fixed-point value.
func (m Money) String() string {
	var absolute big.Int
	absolute.Abs(&m.scaled)

	digits := absolute.String()
	if len(digits) <= moneyScale {
		digits = strings.Repeat("0", moneyScale-len(digits)+1) + digits
	}

	splitAt := len(digits) - moneyScale
	whole := digits[:splitAt]
	fractional := digits[splitAt:]
	if whole == "" {
		whole = "0"
	}

	prefix := ""
	if m.Sign() < 0 {
		prefix = "-"
	}

	return prefix + whole + "." + fractional
}

// MarshalJSON keeps monetary values exact by serializing them as strings.
func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

// UnmarshalJSON accepts a quoted decimal string or a raw JSON number.
func (m *Money) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return errors.New("money payload is empty")
	}

	if string(data) == "null" {
		*m = Money{}
		return nil
	}

	var literal string
	if data[0] == '"' {
		if err := json.Unmarshal(data, &literal); err != nil {
			return fmt.Errorf("decode money: %w", err)
		}
	} else {
		literal = string(data)
	}

	parsed, err := ParseMoney(literal)
	if err != nil {
		return err
	}

	*m = parsed
	return nil
}

func isDigitsOnly(value string) bool {
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// GoString keeps debugging output readable in tests.
func (m Money) GoString() string {
	return strconv.Quote(m.String())
}

// Split divides a money amount into n parts, keeping NUMERIC(20,4) precision and distributing any remainder from the start.
func (m Money) Split(parts int) ([]Money, error) {
	if parts <= 0 {
		return nil, errors.New("split parts must be greater than zero")
	}

	divisor := big.NewInt(int64(parts))
	var quotient big.Int
	var remainder big.Int
	quotient.QuoRem(&m.scaled, divisor, &remainder)

	remainderCount := int(remainder.Int64())
	if remainderCount < 0 {
		remainderCount = -remainderCount
	}

	partsSlice := make([]Money, parts)
	for index := 0; index < parts; index++ {
		part := cloneBigInt(&quotient)
		if remainderCount > 0 && index < remainderCount {
			if m.Sign() >= 0 {
				part.Add(part, big.NewInt(1))
			} else {
				part.Sub(part, big.NewInt(1))
			}
		}
		partsSlice[index] = newMoneyFromBigInt(part)
	}

	return partsSlice, nil
}

func newMoneyFromBigInt(source *big.Int) Money {
	var scaled big.Int
	scaled.Set(source)
	return Money{scaled: scaled}
}

func cloneBigInt(source *big.Int) *big.Int {
	var cloned big.Int
	cloned.Set(source)
	return &cloned
}
