{{ define "header" }}
<div class="header-content">
    <a href=/><img src="/static/images/santa.png" alt="logo"></a>
    <a href=/ style="text-decoration: none;color: black"><h1>TechCareer Talk</h1></a>
    <nav>
        <div style="padding-right: 15px">
            <a href="/about">About Me</a> <!--このサイトについて、なぜこのようなサイトを作ったのかなど-->
        </div>
        <div style="padding-right: 100px">
            <a href="/blog">Blog</a> <!--このサイトのサーバ構成や個人のブログとして使う-->
        </div>
            {{/* <input type="text" placeholder="Search this site..."> */}}
        {{ if .IsLoggedIn }}
        <div>
            <span class="username">{{ .Username }}さん</span>
        </div>
        <div>
            <a href="/logout" class="auth-link">Logout</a>
        </div>
        {{ else }}
        <div>
            <a href="/signup" class="auth-link">SignUp</a>
        </div>
        <div>
            <a href="/login" class="auth-link">Login</a>
        </div>
        {{ end }}
    </nav>
</div>
{{ end }}