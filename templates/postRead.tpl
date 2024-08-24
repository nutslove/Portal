{{ define "postRead" }}
<div class="post-read">
    <h1 class="post-title">{{ .PostTitle }}</h1>
    <div id="preview"></div>
    {{ if and (.IsLoggedIn) (eq .Username .Author) }}
    <div class=post-read-footer>
        <span>
            <a href=/>編集</a>
        </span>
        <span>
            <a href=/>削除</a>
        </span>
    </div>
    {{ end }}
    <script src="https://cdnjs.cloudflare.com/ajax/libs/marked/4.0.2/marked.min.js"></script>
    <script>
        function updatePreview() {
            var content = '{{ .PostContent }}';
            document.getElementById('preview').innerHTML = marked.parse(content);
        }
        updatePreview();
    </script>
    <style>
    #preview { 
      overflow-y: auto;
      word-wrap: break-word;
      overflow-wrap: break-word;
    }
    </style>
</div>
{{ end }}