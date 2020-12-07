package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	ds "go-samples/kms/datastruct"
	"go-samples/kms/kmsapi"

	log "github.com/sirupsen/logrus"
)

const (
	maxCacheSize = 10000
)

var (
	// cache for validation
	cacheData = make(map[string]cacheInfo, maxCacheSize)

	// mutex lock for sync access to cache
	cacheMutex *sync.Mutex
)

type cacheInfo struct {
	ClientID   string
	Expiration int64
}

func init() {
	cacheMutex = &sync.Mutex{}
}

// Health API
func Health(w http.ResponseWriter, r *http.Request) {
	log.Infoln("Received KMS Health request")
	w.Write([]byte("Welcome! KMS Service Wrapper"))
	json.NewEncoder(w)
}

// RequestKMS API
func RequestKMS(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received KMS request")
	fmt.Printf("KMS raw request: %v\n", r)
	log.Infoln("Received KMS request")
	log.Debugf("KMS raw request: %v\n", r)

	// Read Request body
	reqBody := ds.RequestKMS{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		errStr := fmt.Sprintf("Error: in reading request bosdy : %+v", err)
		log.Errorln(errStr)
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}
	log.Debugf("KMS request body: %+v\n", reqBody)

	// Read Operation Param
	reqOpr := r.URL.Query().Get("opr")
	if reqOpr == "" {
		errStr := "Error: Url Param 'opr' is missing"
		log.Errorln(errStr)
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}
	log.Debugf("KMS request operation: %+v\n", reqOpr)

	// Check request Body and Params are valid
	if err := validateKMSRequestBody(reqBody, reqOpr); err != nil {
		errStr := err.Error()
		log.Errorln(errStr)
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}

	kmsLogInfo := fmt.Sprintf("App Log :: clientid:%s, operation:%s", reqBody.ClientID, reqOpr)
	log.Infof("%s, Received KMS request\n", kmsLogInfo)

	ct := time.Now().Unix()

	//Check token cache
	if ctk, ok := cacheData[reqBody.ClientID]; ok {
		// if client ID are same, check the expiration
		if strings.Trim(ctk.ClientID, "") != strings.Trim(reqBody.ClientID, "") || ctk.Expiration < ct {
			if strings.Trim(ctk.ClientID, "") != strings.Trim(reqBody.ClientID, "") {
				log.Debugln("received client id diff from cache client id")
			}
			if ctk.Expiration < ct {
				log.Debugln("token expired")
			}
			//Cleanup Cache entry
			log.Debugf("KMS cache cleanup for clientID: %s\n", ctk.ClientID)
			cacheMutex.Lock()
			delete(cacheData, ctk.ClientID)
			cacheMutex.Unlock()

			errStr := fmt.Sprintf("Error : Operation request client is not Authorizied or token is not valid")
			log.Errorln(errStr)
			http.Error(w, errStr, http.StatusUnauthorized)
			return
		}
	} else { // validate the token

		// token valid, add to cache
		cd := cacheInfo{}
		cd.ClientID = reqBody.ClientID
		cd.Expiration = 3 * 60 * 60 //3 hours

		//before update cache data
		cacheMutex.Lock()
		cacheData[reqBody.ClientID] = cd
		cacheMutex.Unlock()
	}

	var wg sync.WaitGroup
	var output string
	// Swtich based on operation
	switch reqOpr {
	case "createkeyringcrypto":
		log.Debugln("Operation:create key ring and crypto key")
		wg.Add(1)
		go kmsapi.KeyRingAndCryptoKey(&wg, kmsLogInfo, reqBody, &output)
		wg.Wait()
	case "encrypt":
		log.Debugln("Operation:encrypt")
		wg.Add(1)
		go kmsapi.Encrypt(&wg, kmsLogInfo, reqBody, &output)
		wg.Wait()
	case "decrypt":
		log.Debugln("Operation:decrypt")
		wg.Add(1)
		go kmsapi.Decrypt(&wg, kmsLogInfo, reqBody, &output)
		wg.Wait()
	case "keyringrotation":
		log.Debugln("Operation:key rotation")
		wg.Add(1)
		go kmsapi.KeyRotation(&wg, kmsLogInfo, reqBody, &output)
		wg.Wait()
	default:
		errStr := fmt.Sprintf("KMS request is not supported :%s\n", reqOpr)
		log.Errorln(errStr)
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}

	if strings.Contains(output, "Error") {
		log.Errorln(output)
		http.Error(w, output, http.StatusInternalServerError)
	} else {
		w.Write([]byte(output))
		json.NewEncoder(w).Encode(http.StatusOK)
	}
	return
}

// validate KMS Request Body
func validateKMSRequestBody(rd ds.RequestKMS, operation string) error {

	// generic params
	if rd.KeyRingID == "" || rd.ClientID == "" || rd.KeyID == "" {
		return fmt.Errorf("Error: KMS request, missing json body params")
	}

	switch operation {
	case "createkeyringcrypto":
		if len(rd.Labels) == 0 {
			return fmt.Errorf("Error: KMS request, missing labels body params")
		}
	case "encrypt", "decrypt":
		if len(rd.DataSet) == 0 {
			return fmt.Errorf("Error: KMS request, missing dataset body params")
		}
	case "keyringrotation":
	default:
		return fmt.Errorf("Error: request for operation is not supported :%s", operation)
	}
	return nil
}
