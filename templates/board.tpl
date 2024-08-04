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
    <span><<</span>
    <span style="margin-right: 30px">Previous</span>
    {{ end }}
    {{ range .pageSlice }}
    <span>{{ . }}</span>
    {{ end }}
    {{ if ne .page .pageNum }}
    <span style="margin-left: 30px">Next</span>
    <span>>></span>
    {{ end }}
  </div>
</div>
{{ end }}