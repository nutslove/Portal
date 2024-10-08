<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <link rel="icon" type="image/x-icon" href="/static/favicon.ico">
    <title>{{ .title }} - TechCareer Talk</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: white;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
        }
        .login-container {
            background-color: #D2C7AB;
            padding: 20px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
            border-radius: 5px;
            width: 300px;
            text-align: center;
        }
        .login-container h2 {
            margin-bottom: 20px;
        }
        .login-container input[type="text"], 
        .login-container input[type="password"] {
            width: 70%;
            padding: 10px;
            margin: 5px 0;
            border: 1px solid #ccc;
            border-radius: 3px;
        }
        .login-container button {
            background-color: #28a745;
            color: white;
            margin-top: 10px;
            padding: 10px;
            border: none;
            border-radius: 3px;
            cursor: pointer;
            width: 50%;
        }
        .login-container button:hover {
            background-color: #218838;
        }
        .login-container a {
            display: block;
            margin-top: 20px;
            color: #007bff;
            text-decoration: none;
        }
        .login-container a:hover {
            text-decoration: underline;
        }
        .login-container img {
            margin-left: 10px;
            margin-right: 10px;
            width: 30%;  /* 幅を100%に設定 */
            max-width: 70px;  /* 最大幅を600pxに設定 */
            height: auto;  /* 高さを自動調整 */
        }
    </style>
</head>
<body>
    <div class="login-container">
        <img src="/static/images/santa.png" alt="logo">
        <h2>{{ .title }}</h2>
        {{ if eq .title "Login" }}
            {{ if eq .status "loginfailed" }}
            <p style="color: red;font-style: italic;">NicknameまたはPasswordが間違っています。</p>
            {{ end }}        
            <form method="post" action="/login">
        {{ else if eq .title "SignUp" }}
            {{ if eq .status "signupfailed" }}
                {{ if eq .reason "existalready" }}
                <p style="color: red;font-style: italic;">{{ .userid }}はすでに存在しています。</p>
                {{/*  型チェックは後でJavaScriptに移行 */}}
                {{ else if eq .reason "agestring" }}
                <p style="color: red;font-style: italic;">Ageには1~3桁の数字を入力してください。</p>                
                {{ else if eq .reason "internalservererror" }}
                <p style="color: red;font-style: italic;">ユーザ登録に失敗しました。</p>                
                {{ end }}
            {{ end }}
            <form method="post" action="/signup"> 
        {{ end }}
        
        {{ if ne .title "SignUp Success!" }}
            <input type="text" name="username" minlength="2" maxlength="30" placeholder="Nickname" required>
            <input type="password" name="password" minlength="1" maxlength="15" placeholder="Password" required>
        {{ end }}
        
        {{ if eq .title "SignUp" }}
            <input type="text" name="age" minlength="1" maxlength="3" pattern="^\d+$" title="1~3桁の数字を入力してください。" placeholder="Age (Optional)">
            <input type="text" name="company" minlength="1" maxlength="50" placeholder="Company (Optional)">
            <input type="text" name="role" minlength="1" maxlength="50" placeholder="Role (Optional)">
        {{ end }}
        
        {{ if ne .title "SignUp Success!" }}
            <button type="submit">{{ .title }}</button>
        </form>
        {{ end }}
        
        {{ if eq .title "SignUp Success!" }}
        <a href="/login">Proceed to Login</a>
        {{ else }}
        <a href="/">Back to Home</a>
        {{ end }}
    </div>
</body>
</html>
