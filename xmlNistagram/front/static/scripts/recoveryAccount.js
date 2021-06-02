$(document).ready(function(){
    input_password = $('#id_password');
    input_passwordConf = $('#id_passwordConf');
    var btnChange = document.getElementById("btnChange")
    console.log(btnChange)
    btnChange.disabled = true

    input_password.keyup(function () {
        if(!validatePassword(input_password.val())){
            btnChange.disabled = true
            $(this).addClass(`alert-danger`);
            $('#id_password').css('border-color', 'red');
            $("#errorPassword").text("Password must have at least 8 characters, lower case, upper case, digit, special character!")
            $('#errorPassword').css('color', 'red');
        }else {
            $(this).removeClass(`alert-danger`);
            $('#id_password').css('border-color', '');
            $("#errorPassword").text("")
        }
    });
    input_passwordConf.keyup(function () {
        if(input_password.val()!=input_passwordConf.val()){
            btnChange.disabled = true
            $(this).addClass(`alert-danger`);
            $('#id_passwordConf').css('border-color', 'red');
            $("#errorPasswordConf").text("Passwords must match!")
            $('#errorPasswordConf').css('color', 'red');
        }else {

            $(this).removeClass(`alert-danger`);
            $('#id_passwordConf').css('border-color', '');
            $("#errorPasswordConf").text("")
            btnChange.disabled = false;
        }
    });
    $('#btnChange').click(function() {
        console.log("usao u klik")
        var newPassword = $('#id_password').val()
        var confirmPassword = $('#id_passwordConf').val()
        var email = localStorage.getItem('emailAddress')
        obj = JSON.stringify({email:email,newPass:newPassword,confirmPass:confirmPassword});
        customAjax({
            method:'POST',
            url:'http://localhost:80/user-service/changePassword',
            data : obj,
            contentType: 'application/json',
            success: function(){
                localStorage.removeItem('emailAddress');
                alert("Success changed password!")
            },
            error: function() {
                if (localStorage.getItem('emailAddress') == null) {
                    $('#id_p_log').text('Validation link expired. Try again!')
                } else {
                    $('#id_p_log').text("User with that email doesn't exist")
                }
                localStorage.removeItem('emailAddress');
            }

        });

    });

});
function validatePassword(password) {

    var strongRegex = new RegExp("^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#\$%\^&\*])(?=.{8,})");
    if(password.match(strongRegex)) {
        return true;
    }
    else {
        return false;
    }
}