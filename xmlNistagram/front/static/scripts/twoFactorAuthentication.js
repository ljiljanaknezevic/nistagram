$(document).ready(function(){
    $('#id_button').on('click', function(e) {
        input = $('#code').val()
        customAjax({
            url: 'http://localhost:80/user-service/validateToken/'+input,
            method: 'GET',
            contentType: 'application/json',
            success: function (jwt, status, xhr) {
                if(xhr.status == 200){
                    localStorage.setItem('email', jwt.email);
                    localStorage.setItem('jwt', jwt.token);
                    localStorage.setItem('role', jwt.role);
                    authentification()}

            },
            error: function () {
                $('#id_p_log').text('Invalid code try again');
            }
        });
    });


})
function authentification(){
    var role =localStorage.getItem('role');
    if( role == 'user') {
        window.location.href = "userHomepage.html";
    }
}