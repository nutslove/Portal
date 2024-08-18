{{ define "top" }}
<!DOCTYPE html>
<html lang="ja">
<head>
  <meta charset="UTF-8">
  <!--<meta name="viewport" content="width=device-width, initial-scale=1.0">-->
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
    {{ if .PostList }}
      {{ template "board" . }}
    {{ else if .PostRead }}
      {{ template "postRead" . }}
    {{ else if .PostWrite }}
      {{ template "postWrite" . }}
    {{ end }}
    </main>
{{ end }}
{{ define "bottom" }}
</body>
</html>
{{ end }}