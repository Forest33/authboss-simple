<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">

	<title>abserver - login</title>

	<link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap.min.css" />
	<link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/font-awesome/4.3.0/css/font-awesome.min.css" />

	<script type="text/javascript" src="https://code.jquery.com/jquery-2.1.3.min.js"></script>
	<script type="text/javascript" src="//maxcdn.bootstrapcdn.com/bootstrap/3.3.2/js/bootstrap.min.js"></script>

	<style>
	.form-signin {
      max-width: 330px;
      padding: 15px;
      margin: 0 auto;
    }
    .form-signin .form-signin-heading,
    .form-signin .form-control {
      position: relative;
      box-sizing: border-box;
      height: auto;
      padding: 10px;
      font-size: 16px;
      margin-bottom: 10px;
    }
    .form-signin .form-control:focus {
      z-index: 2;
    }
    .form-signin input[type="text"] {
      margin-bottom: 10px;
      border-bottom-right-radius: 0;
      border-bottom-left-radius: 0;
    }
    .form-signin input[type="password"] {
      margin-bottom: 10px;
      border-top-left-radius: 0;
      border-top-right-radius: 0;
    }
	</style>
</head>
<body class="container-fluid" style="padding-top: 15px;">
    <div class="container" style="width:400px;padding-top:50px;">
        {{with .error}}<div class="alert alert-danger">{{.}}</div>{{end}}
        <form class="form-signin" method="POST" action="/login">
            <label for="inputUser" class="sr-only">Логин</label>
            <input type="text" name="username" class="form-control" placeholder="Логин" value="" required autofocus>
            <label for="inputPassword" class="sr-only">Пароль</label>
            <input type="password" name="password" id="inputPassword" class="form-control" placeholder="Пароль" value="" required>
            <input type="hidden" name="csrf_token" value="YD2gp10XtC3Ft26c3C/m2bAcCPokpJQ3ZLs9zLU9MOy71Ngn/aRDLVJXiCGtro&#43;Jdan5Cos&#43;nNXYlsFSSIsYHA==" />
            <button class="btn btn-lg btn-primary btn-block" type="submit">Войти</button>
        </form>
    </div>
</body>
</html>