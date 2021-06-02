$(document).ready(function(e){
	var email = localStorage.getItem('email')
	
	  $("#editProfile").click(function () {
          console.log("usao u klik")
		  customAjax({
		      url: 'http://localhost:80/user-service/getByEmail/' + email,
		      method: 'GET',
		      success: function(data){
                editProfile(data)
		      },
		      error: function(){
		      }
	
	 });
	  });  
	
	$('#logout').click(function(){
		localStorage.removeItem('jwt')		
		location.href = "/";
		});
});

let editProfile = function(user) {
    var json = JSON.parse(user);
    var gender = ``;
    if (json.gender == "female") {
        gender = `<input type="radio"  name="gender" value="male"> Male
               <input type="radio"  name="gender" value="female" checked="checked"> Female`;
    } else {
        gender = `<input type="radio"  name="gender" value="male" checked="checked"> Male
               <input type="radio"  name="gender" value="female"> Female`;
    }
    $("#showData").html(`<table class="ui large table" style="width:50%; margin-left:auto; 
               margin-right:auto; margin-top: 40px;">
                           <thead>
                               <tr class="success">
                                   <th colspan="2" class = "text-info" style= "text-align:center;">Edit profile</th>
                               </tr>
                           </thead>
                           <tbody>
                               <tr>
                                   <td>Email:</td>
                                   <td class="ui input small"> <input type="text" id="txtEmail" disabled="disabled" value="`+ ((json.email != null) ? json.email:`` ) + `"/></td>
                                 
                               </tr>
                               <tr>
                                   <td>Name:</td>
                                   <td class="ui input small"> <input type="text" id="txtName" value="`+ ((json.name != null) ? json.name:`` ) + `"/></td>
                                 
                               </tr>
                               <tr>
                                   <td>Username:</td>
                                   <td class="ui input small"> <input type="text" id="txtUsername" value="`+ ((json.username != null) ? json.username:`` ) + `"/></td>
                                 
                               </tr>
                               <tr>
                               <td>Phone number:</td>
                               <td class="ui input small"> <input type="text" id="txtPhone" value="`+ ((json.phone != null) ? json.phone:`` ) + `"/></td>
                             
                           </tr>
                           <tr>
                           <td>Birthday:</td>
                           <td class="ui input small"> <input type="date" id="txtBirthday" value="`+ ((json.birthday != null) ? json.birthday:`` ) + `"/></td>
                         
                       </tr>
                       <tr>
                           <td>Website:</td>
                           <td class="ui input small"> <input type="text" id="txtWebsite" value="`+ ((json.website != null) ? json.website:`` ) + `"/></td>
                         
                       </tr>
                       <tr>
                       <td>Biography:</td>
                       <td class="ui input small"> <input type="text" id="txtBiography" value="`+ ((json.biography != null) ? json.biography:`` ) + `"/></td>
                     
                   </tr>
                   <tr>
                   <td>Gender:</td>
                   <td>`
                    + gender +
                    `</td>
                        </tr>
                           
                           </tbody>
                           <tfoot class="full-width">
               <tr>
                 <th></th>
                 <th colspan="2">
                      <input id = "acceptChange" class="ui center floated positive basic button" type = "button" value = "Accept changes"></input>
               
                 </th>
               </tr>
             </tfoot>
                       </table> <p id="er"> </p>`);

        $('#acceptChange').click(function(){

        let email=$('#txtEmail').val()
        let name=$('#txtName').val()
        let username=$('#txtUsername').val()
        let phone=$('#txtPhone').val()
        let birthday=$('#txtBirthday').val()
        let website=$('#txtWebsite').val()
        let biography=$('#txtBiography').val()
        let gender = $("input:radio[name=gender]:checked").val();
        console.log(website)

        obj = JSON.stringify({
        email:email,
        name:name,
        username: username,
        phone:phone,
        birthday:birthday,
        website :website,
        biography:biography,
        gender:gender
        });
        console.log(obj)
        
            customAjax({
        url: 'http://localhost:80/user-service/changeUserData',
        method: 'POST',
        data:obj,
        contentType: 'application/json',
            success: function(){
                alert("Sucess.")
                location.href = "userHomePage.html";
                
            },
                error: function(){
                    alert('Error');
                }
    });
});
                        
}