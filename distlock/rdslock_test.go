package distlock

import "testing"

func TestRdsLock(t *testing.T) {
	result := AquireLock("accoutId", 15)
	t.Log(result)
}
