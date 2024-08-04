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
    <span>&lt;&lt;</span>
    <span>前ページ</span>
    <span>1</span>
    <span>2</span>
    <span>3</span>
    <span>4</span>
    <span>5</span>
    <span>6</span>
    <span>7</span>
    <span>8</span>
    <span>9</span>
    <span>10</span>
    <span>次ページ</span>
    <span>&gt;&gt;</span>
  </div>
</div>
{{ end }}