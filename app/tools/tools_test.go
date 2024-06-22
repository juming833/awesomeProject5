package tools

import (
	"fmt"
	"github.com/spf13/viper"
	"testing"
)

func TestGetUID(t *testing.T) {
	appID := viper.GetString("AppID")
	alipayPublicKey := viper.GetString("AlipayPublickey")
	privateKey := viper.GetString("Privatekey")
	fmt.Println(appID, alipayPublicKey, privateKey)
}
