package cos

import (
	"context"
	"testing"
)

func TestGetTempSecret(t *testing.T) {
	client, err := NewCosApiClient(WithSecretId("AKID24NuIVRUAhFqCUqR6WillaAXszsVs67v"), WithSecretKey("jYGlAh5GCsSr40h8ICWjgN4oKgmpCKVU"))
	if err != nil {
		t.Errorf("NewCosApiClient err:%s", err)
		return
	}
	result, err := client.GetTempSecret(context.Background(), "ap-guangzhou", "gz-gaas-dev-1300767139", "/minigameWeb/nbBankTownHead/")
	if err != nil {
		t.Errorf("GetTempSecret err:%s", err)
		return
	} else if result.Error != nil {
		t.Errorf("GetTempSecret result err:%s", result.Error)
	}
}
