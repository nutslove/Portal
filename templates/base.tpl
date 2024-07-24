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
  <div class="total">
    <header>
      {{ template "header" . }}
    </header>
    {{ template "sidebar" . }}
  </div>
{{ end }}
{{ define "bottom" }}
</body>
</html>
{{ end }}