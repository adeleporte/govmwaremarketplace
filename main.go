package govmwaremarketplace

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
)

func GetAccessToken(myKey string) (string, error) {
	client := &http.Client{}

	res, err := client.Post(fmt.Sprintf("https://console.cloud.vmware.com/csp/gateway/am/api/auth/api-tokens/authorize?refresh_token=%s", myKey), "application/json", nil)

	token := GetAccessTokenResults{}

	body, err := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &token); err != nil {
		return "", err
	}

	if token.AccessToken == "" {
		return "", errors.New("Unable to get access token")
	}

	return token.AccessToken, err
}

func NewClient(token string) (MyClient, error) {

	myclient := MyClient{
		Client: &http.Client{},
	}

	tk, err := GetAccessToken(token)
	if err != nil {
		return myclient, err
	}

	myclient.Token = tk

	return myclient, nil
}

func GetProducts(client MyClient) (GetProductsResponse, error) {

	pl := GetProductsResponse{}

	r, err := http.NewRequest("GET", "https://api.marketplace.cloud.vmware.com/products?pagination={\"page\":1,\"pageSize\":1000}", nil)
	if err != nil {
		return pl, err
	}
	r.Header.Set("csp-auth-token", client.Token)

	res, err := client.Client.Do(r)
	if err != nil {
		return pl, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &pl); err != nil {
		return pl, err
	}

	return pl, nil

}

func GetProductDetail(client MyClient, slug string) (GetProductDetailResponse, error) {

	pl := GetProductDetailResponse{}

	r, err := http.NewRequest("GET", fmt.Sprintf("https://api.marketplace.cloud.vmware.com/products/%s?isSlug=true", slug), nil)
	if err != nil {
		return pl, err
	}
	r.Header.Set("csp-auth-token", client.Token)

	res, err := client.Client.Do(r)
	if err != nil {
		return pl, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &pl); err != nil {
		return pl, err
	}

	return pl, nil

}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

func GetDownload(client MyClient, download_body DownloadBody) (GetDownloadResults, error) {
	log.Println(download_body)

	dr := GetDownloadResults{}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(download_body)

	r, err := http.NewRequest("POST", fmt.Sprintf("https://api.marketplace.cloud.vmware.com/products/%s/download", download_body.ProductID), buf)
	if err != nil {
		return dr, err
	}
	r.Header.Set("csp-auth-token", client.Token)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")

	res, err := client.Client.Do(r)
	if err != nil {
		return dr, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &dr); err != nil {
		return dr, err
	}

	return dr, nil

}

func Download(client MyClient, filename string, url string) error {

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filename)
	if err != nil {
		return err
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")

	// Close the file without defer so it can happen before Rename()
	out.Close()

	return nil

}
