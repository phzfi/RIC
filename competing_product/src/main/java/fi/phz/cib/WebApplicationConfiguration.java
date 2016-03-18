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
		String maxmem = prop.getProperty("memory.maxsize", "1024").toUpperCase();
		Long memoryConstraint = null;
		try {
			memoryConstraint = Long.parseLong(maxmem) * MEGABYTE;
		} catch (NumberFormatException e) {
			logger.error(e.getMessage(), e);
			return null;
		}
		return memoryConstraint;
	}

	public String getRoot() {
		final Properties prop = getProperties();
		return prop.getProperty("server.root", "/var/www");
	}

	@Bean
	public GuavaCacheManager getCacheManager() {
		GuavaCacheManager gcm = new GuavaCacheManager();
		gcm.setAllowNullValues(false);
		CacheBuilder<Object, Object> builder = CacheBuilder.newBuilder();
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
		ImageHandler h = new ImageHandler();
		h.setRoot(getRoot());
		return h;
	}
	
}
