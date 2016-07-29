# Work flow for riclib.js

* 1) WP Users must be able to optimize their images to speed up page loading.
	- Simple UI interface for WP Plugin, description of the plugin must state that plugin is
		not altering the original image.
	- Plugin Options / Settings:
		* Set RIC host url (default to our sass)
		* Set jpeg quality
		* Set Max. resolution
		* Do not optimize images below xyz pixels. Done in RIC server-side
		* NOTE: WP has default media (image) resolution settings, which are accessible from RIC plugin code from 	somehow http://example.dev/wp/wp-admin/options-media.php
	- Plugin is adding 'riclib.js' in to page and added plugin configuration
	- riclib.js is attaching to all <img> elements and replaces 'src' attributes with proper settings to RIC saas.

* 2) Developers must be able to use presets to defined scaling settings.
	- Create preset configuration way in to server-side
	- Generate JSON setting from the presets in to website, where riclib.js reads the configuration
		and can use those presets automatically.
	- Add attribute hook in to <img> elements to define preset (i.e. ric-preset)
	- Improve riclib.js to detect image container size, viewport width etc. to intelligently define proper size fo the images.
	- Add placeholder svg image which is shown while actual image is been fetched, is been processed, is been failed to load (404)

* 3) Add best way to utilize <picture> -element, srcset attibute, sizes attribute, <source> -element type attribute (format) and art-direction (READ w3c).

* 4) Improve server to detect current jpeg quality setting, resolution, etc. to decide whether re-compress the image 

* 5) riclib.js must be able to score client device for performance to optimize used images. Used network connection must be detected to be able to provide optimal image. (Header: X-RIC-Client-Connection: [number])

# Some stuff
```
var screens = {
	lg: {
		"min-width": "1200px"
	},
	md: {
		"min-width": "900px",
		"max-width": "1199px"
	},
	sm: {
		"min-width": "480px",
		"max-width": "899px"
	},
	xs: {
		"max-width": "479px"
	}
};

var presets = {
	img-splitbanner": {
		lg:"578x326",
		md:"458x258",
		sm:"348x196",
		xs:"460x259"
	},
	img-product-thumb": {
		lg:"150x150",
		md:"150x150",
		sm:"150x150",
		xs:"150x150"
	}
}


<ric-image preset="img-product-thumb" src="http://path.to.jpg"></ric-image>

<img src="" srcset="">

<picture>
	<source srcset="http://path.to.jpg?ricPreset=578x392@1x 1x 
			http://path.to.jpg?ricPreset=578x392@2x 2x" media="(min-width: 1200px)">
	<source srcset="http://path.to.jpg?ricPreset=458x258@1x 1x 
			http://path.to.jpg?ricPreset=458x258@2x 2x" media="(min-width: 900px max-width: 1199px)">
 
</picture>
```

# Some refs
https://github.com/fawick/speedtest-resize

https://github.com/disintegration/imaging

https://willnorris.com/go/imageproxy

https://github.com/BBC-News/Imager.js

https://html.spec.whatwg.org/multipage/embedded-content.html#embedded-content

http://blog.chromium.org/2016/05/saving-data-with-googles-pagespeed.html

http://www.clickonf5.org/15311/image-compression-tools-wordpress-plugins/