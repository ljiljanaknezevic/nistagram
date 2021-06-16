let showRequest = function (email) {


        $("#showData").html(`<table class="ui large table" style="width:50%; margin-left:auto; 
               margin-right:auto; margin-top: 40px;">
                           <thead>
                               <tr class="success">
                                   <th colspan="2" class = "text-info" style= "text-align:center;">Verification request</th>
                               </tr>
                           </thead>
                           <tbody>
                               <tr>
                                   <td>Email:</td>
                                   <td class="ui input small"> <input type="text" id="txtEmail" disabled="disabled" value="` + ((email != null) ? email : ``) + `"/></td>                       
                               </tr>
                               <tr>
                                   <td>Full name:</td>
                                   <td class="ui input small"> <input type="text" id="txtFullName"/>
                                   <p id="errorName"></p></td>
                                    
                               </tr>
                              
                               <tr>
                                       <td>Category:</td>
                                       <td>
                                       <select class="ui search dropdown" id="category">
                                       <option value="">Category...</option>
                                       <option value="influencer">Influencer</option>
                                       <option value="sports">Sports</option>
                                       <option value="new">New/Media</option>
                                       <option value="business">Business</option>
                                       <option value="brand">Brand</option>
                                       <option value="organization">Organization</option>
                                       </select>
                                       <p id="errorCategory"></p>
                                        </td>
                               </tr>
                               <tr>
                               <td></td>
                               <td>
                               <div class="field">
                                <label for="file">Choose image off official document:</label>
                                <input type="file" id="file" name="file"  multiple required onchange="readURL(this);" accept=".jpg, .jpeg, .png"  >
                                
                            </div>
                             <div class="field">
                                    <img id="blah" height="150px" alt="your image" />
                                </div>
                                <p id="errorImage"></p>
                                </td>
                               </tr>
                     
                           </tbody>
                           <tfoot class="full-width">
               <tr>
                 <th></th>
                 <th colspan="2">
                      <input id = "sendRequest" class="ui right floated positive basic button" type = "button" value = "Send request"></input>
               
                 </th>
               </tr>
             </tfoot>
                       </table> <p id="er"> </p>`);
    input_name = $('#txtFullName');
    var btnSend = document.getElementById("sendRequest")
    btnSend.disabled = true
    input_name.keyup(function () {
        if(validateFullName(input_name.val()) && $('#category').val()!="") {
            btnSend.disabled = false
        }
        if(!validateFullName(input_name.val())){
            btnSend.disabled = true
            $(this).addClass(`alert-danger`);
            $('#txtFullName').css('border-color', 'red');
            $("#errorName").text("You can only use letters for full name!")
            $('#errorName').css('color', 'red');
        }else {
            $(this).removeClass(`alert-danger`);
            $('#txtFullName').css('border-color', '');
            $("#errorName").text("")
        }
    });
    $('#category').on('change',function(e){
        if(validateFullName(input_name.val()) && $('#category').val()!="") {
            btnSend.disabled = false
        }
        if($('#category').val()==""){
            btnSend.disabled = true
            $(this).addClass(`alert-danger`);
            $('#category').css('border-color', 'red');
            $("#errorCategory").text("You need to choose one!")
            $('#errorCategory').css('color', 'red');
        }else {
            $(this).removeClass(`alert-danger`);
            $('#category').css('border-color', '');
            $("#errorCategory").text("")
        }

    })
    $('#sendRequest').click(function() {
        var email = $('#txtEmail').val();
        var fullName =  $('#txtFullName').val();
        var category=  $('#category').val();
        var image = $('#blah').attr('src');
        if(image =='undefined' || image==null || image==""){
            $("#errorImage").text("You need to choose one!")
            $('#errorImage').css('color', 'red');
            return;
        }
        else{
            $("#errorImage").text("")
        }

        obj = JSON.stringify({
            email:email,
            fullName:fullName,
            category:category,
            image:image
        });

        customAjax({
            url: 'http://localhost:80/user-service/createRequest',
            method: 'POST',
            data:obj,
            contentType: 'application/json',
            success: function(){

                alert("Sucess sent request.")

            },
            error: function(){
                alert("You have already created request.")
            }
        });

    })

    }
    function validateFullName(name) {
        const re = /^[a-zA-Z]+[a-zA-Z\s]*$/;
        return re.test(String(name));
    }
