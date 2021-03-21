// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"github.com/nbkit/mdf/log"
	"html/template"
	"runtime"
	"strconv"
	"strings"
)

const ginSupportMinGoVer = 12

// IsDebugging returns true if the framework is running in debug mode.
// Use SetMode(gin.ReleaseMode) to disable debug mode.
func IsDebugging() bool {
	return ginMode == debugCode
}

// DebugPrintRouteFunc indicates debug log output format.
var DebugPrintRouteFunc func(httpMethod, absolutePath, handlerName string, nuHandlers int)

func debugPrintRoute(httpMethod, absolutePath string, handlers HandlersChain) {
	nuHandlers := len(handlers)
	handlerName := nameOfFunction(handlers.Last())
	if DebugPrintRouteFunc == nil {
		debugPrint("%-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
	} else {
		DebugPrintRouteFunc(httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

func debugPrintLoadTemplate(tmpl *template.Template) {
	var buf strings.Builder
	for _, tmpl := range tmpl.Templates() {
		buf.WriteString("\t- ")
		buf.WriteString(tmpl.Name())
		buf.WriteString("\n")
	}
	debugPrint("Loaded HTML Templates (%d) : %s", len(tmpl.Templates()), buf.String())
}

func debugPrint(format string, values ...interface{}) {
	log.Info().Msgf(format, values...)
}

func getMinVer(v string) (uint64, error) {
	first := strings.IndexByte(v, '.')
	last := strings.LastIndexByte(v, '.')
	if first == last {
		return strconv.ParseUint(v[first+1:], 10, 64)
	}
	return strconv.ParseUint(v[first+1:last], 10, 64)
}

func debugPrintWARNINGDefault() {
	if v, e := getMinVer(runtime.Version()); e == nil && v <= ginSupportMinGoVer {
		debugPrint(`Now Gin requires Go 1.12+.`)
	}
	debugPrint(`Creating an Engine instance with the IOutput and Recovery middleware already attached.

`)
}

func debugPrintWARNINGNew() {
	log.Info().Msgf(`Running in %s mode`, modeName)
}

func debugPrintWARNINGSetHTMLTemplate() {
	debugPrint(`Since SetHTMLTemplate() is NOT thread-safe. It should only be called
at initialization. ie. before any route is registered or the router is listening in a socket:

	router := gin.defaultOutput()
	router.SetHTMLTemplate(template) // << good place

`)
}

func debugPrintError(err error) {
	if err != nil {
		if IsDebugging() {
			log.Print(err)
		}
	}
}
