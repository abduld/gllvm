//
// OCCAM
//
// Copyright (c) 2017, SRI International
//
//  All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// * Neither the name of SRI International nor the names of its contributors may
//   be used to endorse or promote products derived from this software without
//   specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package shared

import (
	"fmt"
	"os"
	"strings"
)

const (
	errorV = iota
	warningV
	infoV
	debugV
)

//loggingLevels is the accepted logging levels.
var loggingLevels = map[string]int{
	"ERROR":   errorV,
	"WARNING": warningV,
	"INFO":    infoV,
	"DEBUG":   debugV,
}

//loggingLevel is the user configured level of logging: ERROR, WARNING, INFO, DEBUG
var loggingLevel = errorV

//loggingFilePointer is where the logging is streamed too.
var loggingFilePointer = os.Stderr

func init() {
	if LLVMLoggingLevel != "" {
		if envLevelVal, ok := loggingLevels[LLVMLoggingLevel]; ok {
			loggingLevel = envLevelVal
		}
	}
	if LLVMLoggingFile != "" {
		//FIXME: is it overboard to defer a close? the OS will close when the process gets cleaned up, do we win
		//anything by being OCD?
		if loggingFP, err := os.OpenFile(LLVMLoggingFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600); err == nil {
			loggingFilePointer = loggingFP
		}
	}
}

func makeLogger(lvl int) func(format string, a ...interface{}) {
	return func(format string, a ...interface{}) {
		if loggingLevel >= lvl {
			msg := fmt.Sprintf(format, a...)
			if !strings.HasSuffix(msg, "\n") {
				msg += "\n"
			}
			//FIXME: (?) if loggingFilePointer != os.Stderr, we could multiplex here
			//and send output to both os.Stderr and loggingFilePointer. We wouldn't
			//want the user to miss any excitement. Then the sanity checker would not
			//need to futz with it.
			loggingFilePointer.WriteString(msg)
		}
	}
}

//LogDebug logs to the configured stream if the logging level is DEBUG.
var LogDebug = makeLogger(debugV)

//LogInfo logs to the configured stream if the logging level is INFO or lower.
var LogInfo = makeLogger(infoV)

//LogWarning logs to the configured stream if the logging level is WARNING or lower.
var LogWarning = makeLogger(warningV)

//LogError logs to the configured stream if the logging level is ERROR or lower.
var LogError = makeLogger(errorV)

//LogFatal logs to the configured stream and then exits.
func LogFatal(format string, a ...interface{}) {
	LogError(format, a...)
	os.Exit(1)
}

//LogWrite writes to the logging stream, irregardless of levels.
var LogWrite = makeLogger(-1)
