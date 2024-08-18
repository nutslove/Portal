{{ define "board" }}
<div class="board">
  <h1>{{ .BoardName }}</h1>
  <nav class="search">
  <form action="/{{ .BoardType }}/searching" method="post">
    {{/* <input type="text" name="search" placeholder="Search post..."> */}}
    <input type="text" placeholder="Search post...">
    <button type="submit">検索</button>
  </form>
  <a href="/{{ .BoardType }}/posting" class="button">新規投稿</a>
  </nav>
  <table>
    <thead>
      <tr>
        <th style="width: 10%;text-align: center;">No.</th>
        <th style="width: 47%;text-align: center;">題名</th>
        <th style="width: 15%;text-align: center;">作成者</th>
        <th style="width: 18%;text-align: center;">作成日</th>
        <th style="width: 10%;text-align: center;">閲覧数</th>
      </tr>
    </thead>
    <tbody>
      {{ range .posts }}
      <tr>
        <td style="width: 10%;text-align: center;">{{ .Number }}</td>
        <td style="width: 47%;">
          <a class="postlist" href="/{{ $.BoardType }}/{{ $.page }}?post={{ .Number }}">
          {{ .Title }}
          </a>
        </td>
        <td style="width: 15%;">{{ .Author }}</td>
        <td style="width: 18%;">{{ .Date }}</td>
        <td style="width: 10%;text-align: center;">{{ .Count }}</td>
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
      {{ . }}
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
{{ end }}