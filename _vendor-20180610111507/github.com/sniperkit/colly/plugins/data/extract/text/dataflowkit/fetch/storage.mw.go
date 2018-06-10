package fetch

import (
	"encoding/base32"
	"encoding/json"
	"errors"
	"time"

	"github.com/slotix/dataflowkit/splash"
	"github.com/slotix/dataflowkit/storage"
	"github.com/slotix/dataflowkit/utils"
	"github.com/spf13/viper"
)

type storageMiddleware struct {
	//storage instance puts fetching results to a cache
	storage storage.Store
	Service
}

// StorageMiddleware caches web pages to be parsed.
func StorageMiddleware(storage storage.Store) ServiceMiddleware {
	return func(next Service) Service {
		return storageMiddleware{storage, next}
	}
}

//get fetches web page content from a storage
func (mw storageMiddleware) get(req FetchRequester) (resp FetchResponser, err error) {
	var fetchResponse FetchResponser
	url := req.GetURL()

	switch req.(type) {
	case BaseFetcherRequest:
		fetchResponse = &BaseFetcherResponse{}
	case splash.Request:
		fetchResponse = &splash.Response{}
	default:
		panic("invalid fetcher request")
	}

	//URL Conversion MD5 Reduces file name length to avoid the error like file name too long.
	storageKey := string(utils.GenerateMD5([]byte(url)))
	//Base32 encoded values are 100% safe for file/uri usage without replacing any characters and guarantees 1-to-1 mapping
	sKey := base32.StdEncoding.EncodeToString([]byte(storageKey))
	value, err := mw.storage.Read(sKey)
	if err == nil {
		if err := json.Unmarshal(value, &fetchResponse); err != nil {
			logger.Error(err)
		}
		//Error responses: a 404 (Not Found) may be cached.
		//if sResponse.Response.Status == 404 {
		//	return nil, &errs.NotFound{URL: url}
		//}
		//check if item is expired.
		expired := fetchResponse.GetExpires()
		untilExpired := -time.Since(expired)
		logger.Infof("%s: Time until expired: %+v", url, untilExpired)
		//If a website is not cachable by some reason, ignore this and use cached copy if any
		ignoreCacheInfo := viper.GetBool("IGNORE_CACHE_INFO")
		if untilExpired > 0 || ignoreCacheInfo { //if cached value is not expired return it
			return fetchResponse, nil
		}
		err = errors.New("Cached item is expired or not cacheable")
	}
	return nil, err
}

//put saves web page content to the storage
func (mw storageMiddleware) put(req FetchRequester, resp FetchResponser) error {
	url := req.GetURL()
	//URL Conversion MD5 Reduces file name length to avoid the error like file name too long.
	storageKey := string(utils.GenerateMD5([]byte(url)))
	sKey := base32.StdEncoding.EncodeToString([]byte(storageKey))

	reasons := resp.GetReasonsNotToCache()
	if len(reasons) == 0 {
		logger.Info(url, " is Cachable")
	} else {
		logger.Info(url, " is not Cachable.", "Reasons to not cache:", resp.GetReasonsNotToCache())
	}
	expired := resp.GetExpires()

	//logger.Info("Expires:", expired)

	r, err := json.Marshal(resp)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	//calculate expiration time. This is actual for Redis only.
	expTime := expired.Unix()
	err = mw.storage.Write(sKey, r, expTime)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

//Fetch returns content either from storage or directly from web.
// func (mw storageMiddleware) Fetch(req FetchRequester) (io.ReadCloser, error) {
// 	resp, err := mw.Response(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return resp.GetHTML()
// }

//Response returns Fetch Response either from storage or directly from web.
//This middleware method is used by Parse service.
func (mw storageMiddleware) Response(req FetchRequester) (FetchResponser, error) {

	//if form Data is emtpy try to get cached data from storage first
	if req.GetFormData() != "" {
		//if form Data exists don't use storage!!!
		return mw.Service.Response(req)
	}
	//loads content from a storage if any
	fromStorage, err := mw.get(req)
	if err == nil {
		return fromStorage, nil
	}
	//logger.Error(err)
	//Get fetch response directly from web if there is nothing in storage
	resp, err := mw.Service.Response(req)
	if err != nil {
		return nil, err
	}

	var fetchResponse FetchResponser
	switch req.Type() {
	case "base":
		fetchResponse = resp.(*BaseFetcherResponse)
	case "splash":
		fetchResponse = resp.(*splash.Response)
	}
	//save fetched content to storage
	err = mw.put(req, fetchResponse)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
