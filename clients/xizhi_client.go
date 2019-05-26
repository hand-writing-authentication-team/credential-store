package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/hand-writing-authentication-team/credential-store/models"
	"github.com/labstack/gommon/log"
)

const (
	REL_ANALYZE_URL  = "analyze"
	REL_VALIDATE_URL = "validate"

	TYPE_JSON = "application/json"
)

type XizhiClient struct {
	Client *http.Client
	Url    *url.URL
}

func NewXizhiClient(xizhiDest string, timeout time.Duration) (xzClient *XizhiClient, err error) {
	httpClient := &http.Client{
		Timeout: timeout,
	}
	xizhiUrl, err := url.Parse(xizhiDest)
	if err != nil {
		log.Errorf("met error %s when parse url", err)
		return nil, err
	}

	xzClient.Client = httpClient
	xzClient.Url = xizhiUrl

	return xzClient, err

}

func (xz *XizhiClient) Analyze(handwriting string, prevFeature []models.Feature) (features []models.Feature, err error) {
	relUrl, _ := url.Parse(REL_ANALYZE_URL)
	fullUrl := xz.Url.ResolveReference(relUrl)
	if err != nil {
		log.Errorf("error when forming analyze url")
		return nil, err
	}
	reqBody := models.FeatureReq{
		UserHandwriting: handwriting,
		PrevDataPoints:  prevFeature,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	resp, err := xz.Client.Post(fullUrl.String(), TYPE_JSON, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		log.Error("error %s occured when sending req to xizhi", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("xizhi returned non 200 resp [%v]", resp.StatusCode)
		return nil, errors.New(fmt.Sprintf("status code [%v]", resp.StatusCode))
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("error happen when reading response from xizhi")
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, &features)
	if err != nil {
		log.Error("error happen when unmarshal response body to features")
		return nil, err
	}

	log.Info("Successfully analyzed features")
	return features, nil
}

func (xz *XizhiClient) Validate(handwriting string, prevFeature []models.Feature) (status bool, err error) {
	relUrl, _ := url.Parse(REL_VALIDATE_URL)
	fullUrl := xz.Url.ResolveReference(relUrl)
	if err != nil {
		log.Errorf("error when forming validate url")
		return false, err
	}
	reqBody := models.FeatureReq{
		UserHandwriting: handwriting,
		PrevDataPoints:  prevFeature,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	resp, err := xz.Client.Post(fullUrl.String(), TYPE_JSON, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		log.Error("error %s occured when sending req to xizhi", err)
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("xizhi returned non 200 resp [%v]", resp.StatusCode)
		return false, errors.New(fmt.Sprintf("status code [%v]", resp.StatusCode))
	}

	var statusObj models.ValidateStatus

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("error happen when reading response from xizhi")
		return false, err
	}

	err = json.Unmarshal(bodyBytes, &statusObj)
	if err != nil {
		log.Error("error happen when unmarshal response body to features")
		return false, err
	}

	log.Info("Successfully analyzed features")
	return statusObj.Status, nil
}
