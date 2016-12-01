package kong

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Kong struct {
	Endpoint string
}

const (
	CreateConsumerPath string = "/consumers/"
	DeleteConsumerPath string = "/consumers/"
	CreateAPIKeyPath string = "/consumers/"
)

// NewKong constructs a new Kong client
// endpoint: API endpoint of your Kong server
func NewKong(endpoint string) *Kong {
	ep := strings.TrimSpace(endpoint)
	ep = strings.TrimSuffix(ep, "/")
	return &Kong{Endpoint: ep}
}

// CreateConsumer creates a new API consumer
// username: (semi-option) The username of the consumer. You must send either this field or custom_id with the request.
// custom_id: (semi-option) Field for storing an existing ID for the consumer, useful for mapping Kong with users in your existing database. 
// You must send either username or custom_id with the request.
func (k *Kong) CreateConsumer(username, custom_id string) (map[string]interface{}, error) {
	var data = url.Values{}
	if username == "" && custom_id == "" {
		return nil, ErrMissingParameter
	}
	if username != "" {
		data.Add("username", username)
	}
	if custom_id != "" {
		data.Add("custom_id", custom_id)
	}

	resp, err := k.post(CreateConsumerPath, data)
	if err != nil {
		fmt.Println(err)
		return nil, ErrCreateConsumerFailed
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		fmt.Println(err)
        return nil, ErrCreateConsumerFailed
	}

	return result, nil
}

// DeleteConsumer deletes an existing API consumer
// identifier: username or id, The unique identifier or the name of the consumer to delete
func (k *Kong) DeleteConsumer(identifier string) error {
	_, err := k.delete(fmt.Sprintf("%s%s", DeleteConsumerPath, identifier))
	if err != nil {
		fmt.Println(err)
		return ErrDeleteConsumerFailed
	}
	return nil
}

// CreateAPIKey creates an API key for a customr
// consumer: The id or username property of the Consumer entity to associate the credentials to.
func (k *Kong) CreateAPIKey(consumer string) (map[string]interface{}, error) {
	path := fmt.Sprintf("%s%s/key-auth", CreateAPIKeyPath, consumer)
	fmt.Println(path)
	var data = url.Values{}
	resp, err := k.post(path, data)
	if err != nil {
		return nil, ErrCreateAPIKeyFailed
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, ErrCreateAPIKeyFailed
	}

	return result, nil
}

func (k *Kong) get(path string) ([]byte, error) {
	return k.request("GET", path, nil)
}

func (k *Kong) post(path string, data url.Values) ([]byte, error) {
	return k.request("POST", path, data)
}

func (k *Kong) delete(path string) ([]byte, error) {
	return k.request("DELETE", path, nil)
}

func (k *Kong) request(method string, path string, data url.Values) ([]byte, error) {
	var resp *http.Response
	var err error
	switch method {
		case "GET":
			resp, err = http.Get(fmt.Sprintf("%s%s", k.Endpoint, path))
		case "POST":
			resp, err = http.PostForm(fmt.Sprintf("%s%s", k.Endpoint, path), data)
		case "DELETE":
			req, err1 := http.NewRequest("DELETE", fmt.Sprintf("%s%s", k.Endpoint, path), nil)
			if err1 != nil {
				return nil, err1
			}
			resp, err = http.DefaultClient.Do(req)
		default:
			err = ErrUnknowMedthod
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return d, nil
}