package server

import (
	"encoding/base64"
	"html/template"
	"net/http"
	"strings"

	"github.com/skip2/go-qrcode"
)

func urlToLabel(url string) string {
	if strings.Contains(url, "localhost") {
		return "Local"
	}
	// Extract host from https://192.168.1.42:8443
	if idx := strings.Index(url[8:], ":"); idx >= 0 {
		return url[8 : 8+idx]
	}
	return url
}

type urlEntry struct {
	URL     string
	Label   string
	QRImage string // base64-encoded PNG
}

type qrPageData struct {
	URLs []urlEntry
}

const qrPageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>gdweb - QR Codes</title>
	<style>
		* { box-sizing: border-box; }
		body {
			margin: 0;
			min-height: 100vh;
			display: flex;
			align-items: center;
			justify-content: center;
			background: #1a1a1a;
			color: #e0e0e0;
			font-family: system-ui, -apple-system, sans-serif;
		}
		.card {
			background: #252525;
			border-radius: 12px;
			padding: 2rem;
			box-shadow: 0 8px 32px rgba(0,0,0,0.4);
			max-width: 400px;
		}
		h1 {
			margin: 0 0 1.5rem 0;
			font-size: 1.25rem;
			font-weight: 600;
		}
		.tabs {
			display: flex;
			flex-wrap: wrap;
			gap: 0.5rem;
			margin-bottom: 1.5rem;
		}
		.tabs input { display: none; }
		.tabs label {
			padding: 0.5rem 1rem;
			background: #333;
			border-radius: 8px;
			cursor: pointer;
			font-size: 0.9rem;
			transition: background 0.2s;
		}
		.tabs label:hover { background: #404040; }
		.tabs input:checked + label { background: #4a7c59; }
		.tab-panel {
			display: none;
			text-align: center;
			padding: 1rem 0;
		}
		.tab-panel.active { display: block; }
		.qr-container {
			display: inline-block;
			background: #fff;
			padding: 16px;
			border-radius: 8px;
			margin-bottom: 1rem;
		}
		.qr-container img {
			display: block;
			width: 256px;
			height: 256px;
		}
		.url-text {
			font-size: 0.9rem;
			word-break: break-all;
			color: #888;
		}
	</style>
</head>
<body>
	<div class="card">
		<h1>Scan to open game</h1>
		<div class="tabs">
			{{range $i, $u := .URLs}}
			<input type="radio" name="tab" id="tab-{{$i}}" {{if eq $i 0}}checked{{end}}>
			<label for="tab-{{$i}}">{{$u.Label}}</label>
			{{end}}
		</div>
		{{range $i, $u := .URLs}}
		<div class="tab-panel {{if eq $i 0}}active{{end}}" id="panel-{{$i}}">
			<div class="qr-container">
				<img src="data:image/png;base64,{{$u.QRImage}}" alt="QR code for {{$u.URL}}" width="256" height="256">
			</div>
			<div class="url-text">{{$u.URL}}</div>
		</div>
		{{end}}
	</div>
	<script>
		(function() {
			var panels = document.querySelectorAll('.tab-panel');
			var radios = document.querySelectorAll('input[name="tab"]');
			for (var i = 0; i < radios.length; i++) {
				radios[i].addEventListener('change', function() {
					panels.forEach(function(p) { p.classList.remove('active'); });
					var idx = Array.from(radios).indexOf(document.querySelector('input[name="tab"]:checked'));
					panels[idx].classList.add('active');
				});
			}
		})();
	</script>
</body>
</html>`

func qrHandler(urls []string) http.Handler {
	var data qrPageData
	for _, u := range urls {
		png, err := qrcode.Encode(u, qrcode.Medium, 256)
		if err != nil {
			continue // skip URLs that fail to encode
		}
		data.URLs = append(data.URLs, urlEntry{
			URL:     u,
			Label:   urlToLabel(u),
			QRImage: base64.StdEncoding.EncodeToString(png),
		})
	}

	tmpl := template.Must(template.New("qr").Parse(qrPageHTML))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/qr" && r.URL.Path != "/qr/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, data)
	})
}
