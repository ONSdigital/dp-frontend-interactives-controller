# \PublicApi

All URIs are relative to *http://localhost:23600*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionCsvGet**](PublicApi.md#DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionCsvGet) | **Get** /downloads/datasets/{datasetID}/editions/{edition}/versions/{version}.csv | Download the full csv for a given datasetID, edition and version
[**DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionCsvMetadataJsonGet**](PublicApi.md#DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionCsvMetadataJsonGet) | **Get** /downloads/datasets/{datasetID}/editions/{edition}/versions/{version}.csv-metadata.json | Download the csv metadata for a given datasetID, edition and version
[**DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionXlsxGet**](PublicApi.md#DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionXlsxGet) | **Get** /downloads/datasets/{datasetID}/editions/{edition}/versions/{version}.xlsx | Download the full excel file for a given datasetID, edition and version
[**DownloadsFilterOutputsFilterOutputIDCsvGet**](PublicApi.md#DownloadsFilterOutputsFilterOutputIDCsvGet) | **Get** /downloads/filter-outputs/{filterOutputID}.csv | Download a filtered csv file for a given filter output id
[**DownloadsFilterOutputsFilterOutputIDXlsxGet**](PublicApi.md#DownloadsFilterOutputsFilterOutputIDXlsxGet) | **Get** /downloads/filter-outputs/{filterOutputID}.xlsx | Download a filtered excel file for a given filter output id
[**ImagesImageIDVariantFilenameGet**](PublicApi.md#ImagesImageIDVariantFilenameGet) | **Get** /images/{imageID}/{variant}/{filename} | Download an image variant



## DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionCsvGet

> *os.File DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionCsvGet(ctx, datasetID, edition, version).Execute()

Download the full csv for a given datasetID, edition and version



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    datasetID := "datasetID_example" // string | The unique identifier for a dataset.
    edition := "edition_example" // string | An edition of a dataset
    version := "version_example" // string | A version of a dataset

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PublicApi.DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionCsvGet(context.Background(), datasetID, edition, version).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PublicApi.DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionCsvGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionCsvGet`: *os.File
    fmt.Fprintf(os.Stdout, "Response from `PublicApi.DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionCsvGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**datasetID** | **string** | The unique identifier for a dataset. | 
**edition** | **string** | An edition of a dataset | 
**version** | **string** | A version of a dataset | 

### Other Parameters

Other parameters are passed through a pointer to a apiDownloadsDatasetsDatasetIDEditionsEditionVersionsVersionCsvGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




### Return type

[***os.File**](*os.File.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/csv

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionCsvMetadataJsonGet

> *os.File DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionCsvMetadataJsonGet(ctx, datasetID, edition, version).Execute()

Download the csv metadata for a given datasetID, edition and version



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    datasetID := "datasetID_example" // string | The unique identifier for a dataset.
    edition := "edition_example" // string | An edition of a dataset
    version := "version_example" // string | A version of a dataset

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PublicApi.DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionCsvMetadataJsonGet(context.Background(), datasetID, edition, version).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PublicApi.DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionCsvMetadataJsonGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionCsvMetadataJsonGet`: *os.File
    fmt.Fprintf(os.Stdout, "Response from `PublicApi.DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionCsvMetadataJsonGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**datasetID** | **string** | The unique identifier for a dataset. | 
**edition** | **string** | An edition of a dataset | 
**version** | **string** | A version of a dataset | 

### Other Parameters

Other parameters are passed through a pointer to a apiDownloadsDatasetsDatasetIDEditionsEditionVersionsVersionCsvMetadataJsonGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




### Return type

[***os.File**](*os.File.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/csvm+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionXlsxGet

> *os.File DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionXlsxGet(ctx, datasetID, edition, version).Execute()

Download the full excel file for a given datasetID, edition and version



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    datasetID := "datasetID_example" // string | The unique identifier for a dataset.
    edition := "edition_example" // string | An edition of a dataset
    version := "version_example" // string | A version of a dataset

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PublicApi.DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionXlsxGet(context.Background(), datasetID, edition, version).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PublicApi.DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionXlsxGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionXlsxGet`: *os.File
    fmt.Fprintf(os.Stdout, "Response from `PublicApi.DownloadsDatasetsDatasetIDEditionsEditionVersionsVersionXlsxGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**datasetID** | **string** | The unique identifier for a dataset. | 
**edition** | **string** | An edition of a dataset | 
**version** | **string** | A version of a dataset | 

### Other Parameters

Other parameters are passed through a pointer to a apiDownloadsDatasetsDatasetIDEditionsEditionVersionsVersionXlsxGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




### Return type

[***os.File**](*os.File.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/vnd.ms-excel

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DownloadsFilterOutputsFilterOutputIDCsvGet

> *os.File DownloadsFilterOutputsFilterOutputIDCsvGet(ctx, filterOutputID).Execute()

Download a filtered csv file for a given filter output id



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    filterOutputID := "filterOutputID_example" // string | The unique identifier for a filter output job

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PublicApi.DownloadsFilterOutputsFilterOutputIDCsvGet(context.Background(), filterOutputID).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PublicApi.DownloadsFilterOutputsFilterOutputIDCsvGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DownloadsFilterOutputsFilterOutputIDCsvGet`: *os.File
    fmt.Fprintf(os.Stdout, "Response from `PublicApi.DownloadsFilterOutputsFilterOutputIDCsvGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**filterOutputID** | **string** | The unique identifier for a filter output job | 

### Other Parameters

Other parameters are passed through a pointer to a apiDownloadsFilterOutputsFilterOutputIDCsvGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[***os.File**](*os.File.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/csv

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DownloadsFilterOutputsFilterOutputIDXlsxGet

> *os.File DownloadsFilterOutputsFilterOutputIDXlsxGet(ctx, filterOutputID).Execute()

Download a filtered excel file for a given filter output id



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    filterOutputID := "filterOutputID_example" // string | The unique identifier for a filter output job

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PublicApi.DownloadsFilterOutputsFilterOutputIDXlsxGet(context.Background(), filterOutputID).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PublicApi.DownloadsFilterOutputsFilterOutputIDXlsxGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DownloadsFilterOutputsFilterOutputIDXlsxGet`: *os.File
    fmt.Fprintf(os.Stdout, "Response from `PublicApi.DownloadsFilterOutputsFilterOutputIDXlsxGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**filterOutputID** | **string** | The unique identifier for a filter output job | 

### Other Parameters

Other parameters are passed through a pointer to a apiDownloadsFilterOutputsFilterOutputIDXlsxGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[***os.File**](*os.File.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/csv

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ImagesImageIDVariantFilenameGet

> *os.File ImagesImageIDVariantFilenameGet(ctx, imageID, variant, filename).Execute()

Download an image variant



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    imageID := "imageID_example" // string | The unique identifier for an image.
    variant := "variant_example" // string | The variant of an image to download
    filename := "filename_example" // string | The filename of the file to download

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PublicApi.ImagesImageIDVariantFilenameGet(context.Background(), imageID, variant, filename).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PublicApi.ImagesImageIDVariantFilenameGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ImagesImageIDVariantFilenameGet`: *os.File
    fmt.Fprintf(os.Stdout, "Response from `PublicApi.ImagesImageIDVariantFilenameGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**imageID** | **string** | The unique identifier for an image. | 
**variant** | **string** | The variant of an image to download | 
**filename** | **string** | The filename of the file to download | 

### Other Parameters

Other parameters are passed through a pointer to a apiImagesImageIDVariantFilenameGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




### Return type

[***os.File**](*os.File.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: image/png

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

