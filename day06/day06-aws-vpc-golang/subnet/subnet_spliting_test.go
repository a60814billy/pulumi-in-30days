package subnet

import "testing"

func TestSplit(t *testing.T) {
	tests := []struct {
		cidr                string
		newPrefix           int
		expectSubnetsLength int
		expectFirstSubnets  []string
	}{
		{
			cidr:                "10.120.0.0/16",
			newPrefix:           22,
			expectSubnetsLength: 64,
			expectFirstSubnets:  []string{"10.120.0.0/22", "10.120.4.0/22", "10.120.8.0/22", "10.120.12.0/22", "10.120.16.0/22"},
		},
		{
			cidr:                "10.120.0.0/16",
			newPrefix:           24,
			expectSubnetsLength: 256,
			expectFirstSubnets:  []string{"10.120.0.0/24", "10.120.1.0/24", "10.120.2.0/24", "10.120.3.0/24"},
		},
		{
			cidr:                "192.168.100.0/24",
			newPrefix:           25,
			expectSubnetsLength: 2,
			expectFirstSubnets:  []string{"192.168.100.0/25", "192.168.100.128/25"},
		},
	}

	for _, tc := range tests {
		newSubnets := Split(tc.cidr, tc.newPrefix)
		if len(newSubnets) != tc.expectSubnetsLength {
			t.Errorf("Expected %d subnets, got %d", tc.expectSubnetsLength, len(newSubnets))
		}
		for i, subnet := range tc.expectFirstSubnets {
			if newSubnets[i] != subnet {
				t.Errorf("Expected %s, got %s", subnet, newSubnets[i])
			}
		}
	}
}
