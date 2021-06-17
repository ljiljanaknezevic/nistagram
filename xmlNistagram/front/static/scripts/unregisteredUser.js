$(document).ready(function(e) {
    $("#search").click(function () {
        console.log("usao u klik")
        var username= $("#userSearch").val();
        customAjax({
            url: 'http://localhost:80/search-service/searchUserByUsernameForUnregistredUser/' + username,
            method: 'GET',
            success: function (data) {
                showProfile(data);
                console.log(data)
            },
            error: function () {
            }

        });

    });
    $("#searchTags").click(function () {
        var tag= $("#tagsSearch").val();
        console.log(tag)
        customAjax({
            url: 'http://localhost:80/search-service/searchPostByTagUnregistered/' + tag ,
            method: 'GET',
            success: function (data) {
                showPosts(data);
            },
            error: function () {
            }

        });

    });
    $("#searchLocation").click(function () {
        var location = $("#locationSearch").val();
        console.log(location)
        customAjax({
            url: 'http://localhost:80/search-service/searchPostByLocationUnregistered/' + location,
            method: 'GET',
            success: function (data) {
                showPosts(data);
            },
            error: function () {
            }

        });
    })
})
let showProfile = function(user) {
    var json = JSON.parse(user);
    var pomocna ="";

    pomocna +=`<div style="margin-top: 50px" ><div class="ui link cards">`;
    for( i in json) {
        var verified = ""
        if(json[i].isVerified) {
            verified += `<i class="check circle icon" style="color: dodgerblue"></i>`
        } else {
            verified = ``
        }

        var pom = '';
        if (json[i].gender == "female") {
            pom += "<img src=\"https://avataaars.io/?avatarStyle=Transparent&topType=LongHairStraightStrand&accessoriesType=Blank&hairColor=BrownDark&facialHairType=Blank&clotheType=BlazerShirt&eyeType=Happy&eyebrowType=Default&mouthType=Smile&skinColor=Light\">";
        } else {
            pom += "<img src='https://avataaars.io/?avatarStyle=Transparent&topType=ShortHairShortFlat&accessoriesType=Blank&hairColor=BrownDark&facialHairType=BeardLight&facialHairColor=BrownDark&clotheType=BlazerShirt&eyeType=Happy&eyebrowType=Default&mouthType=Smile&skinColor=Light'\n" +
                "/>"
        }

        pomocna += `<br><div class="ui card" >
  <div class="image">` + pom + `
  </div>
  <div class="content">
    <a class="header" name="profile" id="`+json[i].email+`">` + json[i].username + ``+verified+`</a>
    <div class="meta">
      <span class="date">Birthday: ` + json[i].birthday + `</span>
    </div>
    <div class="description">
     Biography:  ` + ((json[i].biography != '') ? json[i].biography : `-`) + `
    </div>
     <div class="description">Website:   
     <a href="` + ((json[i].website != '') ? json[i].website : `-`) + `">
     My website
    </a>
    </div>
  </div>
  <div class="extra content">
    <div class="right floated author">` + json[i].name + `
    </div>   
    <div id="error` + json[i].username + `" style="color:red"></div> 
  </div>
  
</div>`;

    }

    pomocna+=`</div></div>`;

    $("#showData").html(pomocna);
    $("a[name=profile]").click(function () {
        console.log("Usao u profile")
        id = this.id
        console.log(id)
        customAjax({
            url: 'http://localhost:80/search-service/getPostsForSearchedUserUnregistered/' + this.id ,
            method: 'GET',
            success: function (data) {
                showPosts(data)
            },
            error: function () {
                alert("Locked")
            }
        })
    })
    }
let showPostsForSearchedUser = function(posts) {
    var json = JSON.parse(posts)
    console.log(json)
    var slika
    result=""
    result +=`<div style="margin-top: 50px" ><div class="ui cards">`;

    for( i in json) {

        slika = ""
        customAjax({
            url: 'http://localhost:80/search-service/getMedia/' + json[i].ImageID,
            method: 'GET',
            async:false,
            success: function (data) {

                slika = data
                console.log(slika)
                console.log("slka slika slika" + slika)

                if (slika.type == "video") {

                    pom1 = `<video id="output" height="150px" alt="slika" autoplay src ="` + 'data:video/mp4;base64,' + slika.path + ` ">`;
                } else {
                    pom1 = `<img id="output" height="150px" alt="slika" src ="` + 'data:image/png;base64,' + slika.path + ` ">`;
                }
                result += `<br><div class="ui card">

  <div class="content">
     <div class="right floated meta">` + json[i].CreatedAt.split("T")[0] + `</div>  
  </div>
   <div class="image">
    ` + pom1 + `
  </div>
  <div class="content">
  <div class="description">
      `+ json[i].description+`
    </div>
  </div> 
  <div class="extra content">
    <div class="ui large transparent left icon input">
      <i class="heart outline icon"></i>
      <input type="text" placeholder="Add Comment...">
    </div>
  </div>
</div>`;

            },
            error: function () {


            }
        })


    }

    result+=`</div></div>`;
    $("#showData").html(result);



}



var pomocna

let showPosts = function(posts) {
    var json = JSON.parse(posts);
    var jsonParse;
    console.log(json)
    var slika
    pomocna=""
    pomocna +=`<div style="margin-top: 50px" ><div class="ui cards">`;

    for( i in json) {

        slika = ""
        customAjax({
            url: 'http://localhost:80/search-service/getMedia/' + json[i].ImageID,
            method: 'GET',
            async:false,
            success: function (data) {
                customAjax({
                    url: 'http://localhost:80/user-service/getByEmail/' + json[i].email,
                    method: 'GET',
                    async:false,
                    success: function (user) {
                        jsonParse = JSON.parse(user);
                        console.log(jsonParse.username)
                    },
                    error: function () {

                    }
                })
                slika = data
                console.log(slika)
                console.log("slka slika slika" + slika)

                if (slika.type == "video") {

                    pom1 = `<video id="output" height="150px" alt="slika" autoplay src ="` + 'data:video/mp4;base64,' + slika.path + ` ">`;
                } else {
                    pom1 = `<img id="output" height="150px" alt="slika" src ="` + 'data:image/png;base64,' + slika.path + ` ">`;
                }
                pomocna += `<br><div class="ui card">

  <div class="content">
     <div class="left floated meta">` + jsonParse.username + `</div>
     <div class="right floated meta">` + json[i].CreatedAt.split("T")[0] + `</div>
         <br>
    <div class="description" style="color:cornflowerblue"><i class="location arrow icon" ></i>
      `+ json[i].Location+`
    </div>
    
  </div>
  
   <div class="image">
    ` + pom1 + `
  </div>
  <div class="content">
  <div class="description">
      `+ json[i].description+`
    </div>
     <br>
    <div class="description"><i class="tags icon"></i>
      `+ json[i].tags+`
    </div>
  </div>
  
  
</div>`;

            },
            error: function () {

            }
        })


    }

    pomocna+=`</div></div>`;
    $("#showData").html(pomocna);



}


