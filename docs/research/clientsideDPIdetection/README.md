# Client DPI/PPI Detecetion

## Device pixel ratio

It is not possible to accurately calculate the device's pixels-per-inch or dots-per-inch ratio. With JavaScript it is possible to query for device pixel ratio and logical screen resolution. The physical resolution can be approximated by multiplying the logical resolution by the device pixel ratio.

We gathered data for dozen or so common devices (phones, tablets, laptops and desktop computers). Typical values for device pixel ratio are 1, 1.5, and 2 but as high as 3 and 4 are in use in high-end phones. On Lumia 930 we encountered device pixel ratio of 2.68.

Notably some common Samsung devices report the logical screen resolution incorrectly depending on the browser in use. For example Samsung Galaxy S III has physical resolution of 720x1280. Different browsers reported either 720x1280 or 360x640 logical resolution even though both reported device pixel ratio of 2.

Also some browsers report logical resolution for full display area, others may report only the visible viewport area.

## Estimating DPI/PPI

Different device categories seem to have following qualities:

* Phones: Smaller logical dimension within range of 320-400. Physical screen size in range from 3.5" to 5.5".

* Tablets: Smaller logical dimension within range of 600-800. Physical screen size in range from 7" to 10".

* Computers: Smaller logical dimension is 900 or higher. Physical screen size varies a lot but in general is 12" or higher.

There are some devices that might not fit this broad categorization but mostly devices do map to one of these. User Agent string could be used on server side to improve the accuracy.

The device DPI/PPI can now be approximated by using the abovementioned device category and the reported device pixel ratio.

## Other techniques

We also investigated srcset tag of img element but that wasn't widely supported. HTML picture element is also in experimental stage. Both these techniques do not provide DPI/PPI value but rather instruct the browser to select the best-fit version of the content.
