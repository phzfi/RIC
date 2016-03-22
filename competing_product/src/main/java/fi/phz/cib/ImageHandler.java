package fi.phz.cib;

import java.util.concurrent.ConcurrentHashMap;

import org.apache.log4j.Logger;
import org.springframework.cache.annotation.CacheConfig;
import org.springframework.cache.annotation.Cacheable;

import magick.ImageInfo;
import magick.Magick;
import magick.MagickException;
import magick.MagickImage;


@CacheConfig
public class ImageHandler {

	static Magick magick = new Magick();
	private static Logger logger = Logger.getLogger(ImageHandler.class);

	public static class ExtraInfo extends ImageInfo {
		
		Integer width;
		Integer height;

		public ExtraInfo() throws MagickException {
			super();
		}

		public ExtraInfo(String fileName) throws MagickException {
			super(fileName);
		}

		public final Integer getWidth() {
			return width;
		}

		public final Integer getHeight() {
			return height;
		}

		public final void setWidth(final Integer width) {
			this.width = width;
		}

		public final void setHeight(final Integer height) {
			this.height = height;
		}
	}

	private String root;
	private ConcurrentHashMap<String, ExtraInfo> infos;
	
	public ImageHandler() {
		this.infos = new ConcurrentHashMap<String, ExtraInfo>();
		logger.info("Initialized Competing Image Bank cache");
	}

	public void setRoot(String root) {
		this.root = root;
		if (!root.endsWith("/")) {
			this.root += "/";
		}
	}

	public ExtraInfo getInfo(String imageId) {
		ExtraInfo info = infos.get(imageId);
		if (info != null) {
			return info;
		}
		final String target = root + imageId + ".jpg";
		try {
			info = new ExtraInfo(target);
			MagickImage image = new MagickImage(info, true);
			info.setWidth(image.getDimension().width);
			info.setHeight(image.getDimension().height);
		} catch (MagickException e) {
			logger.error(e.getMessage(), e);
			return null;
		}
		infos.put(imageId, info);
		return info;
	}

	protected MagickImage toMagick(final ExtraInfo info) {
		MagickImage image = new MagickImage();
		try {
			image.readImage(info);
			info.setWidth(image.getDimension().width);
			info.setHeight(image.getDimension().height);
		} catch (MagickException e) {
			logger.error(e.getMessage(), e);
			return null;
		}
		return image;
	}

	@Cacheable(cacheNames={"images"}, key="#a0")
	public byte[] loadOriginal(final String imageId) {
		ExtraInfo info = getInfo(imageId);
		if (info == null) {
			return null;
		}
		MagickImage img = toMagick(info);
		if (img == null) {
			return null;
		}
		return img.imageToBlob(info);
	}

	@Cacheable(cacheNames={"images"}, key="#a1+'-'+#a2+'-'+#a3")
	public byte[] loadResized(String imageId, int width, int height) {
		ExtraInfo info = getInfo(imageId);
		if (info == null) {
			return null;
		}
		MagickImage img = toMagick(info);
		if (img == null) {
			return null;
		}
		try {
			img = img.scaleImage(width, height);
		} catch (MagickException e) {
			logger.error(e.getMessage(), e);
			return null;
		}
		return img.imageToBlob(info);
	}
}
