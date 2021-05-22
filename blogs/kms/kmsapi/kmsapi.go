package kmsapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	ds "go-samples/kms/datastruct"

	cloudkms "cloud.google.com/go/kms/apiv1"
	log "github.com/sirupsen/logrus"

	"google.golang.org/api/option"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

var (
	kmsProjectID   string
	kmsLocation    string
	credentialPath string
	crendialBytes  []byte
)

type createKMSInfo struct {
	context   context.Context
	kmsClient *cloudkms.KeyManagementClient
}

// Init initialize KMS required params
func Init() error {
	var err error
	kmsProjectID = os.Getenv("KMS_PROJECT_ID")
	if kmsProjectID == "" {
		return fmt.Errorf("missing KMS project id env")
	}
	//log.Debugf("kmsapi init kmsProjectID: %+v\n", kmsProjectID)

	kmsLocation = os.Getenv("KMS_LOCATION")
	if kmsLocation == "" {
		return fmt.Errorf("missing KMS location env")
	}
	//log.Debugf("kmsapi init kmsLocation: %+v\n", kmsLocation)

	credentialPath = os.Getenv("KMS_CREDENTIAL_PATH_FILENAME")
	if credentialPath == "" {
		return fmt.Errorf("missing credential path env")
	}
	//log.Debugf("kmsapi init credentialPath: %+v\n", credentialPath)

	crendialBytes, err = ioutil.ReadFile(credentialPath)
	if err != nil {
		return fmt.Errorf("unable to read project service credential file: %v", err)
	}

	log.Debugln("kmsapi Init success")
	return nil
}

// KeyRingAndCryptoKey create KMS key ring and key
func KeyRingAndCryptoKey(wg *sync.WaitGroup, kmsInfo string, rkms ds.RequestKMS, outstr *string) {
	defer wg.Done()

	// create KMS Client
	ctx := context.Background()
	kmsclient, err := cloudkms.NewKeyManagementClient(ctx, option.WithCredentialsJSON(crendialBytes))
	if err != nil {
		log.Errorf("Error: KeyRingAndCryptoKey failed cloudkms.NewKeyManagementClient: %v\n", err)
		*outstr = getJSONRespStr(kmsInfo, "Falied", "Error: KeyRingAndCryptoKey failed cloudkms.NewKeyManagementClient")
		return
	}

	// create kms information
	ckmsInfo := createKMSInfo{}
	ckmsInfo.context = ctx
	ckmsInfo.kmsClient = kmsclient

	// Create key ring
	krStr, err := createKeyRing(ckmsInfo, rkms)
	if err != nil {
		log.Errorf("Error: createKeyRing failed: %+v\n", err)
		*outstr = getJSONRespStr(kmsInfo, "Falied", "Error: create key ring or crypto key failed")
		return
	}

	// Create Crypto Key
	ckStr, err := createCryptoKey(ckmsInfo, rkms)
	if err != nil {
		log.Errorf("Error: createCryptoKey failed: %+v\n", err)
		*outstr = getJSONRespStr(kmsInfo, "Falied", "Error: create key ring or crypto key failed")
		return
	}
	log.Debugf("Success: created keyRing: %s and CryptoKey: %s\n", krStr, ckStr)
	*outstr = getJSONRespStr(kmsInfo, "Success", ckStr)
	return
}

// createKeyRing creates a new ring to store keys on KMS.
func createKeyRing(ckmsInfo createKMSInfo, rkms ds.RequestKMS) (string, error) {
	parent := fmt.Sprintf("projects/%s/locations/%s", kmsProjectID, kmsLocation)

	keyring := kmspb.KeyRing{}
	// Build the request.
	req := &kmspb.CreateKeyRingRequest{
		Parent:    parent,
		KeyRingId: rkms.KeyRingID,
		KeyRing:   &keyring,
	}
	log.Debugf("createKeyRing request params: %+v\n", req)

	// Call the API.
	result, err := ckmsInfo.kmsClient.CreateKeyRing(ckmsInfo.context, req)
	if err != nil {
		return "", fmt.Errorf("Error: Failed createKeyRing: %v", err)
	}
	log.Debugf("Successfully created key ring :%s", result)
	return "Successfully created key ring", nil
}

type dataset struct {
	KeyVersion string `json:"KeyVersion"`
}

type datasetVal struct {
	DataSet dataset `json:"DataSet"`
}

// createCryptoKey creates a new symmetric encrypt/decrypt key on KMS.
func createCryptoKey(ckmsInfo createKMSInfo, rkms ds.RequestKMS) (string, error) {
	// keyRingName := "projects/PROJECT_ID/locations/global/keyRings/RING_ID"
	// keyID := "key-" + strconv.Itoa(int(time.Now().Unix()))
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", kmsProjectID, kmsLocation, rkms.KeyRingID)

	// Build the request.
	req := &kmspb.CreateCryptoKeyRequest{
		Parent:      parent,
		CryptoKeyId: rkms.KeyID,
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_ENCRYPT_DECRYPT,
			VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
				Algorithm: kmspb.CryptoKeyVersion_GOOGLE_SYMMETRIC_ENCRYPTION,
			},
			Labels: rkms.Labels,
		},
	}
	log.Debugf("createCryptoKey request param: %+v\n", req)

	// Call the API.
	result, err := ckmsInfo.kmsClient.CreateCryptoKey(ckmsInfo.context, req)
	if err != nil {
		return "", fmt.Errorf("Error: Falied createCryptoKey: %v", err)
	}
	log.Debugf("Successfully created crypto key: %s\n", result)
	respStr := strings.Split(result.Primary.Name, "/")

	outval := &datasetVal{}
	outval.DataSet.KeyVersion = fmt.Sprintf("%s/%s/%s/%s", respStr[6], respStr[7], respStr[8], respStr[9])
	tmpStr, _ := json.Marshal(outval)
	return string(tmpStr), nil
}

type encResp struct {
	DataSet dataset           `json:"DataSet"`
	Values  map[string]string `json:"Values"`
}

// Encrypt encrypt given key
// name := "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID"
// plaintext := []byte("Sample message")
func Encrypt(wg *sync.WaitGroup, kmsInfo string, rkms ds.RequestKMS, outstr *string) {
	var err error

	defer wg.Done()
	// create KMS Client
	ctx := context.Background()
	kmsClient, err := cloudkms.NewKeyManagementClient(ctx, option.WithCredentialsJSON(crendialBytes))
	if err != nil {
		log.Errorf("Error: encrypt failed cloudkms.NewKeyManagementClient: %v\n", err)
		*outstr = getJSONRespStr(kmsInfo, "Falied", "Error: encrypt failed")
		return
	}

	// The resource name of the [CryptoKey][google.cloud.kms.v1.CryptoKey]
	name := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", kmsProjectID, kmsLocation, rkms.KeyRingID, rkms.KeyID)

	// response result
	var respName string
	result := &encResp{}
	result.Values = make(map[string]string)

	// Loop for dataset
	for key, value := range rkms.DataSet {
		// Build the request.
		req := &kmspb.EncryptRequest{
			Name:      name,
			Plaintext: []byte(value),
		}
		log.Debugf("encrpyt request params: %+v\n", req)

		// Call the API.
		resp, err := kmsClient.Encrypt(ctx, req)
		if err != nil {
			result.Values[key] = "Failed to encrypt"
			log.Debugf("Error: encrypt failed for key(%s) value(%s) :%s\n", key, value, err.Error())
			continue
		}
		respName = resp.Name
		strB64 := base64.URLEncoding.EncodeToString(resp.Ciphertext)
		rs := fmt.Sprintf("%s", strB64)
		result.Values[key] = fmt.Sprintf("%s", rs)
	}

	var jRes []byte
	nSp := strings.Split(respName, "/")
	result.DataSet.KeyVersion = fmt.Sprintf("%s/%s/%s/%s", nSp[6], nSp[7], nSp[8], nSp[9])
	jRes, err = json.Marshal(result)
	if err != nil {
		log.Infoln("ksmapi encrypt marshal error :", err.Error())
		jRes = []byte("")
	} else {
		log.Infoln("kms encrypt success")
	}
	*outstr = getJSONRespStr(kmsInfo, "Success", string(jRes))

	//clean up
	for k := range result.Values {
		delete(result.Values, k)
	}
	return
}

type decResp struct {
	DataSet map[string]string `json:"DataSet"`
}

// Decrypt decrypt given key
// name := "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID"
func Decrypt(wg *sync.WaitGroup, kmsInfo string, rkms ds.RequestKMS, outstr *string) {
	defer wg.Done()

	// create KMS Client
	ctx := context.Background()
	kmsClient, err := cloudkms.NewKeyManagementClient(ctx, option.WithCredentialsJSON(crendialBytes))
	if err != nil {
		log.Debugf(" Error: decrypt failed cloudkms.NewKeyManagementClient: %v", err)
		*outstr = getJSONRespStr(kmsInfo, "Falied", "Error: decrypt failed")
		return
	}

	// The resource name of the [CryptoKey][google.cloud.kms.v1.CryptoKey]
	name := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", kmsProjectID, kmsLocation, rkms.KeyRingID, rkms.KeyID)

	// response result
	result := &decResp{}
	result.DataSet = make(map[string]string)

	// Loop for dataset
	for key, value := range rkms.DataSet {
		// Build the request.
		ct, _ := base64.URLEncoding.DecodeString(value)
		req := &kmspb.DecryptRequest{
			Name:       name,
			Ciphertext: []byte(ct),
		}
		log.Debugf("kms request for decrpyt key: %+v\n", req)

		// Call the API.
		resp, err := kmsClient.Decrypt(ctx, req)
		if err != nil {
			result.DataSet[key] = "Failed to decrypt"
			log.Debugf("Error: decrypt failed for key(%s) value(%s) :%s\n", key, value, err.Error())
			continue
		}
		result.DataSet[key] = string(resp.Plaintext)
	}

	var jRes []byte
	jRes, err = json.Marshal(result)
	if err != nil {
		log.Infoln("ksmapi decrypt marshal error :", err.Error())
		jRes = []byte("")
	} else {
		log.Infoln("kms decrypt success")
	}
	*outstr = getJSONRespStr(kmsInfo, "Success", string(jRes))

	//clean up
	for k := range result.DataSet {
		delete(result.DataSet, k)
	}
	return
}

// KeyRotation new key version
func KeyRotation(wg *sync.WaitGroup, kmsInfo string, rkms ds.RequestKMS, outstr *string) {
	defer wg.Done()

	// create KMS Client
	ctx := context.Background()
	kmsClient, err := cloudkms.NewKeyManagementClient(ctx, option.WithCredentialsJSON(crendialBytes))
	if err != nil {
		log.Errorf("Error: cryptoKey falied cloudkms.NewKeyManagementClient: %v\n", err)
		*outstr = getJSONRespStr(kmsInfo, "Failed", "Error: KeyRotation failed")
		return
	}

	// The resource name of the [CryptoKey][google.cloud.kms.v1.CryptoKey]
	name := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", kmsProjectID, kmsLocation, rkms.KeyRingID, rkms.KeyID)

	// crypto key version request
	ckvReq := &kmspb.CreateCryptoKeyVersionRequest{
		Parent: name,
	}

	ckvResp, err := kmsClient.CreateCryptoKeyVersion(ctx, ckvReq)
	if err != nil {
		log.Errorf("Error: KeyRotation falied CreateCryptoKeyVersion: %v\n", err)
		*outstr = getJSONRespStr(kmsInfo, "Failed", "Error: key rotation failed")
		return
	}

	log.Debugf("key rotation CreateCryptoKeyVersion response: %+v\n", ckvResp)

	sResp := strings.Split(ckvResp.Name, "/")
	if len(sResp) < 10 {
		log.Errorf("Error: KeyRotation falied, new key version not found: %v\n", err)
		*outstr = getJSONRespStr(kmsInfo, "Failed", "Error: key rotation failed")
		return
	}

	// Build update key version request.
	ukvReq := &kmspb.UpdateCryptoKeyPrimaryVersionRequest{
		Name:               name,
		CryptoKeyVersionId: sResp[9],
	}

	// Call UpdateCryptoKeyPrimaryVersion API.
	ukvResp, err := kmsClient.UpdateCryptoKeyPrimaryVersion(ctx, ukvReq)
	if err != nil {
		log.Errorf("Error: KeyRotation falied UpdateCryptoKeyPrimaryVersion: %v\n", err)
		*outstr = getJSONRespStr(kmsInfo, "Failed", "Error: key rotation failed")
		return
	}

	log.Debugf("key rotation UpdateCryptoKeyPrimaryVersion response: %+v\n", ukvResp)
	respStr := strings.Split(ukvResp.Primary.Name, "/")

	outval := &datasetVal{}
	outval.DataSet.KeyVersion = fmt.Sprintf("%s/%s/%s/%s", respStr[6], respStr[7], respStr[8], respStr[9])
	tmpStr, _ := json.Marshal(outval)
	*outstr = getJSONRespStr(kmsInfo, "Success", string(tmpStr))
	return
}

type respData struct {
	ClientInfo string `json:"clientinfo"`
	Status     string `json:"status"`
	OutText    string `json:"outtext"`
}

func getJSONRespStr(kmsLogInfo string, stat string, txt string) string {
	tmpResp := respData{}
	tmpResp.ClientInfo = kmsLogInfo
	tmpResp.Status = stat
	tmpResp.OutText = txt

	respStr, err := json.Marshal(tmpResp)
	if err != nil {
		log.Infoln("ksmapi getJSONRespStr marshal error :", err.Error())
		return kmsLogInfo + "Status: failed"
	}
	log.Debugln(string(respStr))
	log.Infoln(fmt.Sprintf("%s\n%s\n", tmpResp.ClientInfo, tmpResp.Status))
	return fmt.Sprintf("%s\n%s\n", tmpResp.Status, tmpResp.OutText)
}
