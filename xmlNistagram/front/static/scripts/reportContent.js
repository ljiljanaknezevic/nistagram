let showReport = function(id,email){
    $('#showData').html(`<table class="ui large table" style="width:50%; margin-left:auto; 
                margin-right:auto; margin-top: 40px;">
                            <thead>
                                <tr class="success">
                                    <th colspan="2" class = "text-info" style= "text-align:center;">Report post</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr>
                                <td colspan="2">
                                 <div className="ui form">
         <div className="grouped fields">
             <h3 style="color:black">Why are you reporting this ad?</h3>
             <div className="field">
                 <div className="ui radio checkbox">
                     <input type="radio"  checked="checked" name="exampleReport" value="I find it offensive">
                         <label style="color:black">I find it offensive</label>
                 </div>
             </div>
             <div className="field">
                 <div className="ui radio checkbox">
                     <input type="radio" name="exampleReport" value="It is spam">
                         <label style="color:black" >It is spam</label>
                 </div>
             </div>
             <div className="field">
                 <div className="ui radio checkbox">
                     <input type="radio" name="exampleReport" value="It is a scam or it is misleading">
                         <label style="color:black">It is a scam or it is misleading</label>
                 </div>
             </div>
             <div className="field">
                 <div className="ui radio checkbox">
                     <input type="radio" name="exampleReport" value="It is violent or prohibited content">
                         <label style="color:black">It is violent or prohibited content</label>
                 </div>
             </div>
              <div className="field">
                 <div className="ui radio checkbox">
                     <input type="radio" name="exampleReport" value="It violetes my intellectual property rights">
                         <label style="color:black">It violetes my intellectual property rights</label>
                 </div>
             </div>
         </div>
     </div></td>
                                </tr>
                                
                            </tbody>
                            <tfoot class="full-width">
                <tr>
                  <th></th>
                  <th colspan="2">
                       <input id = "reportAd" class="ui center floated negative basic button" type = "button" value = "Report"></input>
                
                  </th>
                </tr>
              </tfoot>
                        </table> <p id="er"> </p>
 `)
    $('#reportAd').click(function(){
       var reportContent =$("input:radio[name=exampleReport]:checked").val()
       obj = JSON.stringify({
          email:email,
          postId:id,
          reason: reportContent
       });
       console.log(obj)
 
       customAjax({
          url: 'http://localhost:80/post-service/reportPost',
          method: 'POST',
          data:obj,
          contentType: 'application/json',
          success: function(){
             alert("Sucess.")
             location.href = "userHomePage.html";
 
          },
          error: function(){
             alert('You have alreadu reported this post');
          }
       });
    })
 }