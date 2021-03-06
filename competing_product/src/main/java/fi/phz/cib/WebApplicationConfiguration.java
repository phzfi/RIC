package fi.phz.cib;

import java.io.IOException;
import java.io.InputStream;
import java.util.Properties;
import java.util.concurrent.TimeUnit;

import org.apache.log4j.Logger;
import org.springframework.cache.annotation.EnableCaching;
import org.springframework.cache.guava.GuavaCacheManager;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.PropertySource;

import com.google.common.cache.CacheBuilder;
import com.google.common.cache.Weigher;

@Configuration
@EnableCaching
@PropertySource("/log4j.properties")
public class WebApplicationConfiguration {

	private static final Logger logger = Logger.getLogger(WebApplicationConfiguration.class);
	private static final Class<?> clazz = WebApplicationConfiguration.class;
	private static final Long KILOBYTE = 1024L;
	private static final Long MEGABYTE = 1024L * KILOBYTE;
	//private static final Long GIGABYTE = 1024L * MEGABYTE;
	
	public static class ByteSizeWeight implements Weigher<String, byte[]> {
		@Override
		public final int weigh(final String key, final byte[] value) {
			return value.length;
		}
	}

	@Bean
	public Properties getProperties() {
		Properties result = new Properties();
		try (InputStream is = clazz.getResourceAsStream("/web.properties")) {
			result.load(is);
		} catch (IOException e) {
			logger.error(e.getMessage(), e);
			return null;
		}
		logger.info("Settings: " + result.toString());
		return result;
	}

	public Long getMaxMemory() {
		final Properties prop = getProperties();
		final String maxMemory = prop.getProperty("memory.maxsize", "1024").toUpperCase();
		final Long memoryConstraint;
		try {
			memoryConstraint = Long.parseLong(maxMemory) * MEGABYTE;
		} catch (NumberFormatException e) {
			logger.error(e.getMessage(), e);
			return 1024L;
		}
		return memoryConstraint;
	}

	public Integer getMaxConcurrentResizes() {
		final Properties prop = getProperties();
		final String maxResizes = prop.getProperty("imagick.maxresizes", "10");
		final Integer resizeConstraint;
		try {
			resizeConstraint = Integer.parseInt(maxResizes);
		} catch (NumberFormatException e) {
			logger.error(e.getMessage(), e);
			return 10;
		}
		return resizeConstraint;
	}

	public String getRoot() {
		final Properties prop = getProperties();
		return prop.getProperty("server.root", "/var/www");
	}

	@Bean
	public GuavaCacheManager getCacheManager() {
		final GuavaCacheManager gcm = new GuavaCacheManager();
		final CacheBuilder<Object, Object> builder = CacheBuilder.newBuilder();
		gcm.setAllowNullValues(false);
		builder.concurrencyLevel(16);
		builder.expireAfterAccess(15, TimeUnit.MINUTES);
		builder.initialCapacity(256);
		builder.maximumWeight(getMaxMemory());
		builder.weigher(new ByteSizeWeight());
		gcm.setCacheBuilder(builder);
		return gcm;
	}

	@Bean
	public ImageHandler getImageHandler() {
		final ImageHandler h = new ImageHandler();
		h.setRoot(getRoot());
		h.setConcurrentResizeLimit(getMaxConcurrentResizes());
		return h;
	}
	
}
