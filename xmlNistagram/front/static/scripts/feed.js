let showFeed = function () {
    customAjax({
        url: 'http://localhost:80/search-service/getAllPosts/' + localStorage.getItem("email"),
        method: 'GET',
        success: function (data) {
            showPosts(data)
        },
        error: function () {
        }
    });
    $('#showData').html(`<div  style="width:80%; margin-left:auto; 
                             margin-right:auto;">
            
           <div class="ui secondary pointing menu">
            <a href="#" class="item active" id="posts_myprofile">
                Posts
            </a>
            <a href="#" class="item" id="stories_myprofile">
                Stories
            </a>
            </div>
            <h3 class="ui header"></h3>
                <div class="ui four cards" id='posts'></div>
        </div> `)

    $("#stories_myprofile.item").click(function () {
        $(this).addClass('active');
        $("#posts_myprofile").removeClass('active');
        customAjax({
            url: 'http://localhost:80/search-service/getAllStories/' + localStorage.getItem("email"),
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
            url: 'http://localhost:80/search-service/getAllPosts/' + localStorage.getItem("email"),
            method: 'GET',
            success: function (data) {
                showPosts(data)
            },
            error: function () {
            }
        });
    });

    function showPosts(data) {
        json = JSON.parse(data)

        var slika

        result = ""
        result += `<div style="margin-top: 50px" ><div class="ui cards">`;
        console.log(json)
        for (i in json) {
            var postID = json[i].ID

            slika = ""
            customAjax({
                url: 'http://localhost:80/search-service/getMedia/' + json[i].ImageID,
                method: 'GET',
                async: false,
                success: function (data) {
                    customAjax({
                        url: 'http://localhost:80/user-service/getByEmail/' + json[i].email,
                        method: 'GET',
                        async: false,
                        success: function (user) {
                            jsonParse = JSON.parse(user);
                        },
                        error: function () {

                        }
                    })
                    slika = data

                    if (slika.type == "video") {

                        pom1 = `<video id="output" height="150px" alt="slika" autoplay src ="` + 'data:video/mp4;base64,' + slika.path + ` ">`;
                    } else {
                        pom1 = `<img id="output" height="150px" alt="slika" src ="` + 'data:image/png;base64,' + slika.path + ` ">`;
                    }
                    result += `<br><div class="ui card">

  <div class="content">
  
  <div class="left floated meta">` + jsonParse.username + `</div>
      <div class="right floated meta"><button class="ui basic icon button">
        <i class="bookmark outline icon"></i>
        </button>
    </div> 
     <div class="right floated meta">` + json[i].CreatedAt.split("T")[0] + `</div>  
  
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
            <div class="ui large transparent left icon input"><button name="like" id="` + json[i].ID + `">
            <i class="heart outline  icon"></i>
            </button>
            <label name="` + json[i].ID + `"></label>
            </div>
             <div class="ui large transparent  input"> 
                        <input type="text" placeholder="Add Comment..."  name = "` + json[i].ID+`" />
                        <button class="ui primary basic button" name = "add-comment-button" id="` + json[i].ID+`">Post</button>
                    </div>
            <div class="ui divider"></div>
            <div class="ui comments">
                            <div class="comment">
                                        <div class="content">
                                            <a class="author">Matt</a>
                                            <div class="metadata">
                                                <span class="date">Today at 5:42PM</span>
                                            </div>
                                            <div class="text">
                                                How artistic!
                                            </div>
                                            <div class="actions">
                                                <a class="reply">Reply</a>
                                            </div>
                                        </div>
                            </div>

                            <div class="comment">
                                <div class="content">
                                    <a class="author">Matt</a>
                                    <div class="metadata">
                                        <span class="date">Today at 5:42PM</span>
                                    </div>
                                    <div class="text">
                                        How artistic!
                                    </div>
                                    <div class="actions">
                                        <a class="reply">Reply</a>
                                    </div>
                                </div>
                            </div>
                        </div>
      </div>
  </div>`;

                },
                error: function () {
                }
            })
        }


        result += `</div></div>`;
        $('#posts').html(result);

        $("button[name=add-comment-button]").click(function(){
            var ideic= this.id
            var postID= this.id
            console.log(ideic)

            var ispis= $('input[name='+ ideic +']' ).val()
            console.log(ispis)

            var formData = new FormData();
            var text = ispis;
            var postID = ideic;
            var email = localStorage.getItem('email');

            formData.append("text", text)
            formData.append("postID", postID)
            formData.append("email", email)
            customAjax({
                url: 'http://localhost:80/post-service/saveComment',
                method: 'POST',
                data: formData,
                processData: false,
                contentType: false,
                success: function () {
                    alert("success post comment")
                    $('input[name='+ ideic +']' ).val('')

                },
                error: function (e) {
                    alert('Error uploading new post.')
                }
            });
        })


        var action = 1;
        $("button[name=like]").click(function () {
            var postID = this.id
            var userWhoLiked = localStorage.getItem('email');


            customAjax({
                url: 'http://localhost:80/post-service/liked/' + postID + "/" + userWhoLiked,
                method: 'POST',
                success: function (data) {
                    /*
                    var action = 1;

                    $("input").on("click", viewSomething);

                    function viewSomething() {
                        if ( action == 1 ) {
                            $("body").css("background", "honeydew");
                            action = 2;
                        } else {
                            $("body").css("background", "beige");
                            action = 1;
                        }
                    }
                     */

                    //$('label[name=' + postID + ']').text("Liked")
                },
                error: function () {
                }
            })

            if ( action == 1 ) {
                $('label[name=' + postID + ']').text("Liked")
                action = 2;
            } else {
                $('label[name=' + postID + ']').text("")
                action = 1;
            }

        });


    }

}

function showStories(data)
{
    json = JSON.parse(data)
    var slika
    result=""
    result +=`<div style="margin-top: 50px" ><div class="ui five doubling cards">`;

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
                    async: false,
                    success: function (user) {
                        jsonParse = JSON.parse(user);
                    },
                    error: function () {

                    }
                })

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

    result+=`</div></div>`;
    $('#posts').html(result);




}


