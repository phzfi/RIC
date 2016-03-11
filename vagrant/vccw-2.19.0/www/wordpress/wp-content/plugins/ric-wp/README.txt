=== Plugin Name ===
Contributors: fubla, LKeronen
Donate link: https://phz.fi/
Tags: comments, spam
Requires at least: 3.0.1
Tested up to: 3.4
Stable tag: 4.3
License: GPLv2 or later
License URI: http://www.gnu.org/licenses/gpl-2.0.html

== Description ==

This plugin integrates RIC image server with your wordpress site allowing for faster and more responsive media-rich web sites.

== Installation ==

This section describes how to install the plugin and get it working.

1. Upload `ric-wp` folder to the `/wp-content/plugins/` directory
2. Activate the plugin through the 'Plugins' menu in WordPress

== Instructions ==

Define RIC image server url in wordpress admin site for ric-wp plugin settings.
When creating your image divs, use the following format: 

<div class="ricdiv" style="height:[your_image_height]px; width:[your_image_width]px;">
	<img data-id="[your_image_id]" data-mode="nofit" class="ricimg"></img>
</div>

Where id is given with or without the image format ending. The height and width are the dimensions of the image queried from the RIC image server.
Resize mode is optional and aspect ratio is preserved by default.
Width and height can be left blank, in which case the image is fit to the content width.

== Changelog ==

= 1.0 =

Initial release
