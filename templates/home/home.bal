{{define "home"}}
{{template "header"}}
{{template "side"}}
<div class="col-md-9" role="main">
  <div class="panel panel-success">
	<div class="panel-heading">
	  <h2 class="panel-title">Hello User: {{.Who}}</h2>
	</div>
	<div class="panel-body">
	  <p>你借了 {{.Num}} 本书：<p>
		  {{if .Nil}}
		  {{range $idx,$bookid := .Book}}
			<form   role="form"    method="post"  action="goreturn/">
			  <div class="form-group">
				<label for="bid">书ID：</label>
				<label for="Bookid">{{$bookid}}</label>
			  </div>
			  <div class="form-group">
				<label for="bn">书名：</label>
				<label for="Bookname">{{index $.BookName $idx}}</label>
			  </div>
			  <input type="submit" name="borrow" class="btn btn-primary"  value="还书" >
			  <input type="hidden" name="bookid" value={{$bookid}}>
			  <input type="hidden" name="bookname" value={{index $.BookName $idx}}>
			</form>
		  {{end}}
	</div>
  </div>
</div>
{{template  "divfooter"}}
{{template "footer"}}
{{end}}
