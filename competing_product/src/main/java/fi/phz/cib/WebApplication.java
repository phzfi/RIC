package fi.phz.cib;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.ApplicationContext;

@SpringBootApplication
public class WebApplication {

	public static ApplicationContext context;
	
	public static void main(String[] args) {
		context = SpringApplication.run(WebApplication.class, args);
	}
}
