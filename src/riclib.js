"use strict";
var riclib = (function(window, document, undefined) {
	var RICLib = function() {};
	var _config = {};

	// Feature detection methods
	function getDevicePixelRatio() {
		return window.devicePixelRatio || 1;
	};

	// RIC private tasks start
	function getQueryString(obj) {
		var queryStr = '';
		var first = true;

		for(var key in obj) {
			queryStr += (first ? '': '&') + encodeURIComponent(key) + '=' + encodeURIComponent(obj[key])
			first = false;
		}

		return queryStr;
	};

	function handleSingleImage(img) {
		var queryParams = {
			maxres: _config.maxres,
			format: 'jpeg',
			quality: _config.quality,
			url: img.src
		};

		console.debug('riclib::handleSingleImage(), query params', queryParams, img);

		img.src = _config.server_path + '?' + getQueryString(queryParams);
	};

	function processAllImages() {
		// images is 'HTMLCollection'
		var images = document.getElementsByTagName('img');

		var length = images.length;
		for(var i = 0; i<length ; i++) {
			handleSingleImage(images[i]);
		}
	};
	// RIC private tasks end

	// riclib interface start
	RICLib.prototype.init = function(config) {
		_config = config || {};
		console.debug('riclib::init(), called with config', config);

		processAllImages();
	};

	RICLib.prototype.devicePixelRatio = getDevicePixelRatio();
	// riclib interface end

	return new RICLib();
})(window, document);

// initialize riclib and pass potential config
document.addEventListener("DOMContentLoaded", function(event) {
	riclib.init(window.RICConfig);
});