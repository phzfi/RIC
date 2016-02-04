
window.onload = function() {
  var formats = [
    ['webp', 'webp'],
    ['jpeg', 'jpg'],
    ['png', 'png'],
    ['bmp', 'bmp']
  ];
  fmt = '';
  for (var i = 0; i < formats.length; i++) {
    if (supportsIMG(formats[i][0])) {
      fmt = '.' + formats[i][1];
      break;
    }
  }
  var elements = document.getElementsByClassName('ricimg');
  
  for (var i = 0; i < elements.length; i++){
    var id = elements[i].id;
    console.log(id);
    var h = elements[i].parentElement.clientHeight;
    var w = elements[i].parentElement.clientWidth;
    var url = 'http://ric.phz.fi:8005/'+id+fmt+'?width='+w+'&height='+h;
    elements[i].src = url;
  }

  //}
  //document.getElementById('buttonID').onclick = function() {
    /*var id = document.getElementById('imgname').value;
    var imgW = document.getElementById('imdiv').clientWidth;
    var imgH = document.getElementById('imdiv').clientHeight;
    
    if (id.indexOf('.') === -1) {
      for (var i = 0; i < formats.length; i++) {
        if (supportsIMG(formats[i][0])) {
          id += '.' + formats[i][1];
          break;
        }
      }
    }
    */
    
    //document.getElementById('I').src = url;
  //};


  //console.log("testing...");
}

function supportsIMG(img) {
  var canvas = document.createElement('canvas');
  canvas.width = canvas.height = 1;
  var uri = canvas.toDataURL('image/' + img);

  return uri.match('image/' + img) !== null;
}


