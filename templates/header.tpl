{{ define "header" }}
<div class="header-content">
    <img src="/static/images/santa.png" alt="logo">
    <h1>TechCareer Talk</h1>
    <nav>
        <a href="/about">About</a> <!--このサイトについて、なぜこのようなサイトを作ったのかなど-->
        <a href="/blog">Blog</a> <!--このサイトのサーバ構成や個人のブログとして使う-->
        {{/* <input type="text" placeholder="Search this site..."> */}}
        {{ if .IsLoggedIn }}
            <span class="username">{{ .Username }}さん</span>
            <a href="/logout" class="auth-link">Logout</a>
        {{ else }}
            <a href="/signup" class="auth-link">SignUp</a>
            <a href="/login" class="auth-link">Login</a>
        {{ end }}
    </nav>
</div>
{{ end }}