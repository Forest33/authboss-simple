{{define "pagetitle"}}abserver - index{{end}}

{{$logged := .logged}}
{{if $logged}}
<div class="row" style="margin-bottom: 20px;">
	<div class="col-md-12">
		index page
	</div>
</div>
{{else}}
    <div class="container" style="width:400px;padding-top:50px;">
        <form class="form-signin" method="POST" action="/login">
            <label for="inputUser" class="sr-only">Логин</label>
            <input type="text" name="username" class="form-control" placeholder="Логин" value="" required autofocus>

            <label for="inputPassword" class="sr-only">Пароль</label>
            <input type="password" name="password" id="inputPassword" class="form-control" placeholder="Пароль" value="" required>
            <input type="hidden" name="csrf_token" value="{{.csrf_token}}" />

            <button class="btn btn-lg btn-primary btn-block" type="submit">Войти</button>
        </form>
    </div>
{{end}}
