{{define "pagetitle"}}abserver - index{{end}}

{{$logged := .logged}}
{{if $logged}}
<div class="row" style="margin-bottom: 20px;">
	<div class="col-md-12">
		bar page
	</div>
</div>
{{end}}
