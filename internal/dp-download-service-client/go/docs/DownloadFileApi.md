# \DownloadFileApi

All URIs are relative to *http://localhost:23600*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DownloadsNewFilepathGet**](DownloadFileApi.md#DownloadsNewFilepathGet) | **Get** /downloads-new/{filepath} | Download a file



## DownloadsNewFilepathGet

> *os.File DownloadsNewFilepathGet(ctx, filepath).Execute()

Download a file



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
    filepath := "filepath_example" // string | filepath of required file

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DownloadFileApi.DownloadsNewFilepathGet(context.Background(), filepath).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DownloadFileApi.DownloadsNewFilepathGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DownloadsNewFilepathGet`: *os.File
    fmt.Fprintf(os.Stdout, "Response from `DownloadFileApi.DownloadsNewFilepathGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**filepath** | **string** | filepath of required file | 

### Other Parameters

Other parameters are passed through a pointer to a apiDownloadsNewFilepathGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[***os.File**](*os.File.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

