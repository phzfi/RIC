
window.onload = function() {
  URI = php_vars.URI.url;

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

  // loop over all elements of page and inject to image divs RIC URL's
  for (var i = 0; i < elements.length; i++){
    var id = elements[i].id;
    // if no explicit format has been given in the id, use the most appropriate
    if(id.indexOf('.') == -1){
      id += fmt;
    }
	
	// read width and height from parent element dimensions
	// and multiply them with device pixel ratio for high dpi screens
	var dpr = window.devicePixelRatio;
	var width = 'width=' + elements[i].parentElement.clientWidth * dpr;
	var height = 'height=' + elements[i].parentElement.clientHeight * dpr;
	
	// if no height specified (in HTML!), request image using width only
	if (elements[i].parentElement.outerHTML.indexOf('height') == -1) {
		var img = new Image();
		img.src = URI + id + '?' + width;
		elements[i].parentElement.clientHeight = img.height;
		elements[i].src = img.src;
	} else {
		var params = [width, height];
		
		// set image resize mode parameter - fit to parent dimensions as default
		var mode;
		if (elements[i].classList.contains('mode-nofit')) {}
		else if (elements[i].classList.contains('mode-liquid')) {
			mode = 'mode=liquid';
			params.push(mode);
		}
		else {
			mode = 'mode=fit';
			params.push(mode);
		}
		
		// extend uri with the parameters and set image source
		var url = URI + id;
		if (params.length > 0) {
			url += '?' + params[0];
			for (var j = 1; j < params.length; j++) {
				url += '&' + params[j];
			}
		}
		elements[i].src = url;
	}
  }
}

//check if browser supports image format
function supportsIMG(format) {
  var canvas = document.createElement('canvas');
  canvas.width = canvas.height = 1;
  var uri = canvas.toDataURL('image/' + format);

  return uri.match('image/' + format) !== null;
}
