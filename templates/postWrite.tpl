{{ define "postWrite" }}
<textarea id="editor" placeholder="ここにMarkdownを入力してください"></textarea>
<div id="preview"></div>
<button id="submit">投稿する</button>

<script>
    const editor = document.getElementById('editor');
    const preview = document.getElementById('preview');
    const submit = document.getElementById('submit');

    function updatePreview() {
        preview.innerHTML = marked.parse(editor.value);
    }

    editor.addEventListener('input', updatePreview);

    submit.addEventListener('click', function() {
        // ここで投稿処理を行います
        // 例: サーバーにPOSTリクエストを送信
        console.log('投稿内容:', editor.value);
        alert('投稿されました！（実際の送信は行っていません）');
    });

    // 初期プレビューの更新
    updatePreview();
</script>
{{ end }}