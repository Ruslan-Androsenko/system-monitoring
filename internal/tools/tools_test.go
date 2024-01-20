package tools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsingThreeNumbers(t *testing.T) {
	tests := []struct {
		name   string
		format string
		buffer []byte
	}{
		{
			name:   "test cpu load with coma number format",
			format: cpuLoadPatternFormat,
			buffer: []byte("%Cpu(s):  3,6 us,  4,8 sy,  0,0 ni, 92,7 id,  0,0 wa,  0,0 hi,  0,0 si,  0,0 st"),
		},
		{
			name:   "test cpu load with point number format",
			format: cpuLoadPatternFormat,
			buffer: []byte("%Cpu(s):  1.6 us,  0.0 sy,  0.0 ni, 98.4 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st"),
		},

		{
			name:   "test disk load with coma number format",
			format: diskLoadPatternFormat,
			buffer: []byte("sda    10,64    177,23    167,87    0,00    3979996    3769829    0"),
		},
		{
			name:   "test disk load with point number format",
			format: diskLoadPatternFormat,
			buffer: []byte("sda    251.82    5318.19    1085.91    64.95    893509    182444    10912"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.NotNil(t, tc.buffer)
			require.Greater(t, len(tc.buffer), 0)

			first, second, third, err := parsingThreeNumbers(tc.buffer, tc.format)
			require.NoError(t, err)

			require.GreaterOrEqual(t, first, zeroNumber)
			require.GreaterOrEqual(t, second, zeroNumber)
			require.GreaterOrEqual(t, third, zeroNumber)
		})
	}
}
