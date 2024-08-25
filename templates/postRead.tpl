{{ define "postRead" }}
<div class="post-read">
    <h1 class="post-title">{{ .PostTitle }}</h1>
    <div id="preview"></div>
    {{ if and (.IsLoggedIn) (eq .Username .Author) }}
    <div id="loading" style="display: none; text-align: center;">
        <img src="https://i.gifer.com/ZZ5H.gif" alt="Loading..." height="30">
    </div>
    <div class="post-read-footer">
        <button id="modify" class="modify">編集</a>
        <button id="delete" class="delete">削除</a>
    </div>
    {{ end }}
    <script src="https://cdnjs.cloudflare.com/ajax/libs/marked/4.0.2/marked.min.js"></script>
    <script>
        const loading = document.getElementById('loading')
        const modifySubmit = document.getElementById('modify')
        const deleteSubmit = document.getElementById('delete')

        function updatePreview() {
            const content = '{{ .PostContent }}';
            document.getElementById('preview').innerHTML = marked.parse(content);
        }

        {{/* // 修正ボタン処理
        modifySubmit.addEventListener('click', function() {


        }); */}}

        // 削除ボタン処理
        deleteSubmit.addEventListener('click', function() {
            {{/* console.log("BoardType:",{{ .BoardType }})
            console.log("PostId:",{{ .PostId }}) */}}

            let result = confirm("本当に削除しますか？")

            // OKを押した時のみ削除処理を実行
            if (result) {
                loading.style.display = 'block'; // ローディング表示
                deleteSubmit.disabled = true; // ボタンを無効化

                fetch('/{{ .BoardType }}/posting/{{ .PostId }}', {
                    method: 'DELETE'
                })
                .then(response => response.json())
                .then(result => {
                    loading.style.display = 'none'; // ローディング非表示
                    deleteSubmit.disabled = false; // ボタンを有効化

                    if (result.success) {
                        // 成功時の処理
                        window.location.href = result.redirectUrl; // リダイレクト先のURLを指定
                    } else {
                        // 失敗時の処理
                        alert('削除に失敗しました: ' + result.message);
                    }
                })
                .catch(error => {
                    loading.style.display = 'none'; // ローディング非表示
                    deleteSubmit.disabled = false; // ボタンを有効化

                    alert('リクエストに失敗しました: ' + error.message);
                });
            }
        });

        updatePreview();
    </script>
    <style>
    #preview { 
        overflow-y: auto;
        word-wrap: break-word;
        overflow-wrap: break-word;
        margin-bottom: 50px;
    }

    .post-read-footer {
        margin-top: 15px;
        margin-bottom: 25px;
        text-align: center;
    }

    .modify {
        margin-right: 10px;
        padding: 6px 20px;
        background-color: #2196F3;
        color: white;
        border-radius: 4px;
        border: none;
        cursor: pointer;
    }
    .delete {
        margin-left: 10px;
        padding: 6px 20px;
        background-color: red;
        color: white;
        border-radius: 4px;
        border: none;
        cursor: pointer;
    }
    </style>
</div>
{{ end }}