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
            border: 1px solid #ddd;
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
            border: 1px solid #ba55d3;
            border-radius: 20px;
            display: flex;
            overflow: hidden;
            background-color: #fff;
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
            background: #2196F3;
            color: white;
        }
        .segmented-control input:checked + label {
            background: #2196F3;
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
                previewLabel.style.background = '#ba55d3';
                previewLabel.style.color = 'white';
                markdownLabel.style.background = '#fff';
                markdownLabel.style.color = '#333';
            } else {
                editor.style.display = 'block';
                preview.style.display = 'none';
                markdownLabel.style.background = '#ba55d3';
                markdownLabel.style.color = 'white';
                previewLabel.style.background = '#fff';
                previewLabel.style.color = '#333';
            }
        }

        editor.addEventListener('input', updatePreview);
        markdownRadio.addEventListener('change', toggleView);
        previewRadio.addEventListener('change', toggleView);

        submit.addEventListener('click', function() {
            console.log('投稿内容:', editor.value);
            alert('投稿されました！（実際の送信は行っていません）');
        });

        // 初期プレビューの更新
        updatePreview();
        toggleView(); // 初期状態の設定
    </script>
</div>
{{ end }}