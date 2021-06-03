var input_email;
var input_password;
var input_password_repeat;
var button_register;

var p_log;

$(document).ready(function(e){
    input_name=$('#id_name');
    input_username=$('#id_username');
    input_email = $('#id_email');
    input_password = $('#id_password');
    input_password_repeat = $('#id_password_repeat');
    var btnRegister = document.getElementById("id_button")
    btnRegister.disabled = true

    input_name.keyup(function () {
        if(validateEmail(input_email.val()) && validatePassword(input_password.val()) && validateName(input_name.val()) && validateUsername(input_username.val()) ) {
            btnRegister.disabled = false
        }
        if(!validateName(input_name.val())){
            btnRegister.disabled = true
            $(this).addClass(`alert-danger`);
            $('#id_name').css('border-color', 'red');
            $("#errorName").text("You can only use letters for full name!")
            $('#errorName').css('color', 'red');
        }else {
            $(this).removeClass(`alert-danger`);
            $('#id_name').css('border-color', '');
            $("#errorName").text("")
        }
    });

    input_username.keyup(function () {
        if(validateEmail(input_email.val()) && validatePassword(input_password.val()) && validateName(input_name.val()) && validateUsername(input_username.val())  ) {
            btnRegister.disabled = false
        }
        if(!validateUsername(input_username.val())){
            btnRegister.disabled = true
            $(this).addClass(`alert-danger`);
            $('#id_username').css('border-color', 'red');
            $("#errorUsername").text("You can only use letters and numbers for username!")
            $('#errorUsername').css('color', 'red');
        }else {
            $(this).removeClass(`alert-danger`);
            $('#id_username').css('border-color', '');
            $("#errorUsername").text("")
        }
    });

    input_email.keyup(function () {
        if(validateEmail(input_email.val()) && validatePassword(input_password.val()) && validateName(input_name.val()) && validateUsername(input_username.val()) ) {
            btnRegister.disabled = false
        }
        if(!validateEmail(input_email.val())){
            btnRegister.disabled = true
            $(this).addClass(`alert-danger`);
            $('#id_email').css('border-color', 'red');
            $("#errorEmail").text("Email is in wrong format!")
            $('#errorEmail').css('color', 'red');
        }else {
            $(this).removeClass(`alert-danger`);
            $('#id_email').css('border-color', '');
            $("#errorEmail").text("")
        }
    });

    input_password.keyup(function () {
        if(validateEmail(input_email.val()) && validatePassword(input_password.val()) && validateName(input_name.val()) && validateUsername(input_username.val()) ) {
            btnRegister.disabled = false
        }
        if(!validatePassword(input_password.val())) {
            btnRegister.disabled = true
            $(this).addClass(`alert-danger`);
            $('#id_password').css('border-color', 'red');
            $("#errorPassword").text("Password must have at least 8 characters, lower case, upper case, digit, special character!")
            $('#errorPassword').css('color', 'red');
        } else {
            $(this).removeClass(`alert-danger`);
            $('#id_password').css('border-color', '');
            $("#errorPassword").text("")

        }
    });
    input_password_repeat.keyup(function () {
        if(input_password.val()!=input_password_repeat.val()){
            btnRegister.disabled = true
            $("#errorPasswordRepeat").text("Passwords do not match!")
            $(this).addClass(`alert-danger`);
            $('#id_password_repeat').css('border-color', 'red');
            $('#errorPasswordRepeat').css('color', 'red');
        }
        else {
            btnRegister.disabled = false
            $(this).removeClass(`alert-danger`);
            $('#id_password_repeat').css('border-color', '');
            $("#errorPasswordRepeat").text("")
        }


    });
    button_register = $('#id_button');
    p_log = $('#id_p_log');

    button_register.on('click', function(e){
        var isPrivate = document.getElementById("isPrivate").checked;
        var name=input_name.val();
        var username=input_username.val();
        var email = input_email.val();
        var password = input_password.val();

        obj = JSON.stringify({
            name:name,
            username:username,
            email:email,
            password:password,
            isPrivate:isPrivate
        });

        customAjax({
            url: 'http://localhost:80/user-service/confirmRegistration',
            method: 'POST',
            data:obj,
            contentType: 'application/json',
            success: function(){
                p_log.text('')
                localStorage.setItem("obj" , obj)
                alert("Sucess sent email.")
                setInterval(timer, 60000);

            },
            error: function(){
                p_log.text('Error');
            }
        });


    });


    //[a-zA-Z]+
    function validateName(name) {
        const re = /^[a-zA-Z]+[a-zA-Z\s]*$/;
        return re.test(String(name));
    }

    function validateUsername(name) {
        const re = /^[a-zA-Z0-9]+$/;
        return re.test(String(name));
    }

    function validateEmail(email) {
        const re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
        return re.test(String(email).toLowerCase());
    }

    function validatePassword(password) {

        var strongRegex = new RegExp("^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#\$%\^&\*])(?=.{8,})");
        if(password.match(strongRegex)) {
            return true;
        }
        else {
            return false;
        }
    }

    function sanatize(input) {
        var output = input.replace(/<script[^>]*?>.*?<\/script>/gi, '').
        replace(/<[\/\!]*?[^<>]*?>/gi, '').
        replace(/<style[^>]*?>.*?<\/style>/gi, '').
        replace(/<![\s\S]*?--[ \t\n\r]*>/gi, '');
        return output;
    };


    function timer() {
        localStorage.removeItem('obj')
    }

});