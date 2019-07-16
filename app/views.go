// Copyright 2019 Seth Vargo
// Copyright 2019 Google, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"text/template"
)

const (
	htmlProdErr = `
<!DOCTYPE html>
<html>
	<head>
		<title>An error occurred</title>
	</head>

	<body>
		<p style="color:#dd0000; font-weight:bold;">An error occurred.</p>
	</body>
</html>
`

	htmlDevErr = `
<!DOCTYPE html>
<html>
	<head>
		<title>An error occurred</title>
		<style type="text/css">
			body {
				padding: 20px 40px;
			}

			pre {
				background: #eee;
				border-radius: 3px;
				font-size: 1.2em;
				line-height: 1.4;
				overflow: scroll;
				padding: 15px;
				text-align: left;
				white-space: pre-line;
			}
		</style>
	</head>

	<body>
		<h1>An error occurred!</h1>
		<pre>
		{{ .Error }}
		</pre>

		<h2>Request</h2>
		<pre>
			Host:            {{ .Request.Host }}
			Method:          {{ .Request.Method }}
			Protocol:        {{ .Request.Proto }}
			RemoteAddr:      {{ .Request.RemoteAddr }}
			XForwardedFor:   {{ .XForwardedFor }}
			RequestURI:      {{ .Request.RequestURI }}
		</pre>

		<h2>Environment</h2>
		<pre>{{ range .Env }}
		{{ . }}{{ end }}
		</pre>

		<p style="font-size: smaller;">
			You are seeing this information because you are running in non-production
			mode. Set <tt>ENV=production</tt> to disable detailed error pages.
		</p>
	</body>
</html>
`

	htmlIndex = `
<!DOCTYPE html>
<html>
	<head>
		<title>Welcome</title>
	</head>

	<body>
		<h1>Welcome to my site!</h1>
		<h2>You are visitor number {{ .Count }}.</h2>
		You can also <a href="/reset-counter?count=0">reset the counter</a>.
	</body>
</html>
`
)

func processTemplate(tmpl string, data interface{}) string {
	t := template.Must((template.New("tpl").Parse(tmpl)))
	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		panic(err)
	}
	return b.String()
}
