
window.onload = function() {
  setImdivSize();
  setScreenInfo();
}

window.onresize = function() {
  setImdivSize();
  setScreenInfo();
  document.getElementById('I').src = "";
}

// get image
function buttonf(b) {
  var imgW = document.getElementById('imdiv').clientWidth;
  var imgH = Math.round(imgW / 16 * 9);
  var url = 'http://ric.phz.fi:8005/'+b.innerHTML+'?width='+imgW+'&height='+imgH;
  document.getElementById('I').src = url;
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
  //text += 'imgdiv size: ' + imgW + 'x' + imgH;
  text += '<br>' + 'imgdiv clientWidth: ' + document.getElementById('imdiv').clientWidth;
  text += '<br>' + 'imgdiv offsetWidth: ' + document.getElementById('imdiv').offsetWidth;
  text += '<br>' + 'imgdiv scrollWidth: ' + document.getElementById('imdiv').scrollWidth;
  document.getElementById('infodiv').innerHTML = text;
}
