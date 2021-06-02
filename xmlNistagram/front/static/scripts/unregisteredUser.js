$(document).ready(function(e) {
    $("#search").click(function () {
        console.log("usao u klik")
        var username= $("#userSearch").val();
        customAjax({
            url: 'http://localhost:80/search-service/searchUserByUsername/' + username,
            method: 'GET',
            success: function (data) {
                showProfile(data);
                console.log(data)
            },
            error: function () {
            }

        });

    });
})
let showProfile = function(user) {
    var json = JSON.parse(user);
    console.log(json)
    var pomocna ="";
    pomocna +=`<div style="margin-top: 50px" 
        class="ui container"><div class="ui link cards">`;
    for( i in json)
    {
        var pom = '';
        if (json[i].gender == "female") {
            pom += "<img src=\"https://avataaars.io/?avatarStyle=Transparent&topType=LongHairStraightStrand&accessoriesType=Blank&hairColor=BrownDark&facialHairType=Blank&clotheType=BlazerShirt&eyeType=Happy&eyebrowType=Default&mouthType=Smile&skinColor=Light\">";
        } else {
            pom += "<img src='https://avataaars.io/?avatarStyle=Transparent&topType=ShortHairShortFlat&accessoriesType=Blank&hairColor=BrownDark&facialHairType=BeardLight&facialHairColor=BrownDark&clotheType=BlazerShirt&eyeType=Happy&eyebrowType=Default&mouthType=Smile&skinColor=Light'\n" +
                "/>"
        }
        pomocna += `<br><div class="ui card">
  <div class="image">` + pom + `
  </div>
  <div class="content">
    <a class="header">` + json[i].username + `</a>
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
  </div>
</div>`;

    }
    pomocna+=`</div></div>`;
    $("#showData").html(pomocna);
    }


