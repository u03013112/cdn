package cdn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	// "log"
	"net/http"
)

func isPathExsit(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func jsonReader(reqJson interface{}) io.Reader {
	bytesData, err := json.Marshal(reqJson)
	if err != nil {
		fmt.Println(err.Error())
		log.Errorf(err.Error())
		return nil
	}
	reader := bytes.NewReader(bytesData)
	return reader
}
func JsonPost(method string, url string, respSt interface{}, reqJson interface{}) error {

	request, err := http.NewRequest(method, url, jsonReader(reqJson))
	if err != nil {
		// log.Errorf(err.Error())
		log.Fatalf("%s", err.Error())
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalf("%s", err.Error())
		// fmt.Println(err.Error())
		// log.Errorf("url: %s, err: %s", url, err.Error())
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("%d", resp.StatusCode)
		log.Fatalf("%s", err.Error())
		return err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("%s", err.Error())
		// fmt.Println(err.Error())
		// log.Errorf(err.Error())
		return err
	}

	// fmt.Printf("[%d][%s]\n", resp.StatusCode, respBytes)
	// fmt.Printf("%s\n", string(respBytes))
	if err := json.Unmarshal(respBytes, &respSt); err != nil {
		// fmt.Printf("%s\n", string(respBytes))
		// fmt.Printf("respBytes: %d, [%s],err:%s\n", len(respBytes), string(respBytes), err.Error())
		// return err
		//命令行post出去的就暂时不解json了
	}
	return nil
}
