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
let showVideo = function () {
    customAjax({
        url: 'http://localhost:80/user-service/getAllUsersExceptLogging/' + localStorage.getItem("email"),
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
                                <input type="file" id="file" name="file"  multiple required onchange="readURL(this);" accept="video/*"  >
                            </div>
                            <div class=" two fields">
                                <div class="field">
                                    <video id="blah" height="250px" autoplay>
                                    
                                </video>
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
        var description = $('#description').val();
        var tags = $('#tags').val();
        console.log(tags)
        var location = $('#location').val();
        var email = localStorage.getItem('email');
        formData.append("type","video")
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
}


let showVideoStory = function () {
    customAjax({
        url: 'http://localhost:80/user-service/getAllUsersExceptLogging/' + localStorage.getItem("email"),
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
                                <input type="file" id="file" name="file"  multiple required onchange="readURL(this);" accept="video/*"  >
                            </div>
                            <div class=" two fields">
                                <div class="field">
                                    <video id="blah" height="250px" autoplay>
                                    
                                </video>
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
                      </form>
 `
    );

    $('#save_story').click(function () {
        var formData = new FormData();
        formData.append("file", file);
        var description = $('#description').val();
        var tags = $('#tags').val();
        console.log(tags)
        var location = $('#location').val();
        var email = localStorage.getItem('email');
        formData.append("type","video")
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
                        myProfile(data)
                    },
                    error: function () {
                    }
                });
            },
            error: function (e) {
                alert('Error uploading new story.')
            }
        });
    });
}