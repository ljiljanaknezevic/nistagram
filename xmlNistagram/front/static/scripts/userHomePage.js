var file;
var pomocnaP;

let jsonObjekat;

  let slika
function readURL(input) {

    if (input.files && input.files[0]) {
        var reader = new FileReader()
        file=input.files[0];

        reader.onload = function(e) {
            console.log(e.target.result)
            $('#blah').attr('src', e.target.result);
        }
        reader.readAsDataURL(input.files[0]); // convert to base64 string
    }

}
let users;
$(document).ready(function(e){

    var email = localStorage.getItem('email')

    
     $("#addStory").click(function () {
         $(`#showData`).html(`<div class="ui buttons">
  <button id="addStoryImage" class="ui button">Add image</button>
  <button id="addStoryVideo" class="ui button">Add video</button>
</div><div id="showDataMedia"></div>`);
         $("#addStoryVideo").click(function () {
             showVideoStory()
         })

         $("#addStoryImage").click(function () {
             customAjax({
                 url: 'http://localhost:80/user-service/getAllUsersExceptLogging/' + email,
                 method: 'GET',
                 async: false,
                 success: function (data) {
                     var json = JSON.parse(data);
                     users = json

                 },
                 error: function () {
                 }
             })


             function reverseGeocode(coords) {
                 fetch('https://nominatim.openstreetmap.org/reverse?format=json&lon=' + coords[0] + '&lat=' + coords[1])
                     .then(function (response) {
                         //alert(response);
                         return response.json();
                     }).then(function (json) {
                     let location = json["address"]["road"] + ` ` + json["address"]["house_number"] + ` , ` + json["address"]["city"] + ` , ` + json["address"]["country"];
                     $('#location').val(location)

                     // $('#street-number').val(json["address"]["house_number"])
                     //$('#city').val(json["address"]["city"])
                     //$('#zip-code').val(json["address"]["postcode"])


                     //$('#location-longitude').val(json["lon"]);
                     //$('#location-latitude').val(json["lat"]);


                     jsonObjekat = json;
                 });
             };

             pomocnaP = function () {
                 var map = new ol.Map({

                     target: 'map',
                     layers: [
                         new ol.layer.Tile({
                             source: new ol.source.OSM()
                         })
                     ],
                     view: new ol.View({
                         center: ol.proj.fromLonLat([19.8424, 45.2541]),
                         zoom: 15
                     })
                 });
                 //var jsonObjekat;
                 map.on('click', function (evt) {
                     var coord = ol.proj.toLonLat(evt.coordinate);
                     reverseGeocode(coord);
                     var iconFeatures = [];
                     var lon = coord[0];
                     var lat = coord[1];
                     var icon = "marker.png";
                     var iconGeometry = new ol.geom.Point(ol.proj.transform([lon, lat], 'EPSG:4326', 'EPSG:3857'));
                     var iconFeature = new ol.Feature({
                         geometry: iconGeometry
                     });

                     iconFeatures.push(iconFeature);

                     var vectorSource = new ol.source.Vector({
                         features: iconFeatures //add an array of features
                     });


                     var iconStyle = new ol.style.Style({
                         image: new ol.style.Icon(/** @type {olx.style.IconOptions} */({
                             anchor: [0.5, 46],
                             anchorXUnits: 'fraction',
                             anchorYUnits: 'pixels',
                             opacity: 0.95,
                             src: icon
                         }))
                     });

                     var vectorLayer = new ol.layer.Vector({
                         source: vectorSource,
                         style: iconStyle
                     });

                     map.addLayer(vectorLayer);

                 });
             }
             let temp = ""
             for (i in users) {
                 temp += `<div class="item" data-value="` + users[i].username + `">` + users[i].username + `</div>`
             }

             $("#showDataMedia").html(
                 `<form  class="ui large form" 
                             style="width:80%; margin-left:auto; 
                             margin-right:auto; margin-top: 20px;">         
                          <form method="post" enctype="multipart/form-data">
                            <div class="field">
                                <label for="file">Choose image:</label>
                                <input type="file" id="file" name="file"  multiple required onchange="readURL(this);" accept=".jpg, .jpeg, .png"  >
                            </div>
                            <div class=" two fields">
                                <div class="field">
                                    <img id="blah" height="350px" alt="your image" />
                                </div>
                                <div class="field">
                                    <label for="location">Location:</label>
                                    <input type="text"  id="location" placeholder="place for location" />
                                    <div id="map" class="map" style="height:420px;"></div>
                                            <script>pomocnaP();</script>
                                </div>
                            </div>
                             <div class="field">
                                    <label for="description">Description:</label>
                                    <textarea type="text"  id="description" name="description" placeholder="Description" rows = "2"/>
                            </div>
                           <div class="tag">

            <div class="ui fluid multiple search selection dropdown">
                <input name="tags" id="tags" type="hidden">
                <i class="dropdown icon"></i>
                <div class="default text">Tags</div>
                <div class="menu">
              ` + temp + `
                </div>
            </div>

        </div>
                            <div class="ui grid">
                            <div class="two wide column"></div>
                            <div class="two wide column"></div>
                            <div class="two wide column"></div>
                            <div class="two wide column"></div>
                            <div class="two wide column"></div>
                            <div class="two wide column"></div>
                            <div class="two wide column"></div>
                            <div class="two wide column right">  
                            <button type="button" style = "text-align: center" class="ui primary button" id="save_story" >ADD STORY</button>
                            </div>
                            </div>
                             <script>
                                    $('.tag .ui.dropdown').dropdown();
                                </script>
                          </form>
                      </form>`
             );

             $('#save_story').click(function () {
                 var formData = new FormData();
                 formData.append("file", file);
                 var description = $('#description').val();
                 var tags = $('#tags').val();
                 var location = $('#location').val();
                 var email = localStorage.getItem('email');
                 formData.append("type","image")

                 formData.append("description", description)
                 formData.append("tags", tags)
                 formData.append("location", location)
                 formData.append("email", email)
                 customAjax({
                     url: 'http://localhost:80/story-service/saveStory',
                     method: 'POST',
                     data: formData,
                     processData: false,
                     contentType: false,
                     success: function () {
                         alert("Sucess saved story")
                         customAjax({
                             url: 'http://localhost:80/user-service/getByEmail/' + email,
                             method: 'GET',
                             success: function (data) {
                                 console.log("Success")
                                 // myProfile(data)
                             },
                             error: function () {
                             }
                         });
                     },
                     error: function (e) {
                         alert('Error uploading new post.')
                     }
                 });
             });
         });
     });

    $("#addPost").click(function () {
        $(`#showData`).html(`<div class="ui buttons">
  <button id="addPostImage" class="ui button">Add image</button>
  <button id="addPostVideo" class="ui button">Add video</button>
</div><div id="showDataMedia"></div>`);
        $("#addPostVideo").click(function () {
            showVideo()
        })


        $("#addPostImage").click(function () {
            customAjax({
                url: 'http://localhost:80/user-service/getAllUsersExceptLogging/' + email,
                method: 'GET',
                async: false,
                success: function (data) {
                    var json = JSON.parse(data);
                    users = json

                },
                error: function () {
                }
            })


            function reverseGeocode(coords) {
                fetch('https://nominatim.openstreetmap.org/reverse?format=json&lon=' + coords[0] + '&lat=' + coords[1])
                    .then(function (response) {
                        //alert(response);
                        return response.json();
                    }).then(function (json) {
                    let location = json["address"]["road"] + ` ` + json["address"]["house_number"] + ` , ` + json["address"]["city"] + ` , ` + json["address"]["country"];
                    $('#location').val(location)

                    jsonObjekat = json;
                });
            };

            pomocnaP = function () {
                var map = new ol.Map({

                    target: 'map',
                    layers: [
                        new ol.layer.Tile({
                            source: new ol.source.OSM()
                        })
                    ],
                    view: new ol.View({
                        center: ol.proj.fromLonLat([19.8424, 45.2541]),
                        zoom: 15
                    })
                });
                //var jsonObjekat;
                map.on('click', function (evt) {
                    var coord = ol.proj.toLonLat(evt.coordinate);
                    reverseGeocode(coord);
                    var iconFeatures = [];
                    var lon = coord[0];
                    var lat = coord[1];
                    var icon = "marker.png";
                    var iconGeometry = new ol.geom.Point(ol.proj.transform([lon, lat], 'EPSG:4326', 'EPSG:3857'));
                    var iconFeature = new ol.Feature({
                        geometry: iconGeometry
                    });

                    iconFeatures.push(iconFeature);

                    var vectorSource = new ol.source.Vector({
                        features: iconFeatures //add an array of features
                    });


                    var iconStyle = new ol.style.Style({
                        image: new ol.style.Icon(/** @type {olx.style.IconOptions} */({
                            anchor: [0.5, 46],
                            anchorXUnits: 'fraction',
                            anchorYUnits: 'pixels',
                            opacity: 0.95,
                            src: icon
                        }))
                    });

                    var vectorLayer = new ol.layer.Vector({
                        source: vectorSource,
                        style: iconStyle
                    });

                    map.addLayer(vectorLayer);

                });
            }
            console.log(users)
            let temp = ""
            for (i in users) {
                temp += `<div class="item" data-value="` + users[i].username + `">` + users[i].username + `</div>`
            }
            $("#showDataMedia").html(
                `
     <form  class="ui large form" 
                             style="width:80%; margin-left:auto; 
                             margin-right:auto; margin-top: 20px;">         
                          <form method="post" enctype="multipart/form-data">
                            <div class="field">
                                <label for="file">Choose image:</label>
                                <input type="file" id="file" name="file"  multiple required onchange="readURL(this);" accept=".jpg, .jpeg, .png"  >
                            </div>
                            <div class=" two fields">
                                <div class="field">
                                    <img id="blah" height="350px" alt="your image" />
                                </div>
                                <div class="field">
                                    <label for="location">Location:</label>
                                    <input type="text"  id="location" placeholder="place for location" />
                                    <div id="map" class="map" style="height:420px;"></div>
                                            <script>pomocnaP();</script>
                                </div>
                            </div>
                             <div class="field">
                                    <label for="description">Description:</label>
                                    <textarea type="text"  id="description" name="description" placeholder="Description" rows = "2"/>
                            </div>                      
                            
                   
                    <div class="tag">

            <div class="ui fluid multiple search selection dropdown">
                <input name="tags" id="tags" type="hidden">
                <i class="dropdown icon"></i>
                <div class="default text">Tags</div>
                <div class="menu">
              ` + temp + `
                </div>
            </div>

        </div>
                            
                            <div class="ui grid">
                            <div class="two wide column"></div>
                            <div class="two wide column"></div>
                            <div class="two wide column"></div>
                            <div class="two wide column"></div>
                            <div class="two wide column"></div>
                            <div class="two wide column"></div>
                            <div class="two wide column"></div>
                            <div class="two wide column right">  
                            <button type="button" style = "text-align: center" class="ui primary button" id="save_post" >ADD POST</button>
                            </div>
                            </div>
                               <script>
                                    $('.tag .ui.dropdown').dropdown();
                                </script>
                          </form>
                      </form>
 `
            );

            $('#save_post').click(function () {
                var formData = new FormData();
                formData.append("file", file);
                formData.append("type","image")
                var description = $('#description').val();
                var tags = $('#tags').val();
                console.log(tags)
                var location = $('#location').val();
                var email = localStorage.getItem('email');

                formData.append("description", description)
                formData.append("tags", tags)
                formData.append("location", location)
                formData.append("email", email)
                customAjax({
                    url: 'http://localhost:80/post-service/savePost',
                    method: 'POST',
                    data: formData,
                    processData: false,
                    contentType: false,
                    success: function () {
                        alert("Sucess saved post")
                        customAjax({
                            url: 'http://localhost:80/user-service/getByEmail/' + email,
                            method: 'GET',
                            success: function (data) {
                                myProfile(data)
                            },
                            error: function () {
                            }
                        });
                    },
                    error: function (e) {
                        alert('Error uploading new post.')
                    }
                });
            });
        });
    });
    
    $(window).on('load', function () {
        customAjax({
            url: 'http://localhost:80/user-service/getByEmail/' + email,
            method: 'GET',
            success: function(data){
                var json = JSON.parse(data);
                if(json.isPrivate){
                    $("#requestsForFollowing").html(`<a  id="requests"><i class="user plus icon"  style="color:white"></i></a>`);
                    $("#notifications").html(`<a  id="followers"><i class="bell icon"  style="color:white"></i></a>`);
                    $("#requests").click(function () {
                        customAjax({
                            url: 'http://localhost:80/user-service/getAllRequests/' + email,
                            method: 'GET',
                            success: function(data){
                                showRequests(data)
                            },
                            error: function(){
                            }
                        });
                    });
                }
                $("#notifications").html(`<a  id="followers"><i class="bell icon"  style="color:white"></i></a>`);
                $("#followers").click(function () {
                    customAjax({
                        url: 'http://localhost:80/user-service/getAllFollowers/' + email,
                        method: 'GET',
                        success: function(data){
                            console.log(data)
                            showFollowers(data)
                        },
                        error: function(){
                        }
                    });
                });

            },
            error: function(){
            }
    });
    });

    $("#myProfile").click(function () {
   
		  customAjax({
		      url: 'http://localhost:80/user-service/getByEmail/' + email,
		      method: 'GET',
		      success: function(data){
                myProfile(data)
		      },
		      error: function(){
		      }
        });
	});  
	
	  $("#editProfile").click(function () {
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
        localStorage.removeItem('email')
        localStorage.removeItem('role')
		location.href = "/";
		});

    $("#search").click(function () {
        var username= $("#userSearch").val();
        customAjax({
            url: 'http://localhost:80/search-service/searchUserByUsername/' + username + '/' + localStorage.getItem("email"),
            method: 'GET',
            success: function (data) {
                showProfile(data);
            },
            error: function () {
            }

        });

    });

    $("#searchLocation").click(function () {
        var location= $("#locationSearch").val();
        console.log(location)
        customAjax({
            url: 'http://localhost:80/search-service/searchPostByLocation/' + location + "/" + localStorage.getItem("email"),
            method: 'GET',
            success: function (data) {
                showPosts(data);
            },
            error: function () {
            }

        });

    });
    $("#searchTags").click(function () {
        var tag= $("#tagsSearch").val();
        console.log(tag)
        customAjax({
            url: 'http://localhost:80/search-service/searchPostByTag/' + tag + "/" + localStorage.getItem("email"),
            method: 'GET',
            success: function (data) {
                console.log(data)
                showPosts(data);
            },
            error: function () {
            }

        });

    });

});
let myProfile = function(user){
    var json = JSON.parse(user);
    var email = json.email
    var verified = ""
    if(json.isVerified) {
        verified = `<i class="check circle icon" style="color: dodgerblue"></i>`
    } else {
        verified = ``
    }
    customAjax({
        url: 'http://localhost:80/post-service/getAllPostsByEmail/' + email,
        method: 'GET',
        success: function(data){
            showPosts(data)
        },
        error: function(){
        }
    });

    function showPosts(data) {
        json = JSON.parse(data)

        var slika

        result = ""
        result += `<div style="margin-top: 50px" ><div class="ui cards">`;
        console.log(json)
            for (i in json) {


                slika = ""
                customAjax({
                    url: 'http://localhost:80/search-service/getMedia/' + json[i].ImageID,
                    method: 'GET',
                    async: false,
                    success: function (data) {
                        slika = data

                        if (slika.type == "video") {

                            pom1 = `<video id="output" height="150px" alt="slika" autoplay src ="` + 'data:video/mp4;base64,' + slika.path + ` ">`;
                        } else {
                            pom1 = `<img id="output" height="150px" alt="slika" src ="` + 'data:image/png;base64,' + slika.path + ` ">`;
                        }
                        result += `<br><div class="ui card">

  <div class="content">
     <div class="right floated meta">` + json[i].CreatedAt.split("T")[0] +  `</div>  
     <div class="left floated meta"><button class="ui basic icon button">
        <i class="bookmark outline icon"></i>
        </button>
    </div>  
      <br>
    <div class="description" style="color:cornflowerblue"><i class="location arrow icon" ></i>
      ` + json[i].Location + `
    </div>
  </div>
   <div class="image">
    ` + pom1 + `
  </div>
  <div class="content">
  <div class="description">
      ` + json[i].description + `
    </div>
      <br>
    <div class="description"><i class="tags icon"></i>
      ` + json[i].tags + `
    </div>
  </div> 
  <div class="extra content">
    <div class="ui large transparent left icon input"><button name="like" id="`+json[i].ID+`">
      <i class="heart outline  icon"></i>
      </button>
      <label name="`+json[i].ID+`"></label>
      <input type="text" placeholder="Add Comment...">
    </div>
  </div>
</div>`;

                    },
                    error: function () {

        $("button[name=like]").click(function(){
        var postID=this.id
        var userWhoLiked=localStorage.getItem('email');

      
          customAjax({
                    url: 'http://localhost:80/post-service/liked/' + postID + "/" + userWhoLiked,
                    method: 'POST',
                    success: function (data) {
                       $('label[name='+postID+']').text("Liked")
                    },
                    error: function () {
                    }
                })

    });

            result += `</div></div>`;
            $('#posts').html(result);



        }

        $("#showData").html(
            `<div  style="width:80%; margin-left:auto; 
                             margin-right:auto;">
            <h2 class="ui left aligned header" ></h2>
           <div class="ui clearing segment">
           <h3 class="ui right aligned header" >
          <div class="ui vertical labeled icon buttons"><button id="verification" class="ui button"><i class="check circle icon"></i>Request for verification</button></div></h3>
            <h3 class="ui left floated header">
               ` + json.username + ``+ verified+`
            </h3>
            </div>
           <div class="ui secondary pointing menu">
            <a href="#" class="item active" id="posts_myprofile">
                Posts
            </a>
            <a href="#" class="item" id="stories_myprofile">
                Stories
            </a>
            <a  href="#" class="right item" id="liked_posts_myprofile">
                Liked posts
            </a>
            </div>
            <h3 class="ui header"></h3>
                <div class="ui four cards" id='posts'></div>
        </div> `
        );
     $('#liked_posts_myprofile.item').click(function(){
        $(this).addClass('active')
        $("#posts_myprofile").removeClass('active');
        $("#stories_myprofile").removeClass('active');
        customAjax({
            url: 'http://localhost:80/post-service/getAllLikedPostsByEmail/' + email,
            method: 'GET',
            success: function(data){
                showPosts(data)
            },
            error: function(){
            }
        });
     });
     $("#stories_myprofile.item").click(function () {
         $(this).addClass('active');
        $("#posts_myprofile").removeClass('active');
        $("#liked_posts_myprofile").removeClass('active');
         customAjax({
        url: 'http://localhost:80/story-service/getAllStoriesByEmail/' + email,
        method: 'GET',
        success: function(data){
            showStories(data)
        },
        error: function(){
        }
    });
    });
       $("#posts_myprofile.item").click(function () {
         $(this).addClass('active');
         $("#stories_myprofile").removeClass('active');
         $("#liked_posts_myprofile").removeClass('active');
         customAjax({
        url: 'http://localhost:80/post-service/getAllPostsByEmail/' + email,
        method: 'GET',
        success: function(data){
            showPosts(data)
        },
        error: function(){
        }
    });
    });

        $("#verification").click(function () {
            showRequest(localStorage.getItem("email"));

        })
        $("#stories_myprofile.item").click(function () {
            $(this).addClass('active');
            $("#posts_myprofile").removeClass('active');
            customAjax({
                url: 'http://localhost:80/story-service/getAllStoriesByEmail/' + email,
                method: 'GET',
                success: function (data) {
                    showStories(data)
                },
                error: function () {
                }
            });
        });
        $("#posts_myprofile.item").click(function () {
            $(this).addClass('active');
            $("#stories_myprofile").removeClass('active');
            customAjax({
                url: 'http://localhost:80/post-service/getAllPostsByEmail/' + email,
                method: 'GET',
                success: function (data) {
                    showPosts(data)
                },
                error: function () {
                }
            });
        });

      function showStories(data)
    {
        json = JSON.parse(data)
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

                    if(slika.type == "video") {

                        pom1 = `<video id="output" height="150px" alt="slika" autoplay src ="` + 'data:video/mp4;base64,' + slika.path + ` ">`;
                    } else {
                        pom1 = `<img id="output" height="150px" alt="slika" src ="`+'data:image/png;base64,'+ slika.path + ` ">`;
                    }
                    result += `<br><div class="ui card">

  <div class="content">
     <div class="right floated meta">` + json[i].CreatedAt.split("T")[0] + `</div>  
     <div class="left floated meta"><button class="ui basic icon button">
        <i class="icon gem"></i>
        </button>
    </div> 
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

        result+=`</div></div>`;
        $('#posts').html(result);



    }



}

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

    var isPrivate = ``;
    if (json.isPrivate) {
        isPrivate = `<input type="checkbox" name="private" id="isPrivate" checked="checked">`
    } else {
        isPrivate = `<input type="checkbox" name="private" id="isPrivate">`
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
                                    <td> <p id="errorUsername"></p></td>
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
                     <tr>
                     <td>Private profile</td>
                   <td><div class="ui toggle checkbox ">` + isPrivate +`
                        <label></label>
                    </div></td>
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

        input_username=$('#txtUsername');
        var btnEdit = document.getElementById("acceptChange")
        btnEdit.disabled = true
    input_username.keyup(function () {
        if(validateUsername(input_username.val())) {
            btnEdit.disabled = false
        }
        if(!validateUsername(input_username.val())){
            btnEdit.disabled = true
            $(this).addClass(`alert-danger`);
            $('#txtUsername').css('border-color', 'red');
            $("#errorUsername").text("You can only use letters and numbers for username!")
            $('#errorUsername').css('color', 'red');
        }else {
            $(this).removeClass(`alert-danger`);
            $('#txtUsername').css('border-color', '');
            $("#errorUsername").text("")
        }
    });
        $('#acceptChange').click(function(){

        let email=$('#txtEmail').val()
        let name=$('#txtName').val()
        let username=$('#txtUsername').val()
        let phone=$('#txtPhone').val()
        let birthday=$('#txtBirthday').val()
        let website=$('#txtWebsite').val()
        let biography=$('#txtBiography').val()
        let gender = $("input:radio[name=gender]:checked").val();
        let isPrivate = document.getElementById("isPrivate").checked;
        console.log(website)

        obj = JSON.stringify({
        email:email,
        name:name,
        username: username,
        phone:phone,
        birthday:birthday,
        website :website,
        biography:biography,
        gender:gender,
        isPrivate:isPrivate
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
function sleep(milliseconds) {
    const date = Date.now();
    let currentDate = null;
    do {
        currentDate = Date.now();
    } while (currentDate - date < milliseconds);
}
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
        var pom1 = '';
        if (json[i].isPrivate == true) {
            pom1 += "<i class=\"lock icon\"></i>";
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
    <button class="ui teal button" name = "follow" id = "` + json[i].username + `">Follow  <i class = "user icon"></i></button>
    
    `+pom1+`
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
            url: 'http://localhost:80/search-service/getPostsForSearchedUser/' + this.id +"/"+ localStorage.getItem("email"),
            method: 'GET',
            success: function (data) {
                showPostsForSearchedUser(data)
            },
            error: function () {
                alert("Locked")
            }
        })
    })

    $("button[name=follow]").click(function () {
        id = this.id
        console.log(id)
        customAjax({
            url: 'http://localhost:80/user-service/alreadyFollow/' + id + "/" + localStorage.getItem("email"),
            method: 'GET',
            success: function () {
                customAjax({
                    url: 'http://localhost:80/user-service/follow/' + id + "/" + localStorage.getItem("email"),
                    method: 'POST',
                    success: function (data) {
                        console.log("Success za follow")
                        document.getElementById(id).innerText = "Followed"
                        document.getElementById(id).disabled = true
                    },
                    error: function () {
                        console.log("Error za follow")
                    }
                })

            },
            error: function () {
                console.log("Error za already")
                var errorId = "error"+id
                document.getElementById(id).innerText = "Followed"
                document.getElementById(id).disabled = true
                document.getElementById(errorId).innerText="Already followed"
            }
        })


    })


}
let showRequests = function(user) {
    var json = JSON.parse(user);
    var pomocna ="";
    pomocna +=`<div style="margin-top: 50px" ><div class="ui cards">`;
    for( i in json) {
        var pom = '';
        if (json[i].gender == "female") {
            pom += "<img class=\"right floated mini ui image\"  src=\"https://avataaars.io/?avatarStyle=Transparent&topType=LongHairStraightStrand&accessoriesType=Blank&hairColor=BrownDark&facialHairType=Blank&clotheType=BlazerShirt&eyeType=Happy&eyebrowType=Default&mouthType=Smile&skinColor=Light\">";
        } else {
            pom += "<img class=\"right floated mini ui image\"  src='https://avataaars.io/?avatarStyle=Transparent&topType=ShortHairShortFlat&accessoriesType=Blank&hairColor=BrownDark&facialHairType=BeardLight&facialHairColor=BrownDark&clotheType=BlazerShirt&eyeType=Happy&eyebrowType=Default&mouthType=Smile&skinColor=Light'\n" +
                "/>"
        }
    pomocna += `<br><div class="ui card">
  <div class="content">` + pom + ` 
    <a class="header">` + json[i].username + `</a>
    <div class="meta">
      <span class="date"> ` + json[i].name + `</span>
    </div>
  </div>
  <div class="extra content">
   <div class="ui two buttons">
  
        <button class="ui basic green button" name="approve" id="`+json[i].username+`">Approve</button>
        <button class="ui basic red button" name="decline" id="`+json[i].username+`">Decline</button>
      </div>
  </div>
</div>`;


    }
    pomocna+=`</div></div>`;
    $("#showData").html(pomocna);
    $("button[name=approve]").click(function () {
        customAjax({
            url: 'http://localhost:80/user-service/acceptRequest/' + this.id + "/" + localStorage.getItem("email"),
            method: 'POST',
            success: function () {
                location.href="userHomePage.html";
            },
            error: function () {
            }
        })
    })
    $("button[name=decline]").click(function () {
        customAjax({
            url: 'http://localhost:80/user-service/declineRequest/' + this.id + "/" + localStorage.getItem("email"),
            method: 'POST',
            success: function () {
                location.href="userHomePage.html";
            },
            error: function () {
            }
        })
    })



}

let showFollowers = function(user) {
    var json = JSON.parse(user);
    console.log(json)
    var pomocna ="";
    pomocna +=`<div style="margin-top: 50px" ><div class="ui celled list">`;
    for( i in json) {
        var pom = '';
        if (json[i].gender == "female") {
            pom += "<img class=\"ui avatar image\"  src=\"https://avataaars.io/?avatarStyle=Transparent&topType=LongHairStraightStrand&accessoriesType=Blank&hairColor=BrownDark&facialHairType=Blank&clotheType=BlazerShirt&eyeType=Happy&eyebrowType=Default&mouthType=Smile&skinColor=Light\">";
        } else {
            pom += "<img class=\"ui avatar image\"  src='https://avataaars.io/?avatarStyle=Transparent&topType=ShortHairShortFlat&accessoriesType=Blank&hairColor=BrownDark&facialHairType=BeardLight&facialHairColor=BrownDark&clotheType=BlazerShirt&eyeType=Happy&eyebrowType=Default&mouthType=Smile&skinColor=Light'\n" +
                "/>"
        }
        pomocna += `<div class="item">
         ` + pom + ` 
    <div class="content">
      <div class="header">` +json[i].username + `</div> 
      <div class="description"> Started following you.
        <button class="ui teal active button" name="followBack" id="`+json[i].username+`"><i class="user icon"></i>Follow back
  </button>
      </div>
       <div id="error` + json[i].username + `" style="color:red"></div>
    </div>
  </div>
`;


    }
    pomocna+=`</div></div>`;
    $("#showData").html(pomocna);


    $("button[name=followBack]").click(function () {
        id = this.id
        customAjax({
            url: 'http://localhost:80/user-service/alreadyFollow/' + id + "/" + localStorage.getItem("email"),
            method: 'GET',
            success: function () {
                customAjax({
                    url: 'http://localhost:80/user-service/follow/' + id + "/" + localStorage.getItem("email"),
                    method: 'POST',
                    success: function (data) {
                        console.log("Success za follow")
                        document.getElementById(id).innerText = "Followed"
                        document.getElementById(id).disabled = true
                    },
                    error: function () {
                        console.log("Error za follow")
                    }
                })

            },
            error: function () {
                console.log("Error za already")
                var errorId = "error"+id
                document.getElementById(id).innerText = "Followed"
                document.getElementById(id).disabled = true
                document.getElementById(errorId).innerText="Already followed"
            }
        })

    })


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


                pom1 = `<img id="output" height="150px" alt="slika" src ="`+'data:image/png;base64,'+ slika + ` ">`;


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
  
  
  <div class="extra content">
    <div class="ui large transparent left icon input">
    <button name="like" id="`+json[i].ID+`">
      <i class="heart outline  icon"></i>
      </button>
      <label name="`+json[i].ID+`"></label>
      <input type="text" placeholder="Add Comment...">
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

    $("button[name=like]").click(function(){
        var postID=this.id
        var userWhoLiked=localStorage.getItem('email');

      
          customAjax({
                    url: 'http://localhost:80/post-service/liked/' + postID + "/" + userWhoLiked,
                    method: 'POST',
                    success: function (data) {
                        $('label[name='+postID+']').text("Liked");
                    },
                    error: function () {
                    }
                })

    });


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
  <div class="extra content">
    <div class="ui large transparent left icon input">
      <button name="like" id="`+json[i].ID+`">
      <i class="heart outline  icon"></i>
      </button>
      <label name="`+json[i].ID+`"></label>
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

    $("button[name=like]").click(function(){
        var postID=this.id
        var userWhoLiked=localStorage.getItem('email');
          customAjax({
                    url: 'http://localhost:80/post-service/liked/' + postID + "/" + userWhoLiked,
                    method: 'POST',
                    success: function (data) {
                      $('label[name='+postID+']').text("Liked")
                    },
                    error: function () {
                    }
                })

    });


}
function validateUsername(name) {
    const re = /^[a-zA-Z0-9]+$/;
    return re.test(String(name));
}



