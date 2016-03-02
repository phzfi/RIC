package fi.phz.cib;


import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.context.embedded.EmbeddedServletContainerCustomizer;
import org.springframework.context.annotation.Bean;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestController;

import fi.phz.cib.ImageHandler.ExtraInfo;

@RestController
public class WebController {

	@Bean
	public EmbeddedServletContainerCustomizer containerCustomizer() {
		return (container -> {
			container.setPort(8005);
		});
	}

	@Autowired
	private ImageHandler handler;
	
	@RequestMapping(path="/{id}.{ext}")
	@ResponseBody
	public final byte[] getImage(
			@PathVariable("id") String imageId,
			@PathVariable("ext") String extension,
			@RequestParam(required=false, name="width") Integer width,
			@RequestParam(required=false, name="height") Integer height) {

		if ((width == null) && (height == null)) {
			return handler.loadOriginal(imageId);
		}
		ExtraInfo info = handler.getInfo(imageId);

		if ((width == null) || (width.intValue() > info.getWidth())) {
			width = info.getWidth();
		}
		
		if ((height == null) || (height.intValue() > info.getHeight())) {
			height = info.getHeight();
		}

		return handler.loadResized(imageId, width, height);
	}
}
