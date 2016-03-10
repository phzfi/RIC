
window.onload = function() {
  setImdivSize();
  setScreenInfo();
  formats = [
    ['webp', 'webp'],
    ['jpeg', 'jpg'],
    ['png', 'png'],
    ['bmp', 'bmp']
  ];
}

window.onresize = function() {
  setImdivSize();
  setScreenInfo();
  document.getElementById('I').src = "";
}

function getImage() {
  var id = document.getElementById('imgname').value;
  var imgW = document.getElementById('imdiv').clientWidth;
  var imgH = document.getElementById('imdiv').clientHeight;
  var zoom = Math.round((window.outerWidth / (window.innerWidth + 18 / window.devicePixelRatio))*100) / 100;
  
  // if large screen is zoomed in
  if (screen.width > 1000 && zoom != 1) {
	  // set maximum img-element css size
	  // and "squeeze" HQ image into that
	  var I = document.getElementById('I');
	  I.width = imgW;
	  I.height = imgH;
	  
	  imgW *= zoom;
	  imgH *= zoom;
  }

  if (id.indexOf('.') === -1) {
    for (var i = 0; i < formats.length; i++) {
      if (supportsIMG(formats[i][0])) {
        id += '.' + formats[i][1];
        break;
      }
    }
  }
  var url = 'http://ric.phz.fi:8005/'+id+'?width='+imgW+'&height='+imgH;
  document.getElementById('I').src = url;
}

function supportsIMG(img) {
  var canvas = document.createElement('canvas');
  canvas.width = canvas.height = 1;
  var uri = canvas.toDataURL('image/' + img);

  return uri.match('image/' + img) !== null;
}

function setImdivSize() {
  // use aspect ratio of 16:9 to set image height
  var imgW = document.getElementById('imdiv').clientWidth;
  var imgH = Math.round(imgW / 16 * 9);
  document.getElementById('imdiv').style.height = imgH + 'px';
}

function setScreenInfo() {
  // show screen information
  var w = screen.width;
  var h = screen.height;
  var dpr = window.devicePixelRatio;
  var text = 'Your logical screen size is: ' + w + 'x'+ h + ' (' + dpr + ')<br>';
  text += 'Maximum image size for you: ' + Math.round(w * dpr) + 'x';
  text += Math.round(h * dpr) + '<br>';
  text += '<br>' + 'imgdiv clientWidth: ' + document.getElementById('imdiv').clientWidth;
  text += '<br>' + 'imgdiv clientHeight: ' + document.getElementById('imdiv').clientHeight;
  text += '<br>' + 'window.innerWidth: ' + window.innerWidth;
  text += '<br>' + 'window.outerWidth: ' + window.outerWidth + '<br>';
  text += '<br>' + 'Exact Zoom: ' + window.outerWidth / (window.innerWidth + 18/dpr);
  text += '<br>' + 'Rounded Zoom: ' + Math.round((window.outerWidth / (window.innerWidth + 18/dpr))*100) / 100;
  if (screen.width > 1000)
	text += '<br> Computer screen';
  else
	text += '<br> Mobile screen';
  document.getElementById('infodiv').innerHTML = text;
}
