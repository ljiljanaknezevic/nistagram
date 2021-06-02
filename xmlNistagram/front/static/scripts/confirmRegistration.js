$(document).ready(function(e) {
    button_confirm = $('#id_button');
    p_log = $('#id_p_log');

    button_confirm.on('click', function(e) {
        customAjax({
            url: 'http://localhost:80/user-service/signup',
            method: 'POST',
            data: localStorage.getItem("obj"),
            contentType: 'application/json',
            success: function () {
                p_log.text('')
                if(localStorage.getItem('obj') == null) {
                    alert('Validation link expired. Try again!')
                } else {
                    alert("Sucess registration.")
                }
            },
            error: function () {
                p_log.text('Error');
            }
        });
    });
});