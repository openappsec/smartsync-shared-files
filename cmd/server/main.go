// Copyright (C) 2022 Check Point Software Technologies Ltd. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"os"
	"time"

	"openappsec.io/errors/errorloader"

	"openappsec.io/log"
	"openappsec.io/smartsync-shared-files/internal/app/ingector"
)

const (
	initTimeout    = 15
	stopTimeout    = 15
	exitCauseError = 1
	exitCauseDone  = 0
	// errors json file location
	// please edit this file with your own responses messages
	// (the structure of a single response should remain the same: message, description, code, severity
	errorsFilePath = "configs/error-responses.json"

	// please change the msrv error code to your own (should be 3 digits)
	// add your code to this wiki page - https://wiki.checkpoint.com/confluence/display/PROJECTINFO/Error+responses+-+MSRV+codes
	msrvErrCode = "1111"
)

func main() {
	// context object explanation can be found in this repository's Gitlab Wiki page here:
	// https://openappsec.io/smartsync-shared-files/-/wikis/Go-Template#the-context-object
	initCtx, initCancel := context.WithTimeout(context.Background(), initTimeout*time.Second)

	// loads the errors json file
	// the documentation of the errorloader package can be found here:
	// https://openappsec.io/errors
	err := errorloader.Configure(errorsFilePath, msrvErrCode)

	// Initialize the Application
	// Note that in GoLang the dependency injector is called "ingector"
	// more about the ingector in the Go wire package can be found in this gitlab wiki:
	// https://openappsec.io/smartsync-shared-files/-/wikis/Go-Template#wire-and-dependency-injection
	app, err := ingector.InitializeApp(initCtx)
	if err != nil {
		log.WithContext(initCtx).Error("Failed to inject dependencies! Error: ", err.Error())
		initCancel()
		os.Exit(exitCauseError)
	}

	// cancels the initial context since it is not supposed to be used anymore
	// this cancellation "alerts" every function which uses this context object instance
	// that it is invalid
	// upon starting the application a new context will be created for each request
	// explanation about the context object can be found in the wiki page mentioned in the comment
	// above the context creation command of this main function
	initCancel()

	exitCode := exitCauseDone
	if err := app.Start(); err != nil {
		log.WithContext(initCtx).Error("Failed to start app: ", err.Error())
		exitCode = exitCauseError
	}

	stopCtx, stopCancel := context.WithTimeout(context.Background(), stopTimeout*time.Second)
	if err := app.Stop(stopCtx); err != nil {
		log.WithContext(context.Background()).Error("Failed to stop app: ", err.Error())
		exitCode = exitCauseError
	}

	stopCancel()
	os.Exit(exitCode)
}
