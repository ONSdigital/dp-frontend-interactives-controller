// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks_storage

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/download"
	"github.com/ONSdigital/dp-frontend-interactives-controller/storage"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"io"
	"sync"
)

// Ensure, that ProviderMock does implement storage.Provider.
// If this is not the case, regenerate this file with moq.
var _ storage.Provider = &ProviderMock{}

// ProviderMock is a mock implementation of storage.Provider.
//
// 	func TestSomethingThatUsesProvider(t *testing.T) {
//
// 		// make and configure a mocked storage.Provider
// 		mockedProvider := &ProviderMock{
// 			CheckerFunc: func() func(context.Context, *healthcheck.CheckState) error {
// 				panic("mock out the Checker method")
// 			},
// 			GetFunc: func(contextMoqParam context.Context, s string) (io.ReadCloser, error) {
// 				panic("mock out the Get method")
// 			},
// 		}
//
// 		// use mockedProvider in code that requires storage.Provider
// 		// and then make assertions.
//
// 	}
type ProviderMock struct {
	// CheckerFunc mocks the Checker method.
	CheckerFunc func() func(context.Context, *healthcheck.CheckState) error

	// GetFunc mocks the Get method.
	GetFunc func(contextMoqParam context.Context, s string) (io.ReadCloser, error)

	// calls tracks calls to the methods.
	calls struct {
		// Checker holds details about calls to the Checker method.
		Checker []struct {
		}
		// Get holds details about calls to the Get method.
		Get []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// S is the s argument value.
			S string
		}
	}
	lockChecker sync.RWMutex
	lockGet     sync.RWMutex
}

// Checker calls CheckerFunc.
func (mock *ProviderMock) Checker() func(context.Context, *healthcheck.CheckState) error {
	if mock.CheckerFunc == nil {
		panic("ProviderMock.CheckerFunc: method is nil but Provider.Checker was just called")
	}
	callInfo := struct {
	}{}
	mock.lockChecker.Lock()
	mock.calls.Checker = append(mock.calls.Checker, callInfo)
	mock.lockChecker.Unlock()
	return mock.CheckerFunc()
}

// CheckerCalls gets all the calls that were made to Checker.
// Check the length with:
//     len(mockedProvider.CheckerCalls())
func (mock *ProviderMock) CheckerCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockChecker.RLock()
	calls = mock.calls.Checker
	mock.lockChecker.RUnlock()
	return calls
}

// Get calls GetFunc.
func (mock *ProviderMock) Get(contextMoqParam context.Context, s string) (io.ReadCloser, error) {
	if mock.GetFunc == nil {
		panic("ProviderMock.GetFunc: method is nil but Provider.Get was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		S               string
	}{
		ContextMoqParam: contextMoqParam,
		S:               s,
	}
	mock.lockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	mock.lockGet.Unlock()
	return mock.GetFunc(contextMoqParam, s)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//     len(mockedProvider.GetCalls())
func (mock *ProviderMock) GetCalls() []struct {
	ContextMoqParam context.Context
	S               string
} {
	var calls []struct {
		ContextMoqParam context.Context
		S               string
	}
	mock.lockGet.RLock()
	calls = mock.calls.Get
	mock.lockGet.RUnlock()
	return calls
}

// Ensure, that S3BucketMock does implement storage.S3Bucket.
// If this is not the case, regenerate this file with moq.
var _ storage.S3Bucket = &S3BucketMock{}

// S3BucketMock is a mock implementation of storage.S3Bucket.
//
// 	func TestSomethingThatUsesS3Bucket(t *testing.T) {
//
// 		// make and configure a mocked storage.S3Bucket
// 		mockedS3Bucket := &S3BucketMock{
// 			CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error {
// 				panic("mock out the Checker method")
// 			},
// 			GetFunc: func(key string) (io.ReadCloser, *int64, error) {
// 				panic("mock out the Get method")
// 			},
// 		}
//
// 		// use mockedS3Bucket in code that requires storage.S3Bucket
// 		// and then make assertions.
//
// 	}
type S3BucketMock struct {
	// CheckerFunc mocks the Checker method.
	CheckerFunc func(ctx context.Context, state *healthcheck.CheckState) error

	// GetFunc mocks the Get method.
	GetFunc func(key string) (io.ReadCloser, *int64, error)

	// calls tracks calls to the methods.
	calls struct {
		// Checker holds details about calls to the Checker method.
		Checker []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// State is the state argument value.
			State *healthcheck.CheckState
		}
		// Get holds details about calls to the Get method.
		Get []struct {
			// Key is the key argument value.
			Key string
		}
	}
	lockChecker sync.RWMutex
	lockGet     sync.RWMutex
}

// Checker calls CheckerFunc.
func (mock *S3BucketMock) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	if mock.CheckerFunc == nil {
		panic("S3BucketMock.CheckerFunc: method is nil but S3Bucket.Checker was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		State *healthcheck.CheckState
	}{
		Ctx:   ctx,
		State: state,
	}
	mock.lockChecker.Lock()
	mock.calls.Checker = append(mock.calls.Checker, callInfo)
	mock.lockChecker.Unlock()
	return mock.CheckerFunc(ctx, state)
}

// CheckerCalls gets all the calls that were made to Checker.
// Check the length with:
//     len(mockedS3Bucket.CheckerCalls())
func (mock *S3BucketMock) CheckerCalls() []struct {
	Ctx   context.Context
	State *healthcheck.CheckState
} {
	var calls []struct {
		Ctx   context.Context
		State *healthcheck.CheckState
	}
	mock.lockChecker.RLock()
	calls = mock.calls.Checker
	mock.lockChecker.RUnlock()
	return calls
}

// Get calls GetFunc.
func (mock *S3BucketMock) Get(key string) (io.ReadCloser, *int64, error) {
	if mock.GetFunc == nil {
		panic("S3BucketMock.GetFunc: method is nil but S3Bucket.Get was just called")
	}
	callInfo := struct {
		Key string
	}{
		Key: key,
	}
	mock.lockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	mock.lockGet.Unlock()
	return mock.GetFunc(key)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//     len(mockedS3Bucket.GetCalls())
func (mock *S3BucketMock) GetCalls() []struct {
	Key string
} {
	var calls []struct {
		Key string
	}
	mock.lockGet.RLock()
	calls = mock.calls.Get
	mock.lockGet.RUnlock()
	return calls
}

// Ensure, that DownloadServiceAPIClientMock does implement storage.DownloadServiceAPIClient.
// If this is not the case, regenerate this file with moq.
var _ storage.DownloadServiceAPIClient = &DownloadServiceAPIClientMock{}

// DownloadServiceAPIClientMock is a mock implementation of storage.DownloadServiceAPIClient.
//
// 	func TestSomethingThatUsesDownloadServiceAPIClient(t *testing.T) {
//
// 		// make and configure a mocked storage.DownloadServiceAPIClient
// 		mockedDownloadServiceAPIClient := &DownloadServiceAPIClientMock{
// 			CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error {
// 				panic("mock out the Checker method")
// 			},
// 			DownloadFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, path string) (*download.Response, error) {
// 				panic("mock out the Download method")
// 			},
// 		}
//
// 		// use mockedDownloadServiceAPIClient in code that requires storage.DownloadServiceAPIClient
// 		// and then make assertions.
//
// 	}
type DownloadServiceAPIClientMock struct {
	// CheckerFunc mocks the Checker method.
	CheckerFunc func(ctx context.Context, state *healthcheck.CheckState) error

	// DownloadFunc mocks the Download method.
	DownloadFunc func(ctx context.Context, userAuthToken string, serviceAuthToken string, path string) (*download.Response, error)

	// calls tracks calls to the methods.
	calls struct {
		// Checker holds details about calls to the Checker method.
		Checker []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// State is the state argument value.
			State *healthcheck.CheckState
		}
		// Download holds details about calls to the Download method.
		Download []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserAuthToken is the userAuthToken argument value.
			UserAuthToken string
			// ServiceAuthToken is the serviceAuthToken argument value.
			ServiceAuthToken string
			// Path is the path argument value.
			Path string
		}
	}
	lockChecker  sync.RWMutex
	lockDownload sync.RWMutex
}

// Checker calls CheckerFunc.
func (mock *DownloadServiceAPIClientMock) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	if mock.CheckerFunc == nil {
		panic("DownloadServiceAPIClientMock.CheckerFunc: method is nil but DownloadServiceAPIClient.Checker was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		State *healthcheck.CheckState
	}{
		Ctx:   ctx,
		State: state,
	}
	mock.lockChecker.Lock()
	mock.calls.Checker = append(mock.calls.Checker, callInfo)
	mock.lockChecker.Unlock()
	return mock.CheckerFunc(ctx, state)
}

// CheckerCalls gets all the calls that were made to Checker.
// Check the length with:
//     len(mockedDownloadServiceAPIClient.CheckerCalls())
func (mock *DownloadServiceAPIClientMock) CheckerCalls() []struct {
	Ctx   context.Context
	State *healthcheck.CheckState
} {
	var calls []struct {
		Ctx   context.Context
		State *healthcheck.CheckState
	}
	mock.lockChecker.RLock()
	calls = mock.calls.Checker
	mock.lockChecker.RUnlock()
	return calls
}

// Download calls DownloadFunc.
func (mock *DownloadServiceAPIClientMock) Download(ctx context.Context, userAuthToken string, serviceAuthToken string, path string) (*download.Response, error) {
	if mock.DownloadFunc == nil {
		panic("DownloadServiceAPIClientMock.DownloadFunc: method is nil but DownloadServiceAPIClient.Download was just called")
	}
	callInfo := struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		Path             string
	}{
		Ctx:              ctx,
		UserAuthToken:    userAuthToken,
		ServiceAuthToken: serviceAuthToken,
		Path:             path,
	}
	mock.lockDownload.Lock()
	mock.calls.Download = append(mock.calls.Download, callInfo)
	mock.lockDownload.Unlock()
	return mock.DownloadFunc(ctx, userAuthToken, serviceAuthToken, path)
}

// DownloadCalls gets all the calls that were made to Download.
// Check the length with:
//     len(mockedDownloadServiceAPIClient.DownloadCalls())
func (mock *DownloadServiceAPIClientMock) DownloadCalls() []struct {
	Ctx              context.Context
	UserAuthToken    string
	ServiceAuthToken string
	Path             string
} {
	var calls []struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		Path             string
	}
	mock.lockDownload.RLock()
	calls = mock.calls.Download
	mock.lockDownload.RUnlock()
	return calls
}
