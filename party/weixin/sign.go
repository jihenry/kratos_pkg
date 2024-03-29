package weixin

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	Http "net/http"
	"net/url"
	"sort"
	"strings"
)

//获取约惠圈的签名
func genYHQOpenSign2(data []byte, saltKey string) (ans string) {
	maps := make(map[string]interface{})
	//参数解析错误
	if err := json.Unmarshal(data, &maps); err != nil {
		return
	}
	ans = genYHQOpenSign(maps, saltKey)
	return
}

//获取约惠圈的签名
func genYHQOpenSign(maps map[string]interface{}, saltKey string) (ans string) {
	//获取data 中获取参数的key 进行 ascii码 从小到大排序
	if len(maps) <= 0 {
		return
	}
	keys := make([]string, 0, len(maps))
	for k := range maps {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var buf strings.Builder
	for _, key := range keys {
		if value, ok := maps[key]; ok {
			val := toString(value)
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(key)
			buf.WriteByte('=')
			buf.WriteString(val)
		}
	}
	buf.WriteString("&key=")
	buf.WriteString(saltKey)
	log.Printf("signString:%s", buf.String())
	hm := hmac.New(md5.New, []byte(saltKey))
	hm.Write([]byte(buf.String()))
	return hex.EncodeToString(hm.Sum([]byte("")))
}

func toString(value interface{}) (vv string) {
	switch value.(type) {
	case bool:
		vv = fmt.Sprintf("%t", value)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		vv = fmt.Sprintf("%d", value)
	case float64, float32:
		vv = fmt.Sprintf("%2.f", value)
	default:
		vv = fmt.Sprintf("%s", value)
	}
	vv = strings.Replace(vv, " ", "", -1)
	return
}

// get wxapi http get wrapper, parse to proto message
func get(URL string, data map[string]interface{}, reply interface{}) error {
	req, _ := Http.NewRequest("GET", URL, nil)

	params := url.Values{}
	for k, v := range data {
		// v should be either int or string, this is safe
		params.Add(k, fmt.Sprintf("%v", v))
	}
	req.URL.RawQuery = params.Encode()

	client := &Http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	jsonStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("URL:", URL, "data:", data, "reply:", string(jsonStr))

	err = json.Unmarshal(jsonStr, reply)
	if err != nil {
		log.Println("URL:", URL, "data:", data, "reply:", string(jsonStr), "byte:", jsonStr)
		return err
	}
	log.Println("URL:", URL, "request:", data, "reply:", string(jsonStr), "reply proto:", reply)

	// TODO more elegantly check.
	wxErr := &WXErr{}
	json.Unmarshal(jsonStr, wxErr)
	if wxErr.Errcode != 0 {
		return fmt.Errorf("%s", wxErr.Errmsg)
	}

	return nil
}
