$(document).ready(function(){
	$.fn.exists = function () {
		return this.length !== 0;
	};

    $('#login_form').submit(ajaxsubmit_v(precheck,respond));
	$('#signup_form').submit(ajaxsubmit_v(precheck,respond));
	$('#admin_login_form').submit(ajaxsubmit_v(precheck,respond));
	$('#add_form').submit(ajaxsubmit_vv(respond_vv));
	$('#return_form').submit(ajaxsubmit_vv(respond_vv));
	$('form[id^=isbn]').submit(ajaxsubmit_vv(respond_v));
	function ajaxsubmit_v (x,f) {
		return	function (){
			$("#reply").html("");
			$("#insert").html("");
			var b = x();
			if (b){
				$.ajax({
					type: $(this).attr("method"),
					url: $(this).attr("action"),
					data: $(this).serialize(), // serializes the form's elements.
					success: f
				});
				return false;
			} else{
				return false; // avoid to execute the actual submit of the form.
			}
		};
	}
	function ajaxsubmit_vv (f) {
		return	function (){
			$("#reply").html("");
			$("#insert").html("");
				$.ajax({
					type: $(this).attr("method"),
					url: $(this).attr("action"),
					data: $(this).serialize(), // serializes the form's elements.
					success: f
				});
			return false;
		};
	}
	function jump(id,count) {
        window.setTimeout(function(){
            count--;
            if(count > 0) {
                $('#'+id).html(count+" 秒后将转到首页");
                jump(id,count);
            } else {
				if ($("form[id*=admin]").exists()){
					location.href="/admin/";
				}else{
					location.href="/";
				}
            }
        }, 1000);
    }
	function precheck() {
		var uname = $("#username").val();
		if(uname==""){
			$("#msg").html("姓名不能为空！");
			return false;
		}
		var password = $("#password").val();
		if(password==""){
			$("#msg").html("密码不能为空！");
			return false;
		}
		$("#msg").html("正在提交...");
		return true;
	}
	function respond(responseText)  {
		if (responseText.result == 1){
			jump('msg',5);
		}else if(responseText.result == 0) {
			$('#msg').html("您的输入不合法，操作失败");
		}else  {
			$('#msg').html("出现内部错误");
		}
	}
	function respond_vv(responseText, statusText)  {
		$("#insert").html(responseText.txt);
	}
	function respond_v(responseText, statusText)  {
		$("#reply").html(responseText.txt);
	}
});
