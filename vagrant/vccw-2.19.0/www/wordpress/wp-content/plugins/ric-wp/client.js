
window.onload = function() {
  
  /******** replace with your RIC server URI *********/
  URI = URLI['urli'];
  /***************************************************/
  

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
  
  //loop over all elements of page and inject to image divs RIC URL's
  for (var i = 0; i < elements.length; i++){
    var id = elements[i].id;
    //if no explicit format has been given in the id, use the most appropriate
    if(id.indexOf('.') === -1){
      id += fmt;
    }
    
    //get image dimensions from DOM
    var h = elements[i].parentElement.clientHeight;
    var w = elements[i].parentElement.clientWidth;
    var url = URI+id+'?width='+w+'&height='+h;
    elements[i].src = url;
  }
}

//check if browser supports image format
function supportsIMG(format) {
  var canvas = document.createElement('canvas');
  canvas.width = canvas.height = 1;
  var uri = canvas.toDataURL('image/' + format);

  return uri.match('image/' + format) !== null;
}


