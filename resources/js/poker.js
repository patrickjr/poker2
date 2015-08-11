

















var validate_email = function(email) {
    var re = /^([\w-]+(?:\.[\w-]+)*)@((?:[\w-]+\.)*\w[\w-]{0,66})\.([a-z]{2,6}(?:\.[a-z]{2})?)$/i;
    return re.test(email);
};

var get_screen_xy = function(){
    return [screen.height, screen.width];
};

var get_browser_xy = function(){
    var w = window,
        d = document,
        e = d.documentElement,
        g = d.getElementsByTagName('body')[0],
        x = w.innerWidth || e.clientWidth || g.clientWidth,
        y = w.innerHeight|| e.clientHeight|| g.clientHeight;
        return [x,y];
};

var toggle = function (obj) {
    var el = document.getElementById(obj);
    el.style.display = (el.style.display !== 'none' ? 'none' : '' );
};

var toggle_visibility = function(id){
    if ( $(id).css('visibility') === 'hidden' )
        $(id).css('visibility','visible');
    else
        $(id).css('visibility','hidden');
};

var toggleMe = function(x, y){
	$(x).on('click', function(e){
	    e.preventDefault();
	    $(this).next(y).toggle();
	});
};

var overlay = function (id) {
	el = document.getElementById(id);
	el.style.visibility = (el.style.visibility == "visible") ? "hidden" : "visible";
}


var registerSubmit = function(){	
	$('#reg_form').on('submit', function(e){
		e.preventDefault();
		var email = $('#new_user_email').val();
		if(!validate_email(email)){
			return false;
		}
		toggle_visibility("#loader_here");
		toggle_visibility("#submit_here");
		$.ajax({
			url      : 'register',
			type     : 'post',
			data     : $('#reg_form').serialize(),
			success  : function(data){
				alert("ajax working");
			},
			error: function (xhr, textStatus, err) {
				var err_msg = "* " + xhr.responseText;
				$('#email_error').html(err_msg);
				toggle_visibility("#loader_here");
				toggle_visibility("#submit_here");
			}
		});
	});
};

var loginSubmit = function(){
	$('#login_form').on('submit', function(e){
		e.preventDefault();
		toggle_visibility("#loader_login");
		toggle_visibility("#login_button");
		$.ajax({
			url      : 'login',
			type     : 'post',
			data     : $('#login_form').serialize(),
			dataType : 'html',
			success  : function(data){
				$('#login_register').fadeOut(1000, function() { $('#login_register').remove(); })
				$('#user_nav_info').html(data);
			},
			error: function (xhr, textStatus, err) {
				var err_msg = "* " + xhr.responseText;
				$('#login_error').html(err_msg);
				toggle_visibility("#loader_login");
				toggle_visibility("#login_here");
			}

		}); 
	});
};







































