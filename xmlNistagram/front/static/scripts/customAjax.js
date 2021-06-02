function customAjax(params){
    var jwt = localStorage.getItem('jwt');
    $.ajax({
      async: params.async,
      cache: params.cache,
      complete: function(xhr, status){
        console.log(params.url + ': Server returned ' + xhr.status + '; status is ' + status);
        if(params.complete != null){
          params.complete();
        }
      },
      contentType: params.contentType,
      data: params.data,
      dataType: params.dataType,
      error: function(xhr, status, error){
          
        if(xhr.status == 403) {
            $("#error1").html(`Nemate pravo pristupa!`)
        }
        if(xhr.status == 401){
          //localStorage.removeItem('jwt');
          window.location.href = 'index.html';
        }
        if(params.error != null){
          params.error();
        }
      },
      method: params.method,
      mimeType: params.mimeType,
      processData: params.processData,
      success: params.success,
      url: params.url,
      beforeSend: function (xhr) {
        if(jwt != null){
          xhr.setRequestHeader("Authorization", 'Bearer '+ jwt)
          console.log(params.data)
        }
      }
    });
    
  }