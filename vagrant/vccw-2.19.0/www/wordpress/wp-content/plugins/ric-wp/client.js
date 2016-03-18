
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
    var id = elements[i].dataset.id;
    // if no explicit format has been given in the id, use the most appropriate
    if(id.indexOf('.') == -1){
      id += fmt;
    }
	
	// read width and height from parent element
	var width = elements[i].parentElement.clientWidth;
	var height = elements[i].parentElement.clientHeight;
	
	// Use zoom coefficient for large screens
	if (screen.width > 1000) {
		width = Math.round(width * window.devicePixelRatio);
		height = Math.round(height * window.devicePixelRatio);
	}
	
	var url = URI + id;
	var params = ['width=' + width];
	if (elements[i].parentElement.outerHTML.indexOf('height') != -1)
		params.push('height=' + height);
	
	// set image resize mode parameter - fit to parent dimensions as default
	if (elements[i].dataset.mode == 'nofit') {}
	else if (elements[i].dataset.mode == 'liquid') {
		params.push('mode=liquid');
	}
	else {
		params.push('mode=fit');
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
	elements[i].onload = function() {	// Set maximum image css resolution
		if (screen.width > 1000) {
			this.width = Math.round(this.naturalWidth / window.devicePixelRatio);
			this.height = Math.round(this.naturalHeight / window.devicePixelRatio);
		} else {
			this.width = this.naturalWidth;
			this.height = this.naturalHeight;
		}
		
		// if original image is returned, fit to content
		if (this.height > this.parentElement.clientHeight) {
			var dh = this.height - this.parentElement.clientHeight;
			var ratio = this.width / this.height;
			this.width = Math.round(this.width - dh * ratio);
			this.height = this.parentElement.clientHeight;
		}
	};
  }
}

//check if browser supports image format
function supportsIMG(format) {
  var canvas = document.createElement('canvas');
  canvas.width = canvas.height = 1;
  var uri = canvas.toDataURL('image/' + format);

  return uri.match('image/' + format) !== null;
}
