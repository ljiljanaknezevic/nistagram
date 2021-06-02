$(document).ready(function(){

    input_email = $('#id_email');
    let email;
    var btnForgot = document.getElementById("id_button")
    btnForgot.disabled = true

    input_email.keyup(function () {
        if(validateEmail(input_email.val())) {
            btnForgot.disabled = false
        }
        if(!validateEmail(input_email.val())){
            btnForgot.disabled = true
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
    $('#id_button').click(function() {
        email = $('#id_email').val()
        console.log(email)
        customAjax({
            method:'POST',
            url:'http://localhost:80/user-service/sendEmailForAccountRecovery',
            data: JSON.stringify({email : email}),
            contentType: 'application/json',
            success: function(){
                localStorage.setItem('emailAddress', email);
                alert("Success sent email!")
                setInterval(timer, 60000);

            },
            error: function(){
                alert("User with that email doesn't exist")
            }
        });

    });


});
function validateEmail(email) {
    const re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    return re.test(String(email).toLowerCase());
}
function timer() {
    localStorage.removeItem('emailAddress')
}