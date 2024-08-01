{{ define "top" }}
<!DOCTYPE html>
<html lang="ja">
<head>
  <meta charset="UTF-8">
  <title>TechCareer Talk</title>
  <link rel="icon" type="image/x-icon" href="/static/favicon.ico">
  <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <header>
      {{ template "header" . }}
    </header>
    <main>
    {{ template "sidebar" . }}
    {{ if not .TopPage }}
    {{ template "board" . }}
    {{ end }}
    </main>
{{ end }}
{{ define "bottom" }}
</body>
</html>
{{ end }}