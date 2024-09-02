{{ define "board" }}
<div class="board">
  <h1>{{ .BoardName }}</h1>
  <nav class="search">
  <form action="/{{ .BoardType }}/search" method="get" onsubmit="return validateSearchForm()">
  {{/* <form action="/{{ .BoardType }}/search" method="get"> */}}
    {{/* <input type="text" name="query" placeholder="Search post..." required> */}} 
    {{/* requiredだけでは、スペースだけ入れたパターンを弾けない */}}
    <input type="text" name="query" placeholder="Search post..." id="searchQuery">
    <button type="submit">検索</button>
  </form>
  {{ if .IsLoggedIn }}
  <a href="/{{ .BoardType }}/posting" class="button">新規投稿</a>
  {{ end }}
  </nav>
  <div class="search-result">
    {{ if .query }}
      "{{ .query }}" 検索結果: {{ .SearchResultCount}}件
    {{ end }}
  </div>
  <table>
    <thead>
      <tr>
        <th style="width: 10%;text-align: center;">No.</th>
        <th style="width: 47%;text-align: center;">題名</th>
        <th style="width: 15%;text-align: center;">作成者</th>
        <th style="width: 18%;text-align: center;">作成日</th>
        {{/* <th style="width: 10%;text-align: center;">閲覧数</th> */}}
      </tr>
    </thead>
    <tbody>
      {{ range .posts }}
      <tr>
        <td style="width: 10%;text-align: center;">{{ .Number }}</td>
        <td style="width: 47%;">
          {{/* <a class="postlist" href="/{{ $.BoardType }}/{{ $.page }}?post={{ .Number }}"> */}}
          <a class="postlist" href="/{{ $.BoardType }}/posting/{{ .Number }}">
          {{ .Title }}
          </a>
        </td>
        <td style="width: 15%;">{{ .Author }}</td>
        <td style="width: 18%;">{{ .CreatedAt }}</td>
        {{/* <td style="width: 10%;text-align: center;">{{ .Count }}</td> */}}
      </tr>
      {{ end }}
    </tbody>
  </table>
  <div class="pagination">
    {{ if ne .page 1 }}
    <span>
      <a href=/{{ .BoardType }}/1><<</a>
    </span>
    <span style="margin-right: 30px">
      <a href=/{{ .BoardType }}/{{ subtract .page 1 }}>Previous</a>
    </span>
    {{ end }}
    {{ range .pageSlice }}
    <span>
      {{ if eq $.page . }}
      [{{ . }}]
      {{ else }}
      <a href=/{{ $.BoardType }}/{{ . }}>{{ . }}</a>
      {{ end}}
    </span>
    {{ end }}
    {{ if ne .page .pageNum }}
    <span style="margin-left: 30px">
      <a href=/{{ .BoardType }}/{{ add .page 1 }}>Next</a>
    </span>
    <span>
      <a href=/{{ .BoardType }}/{{ .pageNum }}>>></a>
    </span>
    {{ end }}
  </div>
</div>

<script>
function validateSearchForm() {
  let query = document.getElementById("searchQuery").value;
  if (query == null || query.trim() == "") {
    alert("検索キーワードを入力してください");
    return false;
  }
  return true;
}
</script>
{{ end }}