package govmwaremarketplace

import "net/http"

type MyClient struct {
	Client *http.Client
	Token  string
}

type GetAccessTokenResults struct {
	AccessToken string `json:"access_token"`
}

type GetProductsResponse struct {
	Response DataList `json:"response"`
}

type DataList struct {
	DataList []Product `json:"dataList"`
}

type Product struct {
	DisplayName string `json:"displayname"`
	Slug        string `json:"slug"`
}

type GetProductDetailResponse struct {
	Response Data `json:"response"`
}

type Data struct {
	Data ProductDetail `json:"data"`
}

type ProductDetail struct {
	ProductID                  string                  `json:"productid"`
	DisplayName                string                  `json:"displayname"`
	Slug                       string                  `json:"slug"`
	ProductDeploymentFilesList []ProductDeploymentFile `json:"productdeploymentfilesList"`
	ProductLogo                ProductLogo             `json:"productlogo"`
	Description                ProductDescription      `json:"description"`
	SolutionType               string                  `json:"solutiontype"`
}

type ProductDescription struct {
	Summary       string   `json:"summary"`
	ImageURLsList []string `json:"imageurlsList"`
	Description   string   `json:"description"`
}

type ProductLogo struct {
	URL string `json:"url"`
}

type ProductDeploymentFile struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	FileID string `json:"fileid"`
}

type DownloadBody struct {
	DeploymentFileID string `json:"deploymentFileId"`
	ProductID        string `json:"productId"`
}

type GetDownloadResults struct {
	Response GetDownloadResponse `json:"response"`
}

type GetDownloadResponse struct {
	PresignedURL string `json:"presignedurl"`
	Message      string `json:"message"`
	StatusCode   int    `json:"statuscode"`
}

type WriteCounter struct {
	Total uint64
}
