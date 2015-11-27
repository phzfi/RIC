
window.onload = function() {
  // use aspect ratio of 16:9 to set image height
  imgW = document.getElementById('imdiv').clientWidth;
  imgH = Math.round(imgW / 16 * 9);
  document.getElementById('imdiv').style.height = imgH + 'px';


  // show screen information
  var w = screen.width;
  var h = screen.height;
  var dpr = window.devicePixelRatio;
  var text = 'Your logical screen size is: ' + w + 'x'+ h + ' (' + dpr + ')<br>';
  text += 'Maximum image size for you: ' + Math.round(w * dpr) + 'x';
  text += Math.round(h * dpr) + '<br>';
  text += 'imgdiv size: ' + imgW + 'x' + imgH;
  document.getElementById('infodiv').innerHTML = text;
}


// get image
function buttonf(b) {
  var url = 'http://ric.phz.fi:8005/'+b.innerHTML+'?width='+imgW+'&height='+imgH;
  document.getElementById('I').src = url;
}
