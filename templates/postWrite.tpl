{{ define "postWrite" }}
<div class="post">
    <div class="editor-wrapper">
        <div class="segmented-control">
            <input type="radio" name="sc-1-1" id="sc-1-1-1" checked>
            <input type="radio" name="sc-1-1" id="sc-1-1-2">
            <label for="sc-1-1-1" data-value="Markdown">Markdown</label>
            <label for="sc-1-1-2" data-value="Preview">Preview</label>
        </div>
        <div class="editor-container">
            <textarea id="editor" placeholder="ここにMarkdownを入力してください"></textarea>
            <div id="preview"></div>
        </div>
    </div>
    <div id="loading" style="display: none; text-align: center;">
        <img src="https://i.gifer.com/ZZ5H.gif" alt="Loading..." width="50" height="50">
    </div>
    <button id="submit">投稿</button>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/marked/4.0.2/marked.min.js"></script>
    <style>
        .editor-wrapper {
            width: 90%;
            margin: 20px auto;
            position: relative;
        }
        .editor-container {
            width: 100%;
            border: 2px solid #ddd;
            border-radius: 8px;
            overflow: hidden;
        }
        #editor, #preview {
            width: 100%;
            height: 400px;
            padding: 10px;
            box-sizing: border-box;
            border: none;
        }
        #editor, #preview { resize: vertical; }
        #preview { overflow-y: auto; display: none; }
        #submit {
            display: block;
            margin: 20px auto;
            padding: 10px 20px;
            background-color: #2196F3;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
        }

        /* セグメントコントロールのスタイル */
        .segmented-control {
            position: absolute;
            top: -45px;
            right: 0;
            width: 200px;
            border: 1px solid #deb887;
            border-radius: 20px;
            display: flex;
            overflow: hidden;
            background-color: #fff;
            font-size: 13px;
        }
        .segmented-control input[type="radio"] {
            display: none;
        }
        .segmented-control label {
            flex: 1;
            text-align: center;
            padding: 8px 0;
            cursor: pointer;
            transition: all 0.3s ease;
            border-radius: 20px;
        }
        .segmented-control label[for="sc-1-1-1"] {
            background: #deb887;
            color: white;
        }
        .segmented-control input:checked + label {
            background: #deb887;
            color: white;
        }
        .segmented-control input:not(:checked) + label {
            background: #fff;
            color: #333;
        }
    </style>
    <script>
        const editor = document.getElementById('editor');
        const preview = document.getElementById('preview');
        const submit = document.getElementById('submit');
        const loading = document.getElementById('loading');
        const markdownRadio = document.getElementById('sc-1-1-1');
        const previewRadio = document.getElementById('sc-1-1-2');
        const markdownLabel = document.querySelector('label[for="sc-1-1-1"]');
        const previewLabel = document.querySelector('label[for="sc-1-1-2"]');

        function updatePreview() {
            preview.innerHTML = marked.parse(editor.value);
        }

        function toggleView() {
            if (previewRadio.checked) {
                editor.style.display = 'none';
                preview.style.display = 'block';
                previewLabel.style.background = '#deb887';
                previewLabel.style.color = 'white';
                markdownLabel.style.background = '#fff';
                markdownLabel.style.color = '#333';
            } else {
                editor.style.display = 'block';
                preview.style.display = 'none';
                markdownLabel.style.background = '#deb887';
                markdownLabel.style.color = 'white';
                previewLabel.style.background = '#fff';
                previewLabel.style.color = '#333';
            }
        }

        editor.addEventListener('input', updatePreview);
        markdownRadio.addEventListener('change', toggleView);
        previewRadio.addEventListener('change', toggleView);

        submit.addEventListener('click', function() {
            {{/* console.log('投稿内容:', editor.value); */}}
            loading.style.display = 'block'; // ローディング表示
            submit.disabled = true; // ボタンを無効化

            const data = {
                content: editor.value
            };

            fetch('/{{ .BoardType }}/posting', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            })
            .then(response => response.json())
            .then(result => {
                loading.style.display = 'none'; // ローディング非表示
                submit.disabled = false; // ボタンを有効化

                if (result.success) {
                    // 成功時の処理
                    window.location.href = result.redirectUrl; // リダイレクト先のURLを指定
                } else {
                    // 失敗時の処理
                    alert('エラーが発生しました: ' + result.message);
                }
            })
            .catch(error => {
                loading.style.display = 'none'; // ローディング非表示
                submit.disabled = false; // ボタンを有効化

                alert('リクエストに失敗しました: ' + error.message);
            });

            {{/* alert('投稿に失敗しました。時間をおいて再度試してください。'); */}}
        });

        // 初期プレビューの更新
        updatePreview();
        toggleView(); // 初期状態の設定
    </script>
</div>
{{ end }}

/* 上記処理について
`editor.addEventListener('input', updatePreview);`の部分について詳しく説明します。

### `input`イベントのリスナーについて
`editor.addEventListener('input', updatePreview);`は、`<textarea>` 要素（`editor`）で何か入力があったときに、`updatePreview`関数が実行されるように設定されています。

### `updatePreview`関数の役割
この関数は、Markdownで書かれた内容をHTMLに変換し、それを`<div id="preview">`要素に表示します。

```javascript
function updatePreview() {
    preview.innerHTML = marked.parse(editor.value);
}
```

#### この関数の動作:
1. **`editor.value`**: `textarea`に現在入力されているMarkdownの内容を取得します。
2. **`marked.parse(editor.value)`**: 取得したMarkdownの内容を`marked.parse`を使ってHTMLに変換します。
3. **`preview.innerHTML = ...`**: 変換されたHTMLを`<div id="preview">`の`innerHTML`に設定し、プレビュー領域に表示します。

### Previewボタンを押したときの動作
Previewボタンを押すと、以下の理由でMarkdownのプレビューが表示されます。

1. **`toggleView`関数の動作**:
   - Previewボタンが押されると、`toggleView`関数が実行されます。
   - この関数は、`textarea`要素（Markdownエディタ）を隠し、`<div id="preview">`要素を表示します。
   - これにより、`textarea`に入力されていたMarkdownがHTMLに変換されて表示されます。

2. **`input`イベントによるリアルタイム更新**:
   - 既に`editor.addEventListener('input', updatePreview);`が設定されているため、ユーザーがMarkdownエディタに何か入力すると、`updatePreview`関数が自動的に実行されます。
   - そのため、ユーザーがPreviewボタンを押してプレビューを表示する時点では、`<div id="preview">`には既にMarkdownの内容がHTMLに変換されて表示されています。

このようにして、ユーザーがMarkdownを入力すると、その内容がリアルタイムでプレビュー領域に反映される仕組みが実現されています。プレビューは、あらかじめ`input`イベントを使って更新されているため、Previewボタンを押すとすぐに変換済みの内容が表示されるのです。
*/