{{ define "board" }}
<div class="board">
  <table>
    <thead>
      <tr>
        <th>No.</th>
        <th>題名</th>
        <th>作成者</th>
        <th>作成日</th>
        <th>閲覧数</th>
      </tr>
    </thead>
    <tbody>
      {{ range .posts }}
      <tr>
        <td style="width: 30px">{{ .Number }}</td>
        <td style="width: 300px">
          <a>
          {{ .Title }}
          </a>
        </td>
        <td style="width: 90px">{{ .Author }}</td>
        <td style="width: 100px">{{ .Date }}</td>
        <td style="width: 30px">{{ .Count }}</td>
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