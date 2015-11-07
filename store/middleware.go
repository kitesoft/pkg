// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"net/http"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/ctxjwt"
	"github.com/corestoreio/csfw/net/httputils"
	"github.com/juju/errgo"
	"golang.org/x/net/context"
)

// WithValidateBaseUrl is a middleware which checks if the request base URL
// is equal to the one store in the configuration, if not
// i.e. redirect from http://example.com/store/ to http://www.example.com/store/
// @see app/code/Magento/Store/App/FrontController/Plugin/RequestPreprocessor.php
func WithValidateBaseUrl(cr config.ReaderPubSuber) ctxhttp.Middleware {

	// Having the GetBool command here, means you must restart the app to take
	// changes in effect. @todo refactor and use pub/sub to automatically change
	// the isRedirectToBase value.
	checkBaseURL, err := cr.GetBool(config.Path(PathRedirectToBase)) // scope default
	if config.NotKeyNotFoundError(err) && PkgLog.IsDebug() {
		PkgLog.Debug("ctxhttp.WithValidateBaseUrl.GetBool", "err", err, "path", PathRedirectToBase)
	}

	redirectCode := http.StatusMovedPermanently
	if rc, err := cr.GetInt(config.Path(PathRedirectToBase)); rc != redirectCode && false == config.NotKeyNotFoundError(err) {
		redirectCode = http.StatusFound
	}

	return func(h ctxhttp.Handler) ctxhttp.Handler {
		return ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			if checkBaseURL && r.Method != "POST" {

				_, requestedStore, err := FromContextReader(ctx)
				if err != nil {
					if PkgLog.IsDebug() {
						PkgLog.Debug("ctxhttp.WithValidateBaseUrl.FromContextServiceReader", "err", err, "ctx", ctx)
					}
					return errgo.Mask(err)
				}

				baseURL, err := requestedStore.BaseURL(config.URLTypeWeb, requestedStore.IsCurrentlySecure(r))
				if err != nil {
					if PkgLog.IsDebug() {
						PkgLog.Debug("ctxhttp.WithValidateBaseUrl.requestedStore.BaseURL", "err", err, "ctx", ctx)
					}
					return errgo.Mask(err)
				}

				if err := httputils.IsBaseUrlCorrect(r, &baseURL); err != nil {
					if PkgLog.IsDebug() {
						PkgLog.Debug("store.WithValidateBaseUrl.IsBaseUrlCorrect.error", "err", err, "baseURL", baseURL, "request", r)
					}

					baseURL.Path = r.URL.Path
					baseURL.RawPath = r.URL.RawPath
					baseURL.RawQuery = r.URL.RawQuery
					baseURL.Fragment = r.URL.Fragment
					http.Redirect(w, r, (&baseURL).String(), redirectCode)
					return nil
				}
			}
			return h.ServeHTTPContext(ctx, w, r)
		})
	}
}

// WithInitStoreByToken is a middleware which initializes a request based store
// via a JSON Web Token.
// Extracts the store.Reader and jwt.Token from context.Context. If the requested
// store is different than the initialized requested store than the new requested
// store will be saved in the context.
func WithInitStoreByToken() ctxhttp.Middleware {

	return func(h ctxhttp.Handler) ctxhttp.Handler {
		return ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			storeService, requestedStore, err := FromContextReader(ctx)
			if err != nil {
				if PkgLog.IsDebug() {
					PkgLog.Debug("store.WithInitStoreByToken.FromContextServiceReader", "err", err, "ctx", ctx)
				}
				return errgo.Mask(err)
			}

			token, err := ctxjwt.FromContext(ctx)
			if err != nil {
				if PkgLog.IsDebug() {
					PkgLog.Debug("store.WithInitStoreByToken.ctxjwt.FromContext.err", "err", err, "ctx", ctx)
				}
				return errgo.Mask(err)
			}

			scopeOption, err := StoreCodeFromClaim(token.Claims)
			if err != nil {
				if PkgLog.IsDebug() {
					PkgLog.Debug("store.WithInitStoreByToken.StoreCodeFromClaim", "err", err, "token", token, "ctx", ctx)
				}
				return errgo.Mask(err)
			}

			newRequestedStore, err := storeService.RequestedStore(scopeOption)
			if err != nil {
				if PkgLog.IsDebug() {
					PkgLog.Debug("store.WithInitStoreByToken.RequestedStore", "err", err, "token", token, "scopeOption", scopeOption, "ctx", ctx)
				}
				return errgo.Mask(err)
			}

			if newRequestedStore.StoreID() != requestedStore.StoreID() {
				// this may lead to a bug because the previously set storeService and requestedStore
				// will still exists and have not been removed.
				ctx = NewContextReader(ctx, storeService, newRequestedStore)
			}

			return h.ServeHTTPContext(ctx, w, r)
		})
	}
}

// InitByRequest returns a new Store read from a cookie or HTTP request parameter.
// It calls GetRequestStore() to determine the correct store.
// The internal appStore must be set before hand, call Init() before calling this function.
// 1. check cookie store, always a string and the store code
// 2. check for ___store variable, always a string and the store code
// 3. May return nil,nil if nothing is set.
// This function must be used within an HTTP handler.
// The returned new Store must be used in the HTTP context and overrides the appStore.
func WithInitStoreByRequest(scopeType scope.Scope) ctxhttp.Middleware {

	return func(h ctxhttp.Handler) ctxhttp.Handler {
		return ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			return nil
		})
	}

	// res http.ResponseWriter, req *http.Request
	//	if sm.appStore == nil {
	//		// that means you must call Init() before executing this function.
	//		return nil, ErrAppStoreNotSet
	//	}
	//
	//	var reqStore *Store
	//	var so scope.Option
	//	var err error
	//	so, err = StoreCodeFromForm(req)
	//	if err != nil { // no cookie set, lets try via form to find the store code
	//
	//		if err == ErrStoreCodeInvalid {
	//			return nil, PkgLog.Error("store.Service.InitByRequest.GetCodeFromForm", "err", err, "req", req, "scopeType", scopeType.String())
	//		}
	//
	//		so, err = StoreCodeFromCookie(req)
	//		switch err {
	//		case ErrStoreCodeEmpty, http.ErrNoCookie:
	//			err = nil
	//		case nil:
	//		// do nothing
	//		default: // err != nil
	//			return nil, PkgLog.Error("store.Service.InitByRequest.GetCodeFromCookie", "err", err, "req", req, "scopeType", scopeType.String())
	//		}
	//	}
	//
	//	// @todo reqStoreCode if number ... cast to int64 because then group id if ScopeGroup is group.
	//	if reqStore, err = sm.GetRequestStore(so, scopeType); err != nil {
	//		return nil, PkgLog.Error("store.Service.InitByRequest.GetRequestStore", "err", err)
	//	}
	//	soStoreCode := so.StoreCode()
	//
	//	// also delete and re-set a new cookie
	//	if reqStore != nil && reqStore.Data.Code.String == soStoreCode {
	//		wds, err := reqStore.Website.DefaultStore()
	//		if err != nil {
	//			return nil, PkgLog.Error("store.Service.InitByRequest.Website.DefaultStore", "err", err, "soStoreCode", soStoreCode)
	//		}
	//		if wds.Data.Code.String == soStoreCode {
	//			reqStore.DeleteCookie(res) // cookie not needed anymore
	//		} else {
	//			reqStore.SetCookie(res) // make sure we force set the new store
	//		}
	//	}
	//
	//	return reqStore, nil // can be nil,nil
}
