$(document).ready(function() {

    $('#showRequests').click(function () {
        console.log("USAO U KLIK ZA VER")
        customAjax({
            url: 'http://localhost:80/user-service/getAllRequests',
            method: 'GET',
            success: function (data) {
                console.log(data)
                showAllRequests(data)

            }, error: function () {
            }
        });
    });
    $('#logout').click(function(){
        localStorage.removeItem('jwt')
        localStorage.removeItem('email')
        localStorage.removeItem('role')
        location.href = "/";
    });
});

let showAllRequests = function (data) {
    var json = JSON.parse(data);
    var pomocna ="";
    pomocna +=`<div style="margin-top: 50px" ><div class="ui cards">`;
    for( i in json) {

        pomocna += `<br><div class="ui card">
  <div class="content">
    <div class="header">` + json[i].email + `</div>
     <div class="meta">Sent: 
        ` + json[i].CreatedAt.split("T")[0] + `
      </div>
       <div class="description">Full name:
        <b>` + json[i].fullName + `</b>
      </div>
      <div class="description" style="color: cornflowerblue">Category: 
        <b>` + json[i].category + `</b>
      </div>
      <div class="image"> 
        <img id="output" height="150px" alt="slika" src ="`+json[i].image+`">
      </div>
  </div>
  <div class="extra content">
   <div class="ui two buttons">
  
        <button class="ui basic green button" name="approve" id="`+json[i].email+`">Approve</button>
        <button class="ui basic red button" name="decline" id="`+json[i].email+`">Decline</button>
      </div>
  </div>
</div>`;


    }
    pomocna+=`</div></div>`;
    $("#showData").html(pomocna);
    $("button[name=approve]").click(function () {
        console.log(this.id)
        customAjax({
            url: 'http://localhost:80/user-service/acceptVerification/' + this.id ,
            method: 'POST',
            success: function () {
                location.href="adminHomePage.html";
            },
            error: function () {
            }
        })
    })
    $("button[name=decline]").click(function () {
        customAjax({
            url: 'http://localhost:80/user-service/declineVerification/' + this.id ,
            method: 'POST',
            success: function () {
                location.href="adminHomePage.html";
            },
            error: function () {
            }
        })
    })

}
